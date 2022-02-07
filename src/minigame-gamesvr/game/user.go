package game

import (
	"errors"

	"github.com/shark/minigame-common/db"
	"github.com/shark/minigame-common/model"
	"github.com/shark/minigame-common/tlog"
)

var errDiamondNotEnough = errors.New("diamond not enough")
var errCoinNotEnough = errors.New("coin not enough")

type User struct {
	model.User
	NewHand    bool
	Openid     string
	Token      string
	Gold       int64
	Today      int64
	ScoreToday int64 // 当天分数
	LevelToday int64 // 当天关卡
	MedalToday int64 // 当天金牌
	Dirty      bool
}

func (user *User) AddDiamond(diamond int64, reason string) error {
	if err := db.User_UpdateDiamond(user.Userid, diamond); err != nil {
		return err
	}
	user.Diamond = user.Diamond + diamond
	tlog.AddDiamond(user.Openid, user.Userid, diamond, reason)
	return nil
}

func (user *User) DecDiamond(diamond int64, reason string) error {
	if user.Diamond < diamond {
		return errDiamondNotEnough
	}
	if err := db.User_UpdateDiamond(user.Userid, -diamond); err != nil {
		return err
	}
	user.Diamond = user.Diamond - diamond
	tlog.DecDiamond(user.Openid, user.Userid, diamond, reason)
	return nil
}

func (user *User) AddCoin(coin int64, reason string) error {
	/*if err := db.User_UpdateCoin(user.Userid, coin); err != nil {
		return err
	}*/
	user.Coin = user.Coin + coin
	user.Dirty = true
	//tlog.AddCoin(user.Userid, coin, reason)
	return nil
}

func (user *User) DecCoin(coin int64, reason string) error {
	if user.Coin < coin {
		return errCoinNotEnough
	}
	/*if err := db.User_UpdateCoin(user.Userid, -coin); err != nil {
		return err
	}*/
	user.Coin = user.Coin - coin
	user.Dirty = true
	//tlog.DecCoin(user.Openid, user.Userid, coin, reason)
	return nil
}
