package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type GateStatInsertArgs struct {
	GameId       int32
	NodeName     string
	NodeFullName string
	Host         string
	Port         int32
	OnlineNum    int32
}

type GateStatGetArgs struct {
}

type GateStatGetReply struct {
	ErrNoRows bool
	Gate      model.GateStat
}

func GateStat_Get() (*model.GateStat, error) {
	var args = GateStatGetArgs{}
	var reply = GateStatGetReply{}
	err := roada.Call("db", "GateStat_Get", &args, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Gate, err
}

func (self *DbService) GateStat_Get(r *roada.Request, args *GateStatGetArgs, reply *GateStatGetReply) error {
	//return errors.New("some error")
	now := time.Now().Unix()
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Gate, `SELECT nodename, nodefullname, host, port, onlinenum 
	FROM gatestat 
	WHERE updatetime>? AND host!= ""
	ORDER BY onlinenum ASC 
	LIMIT 1`, now-120)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("GateStat_Get err: %+v\n", err)
		return err
	}
	return nil
}

func GateStat_Insert(gameId int32, nodeName string, nodeFullName string, host string, port int32, onlineNum int32) error {
	var args = GateStatInsertArgs{
		GameId:       gameId,
		NodeName:     nodeName,
		NodeFullName: nodeFullName,
		Host:         host,
		Port:         port,
		OnlineNum:    onlineNum,
	}
	var reply int
	err := roada.Call("db", "GateStat_Insert", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) GateStat_Insert(r *roada.Request, args *GateStatInsertArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO gatestat
		(nodename, nodefullname, gameid, host, port, onlinenum, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		onlinenum=VALUES(onlinenum), gameid=VALUES(gameid), host=VALUES(host), port=VALUES(port), updatetime=VALUES(updatetime)`,
		args.NodeName, args.NodeFullName, args.GameId, args.Host, args.Port, args.OnlineNum, now, now)
	if err != nil {
		log.Printf("db.GateStat_Insert err %+v\n", err)
		return err
	}
	return nil
}
