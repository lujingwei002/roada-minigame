package game

import (
	"encoding/json"
	"log"
	"math"
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

type HuanlemaomibeiData struct {
	model.Huanlemaomibei
	Dirty bool
}

type HuanlemaomibeiHandler struct {
	agent     *Agent
	gamedata  *HuanlemaomibeiData
	CupDict   map[int32]int32
	InkDict   map[int32]int32
	levelDict map[int32]*model.HuanlemaomibeiLevel
}

func newHuanlemaomibeiHandler(agent *Agent) *HuanlemaomibeiHandler {
	handler := &HuanlemaomibeiHandler{
		agent:   agent,
		CupDict: make(map[int32]int32),
		InkDict: make(map[int32]int32),
	}
	return handler
}

func (handler *HuanlemaomibeiHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *HuanlemaomibeiHandler) GetData(r *gat.Request, args *gamepb.HuanlemaomibeiGetDataRequest) error {
	log.Printf("[HuanlemaomibeiHandler] GetData %+v\n", args)
	response := &gamepb.HuanlemaomibeiGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *HuanlemaomibeiData
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
	response.CupArr = gamedata.CupArr
	response.InkArr = gamedata.InkArr
	response.CupId = gamedata.CupId
	response.InkId = gamedata.InkId
	response.FreedrawTime = gamedata.FreedrawTime
	response.OfflineTime = gamedata.OfflineTime
	response.GameId = conf.Ini.Game.Id
	response.LevelArr = make([]*gamepb.HuanlemaomibeiLevel, 0)
	for _, level := range handler.levelDict {
		coin := level.Coin
		if coin > 0 {
			coin = 1
		}
		response.LevelArr = append(response.LevelArr, &gamepb.HuanlemaomibeiLevel{
			Section: level.Section,
			Level:   level.Level,
			Unlock:  level.Unlock,
			Star:    level.Star,
			Coin:    coin,
		})
	}
	response.Hp = gamedata.Hp
	log.Printf("[HuanlemaomibeiHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *HuanlemaomibeiHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *HuanlemaomibeiHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[HuanlemaomibeiHandler] saveData, Dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	if bytes, err := json.Marshal(handler.InkDict); err != nil {
		log.Println(err)
		return err
	} else {
		gamedata.InkArr = string(bytes)
	}
	//log.Println("ccc", gamedata.InkArr)
	if bytes, err := json.Marshal(handler.CupDict); err != nil {
		log.Println(err)
		return err
	} else {
		gamedata.CupArr = string(bytes)
	}
	err := db.Huanlemaomibei_Save(&gamedata.Huanlemaomibei)
	if err != nil {
		log.Printf("[HuanlemaomibeiHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[HuanlemaomibeiHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *HuanlemaomibeiHandler) loadData(userid int64) (*HuanlemaomibeiData, error) {
	dbrow, err := db.Huanlemaomibei_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *HuanlemaomibeiData
	if dbrow == nil {
		gamedata = &HuanlemaomibeiData{
			model.Huanlemaomibei{
				Userid:       userid,
				FreedrawTime: 0,
				InkArr:       "{}",
				CupArr:       "{}",
				LastSignTime: 0,
				SignDay:      0,
				SignChecked:  0,
				OfflineTime:  time.Now().Unix(),
			},
			true, //Dirty
		}
	} else {
		gamedata = &HuanlemaomibeiData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	if err := json.Unmarshal([]byte(gamedata.InkArr), &handler.InkDict); err != nil {
		log.Println(err)
		return nil, err
	}
	if err := json.Unmarshal([]byte(gamedata.CupArr), &handler.CupDict); err != nil {
		log.Println(err)
		return nil, err
	}

	levelArr, err := db.Huanlemaomibei_LevelGet(userid)
	if err != nil {
		return nil, err
	}
	handler.levelDict = make(map[int32]*model.HuanlemaomibeiLevel)
	for _, level := range levelArr {
		handler.levelDict[level.Section*100+level.Level] = level
	}
	//log.Println("gggg", handler.InkDict)
	return gamedata, nil
}

func (handler *HuanlemaomibeiHandler) SignQuery(r *gat.Request, args *gamepb.HuanlemaomibeiSignQueryRequest) error {
	log.Printf("[HuanlemaomibeiHandler] SignQuery %+v\n", args)
	response := &gamepb.HuanlemaomibeiSignQueryResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	var signToday int32 = 0
	if timeutil.IsSameDay(gamedata.LastSignTime, time.Now().Unix()) {
		signToday = 1
	}
	if gamedata.LastSignTime != 0 && signToday == 0 && gamedata.SignChecked == 0 {
		gamedata.SignChecked = 1
		gamedata.SignDay = (gamedata.SignDay + 1) % 7
		gamedata.Dirty = true
	}
	response.SignToday = signToday
	response.SignDay = gamedata.SignDay
	return r.Response(response)
}

func (handler *HuanlemaomibeiHandler) Sign(r *gat.Request, args *gamepb.HuanlemaomibeiSignRequest) error {
	log.Printf("[HuanlemaomibeiHandler] Sign %+v\n", args)
	response := &gamepb.HuanlemaomibeiSignResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	var signToday int32 = 0
	if timeutil.IsSameDay(gamedata.LastSignTime, time.Now().Unix()) {
		signToday = 1
	}
	if signToday == 1 {
		response.Code = errCodeHuanlemaomibeiSign
		return r.Response(response)
	}
	conf := config.Huanlemaomibei.Sign[gamedata.SignDay]
	if conf.Ink != 0 || conf.Cup != 0 {
		if conf.Ink != 0 {
			log.Println("add ink", conf.Ink)
			handler.InkDict[conf.Ink] = -1
			gamedata.Dirty = true
		}
		if conf.Cup != 0 {
			log.Println("add cup", conf.Cup)
			handler.CupDict[conf.Cup] = -1
			gamedata.Dirty = true
		}
	} else {
		log.Println("add gold", conf.Gold)
		user.AddCoin(conf.Gold, tlogHuanlemaomibeiSign)
	}

	gamedata.LastSignTime = time.Now().Unix()
	gamedata.SignChecked = 0
	gamedata.Dirty = true
	return r.Response(response)
}

func (handler *HuanlemaomibeiHandler) getRandomInLuckyPool() int {
	var total int32 = 0
	for _, conf := range config.Huanlemaomibei.Roulette {
		total = total + conf.Weight
	}
	val := rand.Int31n(total)
	for index, conf := range config.Huanlemaomibei.Roulette {
		if val < conf.Weight {
			return index
		}
		val -= conf.Weight
	}
	return 0
}

func (handler *HuanlemaomibeiHandler) FlyQuery(r *gat.Request, args *gamepb.HuanlemaomibeiFlyQueryRequest) error {
	log.Printf("[HuanlemaomibeiHandler] FlyQuery %+v\n", args)
	response := &gamepb.HuanlemaomibeiFlyQueryResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.FlyTime, time.Now().Unix()) <= 0 {
		response.FlyCount = 0
		return r.Response(response)
	} else {
		response.FlyCount = 1
		return r.Response(response)
	}
}

//免费抽奖
func (handler *HuanlemaomibeiHandler) Fly(r *gat.Request, args *gamepb.HuanlemaomibeiFlyRequest) error {
	log.Printf("[HuanlemaomibeiHandler] Fly %+v\n", args)
	response := &gamepb.HuanlemaomibeiFlyResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.FlyTime, time.Now().Unix()) <= 0 {
		response.Code = errCodeHuanlemaomibeiFly
		return r.Response(response)
	}
	index := handler.getRandomInLuckyPool()
	//index = 5
	conf := config.Huanlemaomibei.Roulette[index]
	if conf.Skin == 0 {
		log.Println("add gold", conf.Gold)
		user.AddCoin(conf.Gold, tlogHuanlemaomibeiFly)
	} else if conf.Skin == 1 {
		//皮肤是否已经领完，领完就奖励金币
		pool := make([]int32, 0)
		for _, conf := range config.Huanlemaomibei.Ink {
			if count, ok := handler.InkDict[conf.Id]; !(ok && count == -1) {
				pool = append(pool, conf.Id)
				break
			}
		}
		if len(pool) <= 0 {
			//log.Println("1add gold", conf.Gold)
			//user.AddCoin(conf.Gold, tlogHuanlemaomibeiFly)
		} else {
			index := rand.Intn(len(pool))
			conf := config.Huanlemaomibei.Ink[pool[index]-1]
			handler.InkDict[conf.Id] = -1
			gamedata.Dirty = true
			log.Println("1add ink", conf.Id, pool)
		}
	} else if conf.Skin == 2 {
		//皮肤是否已经领完，领完就奖励金币
		pool := make([]int32, 0)
		for _, conf := range config.Huanlemaomibei.Cup {
			if count, ok := handler.CupDict[conf.Id]; !(ok && count == -1) {
				pool = append(pool, conf.Id)
				break
			}
		}
		if len(pool) <= 0 {
			//log.Println("1add gold", conf.Gold)
			//user.AddCoin(conf.Gold, tlogHuanlemaomibeiFly)
		} else {
			index := rand.Intn(len(pool))
			conf := config.Huanlemaomibei.Cup[pool[index]-1]
			handler.CupDict[conf.Id] = -1
			gamedata.Dirty = true
			log.Println("1add cup", conf.Id, pool)
		}
	}

	gamedata.FlyTime = time.Now().Unix()
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.GiftIndex = int32(index)
	response.FlyTime = gamedata.FlyTime
	return r.Response(response)
}

func (handler *HuanlemaomibeiHandler) FreeDrawQuery(r *gat.Request, args *gamepb.HuanlemaomibeiFreeDrawQueryRequest) error {
	log.Printf("[HuanlemaomibeiHandler] FreeDrawQuery %+v\n", args)
	response := &gamepb.HuanlemaomibeiFreeDrawQueryResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.FreedrawTime, time.Now().Unix()) <= 0 {
		response.FreedrawCount = 0
		return r.Response(response)
	} else {
		response.FreedrawCount = 1
		return r.Response(response)
	}
}

//免费抽奖
func (handler *HuanlemaomibeiHandler) FreeDraw(r *gat.Request, args *gamepb.HuanlemaomibeiFreeDrawRequest) error {
	log.Printf("[HuanlemaomibeiHandler] FreeDraw %+v\n", args)
	response := &gamepb.HuanlemaomibeiFreeDrawResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	if timeutil.PassDays(gamedata.FreedrawTime, time.Now().Unix()) <= 0 {
		response.Code = errCodeHuanlemaomibeiFreeDraw
		return r.Response(response)
	}
	index := handler.getRandomInLuckyPool()
	//index = 5
	conf := config.Huanlemaomibei.Roulette[index]
	if conf.Skin == 0 {
		log.Println("add gold", conf.Gold)
		user.AddCoin(conf.Gold, tlogHuanlemaomibeiFreeDraw)
	} else if conf.Skin == 1 {
		//皮肤是否已经领完，领完就奖励金币
		pool := make([]int32, 0)
		for _, conf := range config.Huanlemaomibei.Ink {
			if conf.Roulette == 1 {
				if count, ok := handler.InkDict[conf.Id]; !(ok && count == -1) {
					pool = append(pool, conf.Id)
					break
				}
			}
		}
		if len(pool) <= 0 {
			log.Println("1add gold", conf.Gold)
			user.AddCoin(conf.Gold, tlogHuanlemaomibeiFreeDraw)
		} else {
			index := rand.Intn(len(pool))
			conf := config.Huanlemaomibei.Ink[pool[index]-1]
			handler.InkDict[conf.Id] = -1
			gamedata.Dirty = true
			log.Println("1add ink", conf.Id, pool)
		}
	} else if conf.Skin == 2 {
		//皮肤是否已经领完，领完就奖励金币
		pool := make([]int32, 0)
		for _, conf := range config.Huanlemaomibei.Cup {
			if conf.Roulette == 1 {
				if count, ok := handler.CupDict[conf.Id]; !(ok && count == -1) {
					pool = append(pool, conf.Id)
					break
				}
			}
		}
		if len(pool) <= 0 {
			log.Println("1add gold", conf.Gold)
			user.AddCoin(conf.Gold, tlogHuanlemaomibeiFreeDraw)
		} else {
			index := rand.Intn(len(pool))
			conf := config.Huanlemaomibei.Cup[pool[index]-1]
			handler.CupDict[conf.Id] = -1
			gamedata.Dirty = true
			log.Println("1add cup", conf.Id, pool)
		}
	}
	//handler.AddItem(conf.Item, conf.Num, tlogHuanlemaomibeiFreeDraw)
	gamedata.FreedrawTime = time.Now().Unix()
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.GiftIndex = int32(index)
	response.FreedrawTime = gamedata.FreedrawTime
	return r.Response(response)
}

//离线奖励
func (handler *HuanlemaomibeiHandler) OfflineCoin(r *gat.Request, args *gamepb.HuanlemaomibeiOfflineCoinRequest) error {
	log.Printf("[HuanlemaomibeiHandler] OfflineCoin %+v\n", args)
	response := &gamepb.HuanlemaomibeiOfflineCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata

	var coin int64 = 0
	// 每小时 20，上限 12 小时
	duration := int64(math.Floor(float64(time.Now().Unix()-gamedata.OfflineTime) / 3600.0))
	if duration < 0 {
		duration = 0
	}
	if duration > 12 {
		duration = 12
	}
	coin = duration * 20

	user.AddCoin(coin, tlogHuanlemaomibeiOfflineCoin)

	gamedata.OfflineTime = time.Now().Unix()
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.Coin = coin
	response.OfflineTime = gamedata.OfflineTime
	return r.Response(response)
}

//使用皮肤
func (handler *HuanlemaomibeiHandler) UseSkin(r *gat.Request, args *gamepb.HuanlemaomibeiUseSkinRequest) error {
	log.Printf("[HuanlemaomibeiHandler] UseSkin %+v\n", args)
	response := &gamepb.HuanlemaomibeiUseSkinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	//user := handler.agent.user
	gamedata := handler.gamedata
	if args.InkId > 0 {
		if args.InkId <= 0 || args.InkId > int32(len(config.Huanlemaomibei.Ink)) {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		//conf := config.Huanlemaomibei.Ink[args.InkId-1]
		if count, ok := handler.InkDict[args.InkId]; !ok || count == 0 {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		gamedata.InkId = args.InkId
	}
	if args.CupId > 0 {
		if args.CupId <= 0 || args.CupId > int32(len(config.Huanlemaomibei.Cup)) {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		// := config.Huanlemaomibei.Cup[args.CupId-1]
		if count, ok := handler.CupDict[args.CupId]; !ok || count == 0 {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		gamedata.CupId = args.CupId
	}
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.InkId = args.InkId
	response.CupId = args.CupId

	return r.Response(response)
}

//使用皮肤
func (handler *HuanlemaomibeiHandler) BuySkin(r *gat.Request, args *gamepb.HuanlemaomibeiBuySkinRequest) error {
	log.Printf("[HuanlemaomibeiHandler] BuySkin %+v\n", args)
	response := &gamepb.HuanlemaomibeiBuySkinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	response.Coin = 0
	if args.InkId > 0 {
		if args.InkId <= 0 || args.InkId > int32(len(config.Huanlemaomibei.Ink)) {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		conf := config.Huanlemaomibei.Ink[args.InkId-1]
		if conf.Coin == 0 {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		var coin int64 = conf.Coin
		if user.Coin < coin {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		user.DecCoin(coin, tlogHuanlemaomibeiBuySkin)
		handler.InkDict[args.InkId] = -1
		gamedata.InkId = args.InkId
		response.Coin += coin
	}
	if args.CupId > 0 {
		if args.CupId <= 0 || args.CupId > int32(len(config.Huanlemaomibei.Cup)) {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		conf := config.Huanlemaomibei.Cup[args.CupId-1]
		if conf.Coin == 0 {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		var coin int64 = conf.Coin
		if user.Coin < coin {
			response.Code = errCodeHuanlemaomibeiUseSkin
			return r.Response(response)
		}
		user.DecCoin(coin, tlogHuanlemaomibeiBuySkin)
		handler.CupDict[args.CupId] = -1
		gamedata.CupId = args.CupId
		response.Coin += coin
	}
	gamedata.Dirty = true
	handler.agent.bgSave()
	response.InkId = args.InkId
	response.CupId = args.CupId

	return r.Response(response)
}

//关卡解锁
func (handler *HuanlemaomibeiHandler) LevelUnlock(r *gat.Request, args *gamepb.HuanlemaomibeiLevelUnlockRequest) error {
	log.Printf("[HuanlemaomibeiHandler] LevelUnlock %+v\n", args)
	response := &gamepb.HuanlemaomibeiLevelUnlockResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	//gamedata := handler.gamedata

	if args.Section < 0 || int(args.Section) >= len(config.Huanlemaomibei.Section) {
		log.Println("1")
		response.Code = errCodeHuanlemaomibeiLevelUnlock
		return r.Response(response)
	}
	sectionConf := config.Huanlemaomibei.Section[args.Section]
	if args.Level < 0 || int(args.Level) >= len(sectionConf.LevelArr) {
		log.Println("2")
		response.Code = errCodeHuanlemaomibeiLevelUnlock
		return r.Response(response)
	}
	levelConf := sectionConf.LevelArr[args.Level]
	coin := levelConf.Price
	if user.Coin < coin {
		log.Println("3")
		response.Code = errCodeHuanlemaomibeiLevelUnlock
		return r.Response(response)
	}
	if args.Section != 0 && args.Level != 0 {
		lastSection := args.Section
		lastLevel := args.Level
		if lastLevel == 0 {
			lastSection = lastSection - 1
			lastLevel = int32(len(config.Huanlemaomibei.Section[lastSection].LevelArr) - 1)
		} else {
			lastLevel = lastLevel - 1
		}
		if _, ok := handler.levelDict[lastSection*100+lastLevel]; !ok {
			log.Println("4", lastSection, lastLevel, handler.levelDict)
			response.Code = errCodeHuanlemaomibeiLevelUnlock
			return r.Response(response)
		}
	}
	user.DecCoin(coin, tlogHuanlemaomibeiLevelUnlock)
	if _, ok := handler.levelDict[args.Section*100+args.Level]; !ok {
		handler.levelDict[args.Section*100+args.Level] = &model.HuanlemaomibeiLevel{
			Section: args.Section,
			Level:   args.Level,
			Unlock:  1,
			Star:    0,
			Coin:    0,
		}
		db.Huanlemaomibei_LevelSave(userid, args.Section, args.Level, 1, 0, 0)
	}
	response.Coin = coin
	response.Section = args.Section
	response.Level = args.Level
	return r.Response(response)
}

//关卡开始
func (handler *HuanlemaomibeiHandler) RoundStart(r *gat.Request, args *gamepb.HuanlemaomibeiStartRequest) error {
	log.Printf("[HuanlemaomibeiHandler] RoundStart %+v\n", args)
	response := &gamepb.HuanlemaomibeiStartResponse{}
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
func (handler *HuanlemaomibeiHandler) RoundResult(r *gat.Request, args *gamepb.HuanlemaomibeiResultRequest) error {
	log.Printf("[HuanlemaomibeiHandler] RoundResult %+v\n", args)
	response := &gamepb.HuanlemaomibeiResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	//gamedata := handler.gamedata
	userid := user.Userid
	roundResultTime := time.Now().Unix()
	//duration := roundResultTime - user.roundStartTime
	//	log.Println("ofafsfasf", duration)
	//保存关卡记录
	var addStar int32 = 0
	if _, ok := handler.levelDict[args.Section*100+args.Level]; !ok {
		if !(args.Section == 0 && args.Level == 0) {
			response.Code = errCodeHacker
			return r.Response(response)
		}
		handler.levelDict[args.Section*100+args.Level] = &model.HuanlemaomibeiLevel{
			Section: args.Section,
			Level:   args.Level,
			Unlock:  1,
			Star:    0,
			Coin:    0,
		}
	}
	if level, ok := handler.levelDict[args.Section*100+args.Level]; ok {
		if args.Star > level.Star {
			addStar = args.Star - level.Star
			level.Star = args.Star
		}
		db.Huanlemaomibei_LevelSave(userid, args.Section, args.Level, level.Unlock, level.Star, level.Coin)
	}
	var totalStar int32 = 0
	for _, level := range handler.levelDict {
		totalStar = totalStar + level.Star
	}
	if int64(totalStar) > user.Score {
		user.Score = int64(totalStar)
		db.ScoreRank_Update(userid, int64(totalStar))
	}
	if addStar > 0 {
		user.ScoreToday += int64(addStar)
		db.ScoreDayRank_Update(user.Today, userid, user.ScoreToday)
	}
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	user.Dirty = true
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, int64(totalStar), 1)
	return r.Response(response)
}

//关卡奖励
func (handler *HuanlemaomibeiHandler) LevelCoin(r *gat.Request, args *gamepb.HuanlemaomibeiLevelCoinRequest) error {
	log.Printf("[HuanlemaomibeiHandler] LevelCoin %+v\n", args)
	response := &gamepb.HuanlemaomibeiLevelCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	//gamedata := handler.gamedata
	userid := user.Userid
	level, ok := handler.levelDict[args.Section*100+args.Level]
	if !ok {
		log.Println("11", level)
		response.Code = errCodeHuanlemaomibeiLevelCoin
		return r.Response(response)
	} else if level.Star == 0 || level.Coin != 0 {
		log.Println("22", level)
		response.Code = errCodeHuanlemaomibeiLevelCoin
		return r.Response(response)
	}
	sectionConf := config.Huanlemaomibei.Section[args.Section]
	levelConf := sectionConf.LevelArr[args.Level]
	coin := levelConf.Coin
	user.AddCoin(coin, tlogHuanlemaomibeiLevelCoin)
	level.Coin = coin
	db.Huanlemaomibei_LevelSave(userid, args.Section, args.Level, level.Unlock, level.Star, level.Coin)
	return r.Response(response)
}

//增加体力
func (handler *HuanlemaomibeiHandler) AddHp(r *gat.Request, args *gamepb.HuanlemaomibeiAddHpRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddHp %+v\n", args)
	response := &gamepb.HuanlemaomibeiAddHpResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	//user := handler.agent.user
	gamedata := handler.gamedata
	gamedata.Hp = gamedata.Hp + args.Hp
	if gamedata.Hp > 3 {
		gamedata.Hp = 3
	}
	gamedata.Dirty = true
	response.Hp = args.Hp
	return r.Response(response)
}

//关卡开始
func (handler *HuanlemaomibeiHandler) LevelStart(r *gat.Request, args *gamepb.HuanlemaomibeiLevelStartRequest) error {
	log.Printf("[HuanlemaomibeiHandler] LevelStart %+v\n", args)
	response := &gamepb.HuanlemaomibeiLevelStartResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	user.RoundStartTime = time.Now().Unix()
	user.Dirty = true
	gamedata.Hp = gamedata.Hp - 1
	if gamedata.Hp < 0 {
		gamedata.Hp = 0
	}
	gamedata.Dirty = true
	tlog.RoundStart(user.Openid, userid)
	return r.Response(response)
}

//关卡结束
func (handler *HuanlemaomibeiHandler) LevelResult(r *gat.Request, args *gamepb.HuanlemaomibeiLevelResultRequest) error {
	log.Printf("[HuanlemaomibeiHandler] LevelResult %+v\n", args)
	response := &gamepb.HuanlemaomibeiLevelResultResponse{}
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
