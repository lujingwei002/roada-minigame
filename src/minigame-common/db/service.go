package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/roada-go/roada"
)

type DbService struct {
	db      *sqlx.DB
	cache   *redisClient
	road    *roada.Road
	service *roada.Service
	chRoad  chan *roada.Request
}

func newDbService(road *roada.Road, addr string, index int) error {
	dbconn, err := sqlx.Open("mysql", addr)
	if err != nil {
		return err
	}
	err = dbconn.Ping()
	if err != nil {
		return err
	}
	cache, err := newRedisClient()
	if err != nil {
		return err
	}
	//log.Printf("[db] connect success\n")
	var self = &DbService{
		road:   road,
		db:     dbconn,
		cache:  cache,
		chRoad: make(chan *roada.Request, 1),
	}
	//注册成服务
	serviceName := fmt.Sprintf("db%d", index)
	if err := road.LocalSet(serviceName); err != nil {
		return err
	}
	if err := road.LocalGroupAdd("db", serviceName); err != nil {
		return err
	}
	/*if err := road.Handle("db", self); err != nil {
		return err
	}*/
	if err := road.Handle(serviceName, self); err != nil {
		return err
	}
	service, err := road.Register(self)
	if err != nil {
		return err
	}
	self.service = service
	go self.loop()
	return nil
}

func (self *DbService) loop() {
	for {
		select {
		case r := <-self.chRoad:
			{
				self.service.ServeRPC(self, r)
			}
		}
	}
}

func (self *DbService) ServeRPC(r *roada.Request) {
	self.chRoad <- r
	r.Wait(5)
}
