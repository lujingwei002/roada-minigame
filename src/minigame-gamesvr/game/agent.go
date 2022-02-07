package game

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/roada-go/gat"
	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/gamepb"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type Agent struct {
	game        *GameService
	road        *roada.Road
	rkey        string
	session     *gat.Session
	chRoad      chan *roada.Request
	chGate      chan *gat.Request
	chQuit      chan bool
	user        *User
	userid      int64
	openid      string
	handlerDict map[string]HandlerInterface
	isClose     bool
}

type HandlerInterface interface {
	onLogout()
	bgSave()
}

func newAgent(game *GameService, session *gat.Session) *Agent {
	agent := &Agent{
		game:        game,
		road:        game.road,
		session:     session,
		handlerDict: make(map[string]HandlerInterface),
		chRoad:      make(chan *roada.Request, 1),
		chGate:      make(chan *gat.Request, 1),
		chQuit:      make(chan bool),
	}
	if conf.Ini.Game.Name == "daxigua" || conf.Ini.Game.Name == "minigame" {
		handler := newDaxiguaHandler(agent)
		agent.handlerDict["daxigua"] = handler
	}
	if conf.Ini.Game.Name == "bpbxq" || conf.Ini.Game.Name == "minigame" {
		handler := newBpbxqHandler(agent)
		agent.handlerDict["bpbxq"] = handler
	}
	if conf.Ini.Game.Name == "caichengyu" || conf.Ini.Game.Name == "minigame" {
		handler := newCaichengyuHandler(agent)
		agent.handlerDict["caichengyu"] = handler
	}
	if conf.Ini.Game.Name == "gongjianchuanshu" || conf.Ini.Game.Name == "minigame" {
		handler := newGongjianchuanshuHandler(agent)
		agent.handlerDict["gongjianchuanshu"] = handler
	}
	if conf.Ini.Game.Name == "fangkuainiao" || conf.Ini.Game.Name == "minigame" {
		handler := newFangkuainiaoHandler(agent)
		agent.handlerDict["fangkuainiao"] = handler
	}
	if conf.Ini.Game.Name == "paopaolong" || conf.Ini.Game.Name == "minigame" {
		handler := newPaopaolongHandler(agent)
		agent.handlerDict["paopaolong"] = handler
	}
	if conf.Ini.Game.Name == "tanchishedazuozhan" || conf.Ini.Game.Name == "minigame" {
		handler := newTanchishedazuozhanHandler(agent)
		agent.handlerDict["tanchishedazuozhan"] = handler
	}
	if conf.Ini.Game.Name == "tiantianpaoku" || conf.Ini.Game.Name == "minigame" {
		handler := newTiantianpaokuHandler(agent)
		agent.handlerDict["tiantianpaoku"] = handler
	}
	if conf.Ini.Game.Name == "huanlemaomibei" || conf.Ini.Game.Name == "minigame" {
		handler := newHuanlemaomibeiHandler(agent)
		agent.handlerDict["huanlemaomibei"] = handler
	}
	if conf.Ini.Game.Name == "yangzhunongchang" || conf.Ini.Game.Name == "minigame" {
		handler := newYangzhunongchangHandler(agent)
		agent.handlerDict["yangzhunongchang"] = handler
	}
	go agent.forever()
	return agent
}

func (agent *Agent) onSessionOpen(session *gat.Session) {
}

func (agent *Agent) onSessionClose(session *gat.Session) {
	close(agent.chQuit)
}

