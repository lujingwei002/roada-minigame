package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type HuanlemaomibeiGetReply struct {
	ErrNoRows bool
	Data      model.Huanlemaomibei
}

type HuanlemaomibeiSaveArgs struct {
	model.Huanlemaomibei
}

type Huanlemaomibei_LevelGetReply struct {
	LevelArr []*model.HuanlemaomibeiLevel
}

type Huanlemaomibei_LevelSaveArgs struct {
	model.HuanlemaomibeiLevel
}

func Huanlemaomibei_Save(data *model.Huanlemaomibei) error {
	args := HuanlemaomibeiSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Huanlemaomibei_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Huanlemaomibei_Save(r *roada.Request, args *HuanlemaomibeiSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO huanlemaomibei
		(userid, freedraw_time, inkarr, cuparr, inkid, cupid, signday, lastsigntime, signchecked, fly_time, offline_time, hp, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		freedraw_time=VALUES(freedraw_time), inkarr=VALUES(inkarr), cuparr=VALUES(cuparr),
		inkid=VALUES(inkid), cupid=VALUES(cupid),
		signday=VALUES(signday), lastsigntime=VALUES(lastsigntime), signchecked=VALUES(signchecked),
		fly_time=VALUES(fly_time), offline_time=VALUES(offline_time), hp=VALUES(hp),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.FreedrawTime, args.InkArr, args.CupArr, args.InkId, args.CupId,
		args.SignDay, args.LastSignTime, args.SignChecked, args.FlyTime, args.OfflineTime, args.Hp, now, now)
	if err != nil {
		log.Printf("Huanlemaomibei_Save err %+v\n", err)
		return err
	}
	return nil
}

func Huanlemaomibei_Get(userid int64) (*model.Huanlemaomibei, error) {
	var reply = HuanlemaomibeiGetReply{}
	err := roada.Call("db", "Huanlemaomibei_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Huanlemaomibei_Get(r *roada.Request, args int64, reply *HuanlemaomibeiGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, `SELECT 
		userid, freedraw_time, inkarr, cuparr, inkid, cupid,
		signday, signchecked, lastsigntime, fly_time, offline_time, hp
		FROM huanlemaomibei WHERE userid=?`, userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Huanlemaomibei_Get err: %+v\n", err)
		return err
	}
	return nil
}

func Huanlemaomibei_LevelGet(userid int64) ([]*model.HuanlemaomibeiLevel, error) {
	var reply = Huanlemaomibei_LevelGetReply{}
	err := roada.Call("db", "Huanlemaomibei_LevelGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.LevelArr, err
}

func (self *DbService) Huanlemaomibei_LevelGet(r *roada.Request, args int64, reply *Huanlemaomibei_LevelGetReply) error {
	var userid int64 = args
	reply.LevelArr = make([]*model.HuanlemaomibeiLevel, 0)
	err := self.db.Select(&reply.LevelArr, `SELECT userid, section, level, unlocked, star, coin		
		FROM huanlemaomibei_level WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Huanlemaomibei_LevelGet err: %+v\n", err)
		return err
	}
	return nil
}

func Huanlemaomibei_LevelSave(userid int64, section int32, level int32, unlock int32, star int32, coin int64) error {
	args := Huanlemaomibei_LevelSaveArgs{
		model.HuanlemaomibeiLevel{
			Userid:  userid,
			Section: section,
			Level:   level,
			Unlock:  unlock,
			Star:    star,
			Coin:    coin,
		},
	}
	var reply int
	err := roada.Call("db", "Huanlemaomibei_LevelSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Huanlemaomibei_LevelSave(r *roada.Request, args *Huanlemaomibei_LevelSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO huanlemaomibei_level
		(userid, section, level, unlocked, star, coin, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		unlocked=VALUES(unlocked), star=VALUES(star), coin=VALUES(coin),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Section, args.Level, args.Unlock, args.Star, args.Coin,
		now, now)
	if err != nil {
		log.Printf("Huanlemaomibei_LevelSave err %+v\n", err)
		return err
	}
	return nil
}
