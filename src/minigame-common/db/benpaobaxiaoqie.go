package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type BpbxqGetReply struct {
	ErrNoRows bool
	Data      model.Benpaobaxiaoqie
}

type BpbxqInsertArgs struct {
	model.Benpaobaxiaoqie
}

type BpbxqSaveArgs struct {
	Userid       int64
	Skin         int32
	SkinArr      string
	LastSignTime int64
	SignTimes    int64
}

func Bpbxq_Save(data *model.Benpaobaxiaoqie) error {
	var args = BpbxqInsertArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Bpbxq_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Bpbxq_Save(r *roada.Request, args *BpbxqInsertArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO benpaobaxiaoqie
		(userid, skin, skinarr, lastsigntime, signtimes, createtime, updatetime) VALUES
		(?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		skin=VALUES(skin), skinarr=VALUES(skinarr), lastsigntime=VALUES(lastsigntime), signtimes=VALUES(signtimes), updatetime=VALUES(updatetime)`,
		args.Userid, args.Skin, args.SkinArr, args.LastSignTime, args.SignTimes, now, now)
	if err != nil {
		log.Printf("Bpbxq_Save err %+v\n", err)
		return err
	}
	return nil
}

func Bpbxq_Get(userid int64) (*model.Benpaobaxiaoqie, error) {
	var reply = BpbxqGetReply{}
	err := roada.Call("db", "Bpbxq_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Bpbxq_Get(r *roada.Request, args int64, reply *BpbxqGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, "SELECT userid, skin, skinarr, lastsigntime, signtimes FROM benpaobaxiaoqie WHERE userid=?", userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Bpbxq_Get err: %+v\n", err)
		return err
	}
	return nil
}
