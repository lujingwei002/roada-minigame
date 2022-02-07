package gat

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/roada-go/gat/crypto"
	"github.com/roada-go/gat/log"
	"github.com/roada-go/gat/message"
	"github.com/roada-go/gat/packet"
)

const (
	_ int32 = iota
	statusStart
	statusHandshake
	statusWorking
	statusClosed
)

const (
	agentWriteBacklog = 16
)

type Agent struct {
	session *Session
	conn    net.Conn
	state   int32
	chDie   chan struct{}
	chSend  chan pendingMessage
	lastAt  int64
	decoder *packet.Decoder
	gate    *Gate
}

type pendingMessage struct {
	typ     message.Type
	route   string
	mid     uint64
	payload interface{}
}

type handShakeRequest struct {
	Sys struct {
		Token   string `json:"token"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"sys"`
}

// Create new agent instance
func newAgent(gate *Gate, conn net.Conn) *Agent {
	self := &Agent{
		gate:    gate,
		conn:    conn,
		state:   statusStart,
		chDie:   make(chan struct{}),
		lastAt:  time.Now().Unix(),
		chSend:  make(chan pendingMessage, agentWriteBacklog),
		decoder: packet.NewDecoder(),
	}
	self.session = newSession(self)
	return self
}

func (a *Agent) send(m pendingMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrBrokenPipe
		}
	}()
	a.chSend <- m
	return
}

func (agent *Agent) syncSend(typ message.Type, mid uint64, route string, v interface{}) error {
	payload, err := agent.serialize(v)
	if err != nil {
		switch typ {
		case message.Push:
			log.Printf("[agent] push: %s error: %s\n", route, err.Error())
		case message.Response:
			log.Printf("[agent] response message(id: %d) error: %s\n", mid, err.Error())
		default:
			// expect
		}
		return err
	}
	m := &message.Message{
		Type:  typ,
		Data:  payload,
		Route: route,
		ID:    mid,
	}
	em, err := m.Encode()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	p, err := packet.Encode(packet.Data, em)
	if err != nil {
		log.Println(err)
		return err
	}
	if _, err := agent.conn.Write(p); err != nil {
		log.Printf("[agent] conn write failed, error:%s\n", err.Error())
		return err
	}
	return nil
}

func (self *Agent) push(route string, v interface{}) error {
	if self.status() == statusClosed {
		return ErrBrokenPipe
	}
	if len(self.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}
	if self.gate.debug {
		switch d := v.(type) {
		case []byte:
			log.Printf("[agent] type=push, sessionid=%d, uid=%d, route=%s, data=%dbytes\n",
				self.session.ID(), self.session.UID(), route, len(d))
		default:
			log.Printf("[agent] type=push, sessionid=%d, uid=%d, route=%s, data=%+v\n",
				self.session.ID(), self.session.UID(), route, v)
		}
	}
	return self.send(pendingMessage{typ: message.Push, route: route, payload: v})
}

func (self *Agent) syncPush(route string, v interface{}) error {
	if self.status() == statusClosed {
		return ErrBrokenPipe
	}
	if len(self.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}
	if self.gate.debug {
		switch d := v.(type) {
		case []byte:
			log.Printf("[agent] type=push, sessionid=%d, uid=%d, route=%s, data=%dbytes\n",
				self.session.ID(), self.session.UID(), route, len(d))
		default:
			log.Printf("[agent] type=push, sessionid=%d, uid=%d, route=%s, data=%+v\n",
				self.session.ID(), self.session.UID(), route, v)
		}
	}
	return self.syncSend(message.Push, 0, route, v)
}

func (self *Agent) response(mid uint64, v interface{}) error {
	if self.status() == statusClosed {
		return ErrBrokenPipe
	}
	if mid <= 0 {
		return ErrSessionOnNotify
	}
	if len(self.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}
	if self.gate.debug {
		switch d := v.(type) {
		case []byte:
			log.Printf("[agent] type=response, sessionid=%d, uid=%d, mid=%d, data=%dbytes\n",
				self.session.ID(), self.session.UID(), mid, len(d))
		default:
			log.Printf("[agent] type=response, sessionid=%d, uid=%d, mid=%d, data=%+v\n",
				self.session.ID(), self.session.UID(), mid, v)
		}
	}
	return self.send(pendingMessage{typ: message.Response, mid: mid, payload: v})
}

