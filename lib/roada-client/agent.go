package cli

import (
	"crypto/rand"
	"encoding/base32"
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
	client         *Client
	session        *Session
	conn           net.Conn
	state          int32
	chDie          chan struct{}
	chSend         chan pendingMessage
	lastAt         int64
	decoder        *packet.Decoder
	mid            uint64
	responseRouter map[uint64]string
	chWrite        chan []byte
}

type handShakeResponse struct {
	Sys struct {
		Heartbeat int   `json:"heartbeat"`
		Session   int64 `json:"session"`
	} `json:"sys"`
	Code int `json:"code"`
}
type pendingMessage struct {
	typ     message.Type
	route   string
	mid     uint64
	payload interface{}
}

func newAgent(client *Client, conn net.Conn) *Agent {
	self := &Agent{
		client:         client,
		conn:           conn,
		decoder:        packet.NewDecoder(),
		state:          statusStart,
		chDie:          make(chan struct{}),
		lastAt:         time.Now().Unix(),
		chSend:         make(chan pendingMessage, agentWriteBacklog),
		responseRouter: make(map[uint64]string),
		chWrite:        make(chan []byte, agentWriteBacklog),
	}
	session := newSession(self)
	self.session = session
	return self
}

func (self *Agent) send(m pendingMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = ErrBrokenPipe
		}
	}()
	self.chSend <- m
	return
}

func (self *Agent) request(route string, v interface{}) error {
	if self.status() == statusClosed {
		return ErrBrokenPipe
	}
	if len(self.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}
	mid := atomic.AddUint64(&self.mid, 1)
	if self.client.debug {
		switch d := v.(type) {
		case []byte:
			log.Println(fmt.Sprintf("[agent] type=request, sessionid=%d, uid=%d, mid=%d, route=%s, data=%dbytes",
				self.session.ID(), self.session.UID(), mid, route, len(d)))
		default:
			log.Println(fmt.Sprintf("[agent] type=request, sessionid=%d, uid=%d, mid=%d, route=%s, data=%+v",
				self.session.ID(), self.session.UID(), mid, route, v))
		}
	}
	self.responseRouter[mid] = route
	return self.send(pendingMessage{typ: message.Request, mid: mid, route: route, payload: v})
}

func (self *Agent) notify(route string, v interface{}) error {
	if self.status() == statusClosed {
		return ErrBrokenPipe
	}
	if len(self.chSend) >= agentWriteBacklog {
		return ErrBufferExceed
	}
	if self.client.debug {
		switch d := v.(type) {
		case []byte:
			log.Println(fmt.Sprintf("Type=Notify, ID=%d, UID=%d, Route=%s, Data=%dbytes",
				self.session.ID(), self.session.UID(), route, len(d)))
		default:
			log.Println(fmt.Sprintf("Type=Notify, ID=%d, UID=%d, Route=%s, Data=%+v",
				self.session.ID(), self.session.UID(), route, v))
		}
	}
	return self.send(pendingMessage{typ: message.Notify, route: route, payload: v})
}

