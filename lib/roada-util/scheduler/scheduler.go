package gat

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"

	"log"
)

const ()

type LocalScheduler interface {
	Schedule(Task)
}

type Task func()

type Hook func()

type Scheduler struct {
	T              chan Task
	chDie          chan struct{}
	chExit         chan struct{}
	started        int32
	closed         int32
	timerPrecision time.Duration
	timerManager   *timerManager
}

func NewScheduler() *Scheduler {
	self := &Scheduler{}
	self.chDie = make(chan struct{})
	self.chExit = make(chan struct{})
	self.T = make(chan Task, 1<<8)
	self.timerPrecision = time.Second
	self.timerManager = newTimerManager()
	return self
}

func (self *Scheduler) Invoke(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle message panic: %+v\n%s", err, debug.Stack()))
		}
	}()
	f()
}

func (self *Scheduler) Sched() {
	self.timerManager.cron()
}

/*func (self *Scheduler) Sched2() {
	if atomic.AddInt32(&self.started, 1) != 1 {
		return
	}
	ticker := time.NewTicker(self.timerPrecision)
	defer func() {
		ticker.Stop()
		close(self.chExit)
	}()
	for {
		select {
		case <-ticker.C:
			self.timerManager.cron()
		case f := <-self.chTasks:
			self.try(f)
		case <-self.chDie:
			return
		}
	}
}*/

func (self *Scheduler) Close() {
	if atomic.AddInt32(&self.closed, 1) != 1 {
		return
	}
	log.Printf("[scheduler] stopped\n")
}

/*func (self *Scheduler) Close2() {
	if atomic.AddInt32(&self.closed, 1) != 1 {
		return
	}
	close(self.chDie)
	<-self.chExit
	log.Printf("[scheduler] stopped\n")
}*/

func (self *Scheduler) PushTask(task Task) {
	self.T <- task
}

func (self *Scheduler) NewTimer(interval time.Duration, fn TimerFunc) *Timer {
	return self.timerManager.newTimer(interval, fn)
}

func (self *Scheduler) NewCountTimer(interval time.Duration, count int, fn TimerFunc) *Timer {
	return self.timerManager.newCountTimer(interval, count, fn)
}

func (self *Scheduler) NewAfterTimer(duration time.Duration, fn TimerFunc) *Timer {
	return self.timerManager.newAfterTimer(duration, fn)
}
