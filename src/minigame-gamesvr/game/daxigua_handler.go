package game

import (
	"log"
	"math"
	"time"

	"github.com/roada-go/gat"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type DaxiguaData struct {
	model.Daxigua
	Dirty bool
}

type DaxiguaHandler struct {
	agent    *Agent
	gamedata *DaxiguaData
}

func newDaxiguaHandler(agent *Agent) *DaxiguaHandler {
	handler := &DaxiguaHandler{
		agent: agent,
	}
	return handler
}

func (handler *DaxiguaHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *DaxiguaHandler) GetData(r *gat.Request, args *gamepb.DaxiguaGetDataRequest) error {
	log.Printf("[DaxiguaHandler] GetData %+v\n", args)
	response := &gamepb.DaxiguaGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *DaxiguaData
	var err error
	if handler.gamedata != nil {
		gamedata = handler.gamedata
	} else {
		gamedata, err = handler.loadData(userid)
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
		handler.gamedata = gamedata
	}
	response.GameId = conf.Ini.Game.Id
	log.Printf("[DaxiguaHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *DaxiguaHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *DaxiguaHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[DaxiguaHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Daxigua_Save(&gamedata.Daxigua)
	if err != nil {
		log.Printf("[DaxiguaHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[DaxiguaHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *DaxiguaHandler) loadData(userid int64) (*DaxiguaData, error) {
	dbrow, err := db.Daxigua_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *DaxiguaData
	if dbrow == nil {
		gamedata = &DaxiguaData{
			Dirty: false, //Dirty
		}
		gamedata.Daxigua = model.Daxigua{
			Userid: userid,
		}
	} else {
		gamedata = &DaxiguaData{
			Dirty: false, //Dirty
		}
		gamedata.Daxigua = *dbrow
	}
	return gamedata, nil
}

func (handler *DaxiguaHandler) RoundStart(r *gat.Request, args *gamepb.DaxiguaStartRequest) error {
	log.Printf("[DaxiguaHandler] RoundStart %+v\n", args)
	response := &gamepb.DaxiguaStartResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	user.RoundStartTime = time.Now().Unix()
	user.Dirty = true
	tlog.RoundStart(user.Openid, userid)
	return r.Response(response)
}

func (handler *DaxiguaHandler) RoundResult(r *gat.Request, args *gamepb.DaxiguaResultRequest) error {
	log.Printf("[DaxiguaHandler] RoundResult %+v\n", args)
	response := &gamepb.DaxiguaResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	score := args.Score
	medal := args.Medal
	user := handler.agent.user
	userid := user.Userid
	roundResultTime := time.Now().Unix()
	duration := float64(roundResultTime - user.RoundStartTime)
	log.Println("ggg", duration, roundResultTime, user.RoundStartTime)
	minutes := int64(math.Ceil(duration / 60))
	maxScore := 1000 * minutes
	log.Println("fffffffff", minutes, maxScore)
	if score > maxScore {
		response.Code = errCodeHacker
		return r.Response(response)
	}
	//	log.Println("ofafsfasf", duration)
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	user.Dirty = true
	if score > user.Score {
		user.Score = score
		db.ScoreRank_Update(userid, score)
	}
	if score > user.ScoreToday {
		user.ScoreToday = score
		db.ScoreDayRank_Update(user.Today, userid, score)
	}
	if medal > user.Medal {
		user.Medal = medal
		db.MedalRank_Update(userid, medal)
	}
	if medal > user.MedalToday {
		user.MedalToday = medal
		db.MedalDayRank_Update(user.Today, userid, medal)
	}
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, score, medal)
	return r.Response(response)
}
