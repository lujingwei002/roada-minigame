package cli

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
	chDie          chan struct{}
	chExit         chan struct{}
	chTasks        chan Task
	started        int32
	closed         int32
	timerPrecision time.Duration
	timerManager   *TimerManager
}

func NewScheduler() *Scheduler {
	self := &Scheduler{}
	self.chDie = make(chan struct{})
	self.chExit = make(chan struct{})
	self.chTasks = make(chan Task, 1<<8)
	self.timerPrecision = time.Second
	self.timerManager = newTimerManager()
	return self
}

func (self *Scheduler) try(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("Handle message panic: %+v\n%s", err, debug.Stack()))
		}
	}()
	f()
}

func (self *Scheduler) Sched() {
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
}

func (self *Scheduler) Close() {
	if atomic.AddInt32(&self.closed, 1) != 1 {
		return
	}
	close(self.chDie)
	<-self.chExit
	log.Println("Scheduler stopped")
}

func (self *Scheduler) PushTask(task Task) {
	self.chTasks <- task
}
