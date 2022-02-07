package game

import (
	"fmt"
	"log"
	"time"

	"github.com/roada-go/gat"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/config"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type CaichengyuData struct {
	model.Caichengyu
	Dirty bool
}

type CaichengyuHandler struct {
	agent    *Agent
	gamedata *CaichengyuData
}

func newCaichengyuHandler(agent *Agent) *CaichengyuHandler {
	handler := &CaichengyuHandler{
		agent: agent,
	}
	return handler
}

func (handler *CaichengyuHandler) onLogout() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *CaichengyuHandler) GetData(r *gat.Request, args *gamepb.CaichengyuGetDataRequest) error {
	log.Printf("[CaichengyuHandler] GetData %+v\n", args)
	response := &gamepb.CaichengyuGetDataResponse{Code: 0}
	user := handler.agent.user
	userid := user.Userid
	var gamedata *CaichengyuData
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
	handler.resetHpData()
	response.GameId = conf.Ini.Game.Id
	response.RoleLevel = gamedata.RoleLevel
	response.BuildLevel = gamedata.BuildLevel
	response.Level = gamedata.Level
	response.LevelType = gamedata.LevelType
	response.LevelTip = gamedata.LevelTip
	response.Hp = gamedata.Hp
	response.HpDate = gamedata.HpDate
	response.GetHpCount = gamedata.GetHpCount
	response.GetHpDay = gamedata.GetHpDay
	log.Printf("[CaichengyuHandler] GetData succ, response=%+v\n", response)
	return r.Response(response)
}

func (handler *CaichengyuHandler) resetHpData() {
	gamedata := handler.gamedata
	now := time.Now().Unix()
	addHp := (now - gamedata.HpDate) / (60 * 20)
	if addHp <= 0 {
		return
	}
	newHp := gamedata.Hp + int32(addHp)
	if newHp > 5 {
		gamedata.Hp = 5
		gamedata.HpDate = now
	} else {
		gamedata.Hp = newHp
		gamedata.HpDate = gamedata.HpDate + addHp*60*20
	}
	gamedata.Dirty = true
}

func (handler *CaichengyuHandler) bgSave() {
	if handler.gamedata == nil {
		return
	}
	handler.saveData()
}

func (handler *CaichengyuHandler) saveData() error {
	gamedata := handler.gamedata
	if gamedata == nil {
		return nil
	}
	log.Printf("[CaichengyuHandler] saveData, dirty=%v\n", gamedata.Dirty)
	if !gamedata.Dirty {
		return nil
	}
	err := db.Caichengyu_Save(&gamedata.Caichengyu)
	if err != nil {
		log.Printf("[CaichengyuHandler] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[CaichengyuHandler] saveData success\n")
	gamedata.Dirty = false
	return nil
}

func (handler *CaichengyuHandler) loadData(userid int64) (*CaichengyuData, error) {
	dbrow, err := db.Caichengyu_Get(userid)
	if err != nil {
		return nil, err
	}
	var gamedata *CaichengyuData
	if dbrow == nil {
		now := time.Now()
		gamedata = &CaichengyuData{
			model.Caichengyu{
				Userid:     userid,
				RoleLevel:  0,
				BuildLevel: 0,
				Level:      1,
				LevelType:  0,
				Hp:         5,
				HpDate:     now.Unix(),
				GetHpCount: 0,
				GetHpDay:   fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day()),
			},
			false, //Dirty
		}
	} else {
		gamedata = &CaichengyuData{
			*dbrow,
			false, //Dirty
		}
	}
	return gamedata, nil
}

//恢复体力
func (handler *CaichengyuHandler) ResetHp(r *gat.Request, args *gamepb.CaichengyuResetHpRequest) error {
	log.Printf("[CaichengyuHandler] ResetHp %+v\n", args)
	response := &gamepb.CaichengyuResetHpResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	handler.resetHpData()
	response.Hp = gamedata.Hp
	response.HpDate = gamedata.HpDate
	return r.Response(response)
}

//升级角色
func (handler *CaichengyuHandler) UpgradeRole(r *gat.Request, args *gamepb.CaichengyuUpgradeRoleRequest) error {
	log.Printf("[CaichengyuHandler] UpgradeRole %+v\n", args)
	response := &gamepb.CaichengyuUpgradeRoleResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	lastRoleLevel := gamedata.RoleLevel
	if gamedata.RoleLevel >= int32(len(config.Caichengyu.Role)) {
		response.Code = errCodeCaichengyuUpgradeRole
		return r.Response(response)
	}
	for i := int(gamedata.RoleLevel + 1); i < len(config.Caichengyu.Role); i++ {
		nextLevel := int64(config.Caichengyu.Role[i])
		if gamedata.Level <= nextLevel {
			break
		} else {
			gamedata.RoleLevel = int32(i)
			gamedata.Dirty = true
		}
	}
	if lastRoleLevel != gamedata.RoleLevel {
		handler.agent.bgSave()
	}
	response.RoleLevel = gamedata.RoleLevel
	return r.Response(response)
}

//升级房子
func (handler *CaichengyuHandler) UpgradeBuild(r *gat.Request, args *gamepb.CaichengyuUpgradeBuildRequest) error {
	log.Printf("[CaichengyuHandler] UpgradeBuild %+v\n", args)
	response := &gamepb.CaichengyuUpgradeBuildResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	lastBuildLevel := gamedata.BuildLevel
	if gamedata.BuildLevel >= int32(len(config.Caichengyu.Build)) {
		response.Code = errCodeCaichengyuUpgradeBuild
		return r.Response(response)
	}
	for i := int(gamedata.BuildLevel + 1); i < len(config.Caichengyu.Build); i++ {
		nextLevel := int64(config.Caichengyu.Build[i])
		if gamedata.Level <= nextLevel {
			break
		} else {
			gamedata.BuildLevel = int32(i)
			gamedata.Dirty = true
		}
	}
	if lastBuildLevel != gamedata.BuildLevel {
		handler.agent.bgSave()
	}
	response.BuildLevel = gamedata.BuildLevel
	return r.Response(response)
}

