package other

import (
	"strconv"
	"time"
)

//IntTime return int time
func IntTime() int {
	return int(time.Now().Unix())
}

//StrTime return str time
func StrTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10) //"1576395361"
}

//StrTime13 return str time13
func StrTime13() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10) //"1576395361"
}

//TI return \r[2006-01-02 15:04:05]
func TI() string {
	return time.Now().In(time.FixedZone("CST", 28800)).Format("\r[01-02 15:04:05] ")
}

//TodayDay return 0102
func TodayDay() string {
	return time.Now().In(time.FixedZone("CST", 28800)).Format("0102")
}

//BeforeDay return 0102
func BeforeDay(day int) string {
	return time.Now().AddDate(0, 0, -day).In(time.FixedZone("CST", 28800)).Format("0102")
}

//BeforeMinute return 0102
func BeforeMinute(minute int) string {
	return time.Now().Add(-time.Minute * time.Duration(minute)).In(time.FixedZone("CST", 28800)).Format("01020304")
}

//NowHour ...
func NowHour() int {
	return time.Now().In(time.FixedZone("CST", 28800)).Hour()
}
