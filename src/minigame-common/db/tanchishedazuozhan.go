package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type TanchishedazuozhanGetReply struct {
	ErrNoRows bool
	Data      model.Tanchishedazuozhan
}

type TanchishedazuozhanSaveArgs struct {
	model.Tanchishedazuozhan
}

func Tanchishedazuozhan_Save(data *model.Tanchishedazuozhan) error {
	args := TanchishedazuozhanSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Tanchishedazuozhan_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Tanchishedazuozhan_Save(r *roada.Request, args *TanchishedazuozhanSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO tanchishedazuozhan
		(userid, skinarr, skinid, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		skinarr=VALUES(skinarr), skinid=VALUES(skinid),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.SkinArr, args.SkinId, now, now)
	if err != nil {
		log.Printf("Tanchishedazuozhan_Save err %+v\n", err)
		return err
	}
	return nil
}

func Tanchishedazuozhan_Get(userid int64) (*model.Tanchishedazuozhan, error) {
	var reply = TanchishedazuozhanGetReply{}
	err := roada.Call("db", "Tanchishedazuozhan_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Tanchishedazuozhan_Get(r *roada.Request, args int64, reply *TanchishedazuozhanGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid, skinarr, skinid FROM tanchishedazuozhan WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Tanchishedazuozhan_Get err: %+v\n", err)
		return err
	}
	return nil
}
