package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/roada-go/roada"
	"github.com/shark/minigame-common/conf"
	"github.com/shark/minigame-common/model"
)

type ScoreRank_UpdateArgs struct {
	Userid int64
	Score  int64
}

type ScoreRank_UpdateReply struct {
}

type ScoreDayRank_UpdateArgs struct {
	Rankid int64
	Userid int64
	Score  int64
}

type ScoreDayRank_UpdateReply struct {
}

type LevelRank_UpdateArgs struct {
	Userid int64
	Score  int64
}

type LevelRank_UpdateReply struct {
}

type LevelDayRank_UpdateArgs struct {
	Rankid int64
	Userid int64
	Score  int64
}

type LevelDayRank_UpdateReply struct {
}

type MedalRank_UpdateArgs struct {
	Userid int64
	Score  int64
}

type MedalRank_UpdateReply struct {
}

type MedalDayRank_UpdateArgs struct {
	Rankid int64
	Userid int64
	Score  int64
}

type MedalDayRank_UpdateReply struct {
}

//分数总榜
type ScoreRank_RankArgs struct {
}

type ScoreRank_RankReply struct {
	Users []*model.User
}

//清理分数总榜
type ScoreRank_ClearArgs struct {
	Limit int64
}

type ScoreRank_ClearReply struct {
	RowsAffected int64
}

//分数日榜
type ScoreDayRank_RankArgs struct {
	Rankid int64
}

type ScoreDayRank_RankReply struct {
	Users []*model.User
}

//分数日榜成绩
type ScoreDayRank_ScoreArgs struct {
	Userid int64
	Rankid int64
}

type ScoreDayRank_ScoreReply struct {
	Score int64
}

//清理分数日榜
type ScoreDayRank_ClearArgs struct {
	Rankid int64
	Limit  int64
}

type ScoreDayRank_ClearReply struct {
	RowsAffected int64
}

//金牌总榜
type MedalRank_RankArgs struct {
}

type MedalRank_RankReply struct {
	Users []*model.User
}

//清理金牌总榜
type MedalRank_ClearArgs struct {
	Limit int64
}

type MedalRank_ClearReply struct {
	RowsAffected int64
}

//金牌日榜
type MedalDayRank_RankArgs struct {
	Rankid int64
}

type MedalDayRank_RankReply struct {
	Users []*model.User
}

//金牌日榜成绩
type MedalDayRank_ScoreArgs struct {
	Userid int64
	Rankid int64
}

type MedalDayRank_ScoreReply struct {
	Score int64
}

//清理金牌日榜
type MedalDayRank_ClearArgs struct {
	Rankid int64
	Limit  int64
}

type MedalDayRank_ClearReply struct {
	RowsAffected int64
}

//关卡总榜
type LevelRank_RankArgs struct {
}

type LevelRank_RankReply struct {
	Users []*model.User
}

//清理关卡总榜
type LevelRank_ClearArgs struct {
	Limit int64
}

type LevelRank_ClearReply struct {
	RowsAffected int64
}

//关卡日榜
type LevelDayRank_RankArgs struct {
	Rankid int64
}

type LevelDayRank_RankReply struct {
	Users []*model.User
}

//关卡日榜成绩
type LevelDayRank_ScoreArgs struct {
	Userid int64
	Rankid int64
}

type LevelDayRank_ScoreReply struct {
	Score int64
}

//清理关卡总榜
type LevelDayRank_ClearArgs struct {
	Rankid int64
	Limit  int64
}

type LevelDayRank_ClearReply struct {
	RowsAffected int64
}

