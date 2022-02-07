package cli

import (
	"fmt"
	"log"
	"math"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	infinite = -1
)

type (
	TimerManager struct {
		incrementID    int64
		timers         map[int64]*Timer
		muClosingTimer sync.RWMutex
		closingTimer   []int64
		muCreatedTimer sync.RWMutex
		createdTimer   []*Timer
	}
	TimerFunc func()

	TimerCondition interface {
		Check(now time.Time) bool
	}

	Timer struct {
		id        int64
		fn        TimerFunc
		createAt  int64
		interval  time.Duration
		condition TimerCondition
		elapse    int64
		closed    int32
		counter   int
	}
)

func newTimerManager() *TimerManager {
	self := &TimerManager{}
	self.timers = map[int64]*Timer{}
	return self
}

func (t *Timer) ID() int64 {
	return t.id
}

func (t *Timer) Stop() {
	if atomic.AddInt32(&t.closed, 1) != 1 {
		return
	}
	t.counter = 0
}

func safecall(id int64, fn TimerFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle timer panic: %+v\n%s", err, debug.Stack()))
		}
	}()
	fn()
}

func (self *TimerManager) cron() {
	if len(self.createdTimer) > 0 {
		self.muCreatedTimer.Lock()
		for _, t := range self.createdTimer {
			self.timers[t.id] = t
		}
		self.createdTimer = self.createdTimer[:0]
		self.muCreatedTimer.Unlock()
	}
	if len(self.timers) < 1 {
		return
	}
	now := time.Now()
	unn := now.UnixNano()
	for id, t := range self.timers {
		if t.counter == infinite || t.counter > 0 {
			// condition timer
			if t.condition != nil {
				if t.condition.Check(now) {
					safecall(id, t.fn)
				}
				continue
			}
			// execute job
			if t.createAt+t.elapse <= unn {
				safecall(id, t.fn)
				t.elapse += int64(t.interval)
				// update timer counter
				if t.counter != infinite && t.counter > 0 {
					t.counter--
				}
			}
		}
		if t.counter == 0 {
			self.muClosingTimer.Lock()
			self.closingTimer = append(self.closingTimer, t.id)
			self.muClosingTimer.Unlock()
			continue
		}
	}
	if len(self.closingTimer) > 0 {
		self.muClosingTimer.Lock()
		for _, id := range self.closingTimer {
			delete(self.timers, id)
		}
		self.closingTimer = self.closingTimer[:0]
		self.muClosingTimer.Unlock()
	}
}

func (self *TimerManager) NewTimer(interval time.Duration, fn TimerFunc) *Timer {
	return self.NewCountTimer(interval, infinite, fn)
}

func (self *TimerManager) NewCountTimer(interval time.Duration, count int, fn TimerFunc) *Timer {
	if fn == nil {
		panic("timer: nil timer function")
	}
	if interval <= 0 {
		panic("timer: non-positive interval for NewTimer")
	}
	t := &Timer{
		id:       atomic.AddInt64(&self.incrementID, 1),
		fn:       fn,
		createAt: time.Now().UnixNano(),
		interval: interval,
		elapse:   int64(interval),
		counter:  count,
	}
	self.muCreatedTimer.Lock()
	self.createdTimer = append(self.createdTimer, t)
	self.muCreatedTimer.Unlock()
	return t
}

func (self *TimerManager) NewAfterTimer(duration time.Duration, fn TimerFunc) *Timer {
	return self.NewCountTimer(duration, 1, fn)
}

func (self *TimerManager) NewCondTimer(condition TimerCondition, fn TimerFunc) *Timer {
	if condition == nil {
		panic("timer: nil condition")
	}
	t := self.NewCountTimer(time.Duration(math.MaxInt64), infinite, fn)
	t.condition = condition
	return t
}