func (agent *Agent) request(r *gat.Request) {
	if r == nil {
		return
	}
	index := strings.LastIndex(r.Route, ".")
	if index < 0 {
		log.Println(fmt.Sprintf("[agent] invalid route, route=%s", r.Route))
		return
	}
	handlerName := r.Route[:index]
	methodName := r.Route[index+1:]
	if agent.user == nil && methodName != "login" {
		log.Println(fmt.Sprintf("[agent] please login first! route=%s, sessionid=%d", r.Route, agent.session.ID()))
		return
	}
	var handler interface{}
	var service *gat.Service
	var ok bool
	service, ok = agent.game.serviceDict[handlerName]
	if !ok {
		log.Printf("[agent] service not found, route=%s\n", r.Route)
		return
	}
	if handlerName == "game" {
		handler = agent
	} else {
		handler, ok = agent.handlerDict[handlerName]
		if !ok {
			log.Printf("[agent] handler not found, route=%s\n", r.Route)
			return
		}
	}
	service.Unpack(handler, r)
}

func (agent *Agent) rpc(r *roada.Request) {
	if r == nil {
		return
	}
	agent.game.rpcService.ServeRPC(agent, r)
}

func (agent *Agent) forever() {
	session := agent.session
	tick := time.NewTicker(300 * time.Second)
	log.Printf("[agent] loop start, sessionid:%d\n", session.ID())
	defer func() {
		log.Printf("[agent] loop end, sessionid:%d\n", session.ID())
		tick.Stop()
		agent.close()
	}()
	for !agent.isClose {
		select {
		case r := <-agent.chGate:
			{
				agent.request(r)
			}
		case r := <-agent.chRoad:
			{
				agent.rpc(r)
			}
		case <-agent.chQuit:
			{
				return
			}
		case <-tick.C:
			{
				agent.bgSave()
			}
		}
	}
}

func (agent *Agent) bgSave() {
	agent.saveData()
	for _, handler := range agent.handlerDict {
		handler.bgSave()
	}
}

func (agent *Agent) kick(reason string) {
	response := &gamepb.KickPush{Reason: reason}
	agent.session.SyncPush("game.kick", response)
	agent.close()
	agent.session.Close()
}

func (agent *Agent) close() {
	if agent.isClose {
		return
	}
	agent.isClose = true
	session := agent.session
	user := agent.user
	userid := agent.userid
	log.Printf("[agent] close, sessionid:%d, userid:%d\n", session.ID(), userid)
	if user != nil {
		//释放user
		atomic.AddInt32(&agent.game.onlineNum, -1)
		user.Logouttime = time.Now().Unix()
		user.Dirty = true
		agent.saveData()
		for _, handler := range agent.handlerDict {
			handler.onLogout()
		}
		tlog.UserLogout(agent.openid, userid, user.Logintime, time.Now().Unix())
		log.Printf("[agent] Logout, sessionid:%d, userid:%d\n", session.ID(), userid)
	}
	//释放agent
	close(agent.chGate)
	close(agent.chRoad)
	if len(agent.rkey) > 0 {
		agent.unregisterPos()
	}
}

func (agent *Agent) ServeMessage(r *gat.Request) error {
	if agent.isClose {
		return nil
	}
	agent.chGate <- r
	return nil
}

func (agent *Agent) ServeRPC(r *roada.Request) {
	log.Printf("[agent] recv ServeRPC, sessionid=%d\n", agent.session.ID())
	agent.chRoad <- r
	r.Wait(5)
}

func (agent *Agent) tokenValidate(openid string, nickname string, avatar string, expected_token string, time int32) bool {
	if !conf.Ini.Game.UrlValidate {
		return true
	}
	if _, err := strconv.ParseInt(openid, 10, 64); err != nil {
		log.Printf("[game] tokenValidate failed, openid=%s", openid)
		return false
	}
	/*s := fmt.Sprintf("%s%s%d%s%s", avatar, nickname, time, openid, conf.Ini.Game.UrlValidateSecret)
	token := fmt.Sprintf("%x", md5.Sum([]byte(s)))
	if expected_token != token {
		log.Printf("[game] tokenValidate failed, token=%s, expected_token=%s", token, expected_token)
		return false
	}*/
	return true
}