func (self *Agent) kick(reason string) error {
	data, err := packet.Encode(packet.Kick, []byte(reason))
	if err != nil {
		return err
	}
	if _, err := self.conn.Write(data); err != nil {
		log.Printf("[agent] kick, conn write failed, error:%s\n", err.Error())
		return err
	}
	return nil
}

func (self *Agent) close() error {
	if self.status() == statusClosed {
		return ErrCloseClosedSession
	}
	self.setStatus(statusClosed)
	if self.gate.debug {
		log.Printf("[agent] session closed, sessionid=%d, uid=%d, ip=%s\n",
			self.session.ID(), self.session.UID(), self.conn.RemoteAddr())
	}
	select {
	case <-self.chDie:
		// expect
	default:
		close(self.chDie)
	}
	return self.conn.Close()
}

func (self *Agent) remoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *Agent) String() string {
	return fmt.Sprintf("[agent] remote=%s, lastTime=%d", self.conn.RemoteAddr().String(), atomic.LoadInt64(&self.lastAt))
}

func (self *Agent) status() int32 {
	return atomic.LoadInt32(&self.state)
}

func (self *Agent) setStatus(state int32) {
	atomic.StoreInt32(&self.state, state)
}

func (self *Agent) serialize(v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}
	data, err := self.gate.serializer.Marshal(v)
	if err != nil {
		return nil, err
	}
	var session = self.session
	if session.getSecret() != "" {
		data, err = crypto.DesEncrypt(data, session.getSecret())
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (self *Agent) write() {
	ticker := time.NewTicker(self.gate.heartbeat)
	chWrite := make(chan []byte, agentWriteBacklog)
	// clean func
	defer func() {
		ticker.Stop()
		close(self.chSend)
		close(chWrite)
		self.close()
		if self.gate.debug {
			log.Printf("[agent] session write goroutine exit, sessionid=%d, uid=%d\n",
				self.session.ID(), self.session.UID())
		}
	}()
	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(-2 * self.gate.heartbeat).Unix()
			if atomic.LoadInt64(&self.lastAt) < deadline {
				log.Printf("[agent] session heartbeat timeout, lastTime=%d, deadline=%d\n",
					atomic.LoadInt64(&self.lastAt), deadline)
				return
			}
			chWrite <- self.gate.heartbeatPacket
		case data := <-chWrite:
			if _, err := self.conn.Write(data); err != nil {
				log.Printf("[agent] conn write failed, error:%s\n", err.Error())
				return
			}
		case data := <-self.chSend:
			payload, err := self.serialize(data.payload)
			if err != nil {
				switch data.typ {
				case message.Push:
					log.Printf("[agent] push: %s error: %s\n", data.route, err.Error())
				case message.Response:
					log.Printf("[agent] response message(id: %d) error: %s\n", data.mid, err.Error())
				default:
					// expect
				}
				break
			}
			m := &message.Message{
				Type:  data.typ,
				Data:  payload,
				Route: data.route,
				ID:    data.mid,
			}
			em, err := m.Encode()
			if err != nil {
				log.Println(err.Error())
				break
			}
			p, err := packet.Encode(packet.Data, em)
			if err != nil {
				log.Println(err)
				break
			}
			chWrite <- p
		case <-self.chDie:
			return
		case <-self.gate.chDie:
			return
		}
	}
}

func (agent *Agent) onSessionOpen() {
	for _, middleware := range agent.gate.middlewareArr {
		middleware.OnSessionOpen(agent.session)
	}
	if agent.gate.handler != nil {
		agent.gate.handler.OnSessionOpen(agent.session)
	}
}

func (agent *Agent) onSessionClose() {
	for _, middleware := range agent.gate.middlewareArr {
		middleware.OnSessionClose(agent.session)
	}
	if agent.gate.handler != nil {
		agent.gate.handler.OnSessionClose(agent.session)
	}
}

