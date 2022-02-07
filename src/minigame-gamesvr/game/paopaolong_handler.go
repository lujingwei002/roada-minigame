package game

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/roada-go/gat"
	"github.com/roada-go/util/timeutil"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/config"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type PaopaolongData struct {
	model.Paopaolong
	Dirty bool
}

type PaopaolongHandler struct {
	agent     *Agent
	gamedata  *PaopaolongData
	ItemArr   []int32
	levelDict map[int64]*model.PaopaolongLevel
}

func newPaopaolongHandler(agent *Agent) *PaopaolongHandler {
	handler := &PaopaolongHandler{
		agent: agent,
	}
	return handler
}

func (handler *PaopaolongHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *PaopaolongHandler) GetData(r *gat.Request, args *gamepb.PaopaolongGetDataRequest) error {
	log.Printf("[PaopaolongHandler] GetData %+v\n", args)
	response := &gamepb.PaopaolongGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *PaopaolongData
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
	handler.checkToday()
	response.ItemArr = gamedata.ItemArr
	response.Level = gamedata.Level
	response.Hp = gamedata.Hp
	response.FreedrawTime = gamedata.FreedrawTime
	response.NewPackRedeemed = gamedata.NewPackRedeemed
	response.ShopFreeDiamondTime = gamedata.ShopFreeDiamondTime
	response.ShopFreeDiamondTime2 = gamedata.ShopFreeDiamondTime2
	response.LastSignTime = gamedata.LastSignTime
	response.SignedTime = gamedata.SignedTime
	response.GameId = conf.Ini.Game.Id
	response.LevelArr = make([]*gamepb.PaopaolongLevel, 0)
	for _, level := range handler.levelDict {
		response.LevelArr = append(response.LevelArr, &gamepb.PaopaolongLevel{
			Level: level.Level,
			Sec:   level.Sec,
			Lose:  level.Lose,
			Score: level.Score,
		})
	}
	log.Printf("[PaopaolongHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *PaopaolongHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *PaopaolongHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[PaopaolongHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	handler.ItemArr[4] = gamedata.Hp
	if bytes, err := json.Marshal(handler.ItemArr); err != nil {
		return err
	} else {
		gamedata.ItemArr = string(bytes)
	}
	err := db.Paopaolong_Save(&gamedata.Paopaolong)
	if err != nil {
		log.Printf("[PaopaolongHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[PaopaolongHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *PaopaolongHandler) loadData(userid int64) (*PaopaolongData, error) {
	dbrow, err := db.Paopaolong_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *PaopaolongData
	if dbrow == nil {
		t := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local)
		now := t.Unix()
		gamedata = &PaopaolongData{
			Dirty: false,
		}
		gamedata.Paopaolong = model.Paopaolong{
			Userid:               userid,
			Level:                1,
			ItemArr:              "[0,0,0,0,0,0,0]",
			Hp:                   5,
			NewPackRedeemed:      0,
			FreedrawTime:         now,
			ShopFreeDiamondTime:  now,
			ShopFreeDiamondTime2: now,
			LastSignTime:         now,
			SignedTime:           0,
		}
	} else {
		gamedata = &PaopaolongData{
			Dirty: false,
		}
		gamedata.Paopaolong = *dbrow
	}
	handler.ItemArr = make([]int32, 0)
	if err := json.Unmarshal([]byte(gamedata.ItemArr), &handler.ItemArr); err != nil {
		return nil, err
	}
	handler.ItemArr[4] = gamedata.Hp

	levelArr, err := db.Paopaolong_LevelGet(userid)
	if err != nil {
		return nil, err
	}
	handler.levelDict = make(map[int64]*model.PaopaolongLevel)
	for _, level := range levelArr {
		handler.levelDict[level.Level] = level
	}
	return gamedata, nil
}

func (handler *PaopaolongHandler) checkToday() {
	gamedata := handler.gamedata
	//七天后清0
	if timeutil.PassDays(gamedata.LastSignTime, time.Now().Unix()) >= 1 && gamedata.SignedTime >= 7 {
		gamedata.SignedTime = 0
		gamedata.Dirty = true
	}
}

func (handler *PaopaolongHandler) getRandomInLuckyPool() int {
	var total int32 = 0
	for _, conf := range config.Paopaolong.Lucky {
		total = total + conf.Chance
	}
	val := rand.Int31n(total)
	for index, conf := range config.Paopaolong.Lucky {
		if val < conf.Chance {
			return index
		}
		val -= conf.Chance
	}
	return 0
}

func (handler *PaopaolongHandler) AddItem(itemId int32, num int32, reason string) error {
	if itemId < 0 || itemId >= int32(len(handler.ItemArr)) {
		return fmt.Errorf("item id format err")
	}
	user := handler.agent.user
	if itemId == 5 {
		user.AddCoin(int64(num), reason)
	} else if itemId == 4 {
		handler.gamedata.Hp += num
	} else {
		handler.ItemArr[itemId] += num
	}
	handler.gamedata.Dirty = true
	return nil
}

func (handler *PaopaolongHandler) DecItem(itemId int32, num int32, reason string) error {
	if itemId < 0 || itemId >= int32(len(handler.ItemArr)) {
		return fmt.Errorf("item id format err")
	}
	user := handler.agent.user
	if itemId == 5 {
		if user.Coin < int64(num) {
			return fmt.Errorf("item num not enough")
		}
		return user.DecCoin(int64(num), reason)
	} else if itemId == 4 {
		if handler.gamedata.Hp < num {
			return fmt.Errorf("item num not enough")
		}
		handler.gamedata.Hp -= num
		return nil
	} else {
		if handler.ItemArr[itemId] < num {
			return fmt.Errorf("item num not enough")
		}
		handler.ItemArr[itemId] -= num
	}
	handler.gamedata.Dirty = true
	return nil
}

func (handler *PaopaolongHandler) GetItemNum(itemId int32) int32 {
	user := handler.agent.user
	if itemId == 5 {
		return int32(user.Coin)
	} else if itemId == 4 {
		return handler.gamedata.Hp
	} else {
		return handler.ItemArr[itemId]
	}
}

func (handler *PaopaolongHandler) ShopBuy(r *gat.Request, args *gamepb.PaopaolongShopBuyRequest) error {
	log.Printf("[PaopaolongHandler] ShopBuy %+v\n", args)
	response := &gamepb.PaopaolongShopBuyResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	index := args.Id
	if index < 0 || index >= int32(len(config.Paopaolong.Shop)) {
		response.Code = errCodePaopaolongShopBuy
		return r.Response(response)
	}
	conf := config.Paopaolong.Shop[index]
	var coin int64 = conf.Cost
	if user.Coin < coin {
		response.Code = errCodePaopaolongShopBuy
		return r.Response(response)
	}
	user.DecCoin(coin, tlogPaopaolongShopBuy)
	handler.AddItem(conf.Item, conf.Num, tlogPaopaolongShopBuy)
	response.Coin = coin
	response.Id = args.Id
	return r.Response(response)
}

//新手礼包
func (handler *PaopaolongHandler) NewPack(r *gat.Request, args *gamepb.PaopaolongNewPackRequest) error {
	log.Printf("[PaopaolongHandler] NewPack %+v\n", args)
	response := &gamepb.PaopaolongNewPackResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if gamedata.NewPackRedeemed != 0 {
		response.Code = errCodePaopaolongNewPack
		return r.Response(response)
	}
	gamedata.NewPackRedeemed = 1
	for _, conf := range config.Paopaolong.NewPack {
		handler.AddItem(conf.Item, conf.Num, tlogPaopaolongNewPack)
	}
	gamedata.Dirty = true
	handler.agent.bgSave()
	return r.Response(response)
}

//领取商店免费金币
func (handler *PaopaolongHandler) ShopFreeCoin(r *gat.Request, args *gamepb.PaopaolongShopFreeCoinRequest) error {
	log.Printf("[PaopaolongHandler] ShopFreeCoin %+v\n", args)
	response := &gamepb.PaopaolongShopFreeCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.ShopFreeDiamondTime, time.Now().Unix()) <= 0 {
		response.Code = errCodePaopaolongShopCoinFree
		return r.Response(response)
	}
	gamedata.ShopFreeDiamondTime = time.Now().Unix()
	handler.AddItem(5, 100, tlogPaopaolongShopFreeCoin)
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.Coin = 100
	return r.Response(response)
}

//领取商店免费金币2
func (handler *PaopaolongHandler) ShopFreeCoin2(r *gat.Request, args *gamepb.PaopaolongShopFreeCoin2Request) error {
	log.Printf("[PaopaolongHandler] ShopFreeCoin2 %+v\n", args)
	response := &gamepb.PaopaolongShopFreeCoin2Response{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.ShopFreeDiamondTime2, time.Now().Unix()) <= 0 {
		response.Code = errCodePaopaolongShopCoinFree
		return r.Response(response)
	}
	gamedata.ShopFreeDiamondTime2 = time.Now().Unix()
	handler.AddItem(5, 100, tlogPaopaolongShopFreeCoin)
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.Coin = 100
	return r.Response(response)
}

//签到
func (handler *PaopaolongHandler) Sign(r *gat.Request, args *gamepb.PaopaolongSignRequest) error {
	log.Printf("[PaopaolongHandler] Sign %+v\n", args)
	response := &gamepb.PaopaolongSignResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.checkToday()
	gamedata := handler.gamedata
	index := gamedata.SignedTime
	if index < 0 || index >= int32(len(config.Paopaolong.Sign)) {
		response.Code = errCodePaopaolongSign
		return r.Response(response)
	}
	conf := config.Paopaolong.Sign[index]
	if timeutil.PassDays(gamedata.LastSignTime, time.Now().Unix()) <= 0 {
		response.Code = errCodePaopaolongFreeDraw
		return r.Response(response)
	}
	gamedata.LastSignTime = time.Now().Unix()
	gamedata.SignedTime++
	handler.AddItem(conf.Item, conf.Num, tlogPaopaolongFreeDraw)
	gamedata.Dirty = true
	handler.agent.bgSave()
	return r.Response(response)
}

//免费抽奖
func (handler *PaopaolongHandler) FreeDraw(r *gat.Request, args *gamepb.PaopaolongFreeDrawRequest) error {
	log.Printf("[PaopaolongHandler] FreeDraw %+v\n", args)
	response := &gamepb.PaopaolongFreeDrawResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.FreedrawTime, time.Now().Unix()) <= 0 {
		response.Code = errCodePaopaolongFreeDraw
		return r.Response(response)
	}
	index := handler.getRandomInLuckyPool()
	conf := config.Paopaolong.Lucky[index]
	handler.AddItem(conf.Item, conf.Num, tlogPaopaolongFreeDraw)
	gamedata.FreedrawTime = time.Now().Unix()
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.LuckyId = int32(index)
	response.FreedrawTime = gamedata.FreedrawTime
	return r.Response(response)
}

