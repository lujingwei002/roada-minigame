package cli

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	//	"github.com/shark/roada-gate/agent"
	"github.com/roada-go/gat/log"
	"github.com/roada-go/gat/packet"
	"github.com/roada-go/gat/serialize"
	"github.com/roada-go/gat/serialize/protobuf"
)

type Client struct {
	isWebsocket        bool
	tslCertificate     string
	tslKey             string
	handshakeValidator func([]byte) error
	serializer         serialize.Serializer
	heartbeat          time.Duration
	debug              bool
	wsPath             string
	rsaPublicKey       string

	serverAddr      string
	scheduler       *Scheduler
	lifetime        *Lifetime
	chDie           chan bool
	mu              sync.RWMutex
	sessions        map[int64]*Session
	handleFunc      HandleFunc
	heartbeatPacket []byte
}

func NewClient() *Client {
	self := &Client{
		scheduler:          NewScheduler(),
		chDie:              make(chan bool),
		heartbeat:          30 * time.Second,
		debug:              false,
		handshakeValidator: func(_ []byte) error { return nil },
		lifetime:           newLifetime(),
		serializer:         protobuf.NewSerializer(),
		sessions:           map[int64]*Session{},
		rsaPublicKey:       "",
	}
	return self
}

func (self *Client) WithHandshakeValidator(fn func([]byte) error) {
	self.handshakeValidator = fn
}

func (self *Client) WithSerializer(serializer serialize.Serializer) {
	self.serializer = serializer
}

func (self *Client) WithDebugMode() {
	self.debug = true
}

func (self *Client) WithHeartbeatInterval(d time.Duration) {
	self.heartbeat = d
}

func (self *Client) WithDictionary(dict map[string]uint16) {

}

func (self *Client) WithWSPath(path string) {
	self.wsPath = path
}

func (self *Client) WithServerAdd(addr string) {
	self.serverAddr = addr
}

func (self *Client) WithTimerPrecision(precision time.Duration) {
	if precision < time.Millisecond {
		panic("time precision can not less than a Millisecond")
	}
	self.scheduler.timerPrecision = precision
}

func (self *Client) WithIsWebsocket(enableWs bool) {
	self.isWebsocket = enableWs
}

func (self *Client) WithTSLConfig(certificate, key string) {
	self.tslCertificate = certificate
	self.tslKey = key
}

func (self *Client) WithLogger(l log.Logger) {
	log.SetLogger(l)
}

func (self *Client) WithRSAPublicKey(keyFile string) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	self.rsaPublicKey = string(data)
}

func (self *Client) Run() {
	err := self.startup()
	if err != nil {
		log.Fatalf("gate startup failed: %v", err)
	}
	log.Println(fmt.Sprintf("[client] startup, server address: %v",
		self.serverAddr))
	go self.scheduler.Sched()
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-self.chDie:
		log.Println("[client] will shutdown in a few seconds")
	case s := <-sg:
		log.Println("[client] got signal", s)
	}
	log.Println("[client] stopping...")
	self.destory()
	self.scheduler.Close()
}

func (self *Client) HandleFunc(h HandleFunc) {
	self.handleFunc = h
}

func (self *Client) Register(svr interface{}) (*Service, error) {
	s := newService(self, svr)
	if err := s.extractHandler(); err != nil {
		return nil, err
	}
	return s, nil
}

func (self *Client) Dial() error {
	if self.isWebsocket {
		if len(self.tslCertificate) != 0 {
			//self.dialWSTLS()
		} else {
			return self.dialWS()
		}
	} else {
		return self.dial()
	}
	return nil
}

func (self *Client) NewGroup(name string) *Group {
	return newGroup(self, name)
}

func (self *Client) Shutdown() {
	close(self.chDie)
}

func (self *Client) OnSessionOpen(h LifetimeHandler) {
	self.lifetime.onOpen(h)
}

func (self *Client) OnSessionClose(h LifetimeHandler) {
	self.lifetime.onClose(h)
}

func (self *Client) startup() error {
	heartbeatPacket, err := packet.Encode(packet.Heartbeat, nil)
	if err != nil {
		return err
	}
	self.heartbeatPacket = heartbeatPacket
	return nil
}

func (self *Client) destory() {

}

func (self *Client) dial() error {
	conn, err := net.Dial("tcp", self.serverAddr)
	if err != nil {
		return err
	}
	go self.handle(conn)
	return nil
}

func (self *Client) dialWS() error {
	u := url.URL{Scheme: "ws", Host: self.serverAddr, Path: self.wsPath}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	self.handleWS(conn)
	return nil
}

func (self *Client) storeSession(s *Session) {
	self.mu.Lock()
	self.sessions[s.ID()] = s
	self.mu.Unlock()
}

func (self *Client) findSession(sid int64) *Session {
	self.mu.RLock()
	s := self.sessions[sid]
	self.mu.RUnlock()
	return s
}

func (self *Client) handle(conn net.Conn) {
	// create a client agent and startup write gorontine
	agent := newAgent(self, conn)
	self.storeSession(agent.session)
	agent.handle()
}

func (self *Client) handleWS(conn *websocket.Conn) {
	c, err := newWSConn(conn)
	if err != nil {
		log.Println(err)
		return
	}
	go self.handle(c)
}

func (self *Client) sessionClosed(s *Session) error {
	self.mu.Lock()
	delete(self.sessions, s.ID())
	self.mu.Unlock()
	return nil
}
