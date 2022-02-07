package middleware

import (
	"fmt"
	"time"

	"github.com/roada-go/gat"
)

type agentStat struct {
	sid              int64
	packetCount      uint64
	packetTotalCount uint64
	packetAvgCount   uint64
	lastTime         time.Time
}

type AgentStat struct {
}

func newAgentStat(sid int64) *agentStat {
	self := &agentStat{
		sid:      sid,
		lastTime: time.Now(),
	}
	return self
}

func (self *agentStat) inpack(r *gat.Request) {
	self.packetCount++
	self.packetTotalCount++
	if self.lastTime.Add(5 * time.Second).Before(time.Now()) {
		duration := time.Now().Sub(self.lastTime)
		self.lastTime = time.Now()
		self.packetAvgCount = self.packetCount / uint64(duration.Seconds())
		self.packetCount = 0
	}
}

func (self *agentStat) String() string {
	return fmt.Sprintf("[agentstat] sid:%d, total:%d, avg:%d", self.sid, self.packetTotalCount, self.packetAvgCount)
}

func (self *AgentStat) ServeMessage(r *gat.Request) {
	session := r.Session
	var stat *agentStat
	if ok := session.HasKey("stat"); !ok {
		stat = newAgentStat(session.ID())
		session.Set("stat", stat)
	} else {
		stat = session.Value("stat").(*agentStat)
	}
	stat.inpack(r)
}