func (agent *Agent) registerPos(openid string) error {
	rkey := fmt.Sprintf("agent%s", openid)
	tryTimes := 0
	for {
		if err := agent.road.LocalSet(rkey); err == nil {
			break
		}
		if tryTimes >= 3 {
			return fmt.Errorf("[agent] registerPos failed, rkey=%s", rkey)
		}
		tryTimes = tryTimes + 1
		if err := UserInstead(rkey); err != nil {
			return err
		}
	}
	agent.road.Handle(rkey, agent)
	agent.rkey = rkey
	return nil
}

func (agent *Agent) unregisterPos() error {
	if err := agent.road.LocalDel(agent.rkey); err != nil {
		log.Printf("[agent] road.LocalDel failed, rkey=%s, error=%s\n", agent.rkey, err.Error())
		return err
	}
	if err := agent.road.DelHandle(agent.rkey, agent); err != nil {
		log.Printf("[agent] road.DelHandle failed, rkey=%s, error=%s\n", agent.rkey, err.Error())
		return err
	}
	return nil
}

func (agent *Agent) Login(r *gat.Request, args *gamepb.LoginRequest) error {
	session := agent.session
	log.Printf("[agent] Login %+v\n", args)
	response := &gamepb.LoginResponse{Code: 0}
	openid := args.Openid
	token := args.Token
	platformRegisterTime := args.PlatformRegisterTime
	addr := session.RemoteAddr().String()
	pats := strings.Split(addr, ":")
	var ip string
	if len(pats) >= 2 {
		ip = pats[0]
	}
	if !agent.tokenValidate(args.Openid, args.Nickname, args.Avatar, token, args.Time) {
		response.Code = errCodeDbErr
		return r.Response(response)
	}
	if err := agent.registerPos(openid); err != nil {
		log.Printf("[agent] register pos failed, openid=%s, err=%s\n", openid, err.Error())
		session.Close()
		response.Code = errCodeLogin
		return r.Response(response)
	} else {
		log.Printf("[agent] register pos success, openid=%s\n", openid)
	}
	user, err := agent.loadData(openid, args.Nickname, args.Avatar, ip)
	if err != nil {
		response.Code = errCodeDbErr
		return r.Response(response)
	}
	log.Printf("[agent] loadData success, user=%+v\n", user)
	user.Openid = openid
	user.Token = token
	userid := user.Userid
	session.Bind(userid)
	agent.userid = userid
	agent.openid = openid
	agent.user = user
	atomic.AddInt32(&agent.game.onlineNum, 1)
	if user.NewHand {
		tlog.UserRegister(userid, openid, user.Nickname, ip, platformRegisterTime)
	}
	lastLoginTime := time.Unix(user.Logintime, 0)
	user.Today = int64(lastLoginTime.Year()*10000 + int(lastLoginTime.Month())*100 + lastLoginTime.Day())
	agent.checkToday()
	now := time.Now().Unix()
	user.Logintime = now
	user.Dirty = true
	tlog.UserLogin(openid, userid, now)

	//更新账号数据
	openInfoDirty := false
	if args.Nickname != user.Nickname {
		user.Nickname = args.Nickname
		openInfoDirty = true
	}
	if args.Avatar != user.Avatar {
		user.Avatar = args.Avatar
		openInfoDirty = true
	}
	if openInfoDirty {
		db.User_UpdateOpenInfo(userid, user.Nickname, user.Avatar)
	}

	response.Userid = userid
	response.Score = user.Score
	response.Medal = user.Medal
	response.Coin = user.Coin
	response.Diamond = user.Diamond
	response.Level = user.Level
	response.Gold = user.Gold
	response.ServerTime = now
	log.Printf("[agent] Login succ, sessionid=%d, response=%+v\n", session.ID(), response)
	return r.Response(response)
}

