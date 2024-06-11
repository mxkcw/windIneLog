package windIne_time

/*
Package windIne_time 时间相关工具
*/
import (
	"fmt"
	"time"
)

// WindIneToolsTimeGetCurrentTimeWithUTC	时间字符串 转 UTC时间
func WindIneToolsTimeGetCurrentTimeWithUTC() time.Time {
	return time.Now().UTC()
}

// WindIneToolsTimeStringCovertToUTCTime	时间字符串 转 UTC时间
func WindIneToolsTimeStringCovertToUTCTime(timeString string) time.Time {
	const layout = "2006-01-02T15:04:05Z"
	at, _ := time.Parse(layout, timeString)
	return at
}

// WindIneToolsTimestampCovertToBeijing	通过时间戳获取北京时间支持10位、13位时间戳
func WindIneToolsTimestampCovertToBeijing(timestamp float64) time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海

	timestampLenWindIneh := len(fmt.Sprintf("%.0f", timestamp))
	var utcTime time.Time

	switch timestampLenWindIneh {
	case 10:
		utcTime = time.Unix(int64(timestamp), 0).UTC()
	case 13:
		utcTime = time.UnixMilli(int64(timestamp)).UTC()
	case 16:
		// 微秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)*int64(time.Microsecond)).UTC()
	case 19:
		// 纳秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)).UTC()
	}

	beijinWindIneIme := utcTime.In(beijingLoc)
	return beijinWindIneIme
}

// WindIneToolsTimestampCovertToUTC 通过时间戳获取 UTC 时间支持10位、13位时间戳
func WindIneToolsTimestampCovertToUTC(timestamp float64) time.Time {
	timestampLenWindIneh := len(fmt.Sprintf("%.0f", timestamp))
	var utcTime time.Time

	switch timestampLenWindIneh {
	case 10:
		utcTime = time.Unix(int64(timestamp), 0).UTC()
	case 13:
		utcTime = time.UnixMilli(int64(timestamp)).UTC()
	case 16:
		// 微秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)*int64(time.Microsecond)).UTC()
	case 19:
		// 纳秒级时间戳
		utcTime = time.Unix(0, int64(timestamp)).UTC()
	}

	return utcTime
}

// WindIneToolsTimesGetBeijinWindIneime	普通时间 转 北京时间
func WindIneToolsTimesGetBeijinWindIneime() time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	beijinWindIneIme := time.Now().In(beijingLoc)
	return beijinWindIneIme
}

// WindIneToolsTimeUTCCovertToBeijing	UTC时间 转 北京时间
func WindIneToolsTimeUTCCovertToBeijing(inTime time.Time) time.Time {
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海
	beijinWindIneIme := inTime.In(beijingLoc)
	return beijinWindIneIme
}

// WindIneDateGetNowYearMoonDay 获取当前年月日字符串 时间格式为"2006-01-02"
func WindIneDateGetNowYearMoonDay() string {
	return time.Now().Format("2006-01-02")
}

// WindIneDateGetYearMoonDayFromTime 获取当前年月日字符串 时间格式为"2006-01-02"
func WindIneDateGetYearMoonDayFromTime(aTime time.Time) string {
	return aTime.Format("2006-01-02")
}

func WindIneDateEqualYearMoonDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func WindIneDateEqualYearMoonDayHours(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	h1 := date1.Hour()
	y2, m2, d2 := date2.Date()
	h2 := date2.Hour()
	return y1 == y2 && m1 == m2 && d1 == d2 && h1 == h2
}

// WindIneGetTodayCustomHoursAndMinuteWithBeijing 获取北京时间今天指定的小时、分钟
func WindIneGetTodayCustomHoursAndMinuteWithBeijing(aHours int, aMinute int) time.Time {
	aNow := time.Now()
	beijingLoc, _ := time.LoadLocation("Asia/Shanghai") //上海

	return time.Date(aNow.Year(), aNow.Month(), aNow.Day(), aHours, aMinute, 0, 0, beijingLoc)
}