//消耗金币获得体力
func (handler *CaichengyuHandler) GetHp(r *gat.Request, args *gamepb.CaichengyuGetHpRequest) error {
	log.Printf("[CaichengyuHandler] GetHp %+v\n", args)
	response := &gamepb.CaichengyuGetHpResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	var addHp int32 = 1
	if gamedata.Hp >= 5 {
		addHp = 0
	}
	gamedata.Hp += addHp
	gamedata.Dirty = true
	response.Hp = addHp
	return r.Response(response)
}

//消耗体力获得提示
func (handler *CaichengyuHandler) Tip(r *gat.Request, args *gamepb.CaichengyuTipRequest) error {
	log.Printf("[CaichengyuHandler] Tip %+v\n", args)
	response := &gamepb.CaichengyuTipResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if gamedata.LevelTip <= 0 && gamedata.Hp <= 0 {
		response.Code = errCodeCaichengyuTip
		return r.Response(response)
	}
	if gamedata.LevelTip > 0 {
		gamedata.LevelTip--
		response.LevelTip = 1
	}
	if gamedata.Hp > 0 {
		gamedata.Hp--
		response.Hp = 1
	}
	gamedata.Dirty = true
	return r.Response(response)
}

//体力购买提示次数（暂时没用到了）
func (handler *CaichengyuHandler) GetTip(r *gat.Request, args *gamepb.CaichengyuGetTipRequest) error {
	log.Printf("[CaichengyuHandler] GetTip %+v\n", args)
	response := &gamepb.CaichengyuGetTipResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	gamedata := handler.gamedata
	if gamedata.Hp <= 0 {
		response.Code = errCodeCaichengyuTip
		return r.Response(response)
	}
	gamedata.LevelTip++
	gamedata.Hp--
	gamedata.Dirty = true
	return r.Response(response)
}

//关卡开始
func (handler *CaichengyuHandler) RoundStart(r *gat.Request, args *gamepb.CaichengyuStartRequest) error {
	log.Printf("[CaichengyuHandler] RoundStart %+v\n", args)
	response := &gamepb.CaichengyuStartResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	rankuser, err := db.Caichengyu_StageIndex(gamedata.Level)
	if err != nil {
		response.Code = errCodeCaichengyuRoundStart
		return r.Response(response)
	}
	if rankuser != nil {
		response.User = &gamepb.RankUser{
			Openid:   rankuser.Openid,
			Nickname: rankuser.Nickname,
			Avatar:   rankuser.Avatar,
			Score:    rankuser.Level,
		}
	}
	var constHp int32 = 0
	if gamedata.LevelType == 0 {
		if gamedata.Hp <= 0 {
			response.Code = errCodeCaichengyuStart
			return r.Response(response)
		}
		constHp = 1
		gamedata.LevelType = 1
		gamedata.LevelTip = 2
		gamedata.Hp = gamedata.Hp - constHp
		gamedata.Dirty = true
	}
	user.RoundStartTime = time.Now().Unix()
	user.Dirty = true
	tlog.RoundStart(user.Openid, userid)
	response.Hp = constHp
	response.LevelTip = gamedata.LevelTip
	return r.Response(response)
}

//关卡结束
func (handler *CaichengyuHandler) RoundResult(r *gat.Request, args *gamepb.CaichengyuResultRequest) error {
	log.Printf("[CaichengyuHandler] RoundResult %+v\n", args)
	response := &gamepb.CaichengyuResultResponse{}
	if handler.gamedata == nil {
		response.Code = errCodeNotLogin
		return r.Response(response)
	}
	handler.agent.checkToday()
	user := handler.agent.user
	gamedata := handler.gamedata
	userid := user.Userid
	if gamedata.Level-1 >= int64(len(config.Caichengyu.Answer)) {
		response.Code = errCodeCaichengyuRoundResult
		return r.Response(response)
	}
	answer := config.Caichengyu.Answer[gamedata.Level-1]
	if answer != args.Answer {
		response.Code = errCodeCaichengyuRoundResult
		return r.Response(response)
	}
	roundResultTime := time.Now().Unix()
	//duration := roundResultTime - user.roundStartTime
	//	log.Println("ofafsfasf", duration)
	passLevel := gamedata.Level
	lastLevel := user.Level
	user.LastRoundStartTime = user.RoundStartTime
	user.LastRoundEndTime = roundResultTime
	user.RoundStartTime = 0
	user.RoundEndTime = 0
	if passLevel > lastLevel {
		user.Level = passLevel
	}
	user.LevelToday = user.LevelToday + 1
	user.Dirty = true
	gamedata.LevelType = 0
	gamedata.Level = passLevel + 1
	gamedata.Dirty = true
	handler.agent.bgSave()
	if passLevel > lastLevel {
		db.LevelRank_Update(userid, passLevel)
	}
	db.LevelDayRank_Update(user.Today, userid, user.LevelToday)
	tlog.RoundResult(user.Openid, userid, user.LastRoundStartTime, user.LastRoundEndTime, passLevel, 0)
	return r.Response(response)
}
