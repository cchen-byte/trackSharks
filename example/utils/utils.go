package utils

import (
	"log"
	"time"
)

/*
时间格式化
在PHP中, 可以使用"Y-m-d H:i:s"、 "yyyy-MM-dd HH:mm:ss"
在golang中为特定的数字 "2006-01-02 15:04:05" 是Go语言的创建时间，且必须为这几个准确的数字。
记忆: 1234567 1月2日3点4分5秒, 2006年, 7时区

ex: exampleUtils.StrTimeFormat("2022-06-28 15:00:01", "2006-01-02")
 */


// StrTimeFormat
//  @Description: 将字符串时间转为其他格式的字符串时间
//  @param timeValue	字符串时间
//  @param timeLayout	当前字符串时间的格式
//  @param timeFmt		希望转换的格式
//  @return string		转换后的字符串时间
func StrTimeFormat(timeValue string, timeLayout string, timeFmt string) string {
	// 时间格式化
	theTime, err := parseStrTime(timeValue, timeLayout)
	if err != nil {
		log.Printf("[utils][StrTimeFormat]parseStrTime error: %s\n", err)
		return ""
	}
	if theTime.IsZero() {
		log.Printf("[utils][StrTimeFormat]StrToTime: %s -> %s; error: zero\n", timeValue, timeFmt)
		return ""
	} else {
		// 日期转
		dataTimeStr := theTime.Format(timeFmt) //使用模板格式化为日期字符串
		return dataTimeStr
	}
}

// Str2TimeStamp
//  @Description: 将字符串时间转为时间戳 10位
//  @param timeValue	字符串时间
//  @param timeLayout	当前字符串时间的格式
//  @return int64		时间戳
func Str2TimeStamp(timeValue string, timeLayout string) int64 {
	theTime, err := parseStrTime(timeValue, timeLayout)
	if err != nil {
		log.Printf("[utils][Str2TimeStamp]parseStrTime error: %s\n", err)
		return 0
	}
	return theTime.Unix()
}

//
// TimeStamp2Str
//  @Description: 	 	时间戳转字符串时间
//  @param timeUnix	 	时间戳
//  @param timeFmt		希望转换的格式
//  @return string		转换后的字符串时间
//
func TimeStamp2Str(timeUnix int64, timeFmt string) string {
	return time.Unix(timeUnix, 0).Format(timeFmt)
}

//  parseStrTime
//  @Description: 		将时间字符串转成 time.Time 类型
//  @param timeValue	时间字符串
//  @param timeLayout	当前字符串时间的格式
//  @return time.Time
//  @return error
func parseStrTime(timeValue string, timeLayout string) (time.Time, error) {
	return time.ParseInLocation(timeLayout, timeValue, time.Local) // 使用模板在对应时区转化为time.time类型
}