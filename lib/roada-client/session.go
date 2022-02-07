package cli

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrIllegalUID = errors.New("illegal uid")
)

type Session struct {
	sync.RWMutex
	id        int64
	uid       int64
	lastTime  int64
	agent     *Agent
	data      map[string]interface{}
	secretKey string
}

func newSession(agent *Agent) *Session {
	return &Session{
		id:        0,
		agent:     agent,
		data:      make(map[string]interface{}),
		lastTime:  time.Now().Unix(),
		secretKey: "",
	}
}

func (s *Session) getSecretKey() string {
	return s.secretKey
}

func (s *Session) setSecretKey(key string) {
	s.secretKey = key
}

func (s *Session) Notify(route string, v interface{}) error {
	return s.agent.notify(route, v)
}

func (s *Session) Request(route string, v interface{}) error {
	return s.agent.request(route, v)
}

func (s *Session) ID() int64 {
	return s.id
}

func (s *Session) UID() int64 {
	return atomic.LoadInt64(&s.uid)
}

func (s *Session) Bind(uid int64) error {
	if uid < 1 {
		return ErrIllegalUID
	}
	atomic.StoreInt64(&s.uid, uid)
	return nil
}

func (s *Session) Close() {
	s.agent.close()
}

func (s *Session) RemoteAddr() net.Addr {
	return s.agent.remoteAddr()
}

func (s *Session) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *Session) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *Session) HasKey(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, has := s.data[key]
	return has
}

func (s *Session) Int(key string) int {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(int)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int8(key string) int8 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(int8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int16(key string) int16 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(int16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int32(key string) int32 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(int32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Int64(key string) int64 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(int64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint(key string) uint {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(uint)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint8(key string) uint8 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(uint8)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint16(key string) uint16 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(uint16)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint32(key string) uint32 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(uint32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Uint64(key string) uint64 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(uint64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float32(key string) float32 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(float32)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) Float64(key string) float64 {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	value, ok := v.(float64)
	if !ok {
		return 0
	}
	return value
}

func (s *Session) String(key string) string {
	s.RLock()
	defer s.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return ""
	}
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

func (s *Session) Value(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.data[key]
}

func (s *Session) State() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

func (s *Session) Restore(data map[string]interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data = data
}

func (s *Session) Clear() {
	s.Lock()
	defer s.Unlock()
	s.uid = 0
	s.data = map[string]interface{}{}
}