func (self *Agent) handle() {
	self.onSessionOpen()
	//开启读协程
	go self.write()
	if self.gate.debug {
		log.Println(fmt.Sprintf("[agent] new session established: %s", self.String()))
	}
	defer func() {
		self.close()
		self.gate.sessionClosed(self.session)
		self.onSessionClose()
		if self.gate.debug {
			log.Printf("[agent] session read goroutine exit, sessionid=%d, uid=%d\n",
				self.session.ID(), self.session.UID())
		}
	}()
	buf := make([]byte, 2048)
	for {
		n, err := self.conn.Read(buf)
		if err != nil {
			log.Printf("[agent] read message error: %s, session will be closed immediately, sessionid=%d\n",
				err.Error(), self.session.ID())
			return
		}
		packets, err := self.decoder.Decode(buf[:n])
		if err != nil {
			log.Println(err.Error())
			return
		}
		if len(packets) < 1 {
			continue
		}
		for i := range packets {
			if err := self.processPacket(packets[i]); err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

func (self *Agent) processPacket(p *packet.Packet) error {
	switch p.Type {
	case packet.Handshake:
		msg := &handShakeRequest{}
		err := json.Unmarshal(p.Data, msg)
		if err != nil {
			return err
		}
		if self.gate.rsaPrivateKey != "" {
			token, err := crypto.RsaDecryptWithSha1Base64(msg.Sys.Token, self.gate.rsaPrivateKey)
			if err != nil {
				return err
			}
			self.session.setSecret(token)
		}
		if err := self.gate.handshakeValidator(p.Data); err != nil {
			return err
		}
		data, err := json.Marshal(map[string]interface{}{
			"code": 200,
			"sys": map[string]interface{}{
				"heartbeat": self.gate.heartbeat.Seconds(),
				"session":   self.session.ID(),
			},
		})
		if err != nil {
			return err
		}
		handsharkResponse, err := packet.Encode(packet.Handshake, data)
		if err != nil {
			return err
		}
		if _, err := self.conn.Write(handsharkResponse); err != nil {
			return err
		}
		self.setStatus(statusHandshake)
		if self.gate.debug {
			log.Printf("[agent] session handshake sessionid=%d, remote=%s, secret=%s\n",
				self.session.ID(), self.conn.RemoteAddr(), self.session.getSecret())
		}
	case packet.HandshakeAck:
		self.setStatus(statusWorking)
		if self.gate.debug {
			log.Printf("[agent] receive handshake ack sessionid=%d, remote=%s\n",
				self.session.ID(), self.conn.RemoteAddr())
		}
	case packet.Data:
		if self.status() < statusWorking {
			return fmt.Errorf("[agent] receive data on socket which not yet ack, session will be closed immediately, sessionid=%d, remote=%s",
				self.session.ID(), self.conn.RemoteAddr().String())
		}
		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		self.processMessage(msg)
	case packet.Heartbeat:
		// expected
	}
	self.lastAt = time.Now().Unix()
	return nil
}

func (self *Agent) processMessage(msg *message.Message) {
	var mid uint64
	switch msg.Type {
	case message.Request:
		mid = msg.ID
	case message.Notify:
		mid = 0
	default:
		log.Println("[agent] invalid message type: " + msg.Type.String())
		return
	}
	var session = self.session
	var payload = msg.Data
	var err error
	if session.getSecret() != "" {
		payload, err = crypto.DesDecrypt(payload, session.getSecret())
		if err != nil {
			log.Printf("[agent] des decrypt failed, error:%s, payload:(%v)\n", err.Error(), payload)
			return
		}
	}
	index := strings.LastIndex(msg.Route, ".")
	if index < 0 {
		log.Printf("[agent] invalid route, route:%s\n", msg.Route)
		return
	}
	if self.gate.handler == nil {
		log.Printf("[agent] handler not found, route:%s\n", msg.Route)
		return
	}
	r := &Request{
		Session: session,
		Route:   msg.Route,
		Payload: payload,
		mid:     mid,
	}
	//task := func() {
	for _, middleware := range self.gate.middlewareArr {
		middleware.ServeMessage(r)
	}
	self.gate.handler.ServeMessage(r)
	//}
	//self.gate.scheduler.PushTask(task)
}
