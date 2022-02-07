package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type TiantianpaokuGetReply struct {
	ErrNoRows bool
	Data      model.Tiantianpaoku
}

type TiantianpaokuSaveArgs struct {
	model.Tiantianpaoku
}

func Tiantianpaoku_Save(data *model.Tiantianpaoku) error {
	args := TiantianpaokuSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Tiantianpaoku_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Tiantianpaoku_Save(r *roada.Request, args *TiantianpaokuSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO tiantianpaoku
		(userid, skinid, skinarr, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		skinid=VALUES(skinid), skinarr=VALUES(skinarr),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.SkinId, args.SkinArr, now, now)
	if err != nil {
		log.Printf("Tiantianpaoku_Save err %+v\n", err)
		return err
	}
	return nil
}

func Tiantianpaoku_Get(userid int64) (*model.Tiantianpaoku, error) {
	var reply = TiantianpaokuGetReply{}
	err := roada.Call("db", "Tiantianpaoku_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Tiantianpaoku_Get(r *roada.Request, args int64, reply *TiantianpaokuGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid, skinid, skinarr FROM tiantianpaoku WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Tiantianpaoku_Get err: %+v\n", err)
		return err
	}
	return nil
}
