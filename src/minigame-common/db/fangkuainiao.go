package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type FangkuainiaoGetReply struct {
	ErrNoRows bool
	Data      model.Fangkuainiao
}

type FangkuainiaoSaveArgs struct {
	model.Fangkuainiao
}

func Fangkuainiao_Save(data *model.Fangkuainiao) error {
	args := FangkuainiaoSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Fangkuainiao_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Fangkuainiao_Save(r *roada.Request, args *FangkuainiaoSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO fangkuainiao
		(userid, level, birdarr, birdid, signtime, signday, getgoldcount, getgoldtime, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		level=VALUES(level),
		birdarr=VALUES(birdarr), birdid=VALUES(birdid), signtime=VALUES(signtime), 
		signday=VALUES(signday), getgoldcount=VALUES(getgoldcount), getgoldtime=VALUES(getgoldtime),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Level, args.BirdArr, args.BirdId, args.SignTime,
		args.SignDay, args.GetGoldCount, args.GetGoldTime, now, now)
	if err != nil {
		log.Printf("Fangkuainiao_Save err %+v\n", err)
		return err
	}
	return nil
}

func Fangkuainiao_Get(userid int64) (*model.Fangkuainiao, error) {
	var reply = FangkuainiaoGetReply{}
	err := roada.Call("db", "Fangkuainiao_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Fangkuainiao_Get(r *roada.Request, args int64, reply *FangkuainiaoGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid, level, birdarr, birdid, signtime, signday, getgoldcount, getgoldtime FROM fangkuainiao WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Fangkuainiao_Get err: %+v\n", err)
		return err
	}
	return nil
}
