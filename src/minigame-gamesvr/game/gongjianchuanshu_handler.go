package game

import (
	"encoding/json"
	"log"
	"time"

	"github.com/roada-go/gat"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type GongjianchuanshuData struct {
	model.Gongjianchuanshu
	Dirty bool
}

type GongjianchuanshuHandler struct {
	agent    *Agent
	gamedata *GongjianchuanshuData
}

func newGongjianchuanshuHandler(agent *Agent) *GongjianchuanshuHandler {
	handler := &GongjianchuanshuHandler{
		agent: agent,
	}
	return handler
}

func (handler *GongjianchuanshuHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *GongjianchuanshuHandler) GetData(r *gat.Request, args *gamepb.GongjianchuanshuGetDataRequest) error {
	log.Printf("[GongjianchuanshuHandler] GetData %+v\n", args)
	response := &gamepb.GongjianchuanshuGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *GongjianchuanshuData
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
	response.SkinArr = gamedata.SkinArr
	response.SkinId = gamedata.SkinId
	response.ShopArr = gamedata.ShopArr
	response.Level = gamedata.Level
	response.GameId = conf.Ini.Game.Id
	log.Printf("[GongjianchuanshuHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *GongjianchuanshuHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *GongjianchuanshuHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[GongjianchuanshuHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Gongjianchuanshu_Save(&gamedata.Gongjianchuanshu)
	if err != nil {
		log.Printf("[GongjianchuanshuHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[GongjianchuanshuHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *GongjianchuanshuHandler) loadData(userid int64) (*GongjianchuanshuData, error) {
	dbrow, err := db.Gongjianchuanshu_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *GongjianchuanshuData
	if dbrow == nil {
		gamedata = &GongjianchuanshuData{
			model.Gongjianchuanshu{
				Userid:  userid,
				Level:   0,
				SkinArr: "[1,0,0,0,0,0,0,0,0,0]",
				SkinId:  1,
				ShopArr: "[2,3,4,5,6,7,8,9,10]",
			},
			false, //Dirty
		}
	} else {
		gamedata = &GongjianchuanshuData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	return gamedata, nil
}

func (handler *GongjianchuanshuHandler) AwardCoin(r *gat.Request, args *gamepb.GongjianchuanshuAwardCoinRequest) error {
	log.Printf("[GongjianchuanshuHandler] AwardCoin %+v\n", args)
	response := &gamepb.GongjianchuanshuAwardCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var coin int64 = args.Coin
	user.AddCoin(coin, tlogGongjianchuanshuGetAwardCoin)
	gamedata.Dirty = true
	response.Coin = coin
	return r.Response(response)
}

//购买皮肤
func (handler *GongjianchuanshuHandler) UnlockSkin(r *gat.Request, args *gamepb.GongjianchuanshuUnlockSkinRequest) error {
	log.Printf("[GongjianchuanshuHandler] UnlockSkin %+v\n", args)
	response := &gamepb.GongjianchuanshuUnlockSkinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	index := int(args.ShopIndex)
	user := handler.agent.user
	gamedata := handler.gamedata
	shopArr := make([]int32, 0)
	if err := json.Unmarshal([]byte(gamedata.ShopArr), &shopArr); err != nil {
		log.Printf("[GongjianchuanshuHandler] UnlockSkin json.Unmarshal failed, error=%s\n", err.Error())
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	}
	skinArr := make([]int32, 0)
	if err := json.Unmarshal([]byte(gamedata.SkinArr), &skinArr); err != nil {
		log.Printf("[GongjianchuanshuHandler] UnlockSkin json.Unmarshal failed, error=%s\n", err.Error())
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	}
	if index < 0 || index >= len(shopArr) {
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	}
	//已解锁
	if skinArr[shopArr[index]-1] == 1 {
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	}
	skinArr[shopArr[index]-1] = 1
	if bytes, err := json.Marshal(skinArr); err != nil {
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	} else {
		gamedata.SkinArr = string(bytes)
		gamedata.Dirty = true
	}
	var coin int64 = 30
	if user.Coin < coin {
		response.Code = errCodeGongjianchuanshuUnlockSkin
		return r.Response(response)
	}
	user.DecCoin(coin, tlogGongjianchuanshuUnlockSkin)
	handler.agent.bgSave()
	return r.Response(response)
}

//使用皮肤
func (handler *GongjianchuanshuHandler) UseSkin(r *gat.Request, args *gamepb.GongjianchuanshuUseSkinRequest) error {
	log.Printf("[GongjianchuanshuHandler] UseSkin %+v\n", args)
	response := &gamepb.GongjianchuanshuUseSkinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	index := int(args.ShopIndex)
	gamedata := handler.gamedata
	shopArr := make([]int32, 0)
	if err := json.Unmarshal([]byte(gamedata.ShopArr), &shopArr); err != nil {
		log.Printf("[GongjianchuanshuHandler] UseSkin json.Unmarshal failed, error=%s\n", err.Error())
		response.Code = errCodeGongjianchuanshuUseSkin
		return r.Response(response)
	}
	skinArr := make([]int32, 0)
	if err := json.Unmarshal([]byte(gamedata.SkinArr), &skinArr); err != nil {
		log.Printf("[GongjianchuanshuHandler] UseSkin json.Unmarshal failed, error=%s\n", err.Error())
		response.Code = errCodeGongjianchuanshuUseSkin
		return r.Response(response)
	}
	if index < 0 || index >= len(shopArr) {
		response.Code = errCodeGongjianchuanshuUseSkin
		return r.Response(response)
	}
	//末解锁
	if skinArr[shopArr[index]-1] == 0 {
		response.Code = errCodeGongjianchuanshuUseSkin
		return r.Response(response)
	}
	//log.Println("aaaaaa", shopArr, gamedata.SkinId)
	lastShopSkin := shopArr[index]
	shopArr[index] = gamedata.SkinId
	if bytes, err := json.Marshal(shopArr); err != nil {
		response.Code = errCodeGongjianchuanshuUseSkin
		return r.Response(response)
	} else {
		gamedata.ShopArr = string(bytes)
		gamedata.SkinId = lastShopSkin
		gamedata.Dirty = true
	}
	//log.Println("aaaaaa", shopArr, gamedata.SkinId, gamedata.ShopArr)
	gamedata.Dirty = true
	//handler.agent.bgSave()
	return r.Response(response)
}

//关卡开始
func (handler *GongjianchuanshuHandler) LevelStart(r *gat.Request, args *gamepb.GongjianchuanshuLevelStartRequest) error {
	log.Printf("[GongjianchuanshuHandler] LevelStart %+v\n", args)
	response := &gamepb.GongjianchuanshuLevelStartResponse{}
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
func (handler *GongjianchuanshuHandler) LevelResult(r *gat.Request, args *gamepb.GongjianchuanshuLevelResultRequest) error {
	log.Printf("[GongjianchuanshuHandler] LevelResult %+v\n", args)
	response := &gamepb.GongjianchuanshuLevelResultResponse{}
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
	if score > user.Medal {
		user.Medal = score
		db.MedalRank_Update(userid, score)
	}
	if score > user.MedalToday {
		user.MedalToday = score
		db.MedalDayRank_Update(user.Today, userid, score)
	}
	return r.Response(response)
}

//关卡过关
func (handler *GongjianchuanshuHandler) LevelPass(r *gat.Request, args *gamepb.GongjianchuanshuLevelPassRequest) error {
	log.Printf("[GongjianchuanshuHandler] LevelPass %+v\n", args)
	response := &gamepb.GongjianchuanshuLevelPassResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	//handler.agent.checkToday()
	level := args.Level
	gamedata := handler.gamedata
	gamedata.Level = level
	gamedata.Dirty = true
	return r.Response(response)
}

func (handler *GongjianchuanshuHandler) RoundStart(r *gat.Request, args *gamepb.GongjianchuanshuStartRequest) error {
	log.Printf("[GongjianchuanshuHandler] RoundStart %+v\n", args)
	response := &gamepb.GongjianchuanshuStartResponse{}
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

func (handler *GongjianchuanshuHandler) RoundResult(r *gat.Request, args *gamepb.GongjianchuanshuResultRequest) error {
	log.Printf("[GongjianchuanshuHandler] RoundResult %+v\n", args)
	response := &gamepb.GongjianchuanshuResultResponse{}
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
