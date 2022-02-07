package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/model"
)

type CaichengyuGetReply struct {
	ErrNoRows bool
	Data      model.Caichengyu
}

type CaichengyuSaveArgs struct {
	model.Caichengyu
}

type CaichengyuStageIndexArgs struct {
	Level int64
}

type CaichengyuStageIndexReply struct {
	ErrNoRows bool
	User      model.User
}

func Caichengyu_Save(data *model.Caichengyu) error {
	var args = CaichengyuSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Caichengyu_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Caichengyu_Save(r *roada.Request, args *CaichengyuSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO caichengyu
		(userid, rolelevel, buildlevel, level, leveltype, leveltip, hp, hpdate, gethpcount, gethpday, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		rolelevel=VALUES(rolelevel), buildlevel=VALUES(buildlevel), level=VALUES(level), 
		leveltype=VALUES(leveltype), leveltip=VALUES(leveltip),hp=VALUES(hp), hpdate=VALUES(hpdate), gethpcount=VALUES(gethpcount),
		gethpday=VALUES(gethpday), updatetime=VALUES(updatetime)`,
		args.Userid, args.RoleLevel, args.BuildLevel, args.Level,
		args.LevelType, args.LevelTip, args.Hp, args.HpDate, args.GetHpCount, args.GetHpDay, now, now)
	if err != nil {
		log.Printf("Caichengyu_Save err %+v\n", err)
		return err
	}
	return nil
}

func Caichengyu_Get(userid int64) (*model.Caichengyu, error) {
	var reply = CaichengyuGetReply{}
	err := roada.Call("db", "Caichengyu_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Caichengyu_Get(r *roada.Request, args int64, reply *CaichengyuGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, `SELECT userid, rolelevel, buildlevel, level, 
		leveltype, leveltip, hp, hpdate, gethpcount, gethpday  
		FROM caichengyu 
		WHERE userid=?`,
		userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Caichengyu_Get err: %+v\n", err)
		return err
	}
	return nil
}

func Caichengyu_StageIndex(level int64) (*model.User, error) {
	var reply = CaichengyuStageIndexReply{}
	var args = CaichengyuStageIndexArgs{
		Level: level,
	}
	err := roada.Call("db", "Caichengyu_StageIndex", &args, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.User, err
}

func (self *DbService) Caichengyu_StageIndex(r *roada.Request, args *CaichengyuStageIndexArgs, reply *CaichengyuStageIndexReply) error {
	level := args.Level

	//先读缓存
	rkey := fmt.Sprintf("levelrank@%s", conf.Ini.Game.Name)
	min := fmt.Sprintf("%d", self.encodeScoreAndTime(level, 0))
	if rows, err := self.cache.ZRangeByScoreWithScores(rkey, min, "+inf", 0, 1); err == nil && len(rows) > 0 {
		if userid, err := strconv.ParseInt(rows[0].Member.(string), 10, 64); err == nil {
			if err := self.db.Get(&reply.User, `SELECT userid, openid, nickname, avatar, level 
				FROM user 
				WHERE userid=? 
				LIMIT 1`, userid); err == nil {
				return nil
			}
		}
	}
	//再读数据库
	err := self.db.Get(&reply.User, `SELECT userid, openid, nickname, avatar, level 
		FROM user 
		WHERE level>=? 
		ORDER BY level ASC, updatetime DESC 
		LIMIT 1`, level)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Caichengyu_StageIndex err: %+v\n", err)
		return err
	}
	return nil
}
