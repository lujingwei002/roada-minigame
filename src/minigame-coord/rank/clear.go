package rank

import (
	"log"
	"time"

	"github.com/shark/minigame-common/db"
)

func schedule_clear() {
	for {
		now := time.Now()
		daySec := now.Hour()*3600 + now.Minute()*60 + now.Second()
		remainSec := 0
		if daySec < 10800 {
			remainSec = 10800 - daySec
		} else {
			remainSec = 86400 - daySec + 10800
		}
		log.Printf("[rank] wait %d second to clear\n", remainSec)
		time.Sleep(time.Duration(remainSec) * time.Second)
		clearScoreRank()
		clearMedalRank()
		clearLevelRank()
		clearScoreDayRank()
		clearMedalDayRank()
		clearLevelDayRank()
	}
}

func clearScoreRank() {
	var total int64 = 0
	for {
		rowsAffected, err := db.ScoreRank_Clear(1000)
		if err != nil {
			log.Printf("[rank] db.ScoreRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	log.Printf("[rank] clearScoreRank finished, total=%d\n", total)
}

func clearMedalRank() {
	var total int64 = 0
	for {
		rowsAffected, err := db.MedalRank_Clear(1000)
		if err != nil {
			log.Printf("[rank] db.MedalRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	log.Printf("[rank] clearMedalRank finished, total=%d\n", total)
}

func clearLevelRank() {
	var total int64 = 0
	for {
		rowsAffected, err := db.LevelRank_Clear(1000)
		if err != nil {
			log.Printf("[rank] db.LevelRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	log.Printf("[rank] clearLevelRank finished, total=%d\n", total)
}

func clearScoreDayRank() {
	var total int64 = 0
	now := time.Now().AddDate(0, 0, -3)
	rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	for {
		rowsAffected, err := db.ScoreDayRank_Clear(rankid, 1000)
		if err != nil {
			log.Printf("[rank] db.ScoreDayRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	for i := 0; i < 3; i++ {
		now := time.Now().AddDate(0, 0, -3-i)
		rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
		err := db.ScoreDayRank_ClearCache(rankid)
		if err != nil {
			log.Printf("[rank] db.ScoreDayRank_ClearCache failed, error=%s\n", err.Error())
		}
	}
	log.Printf("[rank] clearScoreDayRank finished, rankid=%d, total=%d\n", rankid, total)
}

func clearMedalDayRank() {
	var total int64 = 0
	now := time.Now().AddDate(0, 0, -3)
	rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	for {
		rowsAffected, err := db.MedalDayRank_Clear(rankid, 1000)
		if err != nil {
			log.Printf("[rank] db.MedalDayRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	for i := 0; i < 3; i++ {
		now := time.Now().AddDate(0, 0, -3-i)
		rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
		err := db.MedalDayRank_ClearCache(rankid)
		if err != nil {
			log.Printf("[rank] db.MedalDayRank_ClearCache failed, error=%s\n", err.Error())
		}
	}
	log.Printf("[rank] clearMedalDayRank finished, rankid=%d, total=%d\n", rankid, total)
}

func clearLevelDayRank() {
	var total int64 = 0
	now := time.Now().AddDate(0, 0, -3)
	rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	for {
		rowsAffected, err := db.LevelDayRank_Clear(rankid, 1000)
		if err != nil {
			log.Printf("[rank] db.LevelDayRank_Clear failed, error=%s\n", err.Error())
		}
		if rowsAffected <= 0 {
			break
		}
		total = total + rowsAffected
	}
	for i := 0; i < 3; i++ {
		now := time.Now().AddDate(0, 0, -3-i)
		rankid := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
		err := db.LevelDayRank_ClearCache(rankid)
		if err != nil {
			log.Printf("[rank] db.LevelDayRank_ClearCache failed, error=%s\n", err.Error())
		}
	}
	log.Printf("[rank] clearLevelDayRank finished, rankid=%d, total=%d\n", rankid, total)
}