//消耗游戏币抽奖
func (handler *PaopaolongHandler) CostDraw(r *gat.Request, args *gamepb.PaopaolongCostDrawRequest) error {
	log.Printf("[PaopaolongHandler] CostDraw %+v\n", args)
	response := &gamepb.PaopaolongCostDrawResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	//userid := user.Userid
	var coin int64 = 200
	if user.Coin < coin {
		response.Code = errCodePaopaolongCostDraw
		return r.Response(response)
	}
	index := handler.getRandomInLuckyPool()
	conf := config.Paopaolong.Lucky[index]
	user.DecCoin(coin, tlogPaopaolongCostDraw)
	handler.AddItem(conf.Item, conf.Num, tlogPaopaolongCostDraw)
	response.Coin = coin
	response.LuckyId = int32(index)
	return r.Response(response)
}

//使用道具
func (handler *PaopaolongHandler) UseItem(r *gat.Request, args *gamepb.PaopaolongUseItemRequest) error {
	log.Printf("[PaopaolongHandler] UseItem %+v\n", args)
	response := &gamepb.PaopaolongUseItemResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}

	gamedata := handler.gamedata
	itemId := args.ItemId
	if itemId < 0 || itemId >= int32(len(handler.ItemArr)) {
		response.Code = errCodePaopaolongUseItem
		return r.Response(response)
	}
	if handler.ItemArr[itemId] <= 0 {
		response.Code = errCodePaopaolongStart
		return r.Response(response)
	}
	handler.ItemArr[itemId] = handler.ItemArr[itemId] - 1
	gamedata.Dirty = true
	response.ItemId = itemId
	return r.Response(response)
}