func ScoreRank_Update(userid int64, score int64) error {
	var args = ScoreRank_UpdateArgs{Userid: userid, Score: score}
	var reply = ScoreRank_UpdateReply{}
	err := roada.Call("db", "ScoreRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) ScoreRank_Update(r *roada.Request, arg *ScoreRank_UpdateArgs, reply *ScoreRank_UpdateReply) error {
	now := time.Now().Unix()
	userid := arg.Userid
	score := arg.Score
	if score <= 0 {
		return nil
	}
	//总榜
	_, err := self.db.Exec(`UPDATE user SET score=?, updatetime=? WHERE userid=?`,
		score, now, userid)
	if err != nil {
		log.Printf("[db] ScoreRank_Update err %+v\n", err)
		return err
	}
	_, err = self.db.Exec(`INSERT INTO scorerank (userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		userid, score, now, now)
	if err != nil {
		log.Printf("[db] ScoreRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("scorerank@%s", conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] ScoreRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}
func ScoreDayRank_Update(rankid int64, userid int64, score int64) error {
	var args = ScoreDayRank_UpdateArgs{Rankid: rankid, Userid: userid, Score: score}
	var reply = ScoreDayRank_UpdateReply{}
	err := roada.Call("db", "ScoreDayRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) ScoreDayRank_Update(r *roada.Request, arg *ScoreDayRank_UpdateArgs, reply *ScoreDayRank_UpdateReply) error {
	now := time.Now().Unix()
	rankid := arg.Rankid
	userid := arg.Userid
	score := arg.Score
	if score <= 0 {
		return nil
	}
	//日榜
	_, err := self.db.Exec(`INSERT INTO scoredayrank (rankid, userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		rankid, userid, score, now, now)
	if err != nil {
		log.Printf("[db] ScoreRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("scoredayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] ScoreRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}
func LevelRank_Update(userid int64, score int64) error {
	var args = LevelRank_UpdateArgs{Userid: userid, Score: score}
	var reply = LevelRank_UpdateReply{}
	err := roada.Call("db", "LevelRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) LevelRank_Update(r *roada.Request, arg *LevelRank_UpdateArgs, reply *LevelRank_UpdateReply) error {
	now := time.Now().Unix()
	userid := arg.Userid
	score := arg.Score
	if score <= 0 {
		return nil
	}
	//总榜
	_, err := self.db.Exec(`UPDATE user SET level=?, updatetime=? WHERE userid=?`,
		score, now, userid)
	if err != nil {
		log.Printf("[db] LevelRank_Update err %+v\n", err)
		return err
	}
	_, err = self.db.Exec(`INSERT INTO levelrank (userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		userid, score, now, now)
	if err != nil {
		log.Printf("[db] LevelRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("levelrank@%s", conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] LevelRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}

func LevelDayRank_Update(rankid int64, userid int64, score int64) error {
	var args = LevelDayRank_UpdateArgs{Rankid: rankid, Userid: userid, Score: score}
	var reply = LevelDayRank_UpdateReply{}
	err := roada.Call("db", "LevelDayRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) LevelDayRank_Update(r *roada.Request, arg *LevelDayRank_UpdateArgs, reply *LevelDayRank_UpdateReply) error {
	now := time.Now().Unix()
	userid := arg.Userid
	score := arg.Score
	rankid := arg.Rankid
	if score <= 0 {
		return nil
	}
	//日榜
	_, err := self.db.Exec(`INSERT INTO leveldayrank (rankid, userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		rankid, userid, score, now, now)
	if err != nil {
		log.Printf("[db] LevelRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("leveldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] LevelDayRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}

func MedalRank_Update(userid int64, score int64) error {
	var args = MedalRank_UpdateArgs{Userid: userid, Score: score}
	var reply = MedalRank_UpdateReply{}
	err := roada.Call("db", "MedalRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) MedalRank_Update(r *roada.Request, arg *MedalRank_UpdateArgs, reply *MedalRank_UpdateReply) error {
	now := time.Now().Unix()
	userid := arg.Userid
	score := arg.Score
	if score <= 0 {
		return nil
	}
	//总榜
	_, err := self.db.Exec(`UPDATE user SET medal=?, updatetime=? WHERE userid=?`,
		score, now, userid)
	if err != nil {
		log.Printf("[db] MedalRank_Update err %+v\n", err)
		return err
	}
	_, err = self.db.Exec(`INSERT INTO medalrank (userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		userid, score, now, now)
	if err != nil {
		log.Printf("[db] MedalRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("medalrank@%s", conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] MedalRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}

func MedalDayRank_Update(rankid int64, userid int64, score int64) error {
	var args = MedalDayRank_UpdateArgs{Rankid: rankid, Userid: userid, Score: score}
	var reply = MedalDayRank_UpdateReply{}
	err := roada.Call("db", "MedalDayRank_Update", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) MedalDayRank_Update(r *roada.Request, arg *MedalDayRank_UpdateArgs, reply *MedalDayRank_UpdateReply) error {
	now := time.Now().Unix()
	userid := arg.Userid
	score := arg.Score
	rankid := arg.Rankid
	if score <= 0 {
		return nil
	}
	//日榜
	_, err := self.db.Exec(`INSERT INTO medaldayrank (rankid, userid, score, createtime, updatetime) 
		VALUES(?, ?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		score=VALUES(score), updatetime=VALUES(updatetime)`,
		rankid, userid, score, now, now)
	if err != nil {
		log.Printf("[db] MedalRank_Update err %+v\n", err)
		return err
	}
	rkey := fmt.Sprintf("medaldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.ZAdd(rkey, self.encodeScore(userid, score), fmt.Sprintf("%d", userid)); err != nil {
		log.Printf("[db] MedalRank_Update cache.ZAdd failed, err:%s", err.Error())
		return err
	}
	return nil
}

func (self *DbService) encodeScoreAndTime(score int64, t int64) int64 {
	return score*1000000000 + t
}

func (self *DbService) encodeScore(userid int64, score int64) int64 {
	now := time.Now().Unix() - 1625109340
	return score*1000000000 + now
}

func (self *DbService) decodeScore(value int64) int64 {
	//score, _ := strconv.ParseInt(value, 10, 64)
	score := value / 1000000000
	return score
}

//分数总榜
func ScoreRank_Rank() ([]*model.User, error) {
	var args = ScoreRank_RankArgs{}
	var reply = ScoreRank_RankReply{}
	err := roada.Call("db", "ScoreRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//分数总榜
func (self *DbService) ScoreRank_Rank(r *roada.Request, args *ScoreRank_RankArgs, reply *ScoreRank_RankReply) error {
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score 
		FROM scorerank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`)
	if err != nil {
		return err
	}
	return nil
}

//分数日榜
func ScoreDayRank_Rank(rankid int64) ([]*model.User, error) {
	var args = ScoreDayRank_RankArgs{Rankid: rankid}
	var reply = ScoreDayRank_RankReply{}
	err := roada.Call("db", "ScoreDayRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//分数日榜
func (self *DbService) ScoreDayRank_Rank(r *roada.Request, args *ScoreDayRank_RankArgs, reply *ScoreDayRank_RankReply) error {
	rankid := args.Rankid
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score 
		FROM scoredayrank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.rankid=? AND a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`, rankid)
	if err != nil {
		return err
	}
	return nil
}

//分数日榜成绩
func ScoreDayRank_Score(userid int64, rankid int64) (int64, error) {
	var args = ScoreDayRank_ScoreArgs{Rankid: rankid, Userid: userid}
	var reply = ScoreDayRank_ScoreReply{}
	err := roada.Call("db", "ScoreDayRank_Score", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.Score, nil
}

//分数日榜成绩
func (self *DbService) ScoreDayRank_Score(r *roada.Request, args *ScoreDayRank_ScoreArgs, reply *ScoreDayRank_ScoreReply) error {
	rankid := args.Rankid
	userid := args.Userid
	user := model.User{}
	rkey := fmt.Sprintf("scoredayrank_%d@%s", rankid, conf.Ini.Game.Name)
	scoref, err := self.cache.ZScore(rkey, fmt.Sprintf("%d", userid))
	if err == nil {
		score := int64(scoref)
		score = self.decodeScore(score)
		reply.Score = score
		return nil
	}
	err = self.db.Get(&user, "SELECT userid, score FROM scoredayrank WHERE userid=? AND rankid=?", userid, rankid)
	if err == sql.ErrNoRows {
		reply.Score = 0
	} else if err != nil {
		return err
	}
	reply.Score = user.Score
	return nil
}

//金牌总榜
func MedalRank_Rank() ([]*model.User, error) {
	var args = MedalRank_RankArgs{}
	var reply = MedalRank_RankReply{}
	err := roada.Call("db", "MedalRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//金牌总榜
func (self *DbService) MedalRank_Rank(r *roada.Request, args *MedalRank_RankArgs, reply *MedalRank_RankReply) error {
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score
		FROM medalrank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`)
	if err != nil {
		return err
	}
	return nil
}

//金牌日榜
func MedalDayRank_Rank(rankid int64) ([]*model.User, error) {
	var args = MedalDayRank_RankArgs{Rankid: rankid}
	var reply = MedalDayRank_RankReply{}
	err := roada.Call("db", "MedalDayRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//金牌日榜
func (self *DbService) MedalDayRank_Rank(r *roada.Request, args *MedalDayRank_RankArgs, reply *MedalDayRank_RankReply) error {
	rankid := args.Rankid
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score
		FROM medaldayrank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.rankid=? AND a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`, rankid)
	if err != nil {
		return err
	}
	return nil
}

//金牌日榜成绩
func MedalDayRank_Score(userid int64, rankid int64) (int64, error) {
	var args = MedalDayRank_ScoreArgs{Rankid: rankid, Userid: userid}
	var reply = MedalDayRank_ScoreReply{}
	err := roada.Call("db", "MedalDayRank_Score", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.Score, nil
}

//金牌日榜成绩
func (self *DbService) MedalDayRank_Score(r *roada.Request, args *MedalDayRank_ScoreArgs, reply *MedalDayRank_ScoreReply) error {
	rankid := args.Rankid
	userid := args.Userid
	rkey := fmt.Sprintf("medaldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	scoref, err := self.cache.ZScore(rkey, fmt.Sprintf("%d", userid))
	if err == nil {
		score := int64(scoref)
		score = self.decodeScore(score)
		reply.Score = score
		return nil
	}
	user := model.User{}
	err = self.db.Get(&user, "SELECT userid, score FROM medaldayrank WHERE userid=? AND rankid=?", userid, rankid)
	if err == sql.ErrNoRows {
		reply.Score = 0
	} else if err != nil {
		return err
	}
	reply.Score = user.Score
	return nil
}

//关卡总榜
func LevelRank_Rank() ([]*model.User, error) {
	var args = LevelRank_RankArgs{}
	var reply = LevelRank_RankReply{}
	err := roada.Call("db", "LevelRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//关卡总榜
func (self *DbService) LevelRank_Rank(r *roada.Request, args *LevelRank_RankArgs, reply *LevelRank_RankReply) error {
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score
		FROM levelrank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`)
	if err != nil {
		return err
	}
	return nil
}

//关卡日榜
func LevelDayRank_Rank(rankid int64) ([]*model.User, error) {
	var args = LevelDayRank_RankArgs{Rankid: rankid}
	var reply = LevelDayRank_RankReply{}
	err := roada.Call("db", "LevelDayRank_Rank", &args, &reply)
	if err != nil {
		return nil, err
	}
	return reply.Users, nil
}

//关卡日榜
func (self *DbService) LevelDayRank_Rank(r *roada.Request, args *LevelDayRank_RankArgs, reply *LevelDayRank_RankReply) error {
	rankid := args.Rankid
	reply.Users = make([]*model.User, 0)
	err := self.db.Select(&reply.Users, `SELECT b.openid, b.nickname, b.avatar, a.score
		FROM leveldayrank AS a 
		LEFT JOIN user AS b ON a.userid=b.userid 
		WHERE a.rankid=? AND a.score > 0
		ORDER BY a.score desc, a.updatetime desc LIMIT 50`, rankid)
	if err != nil {
		return err
	}
	return nil
}

//关卡日榜成绩
func LevelDayRank_Score(userid int64, rankid int64) (int64, error) {
	var args = LevelDayRank_ScoreArgs{Rankid: rankid, Userid: userid}
	var reply = LevelDayRank_ScoreReply{}
	err := roada.Call("db", "LevelDayRank_Score", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.Score, nil
}

//关卡日榜成绩
func (self *DbService) LevelDayRank_Score(r *roada.Request, args *LevelDayRank_ScoreArgs, reply *LevelDayRank_ScoreReply) error {
	rankid := args.Rankid
	userid := args.Userid
	//先读缓存
	rkey := fmt.Sprintf("leveldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	scoref, err := self.cache.ZScore(rkey, fmt.Sprintf("%d", userid))
	if err == nil {
		score := int64(scoref)
		score = self.decodeScore(score)
		reply.Score = score
		return nil
	}
	//再读数据库
	user := model.User{}
	err = self.db.Get(&user, "SELECT userid, score FROM leveldayrank WHERE userid=? AND rankid=?", userid, rankid)
	if err == sql.ErrNoRows {
		reply.Score = 0
	} else if err != nil {
		return err
	}
	reply.Score = user.Score
	return nil
}

//清理分数总榜
func ScoreRank_Clear(limit int64) (int64, error) {
	var args = ScoreRank_ClearArgs{
		Limit: limit,
	}
	var reply = ScoreRank_ClearReply{}
	err := roada.Call("db", "ScoreRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理分数总榜
func (self *DbService) ScoreRank_Clear(r *roada.Request, args *ScoreRank_ClearArgs, reply *ScoreRank_ClearReply) error {
	user := model.User{}
	err := self.db.Get(&user, "SELECT userid, score FROM scorerank ORDER BY score DESC, updatetime DESC LIMIT 100, 1")
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	result, err := self.db.Exec(`DELETE FROM scorerank WHERE score<? LIMIT ?`, user.Score, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

//清理金牌总榜
func MedalRank_Clear(limit int64) (int64, error) {
	var args = MedalRank_ClearArgs{
		Limit: limit,
	}
	var reply = MedalRank_ClearReply{}
	err := roada.Call("db", "MedalRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理金牌总榜
func (self *DbService) MedalRank_Clear(r *roada.Request, args *MedalRank_ClearArgs, reply *MedalRank_ClearReply) error {
	user := model.User{}
	err := self.db.Get(&user, "SELECT userid, score FROM medalrank ORDER BY score DESC, updatetime DESC LIMIT 100, 1")
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	result, err := self.db.Exec(`DELETE FROM medalrank WHERE score<? LIMIT ?`, user.Score, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

//清理关卡总榜
func LevelRank_Clear(limit int64) (int64, error) {
	var args = LevelRank_ClearArgs{
		Limit: limit,
	}
	var reply = LevelRank_ClearReply{}
	err := roada.Call("db", "LevelRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理关卡总榜
func (self *DbService) LevelRank_Clear(r *roada.Request, args *LevelRank_ClearArgs, reply *LevelRank_ClearReply) error {
	user := model.User{}
	err := self.db.Get(&user, "SELECT userid, score FROM levelrank ORDER BY score DESC, updatetime DESC LIMIT 100, 1")
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	result, err := self.db.Exec(`DELETE FROM levelrank WHERE score<? LIMIT ?`, user.Score, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

//清理分数日榜
func ScoreDayRank_Clear(rankid int64, limit int64) (int64, error) {
	var args = ScoreDayRank_ClearArgs{
		Rankid: rankid,
		Limit:  limit,
	}
	var reply = ScoreDayRank_ClearReply{}
	err := roada.Call("db", "ScoreDayRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理分数日榜
func (self *DbService) ScoreDayRank_Clear(r *roada.Request, args *ScoreDayRank_ClearArgs, reply *ScoreDayRank_ClearReply) error {
	result, err := self.db.Exec(`DELETE FROM scoredayrank WHERE rankid<? LIMIT ?`, args.Rankid, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

//清理金牌日榜
func MedalDayRank_Clear(rankid int64, limit int64) (int64, error) {
	var args = MedalDayRank_ClearArgs{
		Rankid: rankid,
		Limit:  limit,
	}
	var reply = MedalDayRank_ClearReply{}
	err := roada.Call("db", "MedalDayRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理金牌日榜
func (self *DbService) MedalDayRank_Clear(r *roada.Request, args *MedalDayRank_ClearArgs, reply *MedalDayRank_ClearReply) error {
	result, err := self.db.Exec(`DELETE FROM medaldayrank WHERE rankid<? LIMIT ?`, args.Rankid, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

//清理等级日榜
func LevelDayRank_Clear(rankid int64, limit int64) (int64, error) {
	var args = LevelDayRank_ClearArgs{
		Rankid: rankid,
		Limit:  limit,
	}
	var reply = LevelDayRank_ClearReply{}
	err := roada.Call("db", "LevelDayRank_Clear", &args, &reply)
	if err != nil {
		return 0, err
	}
	return reply.RowsAffected, nil
}

//清理等级日榜
func (self *DbService) LevelDayRank_Clear(r *roada.Request, args *LevelDayRank_ClearArgs, reply *LevelDayRank_ClearReply) error {
	result, err := self.db.Exec(`DELETE FROM leveldayrank WHERE rankid<? LIMIT ?`, args.Rankid, args.Limit)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	reply.RowsAffected = rowsAffected
	return nil
}

func ScoreDayRank_ClearCache(rankid int64) error {
	var args = rankid
	var reply int64 = 0
	err := roada.Call("db", "ScoreDayRank_ClearCache", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (self *DbService) ScoreDayRank_ClearCache(r *roada.Request, args *int64, reply *int64) error {
	rankid := *args
	rkey := fmt.Sprintf("scoredayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.Del(rkey); err != nil {
		log.Printf("[db] ScoreDayRank_ClearCache cache.Del failed, err:%s", err.Error())
		return err
	}
	return nil
}

func MedalDayRank_ClearCache(rankid int64) error {
	var args = rankid
	var reply int64 = 0
	err := roada.Call("db", "MedalDayRank_ClearCache", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}
func (self *DbService) MedalDayRank_ClearCache(r *roada.Request, args *int64, reply *int64) error {
	rankid := *args
	rkey := fmt.Sprintf("medaldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.Del(rkey); err != nil {
		log.Printf("[db] MedalDayRank_ClearCache cache.Del failed, err:%s", err.Error())
		return err
	}
	return nil
}

func LevelDayRank_ClearCache(rankid int64) error {
	var args = rankid
	var reply int64 = 0
	err := roada.Call("db", "LevelDayRank_ClearCache", &args, &reply)
	if err != nil {
		return err
	}
	return nil
}
func (self *DbService) LevelDayRank_ClearCache(r *roada.Request, args *int64, reply *int64) error {
	rankid := *args
	rkey := fmt.Sprintf("leveldayrank_%d@%s", rankid, conf.Ini.Game.Name)
	if err := self.cache.Del(rkey); err != nil {
		log.Printf("[db] LevelDayRank_ClearCache cache.Del failed, err:%s", err.Error())
		return err
	}
	return nil
}
