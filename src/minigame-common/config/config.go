package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/shark/minigame-common/conf"
)

var Caichengyu struct {
	Answer []string
	Role   []int32
	Build  []int32
}

type FangkuainiaoSign struct {
	Id           int32 `json:"id"`
	Type         int32 `json:"type"`
	Group        int32 `json:"group"`
	Days         int32 `json:"days"`
	RewardsMoney int32 `json:"rewardsMoney"`
	RxtreRewards int32 `json:"extreRewards"`
}

var Fangkuainiao struct {
	Bird []struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Type     int32  `json:"type"`
		PayMoney int64  `json:"payMoney"`
		BirdIcon string `json:"birdIcon"`
		EggIcon  string `json:"eggIcon"`
	}
	Sign []struct {
		Id           int32 `json:"id"`
		Type         int32 `json:"type"`
		Group        int32 `json:"group"`
		Days         int32 `json:"days"`
		RewardsMoney int32 `json:"rewardsMoney"`
		RxtreRewards int32 `json:"extreRewards"`
	}
}

var Paopaolong struct {
	Lucky []struct {
		Item   int32 `json:"item"`
		Num    int32 `json:"num"`
		Chance int32 `json:chance`
	}
	Sign []struct {
		Item int32 `json:"item"`
		Num  int32 `json:"num"`
		Type int32 `json:type`
	}
	Shop []struct {
		Item int32 `json:"item"`
		Num  int32 `json:"num"`
		Cost int64 `json:cost`
	}
	NewPack []struct {
		Item int32 `json:"item"`
		Num  int32 `json:"num"`
	}
}

var Tanchishedazuozhan struct {
	Skin []struct {
		Id    int32 `json:"id"`
		Type  int32 `json:"type"`
		Price int64 `json:"price"`
	}
}

var Tiantianpaoku struct {
	Skin []struct {
		Id   int32 `json:"id"`
		Coin int64 `json:coin`
	}
}

var Huanlemaomibei struct {
	Roulette []struct {
		Id     int32  `json:"id"`
		Weight int32  `json:"weight"`
		Name   string `json:"name"`
		Gold   int64  `json:"gold"`
		Skin   int32  `json:"skin"`
	}
	Ink []struct {
		Id       int32  `json:"id"`
		Path     string `json:"path"`
		Coin     int64  `json:"coin"`
		Video    int32  `json:"video"`
		Daily    int32  `json:"daily"`
		Time     int32  `json:"time"`
		Like     int32  `json:"like"`
		Roulette int32  `json:"roulette"`
	}
	Cup []struct {
		Id       int32  `json:"id"`
		Path     string `json:"path"`
		Coin     int64  `json:"coin"`
		Video    int32  `json:"video"`
		Daily    int32  `json:"daily"`
		Time     int32  `json:"time"`
		Like     int32  `json:"like"`
		Roulette int32  `json:"roulette"`
	}
	Sign []struct {
		Id   int32 `json:"id"`
		Gold int64 `json:"gold"`
		Ink  int32 `json:"ink"`
		Cup  int32 `json:"cup"`
	}
	Section []struct {
		Name     string `json:"name"`
		LevelArr []struct {
			Name     string `json:"name"`
			Section  int32  `json:"section"`
			Level    int32  `json:"level"`
			Price    int64  `json:"price"`
			Coin     int64  `json:"coin"`
			DrinkArr []struct {
				Name  string `json:"name"`
				Price int64  `json:"price"`
				Ink   int32  `json:"ink"`
			} `json:"array"`
		} `json:"array"`
	}
}

func LoadCaichengyu() error {
	if err := loadFile("/caichengyu/answer.json", &Caichengyu.Answer); err != nil {
		return err
	}
	if err := loadFile("/caichengyu/role.json", &Caichengyu.Role); err != nil {
		return err
	}
	if err := loadFile("/caichengyu/build.json", &Caichengyu.Build); err != nil {
		return err
	}
	//log.Printf("[config] Caichengyu %+v\n", Caichengyu)
	return nil
}

func LoadPaopaolong() error {
	if err := loadFile("/paopaolong/lucky.json", &Paopaolong.Lucky); err != nil {
		return err
	}
	if err := loadFile("/paopaolong/sign.json", &Paopaolong.Sign); err != nil {
		return err
	}
	if err := loadFile("/paopaolong/shop.json", &Paopaolong.Shop); err != nil {
		return err
	}
	if err := loadFile("/paopaolong/newpack.json", &Paopaolong.NewPack); err != nil {
		return err
	}
	log.Printf("[config] Paopaolong %+v\n", Paopaolong)
	return nil
}

func LoadFangkuainiao() error {
	if err := loadFile("/fangkuainiao/bird.json", &Fangkuainiao.Bird); err != nil {
		return err
	}
	if err := loadFile("/fangkuainiao/sign.json", &Fangkuainiao.Sign); err != nil {
		return err
	}
	//log.Printf("[config] Fangkuainiao %+v\n", Fangkuainiao)
	return nil
}

func LoadTanchishedazuozhan() error {
	if err := loadFile("/tanchishedazuozhan/skin.json", &Tanchishedazuozhan.Skin); err != nil {
		return err
	}
	log.Printf("[config] tanchishedazuozhan %+v\n", Tanchishedazuozhan)
	return nil
}

func LoadTiantianpaoku() error {
	if err := loadFile("/tiantianpaoku/skin.json", &Tiantianpaoku.Skin); err != nil {
		return err
	}
	log.Printf("[config] tiantianpaoku %+v\n", Tiantianpaoku)
	return nil
}

func LoadHuanlemaomibei() error {
	if err := loadFile("/huanlemaomibei/roulette.json", &Huanlemaomibei.Roulette); err != nil {
		return err
	}
	if err := loadFile("/huanlemaomibei/cup.json", &Huanlemaomibei.Cup); err != nil {
		return err
	}
	if err := loadFile("/huanlemaomibei/ink.json", &Huanlemaomibei.Ink); err != nil {
		return err
	}
	if err := loadFile("/huanlemaomibei/sign.json", &Huanlemaomibei.Sign); err != nil {
		return err
	}
	if err := loadFile("/huanlemaomibei/level.json", &Huanlemaomibei.Section); err != nil {
		return err
	}
	log.Printf("[config] huanlemaomibei %+v\n", Huanlemaomibei)
	return nil
}

func loadFile(path string, v interface{}) error {
	path = conf.Ini.Game.ConfigDir + path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, v); err != nil {
		return err
	}
	log.Printf("[config] load %s\n", path)
	return nil
}

func Reload() error {
	return nil
}