func (self *Agent) close() error {
	if self.status() == statusClosed {
		return ErrCloseClosedSession
	}
	self.setStatus(statusClosed)
	// prevent closing closed channel
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

// String, implementation for Stringer interface
func (self *Agent) String() string {
	return fmt.Sprintf("remote=%s, LastTime=%d", self.conn.RemoteAddr().String(), atomic.LoadInt64(&self.lastAt))
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
	data, err := self.client.serializer.Marshal(v)
	if err != nil {
		return nil, err
	}
	var session = self.session
	if session.getSecretKey() != "" {
		data, err = crypto.DesEncrypt(data, session.getSecretKey())
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (self *Agent) sendHandShake() error {
	tokenByte := make([]byte, 8)
	_, err := rand.Read(tokenByte)
	if err != nil {
		return err
	}
	token := ""
	if self.client.rsaPublicKey != "" {
		token = base32.StdEncoding.EncodeToString(tokenByte)[:8]
		self.session.setSecretKey(token)
		token, err = crypto.RsaEncryptWithSha1Base64(token, self.client.rsaPublicKey)
		if err != nil {
			return err
		}
	}
	data, err := json.Marshal(map[string]interface{}{
		"sys": map[string]interface{}{
			"type":    "go-websocket",
			"version": "0.0.1",
			"token":   token,
		},
		"user": map[string]interface{}{},
	})
	if err != nil {
		return err
	}
	handsharkRequest, err := packet.Encode(packet.Handshake, data)
	if err != nil {
		return err
	}
	if _, err := self.conn.Write(handsharkRequest); err != nil {
		return err
	}
	if self.client.debug {
		log.Println(fmt.Sprintf("[agent] session handshake sessionid=%d, remote=%s, token=%s", self.session.ID(), self.conn.RemoteAddr(), token))
	}
	return nil
}

func (self *Agent) write() {
	ticker := time.NewTicker(self.client.heartbeat)
	defer func() {
		ticker.Stop()
		close(self.chSend)
		close(self.chWrite)
		self.close()
		if self.client.debug {
			log.Println(fmt.Sprintf("[agent] session write goroutine exit, sessionid=%d, uid=%d", self.session.ID(), self.session.UID()))
		}
	}()
	err := self.sendHandShake()
	if err != nil {
		log.Println(fmt.Sprintf("[agent] client sendHandShake error: %s", err.Error()))
		return
	}
	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(-2 * self.client.heartbeat).Unix()
			if atomic.LoadInt64(&self.lastAt) < deadline {
				log.Println(fmt.Sprintf("[agent] session heartbeat timeout, sessionid=%d, uid=%d, lastTime=%d, deadline=%d",
					self.session.ID(), self.session.UID(), atomic.LoadInt64(&self.lastAt), deadline))
				return
			}
		case data := <-self.chWrite:
			if _, err := self.conn.Write(data); err != nil {
				log.Println(fmt.Sprintf("[agent] conn write failed, error:%s", err.Error()))
				return
			}
		case data := <-self.chSend:
			payload, err := self.serialize(data.payload)
			if err != nil {
				switch data.typ {
				case message.Notify:
					log.Println(fmt.Sprintf("[agent] nontify: %s error: %s", data.route, err.Error()))
				case message.Request:
					log.Println(fmt.Sprintf("[agent] request message(id: %d) error: %s", data.mid, err.Error()))
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
			self.chWrite <- p
		case <-self.chDie:
			return
		case <-self.client.chDie:
			return
		}
	}
}

func (self *Agent) handle() {
	go self.write()
	defer func() {
		self.close()
		self.client.sessionClosed(self.session)
		self.client.lifetime.close(self.session)
		if self.client.debug {
			log.Println(fmt.Sprintf("[agent] session read goroutine exit, sessionid=%d, uid=%d", self.session.ID(), self.session.UID()))
		}
	}()
	buf := make([]byte, 2048)
	for {
		n, err := self.conn.Read(buf)
		if err != nil {
			log.Println(fmt.Sprintf("[agent] read message error: %s, session will be closed immediately, sessionid=%d", err.Error(), self.session.ID()))
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
		msg := &handShakeResponse{}
		err := json.Unmarshal(p.Data, msg)
		if err != nil {
			return err
		}
		if msg.Code != 200 {
			return ErrHandShake
		}
		self.setStatus(statusHandshake)
		data, err := json.Marshal(map[string]interface{}{})
		if err != nil {
			return err
		}
		handsharkAck, err := packet.Encode(packet.HandshakeAck, data)
		if err != nil {
			return err
		}

		self.session.id = msg.Sys.Session
		if self.client.debug {
			log.Println(fmt.Sprintf("session handshake sessionid=%d, remote=%s packget=%+v",
				self.session.ID(), self.conn.RemoteAddr(), msg))
		}
		self.chWrite <- handsharkAck
		self.setStatus(statusWorking)
		self.client.lifetime.open(self.session)
	case packet.Data:
		if self.status() < statusWorking {
			return fmt.Errorf("receive data on socket which not yet ACK, session will be closed immediately, sessionid=%d, remote=%s, status=%d %v",
				self.session.ID(), self.conn.RemoteAddr().String(), self.status())
		}
		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		self.processMessage(msg)
	case packet.Heartbeat:
		self.chWrite <- self.client.heartbeatPacket
	}
	self.lastAt = time.Now().Unix()
	return nil
}

func (self *Agent) processMessage(msg *message.Message) {
	var mid uint64
	var route = msg.Route
	switch msg.Type {
	case message.Response:
		mid = msg.ID
		if s, ok := self.responseRouter[mid]; ok {
			route = s
			delete(self.responseRouter, mid)
		}
	case message.Push:
		mid = 0
	default:
		log.Println("Invalid message type: " + msg.Type.String())
		return
	}
	var session = self.session
	var payload = msg.Data
	var err error
	if session.getSecretKey() != "" {
		payload, err = crypto.DesDecrypt(payload, session.getSecretKey())
		if err != nil {
			log.Println(fmt.Sprintf("crypto.DesDecrypt failed: %+v (%v)", err, payload))
			return
		}
	}
	index := strings.LastIndex(route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("agent invalid route, route:%s", route))
		return
	}
	handleFunc := self.client.handleFunc
	if handleFunc == nil {
		log.Println(fmt.Sprintf("agent handler not found, route:%s", route))
		return
	}
	r := &Request{
		Session: session,
		Route:   route,
		mid:     mid,
		Payload: payload,
	}
	handleFunc(r)
}
