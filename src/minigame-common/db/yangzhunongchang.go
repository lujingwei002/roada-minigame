package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type YangzhunongchangGetReply struct {
	ErrNoRows bool
	Data      model.Yangzhunongchang
}

type YangzhunongchangSaveArgs struct {
	model.Yangzhunongchang
}

type Yangzhunongchang_ItemGetReply struct {
	ItemArr []*model.YangzhunongchangItem
}

type Yangzhunongchang_ItemSaveArgs struct {
	model.YangzhunongchangItem
}

type Yangzhunongchang_PigGetReply struct {
	PigArr []*model.YangzhunongchangPig
}

type Yangzhunongchang_PigSaveArgs struct {
	model.YangzhunongchangPig
}

type Yangzhunongchang_PigDelArgs struct {
	model.YangzhunongchangPig
}

type Yangzhunongchang_BreedPigGetReply struct {
	PigArr []*model.YangzhunongchangBreedPig
}

type Yangzhunongchang_BreedPigSaveArgs struct {
	model.YangzhunongchangBreedPig
}

type Yangzhunongchang_BreedPigDelArgs struct {
	model.YangzhunongchangBreedPig
}

type Yangzhunongchang_UsuGetReply struct {
	UsuArr []*model.YangzhunongchangUsu
}

type Yangzhunongchang_UsuSaveArgs struct {
	model.YangzhunongchangUsu
}

type Yangzhunongchang_FoodGetReply struct {
	FoodArr []*model.YangzhunongchangFood
}

type Yangzhunongchang_FoodSaveArgs struct {
	model.YangzhunongchangFood
}

type Yangzhunongchang_FoodDelArgs struct {
	model.YangzhunongchangFood
}

type Yangzhunongchang_TaskGetReply struct {
	TaskArr []*model.YangzhunongchangTask
}

type Yangzhunongchang_TaskSaveArgs struct {
	model.YangzhunongchangTask
}

