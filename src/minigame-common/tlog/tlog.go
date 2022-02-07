package tlog

import (
	"fmt"
	"strings"
	"time"

	"github.com/roada-go/roada"
	"github.com/roada-go/util/tlogd"
	"github.com/shark/minigame-common/conf"
)

var road *roada.Road

func Register(r *roada.Road) error {
	config := tlogd.Config{
		Dir:       conf.Ini.Tlog.Dir,
		BackupDir: conf.Ini.Tlog.BackupDir,
		TimeLimit: conf.Ini.Tlog.TimeLimit,
		LineLimit: conf.Ini.Tlog.LineLimit,
		TcpAddr:   conf.Ini.Tlog.TcpAddr,
		Console:   conf.Ini.Tlog.Console,
	}
	road = r
	return tlogd.Register(road, &config)
}

func UserLogin(openid string, userid int64, loginTime int64) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"user_login",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", loginTime),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func UserLogout(openid string, userid int64, loginTime int64, logoutTime int64) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"user_logout",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", loginTime),
		fmt.Sprintf("%d", logoutTime),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func UserRegister(userid int64, openid string, nickname string, ip string, registerTime int64) {
	if road == nil {
		return
	}
	nickname = strings.Replace(nickname, "|", "", 1)
	openid = strings.Replace(openid, "|", "", 1)
	//var reply int
	args := []string{
		"user_register",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%s", nickname),
		fmt.Sprintf("%s", ip),
		fmt.Sprintf("%d", registerTime),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func RoundStart(openid string, userid int64) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"round_start",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func RoundResult(openid string, userid int64, roundStartTime int64, roundResultTime int64, score int64, medal int64) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"round_result",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", roundStartTime),
		fmt.Sprintf("%d", roundResultTime),
		fmt.Sprintf("%d", score),
		fmt.Sprintf("%d", medal),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func RoundShare(openid string, userid int64) {
	if road == nil {
		return
	}
	args := []string{
		"round_share",
		"2",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func GateStat(onlineNum int32) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"gate_stat",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%d", onlineNum),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func AddDiamond(openid string, userid int64, diamond int64, reason string) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"diamond_add",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", diamond),
		reason,
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func DecDiamond(openid string, userid int64, diamond int64, reason string) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"diamond_dec",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", diamond),
		reason,
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func DaySign(openid string, userid int64, day int64) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"daysign",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", day),
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func AddCoin(openid string, userid int64, coin int64, reason string) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"coin_add",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", coin),
		reason,
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func DecCoin(openid string, userid int64, coin int64, reason string) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"coin_dec",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", coin),
		reason,
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}

func DecGold(openid string, userid int64, gold int64, reason string) {
	if road == nil {
		return
	}
	//var reply int
	args := []string{
		"gold_dec",
		"1",
		fmt.Sprintf("%d", time.Now().Unix()),
		fmt.Sprintf("%d", conf.Ini.Game.Id),
		fmt.Sprintf("%s", openid),
		fmt.Sprintf("%d", userid),
		fmt.Sprintf("%d", gold),
		reason,
	}
	road.Post("tlog", "Print", strings.Join(args, "|"))
}
