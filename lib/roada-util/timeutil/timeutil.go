package timeutil

import (
	"time"
)

func IsSameDay(s1 int64, s2 int64) bool {
	time1 := time.Unix(s1, 0)
	time2 := time.Unix(s2, 0)
	if time1.Year()*1000+int(time1.YearDay()) == time2.Year()*1000+int(time2.YearDay()) {
		return true
	}
	return false
}

func IsSameWeek(s1 int64, s2 int64) bool {
	time1 := time.Unix(s1, 0)
	time2 := time.Unix(s2, 0)
	y1, m1 := time1.ISOWeek()
	y2, m2 := time2.ISOWeek()
	if y1*1000+m1 == y2*1000+m2 {
		return true
	}
	return false
}

func PassDays(before int64, now int64) int {
	time1 := time.Unix(before, 0)
	d1 := time.Date(time1.Year(), time1.Month(), time1.Day(), 0, 0, 0, 0, time.Local)

	time2 := time.Unix(now, 0)
	d2 := time.Date(time2.Year(), time2.Month(), time2.Day(), 0, 0, 0, 0, time.Local)
	pass := int(d2.Unix()-d1.Unix()) / 86400
	//log.Println("ggggggggg", pass)
	return pass
}
