package game

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/roada-go/gat"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/config"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type TiantianpaokuData struct {
	model.Tiantianpaoku
	Dirty bool
}

type TiantianpaokuHandler struct {
	agent    *Agent
	gamedata *TiantianpaokuData
}

func newTiantianpaokuHandler(agent *Agent) *TiantianpaokuHandler {
	handler := &TiantianpaokuHandler{
		agent: agent,
	}
	return handler
}

func (handler *TiantianpaokuHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *TiantianpaokuHandler) GetData(r *gat.Request, args *gamepb.TiantianpaokuGetDataRequest) error {
	log.Printf("[TiantianpaokuHandler] GetData %+v\n", args)
	response := &gamepb.TiantianpaokuGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *TiantianpaokuData
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
	response.SkinId = gamedata.SkinId
	response.SkinArr = gamedata.SkinArr
	response.GameId = conf.Ini.Game.Id
	log.Printf("[TiantianpaokuHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TiantianpaokuHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *TiantianpaokuHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[TiantianpaokuHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Tiantianpaoku_Save(&gamedata.Tiantianpaoku)
	if err != nil {
		log.Printf("[TiantianpaokuHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[TiantianpaokuHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *TiantianpaokuHandler) loadData(userid int64) (*TiantianpaokuData, error) {
	dbrow, err := db.Tiantianpaoku_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *TiantianpaokuData
	if dbrow == nil {
		gamedata = &TiantianpaokuData{
			model.Tiantianpaoku{
				Userid:  userid,
				SkinArr: "1",
				SkinId:  1,
			},
			false, //Dirty
		}
	} else {
		gamedata = &TiantianpaokuData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	return gamedata, nil
}

func (handler *TiantianpaokuHandler) BuySkin(r *gat.Request, args *gamepb.TiantianpaokuBuySkinRequest) error {
	log.Printf("[TiantianpaokuHandler] BuySkin %+v\n", args)
	response := &gamepb.TiantianpaokuBuySkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	if args.SkinId <= 0 || args.SkinId > int32(len(config.Tiantianpaoku.Skin)) {
		response.Code = errCodeTiantianpaokuBuySkin
		return r.Response(response)
	}
	conf := config.Tiantianpaoku.Skin[args.SkinId-1]
	coin := conf.Coin

	if user.Coin < coin {
		response.Code = errCodeTiantianpaokuBuySkin
		return r.Response(response)
	}

	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeTiantianpaokuBuySkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			response.Code = errCodeTiantianpaokuBuySkin
			return r.Response(response)
		}
	}
	gamedata.SkinArr = fmt.Sprintf("%s,%d", gamedata.SkinArr, args.SkinId)
	gamedata.Dirty = true

	response.Coin = coin
	user.DecCoin(coin, tlogTiantianpaokuBuySkin)

	handler.agent.bgSave()

	log.Printf("[TiantianpaokuHandler] BuySkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TiantianpaokuHandler) UseSkin(r *gat.Request, args *gamepb.TiantianpaokuUseSkinRequest) error {
	log.Printf("[TiantianpaokuHandler] UseSkin %+v\n", args)
	response := &gamepb.TiantianpaokuUseSkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	found := false
	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeTiantianpaokuUseSkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			found = true
			break
		}
	}
	if !found {
		response.Code = errCodeTiantianpaokuUseSkin
		return r.Response(response)
	}
	gamedata.SkinId = args.SkinId
	gamedata.Dirty = true
	log.Printf("[TiantianpaokuHandler] UseSkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TiantianpaokuHandler) GetCoin(r *gat.Request, args *gamepb.TiantianpaokuGetCoinRequest) error {
	log.Printf("[TiantianpaokuHandler] GetCoin %+v\n", args)
	response := &gamepb.TiantianpaokuGetCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var coin int64 = args.Coin
	user.AddCoin(coin, tlogTiantianpaokuGetCoin)
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

//关卡开始
func (handler *TiantianpaokuHandler) RoundStart(r *gat.Request, args *gamepb.TiantianpaokuStartRequest) error {
	log.Printf("[TiantianpaokuHandler] RoundStart %+v\n", args)
	response := &gamepb.TiantianpaokuStartResponse{}
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

//关卡结束
func (handler *TiantianpaokuHandler) RoundResult(r *gat.Request, args *gamepb.TiantianpaokuResultRequest) error {
	log.Printf("[TiantianpaokuHandler] RoundResult %+v\n", args)
	response := &gamepb.TiantianpaokuResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	score := args.Score
	user := handler.agent.user
	userid := user.Userid
	roundResultTime := time.Now().Unix()

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
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, score, 0)
	return r.Response(response)
}
