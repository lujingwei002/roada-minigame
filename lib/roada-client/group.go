package cli

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/roada-go/gat/log"
)

const (
	groupStatusWorking = 0
	groupStatusClosed  = 1
)

type SessionFilter func(*Session) bool

type Group struct {
	client   *Client
	mu       sync.RWMutex
	status   int32
	name     string
	sessions map[int64]*Session
}

func newGroup(client *Client, name string) *Group {
	return &Group{
		client:   client,
		status:   groupStatusWorking,
		name:     name,
		sessions: make(map[int64]*Session),
	}
}

func (self *Group) Member(uid int64) (*Session, error) {
	self.mu.RLock()
	defer self.mu.RUnlock()
	for _, s := range self.sessions {
		if s.UID() == uid {
			return s, nil
		}
	}
	return nil, ErrMemberNotFound
}

func (self *Group) Members() []int64 {
	self.mu.RLock()
	defer self.mu.RUnlock()
	var members []int64
	for _, s := range self.sessions {
		members = append(members, s.UID())
	}
	return members
}

func (self *Group) Multicast(route string, v interface{}, filter SessionFilter) error {
	if self.isClosed() {
		return ErrClosedGroup
	}
	self.mu.RLock()
	defer self.mu.RUnlock()
	for _, s := range self.sessions {
		data, err := s.agent.serialize(v)
		if err != nil {
			return err
		}
		if !filter(s) {
			continue
		}
		if err = s.Notify(route, data); err != nil {
			log.Println(err.Error())
		}
	}
	return nil
}

func (self *Group) Broadcast(route string, v interface{}) error {
	if self.isClosed() {
		return ErrClosedGroup
	}
	self.mu.RLock()
	defer self.mu.RUnlock()
	for _, s := range self.sessions {
		data, err := s.agent.serialize(v)
		if err != nil {
			continue
		}
		if err = s.Notify(route, data); err != nil {
			log.Println(fmt.Sprintf("Session push message error, ID=%d, UID=%d, Error=%s", s.ID(), s.UID(), err.Error()))
		}
	}
	return nil
}

func (self *Group) Contains(uid int64) bool {
	_, err := self.Member(uid)
	return err == nil
}

func (self *Group) Add(session *Session) error {
	if self.isClosed() {
		return ErrClosedGroup
	}
	if self.client.debug {
		log.Println(fmt.Sprintf("Add session to group %s, ID=%d, UID=%d", self.name, session.ID(), session.UID()))
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	id := session.ID()
	_, ok := self.sessions[session.ID()]
	if ok {
		return ErrSessionDuplication
	}
	self.sessions[id] = session
	return nil
}

func (self *Group) Leave(s *Session) error {
	if self.isClosed() {
		return ErrClosedGroup
	}
	if self.client.debug {
		log.Println(fmt.Sprintf("Remove session from group %s, UID=%d", self.name, s.UID()))
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.sessions, s.ID())
	return nil
}

func (self *Group) LeaveAll() error {
	if self.isClosed() {
		return ErrClosedGroup
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	self.sessions = make(map[int64]*Session)
	return nil
}

func (self *Group) Count() int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return len(self.sessions)
}

func (self *Group) isClosed() bool {
	if atomic.LoadInt32(&self.status) == groupStatusClosed {
		return true
	}
	return false
}

func (self *Group) Close() error {
	if self.isClosed() {
		return ErrCloseClosedGroup
	}
	atomic.StoreInt32(&self.status, groupStatusClosed)
	self.sessions = make(map[int64]*Session)
	return nil
}
