package utils

import (
	"time"
)

var loc *time.Location

func init()  {
	loc = time.FixedZone("CST", 8*3600)
}

// 时间格式化 DataTime("2006-01-02 15:04:05",时间戳)
func Date(format string,timestamp int64) string {
	return time.Unix(timestamp, 0).In(loc).Format(format)
}

// 获取当前时间戳和当天开始和结束的时间戳
func GetTodayTime() (t,ts,te int64){
	date := time.Now().In(loc)
	t = GetTime() // 当前时间戳
	ts = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).Unix() // 当日开始时间戳
	te = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, date.Location()).Unix() // 当日结束时间戳
	return
}

// 获取当前时间戳
func GetTime() int64 {
	return time.Now().In(loc).Unix()
}

// 获取当前时间
func GetNow() time.Time {
	return time.Now().In(loc)
}

// 时间转时间戳。传入字符串时间
func StrToTime(timeStr string) int64 {
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)//转换为时间戳
	return tt.Unix()
}

// 时间戳转时间。传入int64时间戳
func TimeToStr(timeInt int64) string{
	return time.Unix(timeInt, 0).In(loc).Format("2006-01-02 15:04:05")
}

// 检查当前时间是否在指定时间段内。开始时间戳，结束时间戳，是否需要判断天
func IsInTime(ts,te int,day bool) bool {
	times := time.Unix(int64(ts), 0).In(loc)  // 开始时间
	timed := time.Unix(int64(te), 0).In(loc) // 结束时间
	result := true
	if time.Now().In(loc).Hour() >= times.Hour() && time.Now().In(loc).Hour() <= timed.Hour() { // 指定时间范围(在当前小时内)条件满足
		//开始前的分钟还没到
		if time.Now().In(loc).Hour() == times.Hour() && time.Now().In(loc).Minute() <= times.Minute() { // 分钟还未到
			result = false
		}
		//结束前的分钟还没到
		if time.Now().In(loc).Hour() == timed.Hour() && time.Now().In(loc).Minute() >= timed.Minute() { // 分钟超过了
			result = false
		}
		//在时间段内，上面的两个条件满足
		if result {
			if day {
				return times.Year() == time.Now().In(loc).Year() && times.Month() == time.Now().In(loc).Month() && times.Day() == time.Now().In(loc).Day()
			}
			return true
		}
	}
	return false
}