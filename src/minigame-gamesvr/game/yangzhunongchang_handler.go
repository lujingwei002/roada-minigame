package game

import (
	"log"
	"time"

	"github.com/roada-go/gat"
	"github.com/roada-go/util/timeutil"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type YangzhunongchangData struct {
	model.Yangzhunongchang
	Dirty bool
}

type YangzhunongchangHandler struct {
	agent        *Agent
	gamedata     *YangzhunongchangData
	itemDict     map[int32]*model.YangzhunongchangItem
	pigDict      map[string]*model.YangzhunongchangPig
	breedPigDict map[string]*model.YangzhunongchangBreedPig
	usuDict      map[int32]*model.YangzhunongchangUsu
	foodDict     map[string]*model.YangzhunongchangFood
	taskDict     map[string]*model.YangzhunongchangTask
}

func newYangzhunongchangHandler(agent *Agent) *YangzhunongchangHandler {
	handler := &YangzhunongchangHandler{
		agent: agent,
	}
	return handler
}

func (handler *YangzhunongchangHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *YangzhunongchangHandler) GetData(r *gat.Request, args *gamepb.YangzhunongchangGetDataRequest) error {
	log.Printf("[YangzhunongchangHandler] GetData %+v\n", args)
	response := &gamepb.YangzhunongchangGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *YangzhunongchangData
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
	response.ItemArr = make([]*gamepb.YangzhunongchangItem, 0)
	for _, item := range handler.itemDict {
		response.ItemArr = append(response.ItemArr, &gamepb.YangzhunongchangItem{
			Id:  item.Id,
			Num: item.Num,
		})
	}
	response.PigArr = make([]*gamepb.YangzhunongchangPig, 0)
	for _, pig := range handler.pigDict {
		response.PigArr = append(response.PigArr, &gamepb.YangzhunongchangPig{
			Id:         pig.Id,
			Data:       pig.Data,
			Createtime: pig.Createtime,
		})
	}
	response.BreedPigArr = make([]*gamepb.YangzhunongchangBreedPig, 0)
	for _, pig := range handler.breedPigDict {
		response.BreedPigArr = append(response.BreedPigArr, &gamepb.YangzhunongchangBreedPig{
			Id:   pig.Id,
			Data: pig.Data,
		})
	}
	response.UsuArr = make([]*gamepb.YangzhunongchangUsu, 0)
	for _, usu := range handler.usuDict {
		response.UsuArr = append(response.UsuArr, &gamepb.YangzhunongchangUsu{
			Id: usu.Id,
		})
	}
	response.FoodArr = make([]*gamepb.YangzhunongchangFood, 0)
	for _, food := range handler.foodDict {
		response.FoodArr = append(response.FoodArr, &gamepb.YangzhunongchangFood{
			Id:   food.Id,
			Data: food.Data,
		})
	}
	response.TaskArr = make([]*gamepb.YangzhunongchangTask, 0)
	for _, task := range handler.taskDict {
		response.TaskArr = append(response.TaskArr, &gamepb.YangzhunongchangTask{
			Id:    task.Id,
			Index: task.Index,
			Count: task.Count,
		})
	}
	response.FarmLv = gamedata.FarmLv
	response.FarmLvName = gamedata.FarmLvName
	response.FarmLvExp = gamedata.FarmLvExp
	response.FarmLvExpCur = gamedata.FarmLvExpCur
	response.AwardNum = gamedata.AwardNum
	response.GameId = conf.Ini.Game.Id
	log.Printf("[YangzhunongchangHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *YangzhunongchangHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[YangzhunongchangHandler] saveData, Dirty=%v, %+v\n", gamedata.Dirty, gamedata)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Yangzhunongchang_Save(&gamedata.Yangzhunongchang)
	if err != nil {
		log.Printf("[YangzhunongchangHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[YangzhunongchangHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *YangzhunongchangHandler) loadData(userid int64) (*YangzhunongchangData, error) {
	dbrow, err := db.Yangzhunongchang_Get(userid)
	if err != nil {
		return nil, err
	}
	newhand := false
	var gamedata *YangzhunongchangData
	if dbrow == nil {
		gamedata = &YangzhunongchangData{
			model.Yangzhunongchang{
				Userid:       userid,
				FarmLv:       0,
				FarmLvName:   "养猪小白",
				FarmLvExp:    0,
				FarmLvExpCur: 0,
				AwardNum:     0,
				AwardTime:    0,
			},
			true, //Dirty
		}
		newhand = true
	} else {
		gamedata = &YangzhunongchangData{
			*dbrow, //model
			false,  //Dirty
		}
	}
	itemArr, err := db.Yangzhunongchang_ItemGet(userid)
	if err != nil {
		return nil, err
	}
	handler.itemDict = make(map[int32]*model.YangzhunongchangItem)
	for _, item := range itemArr {
		handler.itemDict[item.Id] = item
	}
	if newhand {
		var i int32
		for i = 0; i <= 7; i++ {
			item := &model.YangzhunongchangItem{Id: i, Num: 20}
			handler.itemDict[i] = item
			db.Yangzhunongchang_ItemSave(userid, item.Id, item.Num)
		}
	}

	pigArr, err := db.Yangzhunongchang_PigGet(userid)
	if err != nil {
		return nil, err
	}
	handler.pigDict = make(map[string]*model.YangzhunongchangPig)
	for _, pig := range pigArr {
		handler.pigDict[pig.Id] = pig
	}

	breedPigArr, err := db.Yangzhunongchang_BreedPigGet(userid)
	if err != nil {
		return nil, err
	}
	handler.breedPigDict = make(map[string]*model.YangzhunongchangBreedPig)
	for _, pig := range breedPigArr {
		handler.breedPigDict[pig.Id] = pig
	}

	usuArr, err := db.Yangzhunongchang_UsuGet(userid)
	if err != nil {
		return nil, err
	}
	handler.usuDict = make(map[int32]*model.YangzhunongchangUsu)
	for _, usu := range usuArr {
		handler.usuDict[usu.Id] = usu
	}

	foodArr, err := db.Yangzhunongchang_FoodGet(userid)
	if err != nil {
		return nil, err
	}
	handler.foodDict = make(map[string]*model.YangzhunongchangFood)
	for _, food := range foodArr {
		handler.foodDict[food.Id] = food
	}

	taskArr, err := db.Yangzhunongchang_TaskGet(userid)
	if err != nil {
		return nil, err
	}
	handler.taskDict = make(map[string]*model.YangzhunongchangTask)
	for _, task := range taskArr {
		handler.taskDict[task.Id] = task
	}

	if timeutil.PassDays(gamedata.AwardTime, time.Now().Unix()) >= 1 {
		gamedata.AwardNum = 0
		gamedata.AwardTime = time.Now().Unix()
		gamedata.Dirty = true
	}
	return gamedata, nil
}

//金币
func (handler *YangzhunongchangHandler) AddCoin(r *gat.Request, args *gamepb.YangzhunongchangAddCoinRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddCoin %+v\n", args)
	response := &gamepb.YangzhunongchangAddCoinResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}

	user := handler.agent.user
	coin := args.Coin
	if coin > 0 {
		user.AddCoin(coin, tlogYangzhunongchangAddCoin)
	} else if coin < 0 {
		user.DecCoin(-coin, tlogYangzhunongchangAddCoin)
	}
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) AddDiamond(r *gat.Request, args *gamepb.YangzhunongchangAddDiamondRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddDiamond %+v\n", args)
	response := &gamepb.YangzhunongchangAddDiamondResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}

	user := handler.agent.user
	diamond := args.Diamond
	if diamond > 0 {
		user.AddDiamond(diamond, tlogYangzhunongchangAddDiamond)
	} else if diamond < 0 {
		user.DecDiamond(-diamond, tlogYangzhunongchangAddDiamond)
	}
	return r.Response(response)
}

//增加经验
func (handler *YangzhunongchangHandler) AddExp(r *gat.Request, args *gamepb.YangzhunongchangAddExpRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddExp %+v\n", args)
	response := &gamepb.YangzhunongchangAddExpResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	//user := handler.agent.user
	exp := args.Exp
	gamedata.FarmLvExp += exp
	gamedata.FarmLvExpCur += exp
	gamedata.Dirty = true
	return r.Response(response)
}

//增加等级
func (handler *YangzhunongchangHandler) AddLevel(r *gat.Request, args *gamepb.YangzhunongchangAddLevelRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddLevel %+v\n", args)
	response := &gamepb.YangzhunongchangAddLevelResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	//user := handler.agent.user
	gamedata.FarmLv = args.FarmLv
	gamedata.FarmLvExpCur = args.FarmLvExpCur
	gamedata.FarmLvName = args.FarmLvName
	gamedata.Dirty = true
	return r.Response(response)
}

//道具
func (handler *YangzhunongchangHandler) AddItem(r *gat.Request, args *gamepb.YangzhunongchangAddItemRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddItem %+v\n", args)
	response := &gamepb.YangzhunongchangAddItemResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if item, ok := handler.itemDict[args.Id]; ok {
		item.Num = item.Num + args.Num
		if item.Num < 0 {
			item.Num = 0
		}
		db.Yangzhunongchang_ItemSave(userid, item.Id, item.Num)
	} else {
		if args.Num > 0 {
			item := &model.YangzhunongchangItem{Id: args.Id, Num: args.Num}
			handler.itemDict[args.Id] = item
			db.Yangzhunongchang_ItemSave(userid, item.Id, item.Num)
		}
	}
	return r.Response(response)
}

//猪栏
func (handler *YangzhunongchangHandler) AddPig(r *gat.Request, args *gamepb.YangzhunongchangAddPigRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddPig %+v\n", args)
	response := &gamepb.YangzhunongchangAddPigResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if pig, ok := handler.pigDict[args.Id]; ok {
		pig.Data = args.Data
		db.Yangzhunongchang_PigSave(userid, pig.Id, pig.Data)
	} else {
		pig := &model.YangzhunongchangPig{Id: args.Id, Data: args.Data}
		handler.pigDict[args.Id] = pig
		db.Yangzhunongchang_PigSave(userid, pig.Id, pig.Data)
	}
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) DelPig(r *gat.Request, args *gamepb.YangzhunongchangDelPigRequest) error {
	log.Printf("[HuanlemaomibeiHandler] DelPig %+v\n", args)
	response := &gamepb.YangzhunongchangDelPigResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if _, ok := handler.pigDict[args.Id]; !ok {
		return r.Response(response)
	}
	delete(handler.pigDict, args.Id)
	db.Yangzhunongchang_PigDel(userid, args.Id)
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) AddBreedPig(r *gat.Request, args *gamepb.YangzhunongchangAddBreedPigRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddBreedPig %+v\n", args)
	response := &gamepb.YangzhunongchangAddBreedPigResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if pig, ok := handler.breedPigDict[args.Id]; ok {
		pig.Data = args.Data
		db.Yangzhunongchang_BreedPigSave(userid, pig.Id, pig.Data)
	} else {
		pig := &model.YangzhunongchangBreedPig{Id: args.Id, Data: args.Data}
		handler.breedPigDict[args.Id] = pig
		db.Yangzhunongchang_BreedPigSave(userid, pig.Id, pig.Data)
	}
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) DelBreedPig(r *gat.Request, args *gamepb.YangzhunongchangDelBreedPigRequest) error {
	log.Printf("[HuanlemaomibeiHandler] DelBreedPig %+v\n", args)
	response := &gamepb.YangzhunongchangDelBreedPigResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if _, ok := handler.breedPigDict[args.Id]; !ok {
		return r.Response(response)
	}
	delete(handler.breedPigDict, args.Id)
	db.Yangzhunongchang_BreedPigDel(userid, args.Id)
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) AddUsu(r *gat.Request, args *gamepb.YangzhunongchangAddUsuRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddUsu %+v\n", args)
	response := &gamepb.YangzhunongchangAddUsuResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if usu, ok := handler.usuDict[args.Id]; ok {
		db.Yangzhunongchang_UsuSave(userid, usu.Id)
	} else {
		usu := &model.YangzhunongchangUsu{Id: args.Id}
		handler.usuDict[args.Id] = usu
		db.Yangzhunongchang_UsuSave(userid, usu.Id)
	}
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) AddFood(r *gat.Request, args *gamepb.YangzhunongchangAddFoodRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddFood %+v\n", args)
	response := &gamepb.YangzhunongchangAddFoodResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if food, ok := handler.foodDict[args.Id]; ok {
		food.Data = args.Data
		db.Yangzhunongchang_FoodSave(userid, food.Id, food.Data)
	} else {
		food := &model.YangzhunongchangFood{Id: args.Id, Data: args.Data}
		handler.foodDict[args.Id] = food
		db.Yangzhunongchang_FoodSave(userid, food.Id, food.Data)
	}
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) DelFood(r *gat.Request, args *gamepb.YangzhunongchangDelFoodRequest) error {
	log.Printf("[HuanlemaomibeiHandler] DelFood %+v\n", args)
	response := &gamepb.YangzhunongchangDelFoodResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if _, ok := handler.foodDict[args.Id]; !ok {
		return r.Response(response)
	}
	delete(handler.foodDict, args.Id)
	db.Yangzhunongchang_FoodDel(userid, args.Id)
	return r.Response(response)
}

func (handler *YangzhunongchangHandler) AddTask(r *gat.Request, args *gamepb.YangzhunongchangAddTaskRequest) error {
	log.Printf("[HuanlemaomibeiHandler] AddTask %+v\n", args)
	response := &gamepb.YangzhunongchangAddPigResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	userid := user.Userid
	if task, ok := handler.taskDict[args.Id]; ok {
		task.Index = args.Index
		task.Count = args.Count
		db.Yangzhunongchang_TaskSave(userid, task.Id, task.Index, task.Count)
	} else {
		task := &model.YangzhunongchangTask{Id: args.Id, Index: args.Index, Count: args.Count}
		handler.taskDict[args.Id] = task
		db.Yangzhunongchang_TaskSave(userid, task.Id, task.Index, task.Count)
	}
	return r.Response(response)
}

//抽奖
func (handler *YangzhunongchangHandler) Award(r *gat.Request, args *gamepb.YangzhunongchangAwardRequest) error {
	log.Printf("[HuanlemaomibeiHandler] Award %+v\n", args)
	response := &gamepb.YangzhunongchangAwardResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	//user := handler.agent.user
	gamedata.AwardNum++
	if gamedata.AwardNum < 0 {
		gamedata.AwardNum = 0
	}
	gamedata.Dirty = true
	return r.Response(response)
}

//关卡开始
func (handler *YangzhunongchangHandler) RoundStart(r *gat.Request, args *gamepb.YangzhunongchangStartRequest) error {
	log.Printf("[YangzhunongchangHandler] RoundStart %+v\n", args)
	response := &gamepb.YangzhunongchangStartResponse{}
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
func (handler *YangzhunongchangHandler) RoundResult(r *gat.Request, args *gamepb.YangzhunongchangResultRequest) error {
	log.Printf("[YangzhunongchangHandler] RoundResult %+v\n", args)
	response := &gamepb.YangzhunongchangResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	//gamedata := handler.gamedata
	//userid := user.Userid
	roundResultTime := time.Now().Unix()
	//duration := roundResultTime - user.roundStartTime
	//	log.Println("ofafsfasf", duration)
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	user.Dirty = true
	//	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, passLevel, 0)
	return r.Response(response)
}
