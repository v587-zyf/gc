package utils

import (
	"strconv"
	"time"
)

var Time1970 = time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)

func GetNowUTC() time.Time {
	return time.Now().UTC()
}

// GetYesterdayUTC
func GetYesterdayUTC() time.Time {
	now := time.Now().UTC()
	yesterday := now.AddDate(0, 0, -1)
	return yesterday
}

// 20240625
func GetYearMonthDay(t time.Time) int {
	date, _ := strconv.Atoi(t.Format("20060102"))
	return date
}

// 202447
func GetYearWeek(t time.Time) int {
	year, week := t.ISOWeek()
	date := year*100 + week
	return date
}

// 20247
func GetYearMonth(t time.Time) int {
	year, month, _ := t.Date()
	date := year*100 + int(month)
	return date
}

// IsLastDayOfMonth 检查当前时间是否是所在月份的最后一天
func IsLastDayOfMonth(t time.Time) bool {
	year, month, _ := t.Date()
	lastDayOfMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, t.Location())
	return t.YearDay() == lastDayOfMonth.YearDay()
}

// IsToday 检查给定的时间戳是否表示今天的日期
func IsToday(timestamp int64) bool {
	now := time.Now().Unix()
	t := time.Unix(timestamp, 0)
	today := time.Unix(now, 0).Truncate(24 * time.Hour)
	return t.Truncate(24 * time.Hour).Equal(today)
}

// IsYesterday 判断给定的时间是否是昨天
func IsYesterday(t time.Time) bool {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() &&
		t.Month() == yesterday.Month() &&
		t.Day() == yesterday.Day()
}

// 用字符串格式化时间
func GetTimeByData(dateStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)
}

// 获取指定时间到现在多少天
func GetTheDays(startTime time.Time) int {
	startTimeZeroTime := GetZeroTime(startTime).Unix()
	nowTime := time.Now()
	nowZeroTime := GetZeroTime(nowTime).Unix()
	openDays := (nowZeroTime-startTimeZeroTime)/(24*60*60) + 1
	return int(openDays)
}

// 获取0点时间
func GetZeroTime(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// 获取传入时间的零点时刻的时间戳
func GetZeroTimeInt64(t time.Time) int64 {
	ts, _ := time.Parse("2006-01-02", t.Format("2006-01-02"))
	return ts.Unix()
}

// 时间转换为int 如(20230101)
func GetDateInt(t time.Time) int {
	y, m, d := t.Date()
	date := y*10000 + int(m)*100 + d
	return date
}

// 获取两个时间差多少小时
func GetHourDiffer(startTime, endTime string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix()
		hour = diff / 3600
		return hour
	} else {
		return hour
	}
}
