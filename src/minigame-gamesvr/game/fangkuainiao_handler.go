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

type FangkuainiaoData struct {
	model.Fangkuainiao
	Dirty bool
}

type FangkuainiaoHandler struct {
	agent    *Agent
	gamedata *FangkuainiaoData
}

func newFangkuainiaoHandler(agent *Agent) *FangkuainiaoHandler {
	handler := &FangkuainiaoHandler{
		agent: agent,
	}
	return handler
}

func (handler *FangkuainiaoHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (svr *FangkuainiaoHandler) getLocalDate() int64 {
	now := time.Now()
	return int64(now.Year()*10000 + int(now.Month()*100) + now.Day())
}

func (handler *FangkuainiaoHandler) GetData(r *gat.Request, args *gamepb.FangkuainiaoGetDataRequest) error {
	log.Printf("[FangkuainiaoHandler] GetData %+v\n", args)
	response := &gamepb.FangkuainiaoGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *FangkuainiaoData
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
	handler.resetGetCoin()
	response.GameId = conf.Ini.Game.Id
	response.Level = gamedata.Level
	response.BirdId = gamedata.BirdId
	response.BirdArr = gamedata.BirdArr
	response.SignTime = gamedata.SignTime
	response.SignDay = gamedata.SignDay
	response.GetGoldCount = gamedata.GetGoldCount
	response.GetGoldTime = gamedata.GetGoldTime
	log.Printf("[FangkuainiaoHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) resetGetCoin() {
	gamedata := handler.gamedata
	now := handler.getLocalDate()
	if gamedata.GetGoldTime < now {
		gamedata.GetGoldCount = 0
		gamedata.GetGoldTime = now
		gamedata.Dirty = true
	}
}

func (handler *FangkuainiaoHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *FangkuainiaoHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[FangkuainiaoHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Fangkuainiao_Save(&gamedata.Fangkuainiao)
	if err != nil {
		log.Printf("[FangkuainiaoHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[FangkuainiaoHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *FangkuainiaoHandler) loadData(userid int64) (*FangkuainiaoData, error) {
	dbrow, err := db.Fangkuainiao_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *FangkuainiaoData
	if dbrow == nil {
		gamedata = &FangkuainiaoData{
			model.Fangkuainiao{
				Userid:       userid,
				BirdId:       1,
				BirdArr:      "1",
				SignDay:      0,
				SignTime:     0,
				GetGoldCount: 0,
				GetGoldTime:  0,
				Level:        1,
			},
			false, //Dirty
		}
	} else {
		gamedata = &FangkuainiaoData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	return gamedata, nil
}

func (handler *FangkuainiaoHandler) BuySkin(r *gat.Request, args *gamepb.FangkuainiaoBuySkinRequest) error {
	log.Printf("[FangkuainiaoHandler] BuySkin %+v\n", args)
	response := &gamepb.FangkuainiaoBuySkinResponse{Code: 0, BirdId: args.BirdId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	if args.BirdId <= 0 || args.BirdId >= int32(len(config.Fangkuainiao.Bird)) {
		response.Code = errCodeFangkuainiaoBuySkin
		return r.Response(response)
	}
	if config.Fangkuainiao.Bird[args.BirdId].Type != 2 {
		response.Code = errCodeFangkuainiaoBuySkin
		return r.Response(response)
	}
	coin := config.Fangkuainiao.Bird[args.BirdId].PayMoney
	if user.Coin < coin {
		response.Code = errCodeFangkuainiaoBuySkin
		return r.Response(response)
	}
	birdArr := strings.Split(gamedata.BirdArr, ",")
	for _, _birdId := range birdArr {
		if birdId, err := strconv.Atoi(_birdId); err != nil {
			response.Code = errCodeFangkuainiaoBuySkin
			return r.Response(response)
		} else if birdId == int(args.BirdId) {
			response.Code = errCodeFangkuainiaoBuySkin
			return r.Response(response)
		}
	}
	//gamedata.BirdId = args.BirdId
	gamedata.BirdArr = fmt.Sprintf("%s,%d", gamedata.BirdArr, args.BirdId)
	gamedata.Dirty = true
	user.DecCoin(coin, tlogFangkuainiaoBuySkin)
	handler.agent.bgSave()
	response.Coin = coin
	log.Printf("[FangkuainiaoHandler] BuySkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) UseSkin(r *gat.Request, args *gamepb.FangkuainiaoUseSkinRequest) error {
	log.Printf("[FangkuainiaoHandler] UseSkin %+v\n", args)
	response := &gamepb.FangkuainiaoUseSkinResponse{Code: 0, BirdId: args.BirdId}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	found := false
	birdArr := strings.Split(gamedata.BirdArr, ",")
	for _, _birdId := range birdArr {
		if birdId, err := strconv.Atoi(_birdId); err != nil {
			response.Code = errCodeFangkuainiaoUseSkin
			return r.Response(response)
		} else if birdId == int(args.BirdId) {
			found = true
			break
		}
	}
	if !found {
		response.Code = errCodeFangkuainiaoUseSkin
		return r.Response(response)
	}
	gamedata.BirdId = args.BirdId
	gamedata.Dirty = true
	log.Printf("[FangkuainiaoHandler] UseSkin succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) DaySign(r *gat.Request, args *gamepb.FangkuainiaoDaySignRequest) error {
	log.Printf("[FangkuainiaoHandler] DaySign %+v\n", args)
	response := &gamepb.FangkuainiaoDaySignResponse{Code: 0}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	now := handler.getLocalDate()
	if gamedata.SignTime >= now {
		response.Code = errCodeFangkuainiaoDaySign
		return r.Response(response)
	}
	rewardIndex := gamedata.SignDay % 7
	reward := config.Fangkuainiao.Sign[rewardIndex]
	var coin int64 = int64(reward.RewardsMoney)

	gamedata.SignTime = now
	gamedata.SignDay = gamedata.SignDay + 1
	gamedata.Dirty = true
	user.AddCoin(coin, tlogFangkuainiaoDaySign)
	handler.agent.bgSave()
	tlog.DaySign(user.Openid, userid, gamedata.SignDay)
	response.SignTime = now
	response.Coin = coin
	log.Printf("[FangkuainiaoHandler] DaySign succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) GetInvincibleCoin(r *gat.Request, args *gamepb.FangkuainiaoGetInvincibleCoinRequest) error {
	log.Printf("[FangkuainiaoHandler] GetInvincibleCoin %+v\n", args)
	response := &gamepb.FangkuainiaoGetInvincibleCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var coin int64 = args.Coin
	user.AddCoin(coin, tlogFangkuainiaoGetInvincibleCoin)
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) GetLevelCoin(r *gat.Request, args *gamepb.FangkuainiaoGetLevelCoinRequest) error {
	log.Printf("[FangkuainiaoHandler] GetLevelCoin %+v\n", args)
	response := &gamepb.FangkuainiaoGetLevelCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var coin int64 = args.Coin
	user.AddCoin(coin, tlogFangkuainiaoGetLevelCoin)
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

func (handler *FangkuainiaoHandler) GetCoin(r *gat.Request, args *gamepb.FangkuainiaoGetCoinRequest) error {
	log.Printf("[FangkuainiaoHandler] GetCoin %+v\n", args)
	response := &gamepb.FangkuainiaoGetCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	handler.resetGetCoin()
	if gamedata.GetGoldCount >= 1 {
		response.Code = errCodeFangkuainiaoGetCoin
		return r.Response(response)
	}
	var coin int64 = 100
	user.AddCoin(coin, tlogFangkuainiaoGetCoin)
	gamedata.GetGoldCount++
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

//关卡开始
func (handler *FangkuainiaoHandler) RoundStart(r *gat.Request, args *gamepb.CaichengyuStartRequest) error {
	log.Printf("[FangkuainiaoHandler] RoundStart %+v\n", args)
	response := &gamepb.CaichengyuStartResponse{}
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
func (handler *FangkuainiaoHandler) RoundResult(r *gat.Request, args *gamepb.FangkuainiaoResultRequest) error {
	log.Printf("[FangkuainiaoHandler] RoundResult %+v\n", args)
	response := &gamepb.FangkuainiaoResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	roundResultTime := time.Now().Unix()
	//duration := roundResultTime - user.roundStartTime
	//	log.Println("ofafsfasf", duration)
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	user.Dirty = true
	passLevel := gamedata.Level
	lastLevel := user.Level
	if args.Result != 1 {
		tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, passLevel, 0)
		return r.Response(response)
	}
	//通过关卡后奖励皮肤
	for _, birdConf := range config.Fangkuainiao.Bird {
		if birdConf.Type == 5 && passLevel >= birdConf.PayMoney {
			birdArr := strings.Split(gamedata.BirdArr, ",")
			found := false
			for _, _birdId := range birdArr {
				if birdId, err := strconv.Atoi(_birdId); err != nil {
					continue
				} else if int64(birdId) == birdConf.Id {
					found = true
					break
				}
			}
			//log.Println("fffffffffff1")
			if !found {
				gamedata.BirdArr = fmt.Sprintf("%s,%d", gamedata.BirdArr, birdConf.Id)
				log.Printf("fffffffffff1 %s\n", gamedata.BirdArr)
				gamedata.Dirty = true
			}
		}
	}
	if passLevel > lastLevel {
		user.Level = passLevel
	}
	user.LevelToday = user.LevelToday + 1
	user.Dirty = true
	gamedata.Level = passLevel + 1
	gamedata.Dirty = true
	handler.agent.bgSave()
	if passLevel > lastLevel {
		db.LevelRank_Update(userid, passLevel)
	}
	db.LevelDayRank_Update(user.Today, userid, user.LevelToday)
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, passLevel, 1)
	return r.Response(response)
}
