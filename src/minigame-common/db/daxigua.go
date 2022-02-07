package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type DaxiguaGetReply struct {
	ErrNoRows bool
	Data      model.Daxigua
}

type DaxiguaSaveArgs struct {
	model.Daxigua
}

func Daxigua_Save(data *model.Daxigua) error {
	args := DaxiguaSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Daxigua_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Daxigua_Save(r *roada.Request, args *DaxiguaSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO daxigua
		(userid, createtime, updatetime) 
		VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE
		updatetime=VALUES(updatetime)`,
		args.Userid, now, now)
	if err != nil {
		log.Printf("Daxigua_Save err %+v\n", err)
		return err
	}
	return nil
}

func Daxigua_Get(userid int64) (*model.Daxigua, error) {
	var reply = DaxiguaGetReply{}
	err := roada.Call("db", "Daxigua_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Daxigua_Get(r *roada.Request, args int64, reply *DaxiguaGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid FROM daxigua WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Daxigua_Get err: %+v\n", err)
		return err
	}
	return nil
}