func (agent *Agent) RoundShare(r *gat.Request, args *gamepb.RoundShareRequest) error {
	session := agent.session
	log.Printf("[agent] RoundShare sessionid=%d\n", session.ID())
	//user := agent.user
	//userid := user.Userid
	//tlog.RoundShare(agent.openid, userid)
	response := &gamepb.RoundShareResponse{Code: 0}
	return r.Response(response)
}

func (agent *Agent) checkToday() {
	user := agent.user
	now := time.Now()
	today := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	log.Printf("[agent] checkToday, now=%d, today=%d\n", today, user.Today)
	if today != user.Today {
		user.Today = today
		agent.resetToday()
	}
}

func (agent *Agent) resetToday() {
	user := agent.user
	user.ScoreToday = 0
	user.MedalToday = 0
	user.LevelToday = 0
}

func (agent *Agent) loadData(openid string, nickname string, avatar string, clientIp string) (*User, error) {
	dbrow, err := db.User_Get(openid)
	if err != nil {
		return nil, err
	}
	newHand := false
	if dbrow == nil {
		_, err := db.User_Insert(openid, nickname, avatar, clientIp)
		if err != nil {
			return nil, err
		}
		dbrow, err = db.User_Get(openid)
		if err != nil {
			return nil, err
		}
		if dbrow == nil {
			return nil, fmt.Errorf("[agent] loadData failed, user not found, openid=%s", openid)
		}
		newHand = true
	}
	user := &User{
		Dirty:   false,
		NewHand: newHand,
	}
	user.User = *dbrow

	now := time.Now()
	rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	score, err := db.ScoreDayRank_Score(user.Userid, rankid)
	if err != nil {
		return nil, err
	}
	user.ScoreToday = score
	score, err = db.MedalDayRank_Score(user.Userid, rankid)
	if err != nil {
		return nil, err
	}
	user.MedalToday = score
	score, err = db.LevelDayRank_Score(user.Userid, rankid)
	if err != nil {
		return nil, err
	}
	user.LevelToday = score
	return user, nil
}

func (agent *Agent) saveData() error {
	user := agent.user
	if user == nil {
		return nil
	}
	log.Printf("[agent] saveData, dirty=%v\n", user.Dirty)
	if !user.Dirty {
		return nil
	}
	err := db.User_Save(&user.User)
	if err != nil {
		log.Printf("[agent] saveData falied, error=%s\n", err.Error())
		return err
	}
	log.Printf("[agent] saveData success, user=%+v\n", user)
	user.Dirty = false
	return nil
}

func (agent *Agent) GetScoreRank(r *gat.Request, args *gamepb.GetScoreRankRequest) error {
	log.Printf("[agent] GetScoreRank args=%+v\n", args)
	response := &gamepb.GetScoreRankResponse{}
	user := agent.user
	//userid := user.Userid
	var users []*model.User
	var err error
	var myScore int64 = -1
	var myRank int32 = -1
	agent.checkToday()
	rankid := user.Today
	if args.Type == 1 {
		users, err = agent.game.getScoreDayRank(rankid)
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	} else {
		users, err = agent.game.getScoreRank()
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	}
	for index, rankuser := range users {
		if rankuser.Openid == user.Openid {
			myRank = int32(index + 1)
			myScore = rankuser.Score
		}
	}
	if myScore == -1 {
		if args.Type == 1 {
			/*myScore, err = db.ScoreDayRank_Score(userid, rankid)
			if err != nil {
				response.Code = errCodeDbErr
				return r.Response(response)
			}*/
			myScore = user.ScoreToday
		} else {
			myScore = user.Score
		}
	}
	for _, rankuser := range users {
		ruser := &gamepb.RankUser{
			Openid:   rankuser.Openid,
			Nickname: rankuser.Nickname,
			Avatar:   rankuser.Avatar,
			Score:    rankuser.Score,
		}
		if ruser.Avatar == "" {
			ruser.Avatar = "https://pic1.zhimg.com/v2-0e12707217a7c6d403259d9667f5d864_r.jpg"
		}
		if ruser.Nickname == "" {
			ruser.Nickname = "..."
		}
		response.Users = append(response.Users, ruser)
	}
	response.Rank = myRank
	response.Score = myScore
	return r.Response(response)
}

