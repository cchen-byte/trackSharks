package utils

import (
	"errors"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"regexp"
	"strings"
	"time"
)

func partition(trackInfo []*httpobj.TrackInfo, low, high int) int {
	pivotValue := trackInfo[low]
	pivot := Str2TimeStamp(pivotValue.Date, "2006-01-02 15:04:05") //导致 low 位置值为空
	for low < high {
		//high指针值 >= pivot high指针左移
		for low < high && pivot <= Str2TimeStamp(trackInfo[high].Date, "2006-01-02 15:04:05") {
			high--
		}
		//填补low位置空值
		//high指针值 < pivot high值 移到low位置
		//high 位置值空
		trackInfo[low] = trackInfo[high]
		//low指针值 <= pivot low指针右移
		for low < high && pivot >= Str2TimeStamp(trackInfo[low].Date, "2006-01-02 15:04:05") {
			low++
		}
		//填补high位置空值
		//low指针值 > pivot low值 移到high位置
		//low位置值空
		trackInfo[high] = trackInfo[low]
	}
	//pivot 填补 low位置的空值
	trackInfo[low] = pivotValue
	return low
}

// 按照轨迹信息时间降序排序
func sortTrackInfo(trackInfo []*httpobj.TrackInfo,low,high int)  {
	if high > low{
		//位置划分
		pivot := partition(trackInfo,low,high)
		//左边部分排序
		sortTrackInfo(trackInfo,low,pivot-1)
		//右边排序
		sortTrackInfo(trackInfo,pivot+1,high)
	}
}

/*
处理物流状态和按时间排序，是否重查，添加查询状态提示

 */

func DealTrackInfoStatusAndLastInfoByFuc(returnData *httpobj.TrackItem, webRoute []string, normal bool, function []string) {
	// 按查询时间降序排序
	sortTrackInfo(returnData.TrackInfo, 0, len(returnData.TrackInfo)-1)

	// 队列中 是否重查
	hasRouteCanTrack := false
	for _, route := range webRoute{
		if getSwitchBool(route) {
			hasRouteCanTrack = true
			break
		}
	}
	// TrackInfo 为空 && 查询过程不正常 && 有其他线路
	if len(returnData.TrackInfo) == 0 && normal == false && hasRouteCanTrack {
		returnData.NeedReTrack = true
	}

	// function 拼接成字符串
	functionStr := strings.Join(function, " ,")
	if len(returnData.Function) == 0 {
		returnData.Function = functionStr
	}else{
		returnData.Function += " ," + functionStr
	}

	// 设置语言
	if len(returnData.Lang) == 0{
		returnData.Lang = "en"
	}

	// 最新轨迹数据
	lastestData := returnData.TrackInfo[0]
	// 提取最新事件拼接 StatusDescription,Details,Date
	returnData.StatusInfo = strings.Join([]string{lastestData.StatusDescription, lastestData.Details, lastestData.Date}, ",")
	// 最新更新时间
	returnData.LastUpdateTime = lastestData.Date
}

func getSwitchBool(route string) bool {
	return true
}

// PassItem
// @Description: 过滤item
// @param date						处理完成的date
// @param statusDescription			处理完成的statusDescription
// @param detail					处理完成的detail
// @return map[string]interface{}	过滤后返回item
// @return error
func PassItem(date, statusDescription, detail string) (*httpobj.TrackInfo, error) {
	item := &httpobj.TrackInfo{}

	// 序列化数组 || 保持原本数据结构
	//statusDescriptionByte, _ := json.Marshal(statusDescription)
	//statusDescription = string(statusDescriptionByte)

	//detailByte, _ := json.Marshal(detail)
	//detail = string(detailByte)

	item.StatusDescription = replaceFunc(strings.Trim(statusDescription, ""))
	item.Date = replaceFunc(strings.Trim(date, ""))
	item.Details = replaceFunc(strings.Trim(detail, ""))

	// 判断数据时间是否正确
	if len(item.Date) == 0 || strings.Contains(item.Date, "1970-01-01") {
		return nil, errors.New("parseItem date error")
	}
	// 将当前时间和数据时间转成时间戳
	dateTimeStamp := Str2TimeStamp(item.Date, "2006-01-02 15:04:05")
	timeStamp := time.Now().Unix() + 3600*24

	if dateTimeStamp < 0 {
		return nil, errors.New("dateTimeStamp < 0")
	}else if dateTimeStamp >= timeStamp{
		// 减一年
		dateTime1 := dateTimeStamp - 3600*24*365
		// 减一个月
		dateTime2 := dateTimeStamp - 3600*24*30
		if dateTime2 <= timeStamp{
			date = TimeStamp2Str(dateTime2, "2006-01-02 15:04:05")
		}else if dateTime1 <= timeStamp {
			date = TimeStamp2Str(dateTime1, "2006-01-02 15:04:05")
		}else{
			// 时间差别有点大，只能不返回数据
			return nil, errors.New("dateTime error")
		}
		//dateByte, _ = json.Marshal(date)
		//date = string(dateByte)
		item.Date = replaceFunc(strings.Trim(date, ""))
	}

	// 状态描述为空 && 详细为空
	if len(item.StatusDescription) == 0 && len(item.Details) == 0 {
		return nil, errors.New("StatusDescription and Details empty")
	}
	return item, nil
}

// 正则替换字符
// 将 \[ \] \r\n\t 删除
func replaceFunc(str string) string {
	reg := regexp.MustCompile(`[\[\]\r\n\t]`)
	str = reg.ReplaceAllString(str, "")
	return str
}