//关卡开始
func (handler *PaopaolongHandler) RoundStart(r *gat.Request, args *gamepb.PaopaolongStartRequest) error {
	log.Printf("[PaopaolongHandler] RoundStart %+v\n", args)
	response := &gamepb.PaopaolongStartResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	/*if gamedata.Hp <= 0 {
		response.Code = errCodePaopaolongStart
		return r.Response(response)
	}*/
	if args.Level > gamedata.Level {
		response.Code = errCodePaopaolongStart
		return r.Response(response)
	}
	//gamedata.Hp = gamedata.Hp - 1
	gamedata.Dirty = true
	user.RoundStartTime = time.Now().Unix()
	user.Dirty = true
	tlog.RoundStart(user.Openid, userid)
	response.Level = args.Level
	return r.Response(response)
}

//关卡结束
func (handler *PaopaolongHandler) RoundResult(r *gat.Request, args *gamepb.PaopaolongResultRequest) error {
	log.Printf("[PaopaolongHandler] RoundResult %+v\n", args)
	response := &gamepb.PaopaolongResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	roundResultTime := time.Now().Unix()
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	if args.Result == 1 {
		if args.Level == gamedata.Level {
			gamedata.Level = args.Level + 1
			gamedata.Dirty = true
		}
	}
	var addStar int64 = 0
	//保存关卡记录
	if _, ok := handler.levelDict[args.Level]; !ok {
		handler.levelDict[args.Level] = &model.PaopaolongLevel{
			Userid: userid,
			Level:  args.Level,
			Sec:    0,
			Lose:   0,
			Score:  0,
			Star:   0,
		}
	}
	if level, ok := handler.levelDict[args.Level]; ok {
		if args.Result == 0 {
			level.Lose++
		} else {
			level.Sec = args.Sec
			if args.Score > level.Score {
				level.Score = args.Score
			}
			if args.Star > level.Star {
				addStar = args.Star - level.Star
				level.Star = args.Star
			}
		}
		db.Paopaolong_LevelSave(userid, level.Level, level.Sec, level.Lose, level.Score, level.Star)
	}
	if args.Result == 1 {
		var totalStar int64 = 0
		for _, level := range handler.levelDict {
			totalStar = totalStar + level.Star
		}
		if totalStar > user.Score {
			user.Score = totalStar
			db.ScoreRank_Update(userid, totalStar)
		}
		if addStar > 0 {
			user.ScoreToday += addStar
			db.ScoreDayRank_Update(user.Today, userid, user.ScoreToday)
		}
		//log.Println("gggggggggg", totalStar, addStar, user.ScoreToday)
		tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, totalStar, 1)
	} else {
		var totalStar int64 = 0
		for _, level := range handler.levelDict {
			totalStar = totalStar + level.Star
		}
		tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, totalStar, 0)
	}
	return r.Response(response)
}