func (agent *Agent) GetMedalRank(r *gat.Request, args *gamepb.GetMedalRankRequest) error {
	log.Printf("[agent] GetMedalRank args=%+v\n", args)
	response := &gamepb.GetMedalRankResponse{}
	user := agent.user
	//userid := user.Userid
	var users []*model.User
	var err error
	var myScore int64 = -1
	var myRank int32 = -1
	agent.checkToday()
	rankid := user.Today
	if args.Type == 1 {
		users, err = agent.game.getMedalDayRank(rankid)
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	} else {
		users, err = agent.game.getMedalRank()
		if err != nil {
			log.Printf("[agent] game.getMedalRank failed, error=%s\n", err.Error())
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	}
	for index, rankuser := range users {
		if rankuser.Openid == user.Openid {
			myRank = int32(index + 1)
			myScore = rankuser.Score
		}
	}
	if myScore == -1 {
		if args.Type == 1 {
			/*myScore, err = db.MedalDayRank_Score(userid, rankid)
			if err != nil {
				response.Code = errCodeDbErr
				return r.Response(response)
			}*/
			myScore = user.MedalToday
		} else {
			myScore = user.Medal
		}
	}
	for _, rankuser := range users {
		ruser := &gamepb.RankUser{
			Openid:   rankuser.Openid,
			Nickname: rankuser.Nickname,
			Avatar:   rankuser.Avatar,
			Score:    rankuser.Score,
		}
		if ruser.Avatar == "" {
			ruser.Avatar = "https://pic1.zhimg.com/v2-0e12707217a7c6d403259d9667f5d864_r.jpg"
		}
		if ruser.Nickname == "" {
			ruser.Nickname = "..."
		}
		response.Users = append(response.Users, ruser)
	}
	response.Rank = myRank
	response.Score = myScore
	return r.Response(response)
}

func (agent *Agent) GetLevelRank(r *gat.Request, args *gamepb.GetLevelRankRequest) error {
	log.Printf("[agent] GetLevelRank args=%+v\n", args)
	response := &gamepb.GetLevelRankResponse{}
	user := agent.user
	//userid := user.Userid
	var users []*model.User
	var err error
	var myScore int64 = -1
	var myRank int32 = -1
	agent.checkToday()
	rankid := user.Today
	if args.Type == 1 {
		users, err = agent.game.getLevelDayRank(rankid)
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	} else {
		users, err = agent.game.getLevelRank()
		if err != nil {
			response.Code = errCodeDbErr
			return r.Response(response)
		}
	}
	for index, rankuser := range users {
		if rankuser.Openid == user.Openid {
			myRank = int32(index + 1)
			myScore = rankuser.Score
		}
	}
	if myScore == -1 {
		if args.Type == 1 {
			/*myScore, err = db.LevelDayRank_Score(userid, rankid)
			if err != nil {
				response.Code = errCodeDbErr
				return r.Response(response)
			}*/
			myScore = user.LevelToday
		} else {
			myScore = user.Level
		}
	}
	for _, rankuser := range users {
		ruser := &gamepb.RankUser{
			Openid:   rankuser.Openid,
			Nickname: rankuser.Nickname,
			Avatar:   rankuser.Avatar,
			Score:    rankuser.Score,
		}
		if ruser.Avatar == "" {
			ruser.Avatar = "https://pic1.zhimg.com/v2-0e12707217a7c6d403259d9667f5d864_r.jpg"
		}
		if ruser.Nickname == "" {
			ruser.Nickname = "..."
		}
		response.Users = append(response.Users, ruser)
	}
	response.Rank = myRank
	response.Score = myScore
	return r.Response(response)
}
