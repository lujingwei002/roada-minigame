package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type Paopaolong_GetReply struct {
	ErrNoRows bool
	Data      model.Paopaolong
}

type Paopaolong_SaveArgs struct {
	model.Paopaolong
}

type Paopaolong_LevelGetReply struct {
	LevelArr []*model.PaopaolongLevel
}

type Paopaolong_LevelSaveArgs struct {
	model.PaopaolongLevel
}

func Paopaolong_Save(data *model.Paopaolong) error {
	args := Paopaolong_SaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Paopaolong_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Paopaolong_Save(r *roada.Request, args *Paopaolong_SaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO paopaolong
		(userid, level, itemarr, hp, freedraw_time, 
		new_pack_redeemed, shop_free_diamond_time, shop_free_diamond_time2, last_sign_time, signed_time, 
		first_sign_time, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?,  ?, ?, ?) ON DUPLICATE KEY UPDATE
		level=VALUES(level), itemarr=VALUES(itemarr), hp=VALUES(hp), freedraw_time=VALUES(freedraw_time),
		new_pack_redeemed=VALUES(new_pack_redeemed),
		shop_free_diamond_time=VALUES(shop_free_diamond_time), shop_free_diamond_time2=VALUES(shop_free_diamond_time2),
		last_sign_time=VALUES(last_sign_time), signed_time=VALUES(signed_time),first_sign_time=VALUES(first_sign_time),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Level, args.ItemArr, args.Hp, args.FreedrawTime,
		args.NewPackRedeemed, args.ShopFreeDiamondTime, args.ShopFreeDiamondTime2,
		args.LastSignTime, args.SignedTime, args.FirstSignTime,
		now, now)
	if err != nil {
		log.Printf("Paopaolong_Save err %+v\n", err)
		return err
	}
	return nil
}

func Paopaolong_Get(userid int64) (*model.Paopaolong, error) {
	var reply = Paopaolong_GetReply{}
	err := roada.Call("db", "Paopaolong_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Paopaolong_Get(r *roada.Request, args int64, reply *Paopaolong_GetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, `SELECT userid, level, itemarr, hp, 
		freedraw_time, new_pack_redeemed, shop_free_diamond_time, shop_free_diamond_time2,
		last_sign_time, signed_time, first_sign_time
		FROM paopaolong WHERE userid=?`, userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Paopaolong_Get err: %+v\n", err)
		return err
	}
	return nil
}

func Paopaolong_LevelGet(userid int64) ([]*model.PaopaolongLevel, error) {
	var reply = Paopaolong_LevelGetReply{}
	err := roada.Call("db", "Paopaolong_LevelGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.LevelArr, err
}

func (self *DbService) Paopaolong_LevelGet(r *roada.Request, args int64, reply *Paopaolong_LevelGetReply) error {
	var userid int64 = args
	reply.LevelArr = make([]*model.PaopaolongLevel, 0)
	err := self.db.Select(&reply.LevelArr, `SELECT userid, level, sec, lose, score, star		
		FROM paopaolong_level WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Paopaolong_LevelGet err: %+v\n", err)
		return err
	}
	return nil
}

func Paopaolong_LevelSave(userid int64, level int64, sec int64, lose int32, score int64, star int64) error {
	args := Paopaolong_LevelSaveArgs{
		model.PaopaolongLevel{
			Userid: userid,
			Level:  level,
			Sec:    sec,
			Lose:   lose,
			Score:  score,
			Star:   star,
		},
	}
	var reply int
	err := roada.Call("db", "Paopaolong_LevelSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Paopaolong_LevelSave(r *roada.Request, args *Paopaolong_LevelSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO paopaolong_level
		(userid, level, sec, lose, score, star, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		sec=VALUES(sec), lose=VALUES(lose), score=VALUES(score),
		star=VALUES(star),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Level, args.Sec, args.Lose, args.Score, args.Star,
		now, now)
	if err != nil {
		log.Printf("Paopaolong_LevelSave err %+v\n", err)
		return err
	}
	return nil
}
