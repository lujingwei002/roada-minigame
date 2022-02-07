package game

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/roada-go/cli"
	"github.com/shark/minigame-common/gamepb"
)

type DaxiguaHandler struct {
	agent *Agent
}

func newDaxiguaHandler(agent *Agent) *DaxiguaHandler {
	handler := &DaxiguaHandler{
		agent: agent,
	}
	return handler
}

func (self *DaxiguaHandler) GetData(r *cli.Request, reply *gamepb.DaxiguaGetDataResponse) error {
	log.Printf("DaxiguaGetDataResponse, sessionid=%d, response=%+v\n", r.Session.ID(), reply)
	if reply.Code != 0 {
		return errors.New("getdata failed")
	}
	var req = &gamepb.RoundStartRequest{}
	r.Session.Request("daxigua.roundstart", req)
	return nil
}

func (self *DaxiguaHandler) RoundStart(r *cli.Request, reply *gamepb.DaxiguaStartResponse) error {
	log.Printf("RoundStartResponse, sessionid=%d, response=%+v\n", r.Session.ID(), reply)
	time.Sleep(1 * time.Second)
	var req = &gamepb.RoundResultRequest{Score: int64(rand.Intn(99999)), Medal: int64(rand.Intn(99999))}
	r.Session.Request("daxigua.roundresult", req)
	return nil
}

func (self *DaxiguaHandler) RoundResult(r *cli.Request, reply *gamepb.DaxiguaResultResponse) error {
	log.Printf("RoundResultResponse, sessionid=%d, response=%+v\n", r.Session.ID(), reply)
	time.Sleep(1 * time.Second)
	if rand.Intn(100) < 0 {
		r.Session.Close()
		return nil
	}
	var req = &gamepb.RoundStartRequest{}
	r.Session.Request("daxigua.roundstart", req)
	return nil

	/*var req = &gamepb.GetScoreRankRequest{
		Type: 1,
	}
	r.Session.Request("game.getscorerank", req)
	return nil*/
}
