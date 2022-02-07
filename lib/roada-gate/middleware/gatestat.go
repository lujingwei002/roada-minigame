package middleware

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/roada-go/gat"
)

type GateStat struct {
	PacketTotalCount uint64
	PacketAvgCount   uint64
	SessionCount     int64

	packetCount uint64
	lastTime    time.Time
}

func NewGateStat() *GateStat {
	stat := &GateStat{
		lastTime: time.Now(),
	}
	return stat
}

func (stat *GateStat) Record() {
	duration := time.Now().Sub(stat.lastTime)
	if uint64(duration.Seconds()) <= 0 {
		return
	}
	stat.lastTime = time.Now()
	stat.PacketAvgCount = stat.packetCount / uint64(duration.Seconds())
	stat.packetCount = 0
}

func (stat *GateStat) String() string {
	return fmt.Sprintf("[gatestat] session:%d, total:%d, avg:%d",
		stat.SessionCount, stat.PacketTotalCount, stat.PacketAvgCount)
}

func (stat *GateStat) ServeMessage(r *gat.Request) {
	atomic.AddUint64(&stat.packetCount, 1)
	atomic.AddUint64(&stat.PacketTotalCount, 1)
}

func (stat *GateStat) OnSessionOpen(s *gat.Session) {
	atomic.AddInt64(&stat.SessionCount, 1)
}

func (stat *GateStat) OnSessionClose(s *gat.Session) {
	atomic.AddInt64(&stat.SessionCount, -1)
}
