package lib

import "time"

//字符串转时间戳
func String_Unix(str string) int {
	timelayo := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	times, _ := time.ParseInLocation(timelayo, str, loc)
	return int(times.Unix())
}

//时间戳转字符串
func Unix_time(unx int) string  {
	timelayo := "2006-01-02 15:04:05"
	dataTimeStr := time.Unix(int64(unx), 0).Format(timelayo)
	return  dataTimeStr
}