func Yangzhunongchang_Save(data *model.Yangzhunongchang) error {
	args := YangzhunongchangSaveArgs{
		*data,
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_Save(r *roada.Request, args *YangzhunongchangSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang
		(userid, farm_lv, farm_lv_name, farm_lv_exp, farm_lv_exp_cur, award_num, award_time, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		farm_lv=VALUES(farm_lv),farm_lv_name=VALUES(farm_lv_name),farm_lv_exp=VALUES(farm_lv_exp),farm_lv_exp_cur=VALUES(farm_lv_exp_cur),
		award_num=VALUES(award_num),award_time=VALUES(award_time),
		updatetime=VALUES(updatetime)`,
		args.Userid, args.FarmLv, args.FarmLvName, args.FarmLvExp, args.FarmLvExpCur, args.AwardNum, args.AwardTime, now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_Save err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_Get(userid int64) (*model.Yangzhunongchang, error) {
	var reply = YangzhunongchangGetReply{}
	err := roada.Call("db", "Yangzhunongchang_Get", userid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.Data, err
}

func (self *DbService) Yangzhunongchang_Get(r *roada.Request, args int64, reply *YangzhunongchangGetReply) error {
	var userid int64 = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.Data, `SELECT userid,
	 farm_lv, farm_lv_name, farm_lv_exp, farm_lv_exp_cur,
	 award_num, award_time
	 FROM yangzhunongchang WHERE userid=?`, userid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("Yangzhunongchang_Get err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_ItemGet(userid int64) ([]*model.YangzhunongchangItem, error) {
	var reply = Yangzhunongchang_ItemGetReply{}
	err := roada.Call("db", "Yangzhunongchang_ItemGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.ItemArr, err
}

func (self *DbService) Yangzhunongchang_ItemGet(r *roada.Request, args int64, reply *Yangzhunongchang_ItemGetReply) error {
	var userid int64 = args
	reply.ItemArr = make([]*model.YangzhunongchangItem, 0)
	err := self.db.Select(&reply.ItemArr, `SELECT userid, id, num	
		FROM yangzhunongchang_item WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Yangzhunongchang_ItemGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_ItemSave(userid int64, id int32, num int32) error {
	args := Yangzhunongchang_ItemSaveArgs{
		model.YangzhunongchangItem{
			Userid: userid,
			Id:     id,
			Num:    num,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_ItemSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_ItemSave(r *roada.Request, args *Yangzhunongchang_ItemSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang_item
		(userid, id, num, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		num=VALUES(num), updatetime=VALUES(updatetime)`,
		args.Userid, args.Id, args.Num,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_ItemSave err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_PigGet(userid int64) ([]*model.YangzhunongchangPig, error) {
	var reply = Yangzhunongchang_PigGetReply{}
	err := roada.Call("db", "Yangzhunongchang_PigGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.PigArr, err
}

func (self *DbService) Yangzhunongchang_PigGet(r *roada.Request, args int64, reply *Yangzhunongchang_PigGetReply) error {
	var userid int64 = args
	reply.PigArr = make([]*model.YangzhunongchangPig, 0)
	err := self.db.Select(&reply.PigArr, `SELECT userid, id, data, createtime	
		FROM yangzhunongchang_pig WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Yangzhunongchang_PigGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_PigSave(userid int64, id string, data string) error {
	args := Yangzhunongchang_PigSaveArgs{
		model.YangzhunongchangPig{
			Userid: userid,
			Id:     id,
			Data:   data,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_PigSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_PigSave(r *roada.Request, args *Yangzhunongchang_PigSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang_pig
		(userid, id, data, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		data=VALUES(data), updatetime=VALUES(updatetime)`,
		args.Userid, args.Id, args.Data,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_PigSave err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_PigDel(userid int64, id string) error {
	args := Yangzhunongchang_PigDelArgs{
		model.YangzhunongchangPig{
			Userid: userid,
			Id:     id,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_PigDel", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_PigDel(r *roada.Request, args *Yangzhunongchang_PigDelArgs, reply *int) error {
	_, err := self.db.Exec(`DELETE  	
		FROM yangzhunongchang_pig WHERE userid=? AND id=?`, args.Userid, args.Id)
	if err != nil {
		log.Printf("Yangzhunongchang_PigDel err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_BreedPigGet(userid int64) ([]*model.YangzhunongchangBreedPig, error) {
	var reply = Yangzhunongchang_BreedPigGetReply{}
	err := roada.Call("db", "Yangzhunongchang_BreedPigGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.PigArr, err
}

func (self *DbService) Yangzhunongchang_BreedPigGet(r *roada.Request, args int64, reply *Yangzhunongchang_BreedPigGetReply) error {
	var userid int64 = args
	reply.PigArr = make([]*model.YangzhunongchangBreedPig, 0)
	err := self.db.Select(&reply.PigArr, `SELECT userid, id, data	
		FROM yangzhunongchang_breedpig WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Yangzhunongchang_BreedPigGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_BreedPigSave(userid int64, id string, data string) error {
	args := Yangzhunongchang_BreedPigSaveArgs{
		model.YangzhunongchangBreedPig{
			Userid: userid,
			Id:     id,
			Data:   data,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_BreedPigSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_BreedPigSave(r *roada.Request, args *Yangzhunongchang_BreedPigSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang_breedpig
		(userid, id, data, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		data=VALUES(data), updatetime=VALUES(updatetime)`,
		args.Userid, args.Id, args.Data,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_BreedPigSave err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_BreedPigDel(userid int64, id string) error {
	args := Yangzhunongchang_BreedPigDelArgs{
		model.YangzhunongchangBreedPig{
			Userid: userid,
			Id:     id,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_BreedPigDel", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_BreedPigDel(r *roada.Request, args *Yangzhunongchang_BreedPigDelArgs, reply *int) error {
	_, err := self.db.Exec(`DELETE  	
		FROM yangzhunongchang_breedpig WHERE userid=? AND id=?`, args.Userid, args.Id)
	if err != nil {
		log.Printf("Yangzhunongchang_BreedPigDel err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_UsuGet(userid int64) ([]*model.YangzhunongchangUsu, error) {
	var reply = Yangzhunongchang_UsuGetReply{}
	err := roada.Call("db", "Yangzhunongchang_UsuGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.UsuArr, err
}

func (self *DbService) Yangzhunongchang_UsuGet(r *roada.Request, args int64, reply *Yangzhunongchang_UsuGetReply) error {
	var userid int64 = args
	reply.UsuArr = make([]*model.YangzhunongchangUsu, 0)
	err := self.db.Select(&reply.UsuArr, `SELECT userid, id	
		FROM yangzhunongchang_usu WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Yangzhunongchang_UsuGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_UsuSave(userid int64, id int32) error {
	args := Yangzhunongchang_UsuSaveArgs{
		model.YangzhunongchangUsu{
			Userid: userid,
			Id:     id,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_UsuSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_UsuSave(r *roada.Request, args *Yangzhunongchang_UsuSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang_usu
		(userid, id, createtime, updatetime) 
		VALUES(?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		updatetime=VALUES(updatetime)`,
		args.Userid, args.Id,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_UsuSave err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_FoodGet(userid int64) ([]*model.YangzhunongchangFood, error) {
	var reply = Yangzhunongchang_FoodGetReply{}
	err := roada.Call("db", "Yangzhunongchang_FoodGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.FoodArr, err
}

func (self *DbService) Yangzhunongchang_FoodGet(r *roada.Request, args int64, reply *Yangzhunongchang_FoodGetReply) error {
	var userid int64 = args
	reply.FoodArr = make([]*model.YangzhunongchangFood, 0)
	err := self.db.Select(&reply.FoodArr, `SELECT userid, id, data	
		FROM yangzhunongchang_food WHERE userid=?`, userid)
	if err != nil {
		log.Printf("Yangzhunongchang_FoodGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_FoodSave(userid int64, id string, data string) error {
	args := Yangzhunongchang_FoodSaveArgs{
		model.YangzhunongchangFood{
			Userid: userid,
			Id:     id,
			Data:   data,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_FoodSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_FoodSave(r *roada.Request, args *Yangzhunongchang_FoodSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`INSERT INTO yangzhunongchang_food
		(userid, id, data, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		data=VALUES(data), updatetime=VALUES(updatetime)`,
		args.Userid, args.Id, args.Data,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_FoodSave err %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_FoodDel(userid int64, id string) error {
	args := Yangzhunongchang_FoodDelArgs{
		model.YangzhunongchangFood{
			Userid: userid,
			Id:     id,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_FoodDel", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_FoodDel(r *roada.Request, args *Yangzhunongchang_FoodDelArgs, reply *int) error {
	_, err := self.db.Exec(`DELETE  	
		FROM yangzhunongchang_food WHERE userid=? AND id=?`, args.Userid, args.Id)
	if err != nil {
		log.Printf("Yangzhunongchang_FoodDel err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_TaskGet(userid int64) ([]*model.YangzhunongchangTask, error) {
	var reply = Yangzhunongchang_TaskGetReply{}
	err := roada.Call("db", "Yangzhunongchang_TaskGet", userid, &reply)
	if err != nil {
		return nil, err
	}
	return reply.TaskArr, err
}

func (self *DbService) Yangzhunongchang_TaskGet(r *roada.Request, args int64, reply *Yangzhunongchang_TaskGetReply) error {
	var userid int64 = args
	reply.TaskArr = make([]*model.YangzhunongchangTask, 0)
	err := self.db.Select(&reply.TaskArr, "SELECT userid, id, `index`, count FROM yangzhunongchang_task WHERE userid=?", userid)
	if err != nil {
		log.Printf("Yangzhunongchang_TaskGet err: %+v\n", err)
		return err
	}
	return nil
}

func Yangzhunongchang_TaskSave(userid int64, id string, index int32, count string) error {
	args := Yangzhunongchang_TaskSaveArgs{
		model.YangzhunongchangTask{
			Userid: userid,
			Id:     id,
			Index:  index,
			Count:  count,
		},
	}
	var reply int
	err := roada.Call("db", "Yangzhunongchang_TaskSave", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) Yangzhunongchang_TaskSave(r *roada.Request, args *Yangzhunongchang_TaskSaveArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec("INSERT INTO yangzhunongchang_task(userid, id, `index`, count, createtime, updatetime) VALUES(?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `index`=VALUES(`index`), count=VALUES(count), updatetime=VALUES(updatetime)",
		args.Userid, args.Id, args.Index, args.Count,
		now, now)
	if err != nil {
		log.Printf("Yangzhunongchang_TaskSave err %+v\n", err)
		return err
	}
	return nil
}
