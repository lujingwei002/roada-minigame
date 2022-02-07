package game

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/roada-go/cli"
)

type User struct {
	Userid         int64
	Openid         string
	Nickname       string
	Avatar         string
	RoundStartTime int64
}

type GameService struct {
	client      *cli.Client
	userCount   int32
	chGate      chan *cli.Request
	serviceDict map[string]*cli.Service
	agentDict   sync.Map
}

func Register(client *cli.Client, userCount int32) error {
	self := &GameService{
		client:      client,
		agentDict:   sync.Map{},
		serviceDict: make(map[string]*cli.Service),
		chGate:      make(chan *cli.Request, 1),
		userCount:   userCount,
	}
	client.HandleFunc(self.ServeMessage)
	client.OnSessionOpen(self.OnSessionOpen)
	client.OnSessionClose(self.OnSessionClose)
	if s, err := client.Register(&Agent{}); err != nil {
		return err
	} else {
		self.serviceDict["game"] = s
	}
	if s, err := client.Register(&DaxiguaHandler{}); err != nil {
		return err
	} else {
		self.serviceDict["daxigua"] = s
	}
	go self.forever()
	return nil
}

func Connect1() error {
	return nil
}

func (self *GameService) forever() {
	tick := time.NewTicker(100 * time.Millisecond)
	defer func() {
		tick.Stop()
	}()
	for {
		select {
		case <-tick.C:
			{
				self.tick()
			}
		}
	}
}

func (self *GameService) tick() {
	if self.userCount <= 0 {
		return
	}
	atomic.AddInt32(&self.userCount, -1)
	log.Println("bbbbbb", self.userCount)
	if err := self.client.Dial(); err != nil {
		atomic.AddInt32(&self.userCount, 1)
		log.Println("dial error", err.Error())
	}
}

func (game *GameService) ServeMessage(r *cli.Request) {
	session := r.Session
	if !session.HasKey("agent") {
		return
	}
	agent := session.Value("agent").(*Agent)
	agent.ServeMessage(r)
}

func (game *GameService) OnSessionOpen(session *cli.Session) {
	log.Printf("[GameService] session open, sessionid=%d\n", session.ID())
	openid := fmt.Sprintf("%d", rand.Intn(10000))
	agent := newAgent(game, session, openid)
	game.agentDict.Store(session.ID(), agent)
	session.Set("agent", agent)
	agent.onSessionOpen(session)
}

func (game *GameService) OnSessionClose(session *cli.Session) {
	/*log.Printf("OnSessionClose, sessionid=%d\n", s.ID())
	atomic.AddInt32(&self.userCount, -1)*/
	atomic.AddInt32(&game.userCount, 1)
	log.Printf("[GameService] session closed, sessionid=%d", session.ID())
	if !session.HasKey("agent") {
		return
	}
	agent := session.Value("agent").(*Agent)
	agent.onSessionClose(session)
	game.agentDict.Delete(session.ID())
	session.Remove("agent")
}
