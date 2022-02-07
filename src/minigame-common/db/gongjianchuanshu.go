package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type GongjianchuanshuGetReply struct {
	ErrNoRows bool
	Data      model.Gongjianchuanshu
}

type GongjianchuanshuSaveArgs struct {
	model.Gongjianchuanshu
}

func Gongjianchuanshu_Save(data *model.Gongjianchuanshu) error {
	args := GongjianchuanshuSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Gongjianchuanshu_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Gongjianchuanshu_Save(r *roada.Request, args *GongjianchuanshuSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO gongjianchuanshu
		(userid, level, skinarr, skinid, shoparr, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		level=VALUES(level), skinarr=VALUES(skinarr), skinid=VALUES(skinid), shoparr=VALUES(shoparr),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Level, args.SkinArr, args.SkinId, args.ShopArr, now, now)
	if err != nil {
		log.Printf("Gongjianchuanshu_Save err %+v\n", err)
		return err
	}
	return nil
}

func Gongjianchuanshu_Get(userid int64) (*model.Gongjianchuanshu, error) {
	var reply = GongjianchuanshuGetReply{}
	err := roada.Call("db", "Gongjianchuanshu_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Gongjianchuanshu_Get(r *roada.Request, args int64, reply *GongjianchuanshuGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid, level, skinarr, skinid, shoparr FROM gongjianchuanshu WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Gongjianchuanshu_Get err: %+v\n", err)
		return err
	}
	return nil
}
