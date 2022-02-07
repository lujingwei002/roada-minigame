package game

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/roada-go/gat"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type BenpaobaxiaoqieData struct {
	model.Benpaobaxiaoqie
	Dirty bool
}

type BpbxqHandler struct {
	agent    *Agent
	gamedata *BenpaobaxiaoqieData
}

func newBpbxqHandler(agent *Agent) *BpbxqHandler {
	handler := &BpbxqHandler{
		agent: agent,
	}
	return handler
}

func (handler *BpbxqHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (svr *BpbxqHandler) getLocalDate() int64 {
	now := time.Now()
	return int64(now.Year()*1000 + now.YearDay())
}

func (handler *BpbxqHandler) GetData(r *gat.Request, args *gamepb.BpbxqGetDataRequest) error {
	log.Printf("[BpbxqHandler] GetData %+v\n", args)
	response := &gamepb.BpbxqGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *BenpaobaxiaoqieData
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
	handler.resetSignData()
	response.GameId = conf.Ini.Game.Id
	response.Skin = gamedata.Skin
	response.SkinArr = gamedata.SkinArr
	response.LastSignTime = gamedata.LastSignTime
	response.SignTimes = gamedata.SignTimes
	response.SignTimeNow = handler.getLocalDate()
	log.Printf("[BpbxqHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) resetSignData() {
	gamedata := handler.gamedata
	now := handler.getLocalDate()
	if gamedata.SignTimes >= 7 && now-gamedata.LastSignTime >= 1 {
		gamedata.SignTimes = 0
		gamedata.Dirty = true
	}
}

func (handler *BpbxqHandler) DaySign(r *gat.Request, args *gamepb.BpbxqDaySignRequest) error {
	log.Printf("[BpbxqHandler] DaySign %+v\n", args)
	response := &gamepb.BpbxqDaySignResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	now := handler.getLocalDate()
	if gamedata.SignTimes < 0 || gamedata.SignTimes > 7 {
		response.Code = errCodeBpbxqDaySign
		return r.Response(response)
	}
	if now-gamedata.LastSignTime < 1 {
		response.Code = errCodeBpbxqDaySign
		return r.Response(response)
	}
	coinReward := []int64{50, 100, 200, 300, 400, 500, 1000}
	coin := coinReward[gamedata.SignTimes]
	gamedata.LastSignTime = now
	gamedata.SignTimes = gamedata.SignTimes + 1
	gamedata.Dirty = true
	user.AddCoin(coin, tlogBpbxqDaySign)
	handler.agent.bgSave()
	tlog.DaySign(user.Openid, userid, gamedata.SignTimes)
	response.LastSignTime = gamedata.LastSignTime
	response.SignTimes = gamedata.SignTimes
	response.SignTimeNow = now
	response.Coin = coin
	log.Printf("[BpbxqHandler] DaySign succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) BuySkin(r *gat.Request, args *gamepb.BpbxqBuySkinRequest) error {
	log.Printf("[BpbxqHandler] BuySkin %+v\n", args)
	response := &gamepb.BpbxqBuySkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	priceDict := []int64{0, 4000, 5000, 6000, 8000}
	if args.SkinId <= 0 || args.SkinId >= int32(len(priceDict)) {
		response.Code = errCodeBpbxqBuySkin
		return r.Response(response)
	}
	coin := priceDict[args.SkinId]
	if user.Coin < coin {
		response.Code = errCodeBpbxqBuySkin
		return r.Response(response)
	}
	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeBpbxqBuySkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			response.Code = errCodeBpbxqBuySkin
			return r.Response(response)
		}
	}
	gamedata.Skin = args.SkinId
	gamedata.SkinArr = fmt.Sprintf("%s,%d", gamedata.SkinArr, args.SkinId)
	gamedata.Dirty = true
	user.DecCoin(coin, tlogBpbxqBuySkin)
	handler.agent.bgSave()
	response.Coin = coin
	log.Printf("[BpbxqHandler] BuySkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) UseSkin(r *gat.Request, args *gamepb.BpbxqUseSkinRequest) error {
	log.Printf("[BpbxqHandler] UseSkin %+v\n", args)
	response := &gamepb.BpbxqUseSkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	found := false
	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeBpbxqUseSkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			found = true
			break
		}
	}
	if !found {
		response.Code = errCodeBpbxqUseSkin
		return r.Response(response)
	}
	gamedata.Skin = args.SkinId
	gamedata.Dirty = true
	log.Printf("[BpbxqHandler] UseSkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) PickCoin(r *gat.Request, args *gamepb.BpbxqPickCoinRequest) error {
	log.Printf("[BpbxqHandler] PickCoin %+v\n", args)
	response := &gamepb.BpbxqPickCoinResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	args.Coin = 1
	user.AddCoin(args.Coin, tlogBpbxqPick)
	response.Coin = args.Coin
	log.Printf("[BpbxqHandler] PickCoin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) Revive(r *gat.Request, args *gamepb.BpbxqReviveRequest) error {
	log.Printf("[BpbxqHandler] Revive %+v\n", args)
	response := &gamepb.BpbxqReviveResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	var coin int64 = 200
	if user.Coin < coin {
		response.Code = errCodeBpbxqRevive
		return r.Response(response)
	}
	user.RoundStartTime = user.LastRoundStartTime
	user.Dirty = true
	user.DecCoin(coin, tlogBpbxqRevive)
	response.Coin = coin
	log.Printf("[BpbxqHandler] Revive succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) Diefly(r *gat.Request, args *gamepb.BpbxqDieflyRequest) error {
	log.Printf("[BpbxqHandler] Diefly %+v\n", args)
	response := &gamepb.BpbxqDieflyResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	var coin int64 = 100
	if user.Coin < coin {
		response.Code = errCodeBpbxqDiefly
		return r.Response(response)
	}
	user.RoundStartTime = user.LastRoundStartTime
	user.Dirty = true
	user.DecCoin(coin, tlogBpbxqDiefly)
	response.Coin = coin
	log.Printf("[BpbxqHandler] Diefly succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) Reward(r *gat.Request, args *gamepb.BpbxqRewardRequest) error {
	log.Printf("[BpbxqHandler] Reward %+v\n", args)
	response := &gamepb.BpbxqRewardResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	var coin int64 = 300
	user.AddCoin(coin, tlogBpbxqReward)
	response.Coin = coin
	log.Printf("[BpbxqHandler] Reward succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *BpbxqHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *BpbxqHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[BpbxqHandler] saveData, dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Bpbxq_Save(&gamedata.Benpaobaxiaoqie)
	if err != nil {
		log.Printf("[BpbxqHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[BpbxqHandler] saveData %+v success\n", gamedata.Benpaobaxiaoqie)
	gamedata.Dirty = false
	return nil
}

func (handler *BpbxqHandler) loadData(userid int64) (*BenpaobaxiaoqieData, error) {
	dbrow, err := db.Bpbxq_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *BenpaobaxiaoqieData
	if dbrow == nil {
		gamedata = &BenpaobaxiaoqieData{
			model.Benpaobaxiaoqie{
				Userid:       userid,
				Skin:         0,
				SkinArr:      "0",
				LastSignTime: 0,
				SignTimes:    0,
			},
			false, //Dirty
		}
	} else {
		gamedata = &BenpaobaxiaoqieData{
			*dbrow,
			false, //Dirty
		}
	}
	return gamedata, nil
}

func (handler *BpbxqHandler) RoundStart(r *gat.Request, args *gamepb.BpbxqStartRequest) error {
	log.Printf("[BpbxqHandler] RoundStart %+v\n", args)
	response := &gamepb.BpbxqStartResponse{}
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

func (handler *BpbxqHandler) RoundResult(r *gat.Request, args *gamepb.BpbxqResultRequest) error {
	log.Printf("[BpbxqHandler] RoundResult %+v\n", args)
	response := &gamepb.BpbxqResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	score := args.Score
	user := handler.agent.user
	userid := user.Userid
	roundResultTime := time.Now().Unix()
	//duration := roundResultTime - user.roundStartTime
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
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, score, 0)
	return r.Response(response)
}
