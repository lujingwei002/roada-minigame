package game

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/roada-go/cli"
	"github.com/shark/minigame-common/gamepb"
)

type Agent struct {
	game        *GameService
	session     *cli.Session
	chGate      chan *cli.Request
	chQuit      chan bool
	userid      int64
	openid      string
	handlerDict map[string]HandlerInterface
	isClose     bool
}

type HandlerInterface interface {
}

func newAgent(game *GameService, session *cli.Session, openid string) *Agent {
	agent := &Agent{
		game:        game,
		openid:      openid,
		session:     session,
		handlerDict: make(map[string]HandlerInterface),
		chGate:      make(chan *cli.Request, 1),
		chQuit:      make(chan bool),
	}
	handler := newDaxiguaHandler(agent)
	agent.handlerDict["daxigua"] = handler
	go agent.forever()
	return agent
}

func (self *Agent) ServeMessage(r *cli.Request) {
	self.chGate <- r
}

func (agent *Agent) onSessionOpen(session *cli.Session) {
	log.Printf("[agent] OnSessionOpen, sessionid=%d\n", session.ID())
	var req = &gamepb.LoginRequest{}
	req.Openid = agent.openid
	//req.Openid = "100"
	session.Request("game.login", req)
}

func (agent *Agent) onSessionClose(session *cli.Session) {
	close(agent.chQuit)
}

func (agent *Agent) forever() {
	session := agent.session
	tick := time.NewTicker(300 * time.Second)
	log.Printf("[agent] loop start, sessionid:%d\n", session.ID())
	defer func() {
		log.Printf("[agent] loop end, sessionid:%d\n", session.ID())
		tick.Stop()
		agent.close()
	}()
	for !agent.isClose {
		select {
		case r := <-agent.chGate:
			{
				agent.response(r)
			}
		case <-agent.chQuit:
			{
				return
			}
		case <-tick.C:
			{

			}
		}
	}
}

func (agent *Agent) response(r *cli.Request) {
	if r == nil {
		return
	}
	index := strings.LastIndex(r.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("[game] invalid route, route=%s", r.Route))
		return
	}
	handlerName := r.Route[:index]
	//methodName := r.Route[index+1:]
	var handler interface{}
	var service *cli.Service
	var ok bool
	service, ok = agent.game.serviceDict[handlerName]
	if !ok {
		log.Printf("[agent] service not found, route=%s\n", r.Route)
		return
	}
	if handlerName == "game" {
		handler = agent
	} else {
		handler, ok = agent.handlerDict[handlerName]
		if !ok {
			log.Printf("[agent] handler not found, route=%s\n", r.Route)
			return
		}
	}
	service.Unpack(handler, r)
}

func (agent *Agent) close() {
	if agent.isClose {
		return
	}
	agent.isClose = true
	session := agent.session
	userid := agent.userid
	log.Printf("[agent] close, sessionid:%d, userid:%d\n", session.ID(), userid)
	//释放agent
	close(agent.chGate)
}

func (agent *Agent) Login(r *cli.Request, msg *gamepb.LoginResponse) error {
	log.Printf("[agent] LoginResponse, sessionid=%d, response=%+v\n", r.Session.ID(), msg)
	if msg.Code != 0 {
		return errors.New("login failed")
	}
	agent.userid = msg.Userid
	var req = &gamepb.DaxiguaGetDataRequest{}
	r.Session.Request("daxigua.getdata", req)
	return nil
}

func (agent *Agent) Kick(r *cli.Request, msg *gamepb.KickPush) error {
	log.Printf("[agent] KickPush, sessionid=%d, push=%+v\n", r.Session.ID(), msg)
	return nil
}

func (self *Agent) GetScoreRank(r *cli.Request, reply *gamepb.GetScoreRankResponse) error {
	log.Printf("[agent] GetScoreRankResponse, response=%+v\n")
	time.Sleep(1 * time.Second)
	if rand.Intn(100) < 3 {
		r.Session.Close()
		return nil
	}
	var req = &gamepb.RoundStartRequest{}
	r.Session.Request("daxigua.roundstart", req)
	return nil
}
