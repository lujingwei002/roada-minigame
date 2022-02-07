package gat

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/roada-go/gat/log"
	"github.com/roada-go/gat/packet"
	"github.com/roada-go/gat/serialize"
	"github.com/roada-go/gat/serialize/protobuf"
)

const (
	_ int32 = iota
	gateStatusStart
	gateStatusMaintain
	gateStatusWorking
	gateStatusClosed
)

type Gate struct {
	ClientAddr string
	Host       string
	Port       int32
	Wd         string

	isWebsocket        bool
	tslCertificate     string
	tslKey             string
	handshakeValidator func([]byte) error
	serializer         serialize.Serializer
	heartbeat          time.Duration
	checkOrigin        func(*http.Request) bool
	debug              bool
	wsPath             string
	rsaPrivateKey      string

	state           int32
	name            string
	startAt         time.Time
	scheduler       *Scheduler
	chDie           chan bool
	mu              sync.RWMutex
	sessions        map[int64]*Session
	handler         HandlerInterface
	middlewareArr   []MiddleWareInterface
	heartbeatPacket []byte
}

func NewGate() *Gate {
	gate := &Gate{
		state:              gateStatusStart,
		scheduler:          NewScheduler(),
		chDie:              make(chan bool),
		heartbeat:          30 * time.Second,
		debug:              false,
		checkOrigin:        func(_ *http.Request) bool { return true },
		handshakeValidator: func(_ []byte) error { return nil },
		serializer:         protobuf.NewSerializer(),
		middlewareArr:      make([]MiddleWareInterface, 0),
		sessions:           map[int64]*Session{},
	}
	return gate
}

func (gate *Gate) WithHandshakeValidator(fn func([]byte) error) {
	gate.handshakeValidator = fn
}

func (gate *Gate) WithSerializer(serializer serialize.Serializer) {
	gate.serializer = serializer
}

func (gate *Gate) WithDebugMode() {
	gate.debug = true
}

func (gate *Gate) WithCheckOriginFunc(fn func(*http.Request) bool) {
	gate.checkOrigin = fn
}

func (gate *Gate) WithHeartbeatInterval(d time.Duration) {
	gate.heartbeat = d
}

func (gate *Gate) WithDictionary(dict map[string]uint16) {

}

func (gate *Gate) WithWSPath(path string) {
	gate.wsPath = path
}

func (gate *Gate) WithTimerPrecision(precision time.Duration) {
	if precision < time.Millisecond {
		panic("time precision can not less than a Millisecond")
	}
	gate.scheduler.timerPrecision = precision
}

func (gate *Gate) WithIsWebsocket(enableWs bool) {
	gate.isWebsocket = enableWs
}

func (gate *Gate) WithTSLConfig(certificate, key string) {
	gate.tslCertificate = certificate
	gate.tslKey = key
}

func (gate *Gate) WithLogger(l log.Logger) {
	log.SetLogger(l)
}

func (gate *Gate) WithRSAPrivateKey(keyFile string) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	gate.rsaPrivateKey = string(data)
}

func (gate *Gate) Run(addr string) error {
	gate.name = strings.TrimLeft(filepath.Base(os.Args[0]), "/")
	gate.startAt = time.Now()
	if wd, err := os.Getwd(); err != nil {
		return err
	} else {
		gate.Wd, _ = filepath.Abs(wd)
	}
	addrPat := strings.SplitN(addr, ":", 2)
	if gate.isWebsocket {
		if len(gate.tslCertificate) != 0 {
			gate.Host = fmt.Sprintf("wss://%s", addrPat[0])
		} else {
			gate.Host = fmt.Sprintf("ws://%s", addrPat[0])
		}
	} else {
		gate.Host = addrPat[0]
	}
	if port, err := strconv.Atoi(addrPat[1]); err != nil {
		return fmt.Errorf("[gate] Run failed, addr format wrong, addr=%s",
			addr)
	} else {
		gate.Port = int32(port)
	}
	if len(addrPat) < 2 {
		return fmt.Errorf("[gate] Run failed, addr format wrong, addr=%s",
			addr)
	}
	gate.ClientAddr = fmt.Sprintf(":%d", gate.Port)
	err := gate.startup()
	if err != nil {
		return fmt.Errorf("[gate] startup failed ,error=%v",
			err)
	}
	log.Printf("[gate] %s Run, client address=%v",
		gate.name, gate.ClientAddr)
	go gate.scheduler.Sched()
	gate.setStatus(gateStatusWorking)
	return nil
}

