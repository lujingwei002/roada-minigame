package game

import (
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/roada-go/gat"
	"github.com/roada-go/gat/middleware"
	"github.com/roada-go/roada"
	sched "github.com/roada-go/util/scheduler"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/config"
	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

type RankCache struct {
	users []*model.User
}

type GameService struct {
	road                *roada.Road
	gate                *gat.Gate
	onlineNum           int32
	agentDict           sync.Map
	serviceDict         map[string]*gat.Service
	rpcService          *roada.Service
	rankCache           sync.Map
	rankCacheExpireTime sync.Map
	gateStat            *middleware.GateStat
	scheduler           *sched.Scheduler
}

var game *GameService

func Shutdown() {
	log.Println("[game] shutdown")
	game.Shutdown()
}

func Register(road *roada.Road, gate *gat.Gate) error {
	game = &GameService{
		gate:                gate,
		road:                road,
		agentDict:           sync.Map{},
		serviceDict:         make(map[string]*gat.Service),
		rankCache:           sync.Map{},
		rankCacheExpireTime: sync.Map{},
		gateStat:            middleware.NewGateStat(),
		scheduler:           sched.NewScheduler(),
	}
	gate.UseMiddleware(game.gateStat)
	//注册rpc消息回调
	if service, err := road.Register(&Agent{}); err != nil {
		return err
	} else {
		game.rpcService = service
	}
	//注册网关消息回调
	gate.Handle(game)
	if s, err := gate.Register(&Agent{}); err != nil {
		return err
	} else {
		game.serviceDict["game"] = s
	}
	if conf.Ini.Game.Name == "daxigua" || conf.Ini.Game.Name == "minigame" {
		if s, err := gate.Register(&DaxiguaHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["daxigua"] = s
		}
	}
	if conf.Ini.Game.Name == "bpbxq" || conf.Ini.Game.Name == "minigame" {
		if s, err := gate.Register(&BpbxqHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["bpbxq"] = s
		}
	}
	if conf.Ini.Game.Name == "caichengyu" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadCaichengyu(); err != nil {
			return err
		}
		if s, err := gate.Register(&CaichengyuHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["caichengyu"] = s
		}
	}
	if conf.Ini.Game.Name == "fangkuainiao" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadFangkuainiao(); err != nil {
			return err
		}
		if s, err := gate.Register(&FangkuainiaoHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["fangkuainiao"] = s
		}
	}
	if conf.Ini.Game.Name == "gongjianchuanshu" || conf.Ini.Game.Name == "minigame" {
		if s, err := gate.Register(&GongjianchuanshuHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["gongjianchuanshu"] = s
		}
	}
	if conf.Ini.Game.Name == "paopaolong" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadPaopaolong(); err != nil {
			return err
		}
		if s, err := gate.Register(&PaopaolongHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["paopaolong"] = s
		}
	}
	if conf.Ini.Game.Name == "tanchishedazuozhan" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadTanchishedazuozhan(); err != nil {
			return err
		}
		if s, err := gate.Register(&TanchishedazuozhanHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["tanchishedazuozhan"] = s
		}
	}
	if conf.Ini.Game.Name == "tiantianpaoku" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadTiantianpaoku(); err != nil {
			return err
		}
		if s, err := gate.Register(&TiantianpaokuHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["tiantianpaoku"] = s
		}
	}
	if conf.Ini.Game.Name == "huanlemaomibei" || conf.Ini.Game.Name == "minigame" {
		if err := config.LoadHuanlemaomibei(); err != nil {
			return err
		}
		if s, err := gate.Register(&HuanlemaomibeiHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["huanlemaomibei"] = s
		}
	}
	if conf.Ini.Game.Name == "yangzhunongchang" || conf.Ini.Game.Name == "minigame" {
		if s, err := gate.Register(&YangzhunongchangHandler{}); err != nil {
			return err
		} else {
			game.serviceDict["yangzhunongchang"] = s
		}
	}
	go game.forever()
	go game.gc()
	for i := 0; i < 2; i++ {
		user, err := db.User_Get("ljw")
		if err != nil {
			log.Println("User_Get failed", err.Error())
		}
		log.Println("[game] User_Get succ", user)
	}
	return nil
}

