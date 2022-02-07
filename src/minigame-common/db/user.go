package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/model"
)

type UserInsertArgs struct {
	Openid   string
	Nickname string
	Avatar   string
	ClientIp string
}

type UserInsertReply struct {
	Userid int64
}

type User_UpdateOpenInfoArgs struct {
	Userid   int64
	Nickname string
	Avatar   string
}

type UserSaveArgs struct {
	model.User
}

type UserSaveReply struct {
	Userid int64
}

type UserGetReply struct {
	ErrNoRows bool
	User      model.User
}

type UserLoginArgs struct {
	Userid int64
}

type UserLogoutArgs struct {
	Userid int64
}

type UserUpdateCoinArgs struct {
	Userid int64
	Coin   int64
}

type UserUpdateDiamondArgs struct {
	Userid  int64
	Diamond int64
}

//读取user
func User_Get(openid string) (*model.User, error) {
	var reply = UserGetReply{}
	err := roada.Call("db", "User_Get", openid, &reply)
	if err != nil {
		return nil, err
	}
	if reply.ErrNoRows {
		return nil, nil
	}
	return &reply.User, err
}

//读取user
func (self *DbService) User_Get(r *roada.Request, args string, reply *UserGetReply) error {
	//return errors.New("some error")
	var openid string = args
	reply.ErrNoRows = false
	err := self.db.Get(&reply.User, `SELECT 
		userid, openid, nickname, avatar, score, level, medal, coin,
		diamond, logintime, logouttime, roundstarttime, roundendtime, lastroundstarttime, lastroundendtime 
		FROM user WHERE openid=?`, openid)
	if err == sql.ErrNoRows {
		reply.ErrNoRows = true
		return nil
	}
	if err != nil {
		log.Printf("User_Get err: %+v\n", err)
		return err
	}
	return nil
}

//创建账号
func User_Insert(openid string, nickname string, avatar string, clientIp string) (int64, error) {
	var args = UserInsertArgs{
		Openid:   openid,
		Nickname: nickname,
		Avatar:   avatar,
		ClientIp: clientIp,
	}
	var reply = UserInsertReply{}
	err := roada.Call("db", "User_Insert", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.Userid, nil
}

//创建账号
func (self *DbService) User_Insert(r *roada.Request, args *UserInsertArgs, reply *UserInsertReply) error {
	now := time.Now().Unix()
	result, err := self.db.Exec(`INSERT INTO user
		(openid, nickname, avatar, client_ip, logintime, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE
		nickname=VALUES(nickname), avatar=VALUES(avatar), client_ip=VALUES(client_ip), 
		logintime=VALUES(logintime), createtime=VALUES(createtime), updatetime=VALUES(updatetime)`,
		args.Openid, args.Nickname, args.Avatar, args.ClientIp, now, now, now)
	if err != nil {
		log.Printf("User_Insert err %+v\n", err)
		return err
	}
	userid, err := result.LastInsertId()
	if err != nil {
		log.Printf("[db] User_Insert err %+v\n", err)
		return err
	}
	reply.Userid = userid
	return nil
}

//保存账号数据
func User_Save(data *model.User) error {
	var args = UserSaveArgs{
		*data,
	}
	var reply = UserSaveReply{}
	err := roada.Call("db", "User_Save", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

//保存账号数据
func (self *DbService) User_Save(r *roada.Request, args *UserSaveArgs, reply *UserSaveReply) error {
	now := time.Now().Unix()
	result, err := self.db.Exec(`UPDATE user
		SET coin=?, roundstarttime=?, roundendtime=?, lastroundstarttime=?, lastroundendtime=?, 
		logintime=?, logouttime=?, updatetime=?
		WHERE userid=?`,
		args.Coin, args.RoundStartTime, args.RoundEndTime, args.LastRoundStartTime, args.LastRoundEndTime,
		args.Logintime, args.Logouttime, now, args.Userid)
	if err != nil {
		log.Printf("User_Save err %+v\n", err)
		return err
	}
	userid, err := result.LastInsertId()
	if err != nil {
		log.Printf("[db] User_Save err %+v\n", err)
		return err
	}
	reply.Userid = userid
	return nil
}

func User_UpdateOpenInfo(userid int64, nickname string, avatar string) error {
	var args = User_UpdateOpenInfoArgs{
		Userid:   userid,
		Nickname: nickname,
		Avatar:   avatar,
	}
	var reply int
	err := roada.Call("db", "User_UpdateOpenInfo", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) User_UpdateOpenInfo(r *roada.Request, args *User_UpdateOpenInfoArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`UPDATE user
		SET nickname=?, avatar=?, updatetime=?
		WHERE userid=?`,
		args.Nickname, args.Avatar, now, args.Userid)
	if err != nil {
		log.Printf("User_UpdateOpenInfo err %+v\n", err)
		return err
	}
	return nil
}

//记录登录时间
func User_Login(userid int64) error {
	var args = UserLoginArgs{Userid: userid}
	var reply int
	err := roada.Call("db", "User_Login", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

//记录登录时间
func (self *DbService) User_Login(r *roada.Request, args *UserLoginArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`UPDATE user SET logintime=?, updatetime=? WHERE userid=?`,
		now, now, args.Userid)
	if err != nil {
		log.Printf("[db] User_Login err %+v\n", err)
		return err
	}
	return nil
}

//记录登出时间
func User_Logout(userid int64) error {
	var args = UserLogoutArgs{Userid: userid}
	var reply int
	err := roada.Call("db", "User_Logout", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

//记录登出时间
func (self *DbService) User_Logout(r *roada.Request, args *UserLogoutArgs, reply *int) error {
	now := time.Now().Unix()
	_, err := self.db.Exec(`UPDATE user SET logouttime=?, updatetime=? WHERE userid=?`,
		now, now, args.Userid)
	if err != nil {
		log.Printf("User_Logout err %+v\n", err)
		return err
	}
	return nil
}

//更新金币
func User_UpdateCoin(userid int64, coin int64) error {
	var args = UserUpdateCoinArgs{Userid: userid, Coin: coin}
	var reply int
	err := roada.Call("db", "User_UpdateCoin", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

//更新金币
func (self *DbService) User_UpdateCoin(r *roada.Request, arg *UserUpdateCoinArgs, reply *int) error {
	now := time.Now().Unix()
	userid := arg.Userid
	coin := arg.Coin
	_, err := self.db.Exec(`UPDATE user SET updatetime=?, coin=coin+? WHERE userid=?`,
		now, coin, userid)
	if err != nil {
		log.Printf("[db] User_UpdateCoin err %+v\n", err)
		return err
	}
	return nil
}

//更新钻石
func User_UpdateDiamond(userid int64, diamond int64) error {
	var args = UserUpdateDiamondArgs{Userid: userid, Diamond: diamond}
	var reply int
	err := roada.Call("db", "User_UpdateDiamond", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

//更新钻石
func (self *DbService) User_UpdateDiamond(r *roada.Request, arg *UserUpdateDiamondArgs, reply *int) error {
	now := time.Now().Unix()
	userid := arg.Userid
	diamond := arg.Diamond
	_, err := self.db.Exec(`UPDATE user SET updatetime=?, diamond=diamond+? WHERE userid=?`,
		now, diamond, userid)
	if err != nil {
		log.Printf("[db] User_UpdateDiamond err %+v\n", err)
		return err
	}
	return nil
}