/*func (self *Gate) NewTimer(interval time.Duration, fn TimerFunc) *Timer {
	return self.scheduler.timerManager.NewTimer(interval, fn)
}
func (self *Gate) NewAfterTimer(duration time.Duration, fn TimerFunc) *Timer {
	return self.scheduler.timerManager.NewAfterTimer(duration, fn)
}*/
func (gate *Gate) Handle(h HandlerInterface) {
	gate.handler = h
}

func (gate *Gate) UseMiddleware(m MiddleWareInterface) {
	gate.middlewareArr = append(gate.middlewareArr, m)
}

func (gate *Gate) Register(svr interface{}) (*Service, error) {
	s := newService(gate, svr)
	if err := s.extractHandler(); err != nil {
		return nil, err
	}
	return s, nil
}

func (gate *Gate) NewGroup(name string) *Group {
	return newGroup(gate, name)
}

func (gate *Gate) Maintain(m bool) {
	status := gate.status()
	if status != gateStatusMaintain && status != gateStatusWorking {
		return
	}
	if m {
		gate.setStatus(gateStatusMaintain)
	} else {
		gate.setStatus(gateStatusWorking)
	}
}

func (gate *Gate) Shutdown() {
	if gate.status() == gateStatusClosed {
		return
	}
	log.Println("[gate] Shutdown")
	gate.setStatus(gateStatusClosed)
	close(gate.chDie)
	gate.scheduler.Close()
}

func (gate *Gate) status() int32 {
	return atomic.LoadInt32(&gate.state)
}

func (gate *Gate) setStatus(state int32) {
	atomic.StoreInt32(&gate.state, state)
}

func (gate *Gate) startup() error {
	heartbeatPacket, err := packet.Encode(packet.Heartbeat, nil)
	if err != nil {
		return err
	}
	gate.heartbeatPacket = heartbeatPacket
	go func() {
		if gate.isWebsocket {
			if len(gate.tslCertificate) != 0 {
				gate.listenAndServeWSTLS()
			} else {
				gate.listenAndServeWS()
			}
		} else {
			gate.listenAndServe()
		}
	}()
	return nil
}

func (gate *Gate) Kick(reason string) {
	//断开已有的链接
	gate.mu.RLock()
	for _, s := range gate.sessions {
		s.Kick(reason)
	}
	gate.mu.RUnlock()
	now := time.Now().Unix()
	for {
		if len(gate.sessions) <= 0 {
			break
		}
		if time.Now().Unix()-now >= 120 {
			log.Println(fmt.Sprintf("[gate] some session not close, count=%d", len(gate.sessions)))
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (gate *Gate) listenAndServe() {
	listener, err := net.Listen("tcp", gate.ClientAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go gate.handle(conn)
	}
}

func (gate *Gate) listenAndServeWS() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gate.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gate.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gate.handleWS(conn)
	})
	if err := http.ListenAndServe(gate.ClientAddr, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func (gate *Gate) listenAndServeWSTLS() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     gate.checkOrigin,
	}
	http.HandleFunc("/"+strings.TrimPrefix(gate.wsPath, "/"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[gate] Upgrade failure, URI=%s, Error=%s", r.RequestURI, err.Error())
			return
		}
		gate.handleWS(conn)
	})
	if err := http.ListenAndServeTLS(gate.ClientAddr, gate.tslCertificate, gate.tslKey, nil); err != nil {
		log.Fatal(err.Error())
	}
}

func (gate *Gate) storeSession(s *Session) {
	gate.mu.Lock()
	gate.sessions[s.ID()] = s
	gate.mu.Unlock()
}

func (gate *Gate) findSession(sid int64) *Session {
	gate.mu.RLock()
	s := gate.sessions[sid]
	gate.mu.RUnlock()
	return s
}

func (gate *Gate) handle(conn net.Conn) {
	if gate.status() != gateStatusWorking {
		log.Println("[gate] gate is not running")
		conn.Close()
		return
	}
	agent := newAgent(gate, conn)
	gate.storeSession(agent.session)
	agent.handle()
}

func (gate *Gate) handleWS(conn *websocket.Conn) {
	if gate.status() != gateStatusWorking {
		log.Println("[gate] gate is not running")
		conn.Close()
		return
	}
	c, err := newWSConn(conn)
	if err != nil {
		log.Println(err)
		return
	}
	go gate.handle(c)
}

func (gate *Gate) sessionClosed(s *Session) error {
	gate.mu.Lock()
	delete(gate.sessions, s.ID())
	gate.mu.Unlock()
	return nil
}