func (game *GameService) gc() {
	for {
		log.Println("[game] FreeOSMemory")
		debug.FreeOSMemory()
		time.Sleep(300 * time.Second)
	}
}

func (game *GameService) forever() {
	game.scheduler.NewAfterTimer(1*time.Second, game.stat)
	game.scheduler.NewTimer(60*time.Second, game.stat)
	tick := time.NewTicker(1 * time.Second)
	defer func() {
		tick.Stop()
		game.scheduler.Close()
	}()
	for {
		select {
		case <-tick.C:
			{
				game.scheduler.Sched()
			}
		case task := <-game.scheduler.T:
			{
				game.scheduler.Invoke(task)
			}
		}
	}
}

func (game *GameService) Shutdown() {
	game.agentDict.Range(func(key, value interface{}) bool {
		agent := value.(*Agent)
		agent.kick("先休息一会，很快回来")
		return true
	})
}

func (game *GameService) ServeMessage(r *gat.Request) {
	session := r.Session
	if !session.HasKey("agent") {
		return
	}
	agent := session.Value("agent").(*Agent)
	agent.ServeMessage(r)
}

func (game *GameService) stat() {
	game.gateStat.Record()
	log.Println(game.gateStat)
	//运行在主线程里
	tlog.GateStat(game.onlineNum)
	db.GateStat_Insert(conf.Ini.Game.Id, game.road.NodeName, game.road.NodeFullName,
		game.gate.Host, game.gate.Port, game.onlineNum)
}

func (game *GameService) OnSessionOpen(session *gat.Session) {
	log.Printf("[GameService] session open, sessionid=%d\n", session.ID())
	agent := newAgent(game, session)
	game.agentDict.Store(session.ID(), agent)
	session.Set("agent", agent)
	agent.onSessionOpen(session)
}

func (game *GameService) OnSessionClose(session *gat.Session) {
	log.Printf("[GameService] session closed, sessionid=%d", session.ID())
	if !session.HasKey("agent") {
		return
	}
	agent := session.Value("agent").(*Agent)
	agent.onSessionClose(session)
	game.agentDict.Delete(session.ID())
	session.Remove("agent")
}

func (game *GameService) getScoreRank() ([]*model.User, error) {
	var rkey string = "scorerank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.ScoreRank_Rank()
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}

func (game *GameService) getScoreDayRank(rankid int64) ([]*model.User, error) {
	var rkey string = "scoredayrank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.ScoreDayRank_Rank(rankid)
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}

func (game *GameService) getMedalRank() ([]*model.User, error) {
	var rkey string = "medalrank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.MedalRank_Rank()
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}

func (game *GameService) getMedalDayRank(rankid int64) ([]*model.User, error) {
	var rkey string = "medaldayrank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.MedalDayRank_Rank(rankid)
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}

func (game *GameService) getLevelRank() ([]*model.User, error) {
	var rkey string = "levelrank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.LevelRank_Rank()
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}

func (game *GameService) getLevelDayRank(rankid int64) ([]*model.User, error) {
	var rkey string = "leveldayrank"
	var users []*model.User
	var err error
	var expireTime int64
	var cache *RankCache
	expireTimeVal, ok := game.rankCacheExpireTime.Load(rkey)
	if ok {
		expireTime = expireTimeVal.(int64)
	}
	cacheVal, ok := game.rankCache.Load(rkey)
	if ok {
		cache = cacheVal.(*RankCache)
	}
	if cache != nil && expireTime >= time.Now().Unix() {
		users = cache.users
	} else {
		users, err = db.LevelDayRank_Rank(rankid)
		if err != nil {
			return users, err
		}
		game.rankCacheExpireTime.Store(rkey, time.Now().Unix()+60)
		game.rankCache.Store(rkey, &RankCache{
			users: users,
		})
	}
	return users, nil
}
