package game

import (
	"fmt"
	"log"
	"math"
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

type TanchishedazuozhanData struct {
	model.Tanchishedazuozhan
	Dirty bool
}

type TanchishedazuozhanHandler struct {
	agent    *Agent
	gamedata *TanchishedazuozhanData
}

func newTanchishedazuozhanHandler(agent *Agent) *TanchishedazuozhanHandler {
	handler := &TanchishedazuozhanHandler{
		agent: agent,
	}
	return handler
}

func (handler *TanchishedazuozhanHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *TanchishedazuozhanHandler) GetData(r *gat.Request, args *gamepb.TanchishedazuozhanGetDataRequest) error {
	log.Printf("[TanchishedazuozhanHandler] GetData %+v\n", args)
	response := &gamepb.TanchishedazuozhanGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *TanchishedazuozhanData
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
	log.Printf("[TanchishedazuozhanHandler] GetData1 succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TanchishedazuozhanHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *TanchishedazuozhanHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[TanchishedazuozhanHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Tanchishedazuozhan_Save(&gamedata.Tanchishedazuozhan)
	if err != nil {
		log.Printf("[TanchishedazuozhanHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[TanchishedazuozhanHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *TanchishedazuozhanHandler) loadData(userid int64) (*TanchishedazuozhanData, error) {
	dbrow, err := db.Tanchishedazuozhan_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *TanchishedazuozhanData
	if dbrow == nil {
		gamedata = &TanchishedazuozhanData{
			model.Tanchishedazuozhan{
				Userid:  userid,
				SkinArr: "1",
				SkinId:  1,
			},
			false, //Dirty
		}
	} else {
		gamedata = &TanchishedazuozhanData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	return gamedata, nil
}

func (handler *TanchishedazuozhanHandler) BuySkin(r *gat.Request, args *gamepb.TanchishedazuozhanBuySkinRequest) error {
	log.Printf("[TanchishedazuozhanHandler] BuySkin %+v\n", args)
	response := &gamepb.TanchishedazuozhanBuySkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	if args.SkinId <= 0 || args.SkinId > int32(len(config.Tanchishedazuozhan.Skin)) {
		response.Code = errCodeTanchishedazuozhanBuySkin
		return r.Response(response)
	}
	conf := config.Tanchishedazuozhan.Skin[args.SkinId-1]
	price := conf.Price
	if conf.Type == 0 {
		if user.Coin < price {
			response.Code = errCodeTanchishedazuozhanBuySkin
			return r.Response(response)
		}
	} else {
		if user.Diamond < price {
			response.Code = errCodeTanchishedazuozhanBuySkin
			return r.Response(response)
		}
	}
	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeTanchishedazuozhanBuySkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			response.Code = errCodeTanchishedazuozhanBuySkin
			return r.Response(response)
		}
	}
	gamedata.SkinArr = fmt.Sprintf("%s,%d", gamedata.SkinArr, args.SkinId)
	gamedata.Dirty = true
	if conf.Type == 0 {
		response.Coin = price
		user.DecCoin(price, tlogTanchishedazuozhanBuySkin)
	} else {
		response.Diamond = price
		user.DecDiamond(price, tlogTanchishedazuozhanBuySkin)
	}
	handler.agent.bgSave()

	log.Printf("[TanchishedazuozhanHandler] BuySkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TanchishedazuozhanHandler) UseSkin(r *gat.Request, args *gamepb.TanchishedazuozhanUseSkinRequest) error {
	log.Printf("[TanchishedazuozhanHandler] UseSkin %+v\n", args)
	response := &gamepb.TanchishedazuozhanUseSkinResponse{Code: 0, SkinId: args.SkinId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	found := false
	skinArr := strings.Split(gamedata.SkinArr, ",")
	for _, _skinId := range skinArr {
		if skinId, err := strconv.Atoi(_skinId); err != nil {
			response.Code = errCodeTanchishedazuozhanUseSkin
			return r.Response(response)
		} else if skinId == int(args.SkinId) {
			found = true
			break
		}
	}
	if !found {
		response.Code = errCodeTanchishedazuozhanUseSkin
		return r.Response(response)
	}
	gamedata.SkinId = args.SkinId
	gamedata.Dirty = true
	log.Printf("[TanchishedazuozhanHandler] UseSkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *TanchishedazuozhanHandler) GetCoin(r *gat.Request, args *gamepb.TanchishedazuozhanGetCoinRequest) error {
	log.Printf("[TanchishedazuozhanHandler] GetCoin %+v\n", args)
	response := &gamepb.TanchishedazuozhanGetCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var coin int64 = args.Coin
	user.AddCoin(coin, tlogTanchishedazuozhanGetCoin)
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

//关卡开始
func (handler *TanchishedazuozhanHandler) RoundStart(r *gat.Request, args *gamepb.TanchishedazuozhanStartRequest) error {
	log.Printf("[TanchishedazuozhanHandler] RoundStart %+v\n", args)
	response := &gamepb.TanchishedazuozhanStartResponse{}
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
func (handler *TanchishedazuozhanHandler) RoundResult(r *gat.Request, args *gamepb.TanchishedazuozhanResultRequest) error {
	log.Printf("[TanchishedazuozhanHandler] RoundResult %+v\n", args)
	response := &gamepb.TanchishedazuozhanResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	score := args.Score
	user := handler.agent.user
	userid := user.Userid
	roundResultTime := time.Now().Unix()

	duration := float64(roundResultTime - user.RoundStartTime)
	minutes := int64(math.Ceil(duration / 60))
	maxScore := 100 * minutes
	if score > maxScore {
		//	response.Code = errCodeHacker
		//return r.Response(response)
	}
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	user.Dirty = true
	if args.Mode == 1 {
		if score > user.Score {
			user.Score = score
			db.ScoreRank_Update(userid, score)
		}
		if score > user.ScoreToday {
			user.ScoreToday = score
			db.ScoreDayRank_Update(user.Today, userid, score)
		}
		//tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, score, 0)
	} else {
		if score > user.Medal {
			user.Medal = score
			db.MedalRank_Update(userid, score)
		}
		if score > user.MedalToday {
			user.MedalToday = score
			db.MedalDayRank_Update(user.Today, userid, score)
		}
		tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, score, 0)
	}
	return r.Response(response)
}
