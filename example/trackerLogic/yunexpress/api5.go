package yunexpress

import (
	"fmt"
	exampleUtils "github.com/cchen-byte/trackeSharkes/example/utils"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/utils"
	"log"
	"strconv"
	"strings"
)

func api5ConstructFirstRequest(reqMeta *httpobj.RequestMeta, trackData *httpobj.TrackData) *httpobj.Request{
	var trackNumberQ []string
	for _, v := range *trackData{
		trackNumberQ = append(trackNumberQ, v.TrackNumber)
	}

	return &httpobj.Request{
		Url: "https://trackapi.yunexpress.com/LMS.API/api/waybill/gettracklist",
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"Authorization": "Basic QzYyMjQ5Jjd4aTFSRXJPWEg0PQ==",
		},
		Json:  trackNumberQ,
		MetaData: &httpobj.MetaData{
			"trackData": trackData,
		},
		Callback: parseApi5TrackData,
		RequestMeta: reqMeta,
	}
}

func parseApi5TrackData(response *httpobj.Response) (*httpobj.ParseResult, error) {
	result := httpobj.NewParseResult()
	response.Request.RequestMeta.Function = append(response.Request.RequestMeta.Function, fmt.Sprintf("yunexpress_api_5_%d", response.StatusCode))

	// 异常判断
	if strings.Contains(response.Text, "item"){
		// 无其他线路节点
		if !response.Request.RequestMeta.HasNextConstructorNode {
			returnData := &httpobj.TrackItem{}
			returnData.NeedReTrack = true
			returnData.Function = strings.Join(response.Request.RequestMeta.Function, " ,")
			result.AppendItem(&httpobj.Item{
				ItemStatus: &httpobj.ItemStatus{
					RequestId: response.Request.RequestMeta.RequestId,
				},
				Item: returnData,
			})
			return result, nil

		}else{
			// 有其他节点
			returnData := &httpobj.TrackItem{}
			result.AppendItem(&httpobj.Item{
				ItemStatus: &httpobj.ItemStatus{
					RequestId: response.Request.RequestMeta.RequestId,
					IsError: true,
				},
				Item: returnData,
			})
			return result, nil
		}
	}

	// 正常逻辑
	metaData := *(response.Request.MetaData)
	trackData := metaData["trackData"].(*httpobj.TrackData)

	respJsonDom, _ := response.GetJsonDom()
	respItemData := respJsonDom.Xpath("//Item/*")
	for _, itemData := range respItemData{
		yunStatus, _ := strconv.Atoi(itemData.XpathOne("//PackageState").InnerText())

		trackNumberRow1 := itemData.XpathOne("//WayBillNumber").InnerText()
		trackNumberRow2 := itemData.XpathOne("//TrackingNumber").InnerText()
		trackNumberRow3 := itemData.XpathOne("//CustomerOrderNumber").InnerText()
		carrierName := itemData.XpathOne("//CarrierName").InnerText()

		var returnDataTrackInfo []*httpobj.TrackInfo
		returnData := httpobj.NewTrackItem()

		returnData.CountryName = itemData.XpathOne("//OriginCountryCode").InnerText()
		returnData.DestinationCountry = itemData.XpathOne("//CountryCode").InnerText()
		returnData.TrackInfo = returnDataTrackInfo

		trackNumberSet := utils.NewMemorySet()
		for _, v := range *trackData{
			_, _ = trackNumberSet.Add(v.TrackNumber)
		}
		// 判断单号是否在单号列表中（三种单号格式）
		if  (len(trackNumberRow1) != 0 && !trackNumberSet.HasItem(trackNumberRow1)) &&
			(len(trackNumberRow2) != 0 && !trackNumberSet.HasItem(trackNumberRow2)) &&
			(len(trackNumberRow3) != 0 && !trackNumberSet.HasItem(trackNumberRow3)) {
			log.Printf("parse error: the tracking number %s, %s, %s is not in the response\n", trackNumberRow1, trackNumberRow2, trackNumberRow3)
			continue
		}

		itemDataOrderTrackingDetails := itemData.Xpath("//OrderTrackingDetails/*")
		for index, ItemTrackData := range itemDataOrderTrackingDetails{
			itemProcessDate := ItemTrackData.XpathOne("//ProcessDate").InnerText()
			date := exampleUtils.StrTimeFormat(itemProcessDate, "2006-01-02T15:04:05", "2006-01-02 15:04:05")

			event := ItemTrackData.XpathOne("//ProcessContent").InnerText()
			// 轨迹索引为0 && 存在最后一公里单号 -> 添加最后一公里单号及运输商
			if index ==0 && len(trackNumberRow2) != 0 {
				event = fmt.Sprintf("%s Last mile tracking number:%s - Last mile tracking carrier:%s", event, trackNumberRow2, carrierName)
			}

			local := ItemTrackData.XpathOne("//ProcessLocation").InnerText()
			item, err := exampleUtils.PassItem(date, event, local)
			if err != nil{
				log.Printf("parse error: PassItem error %s", err.Error())
				continue
			}
			returnDataTrackInfo = append(returnDataTrackInfo, item)
		}
		returnData.TrackInfo = returnDataTrackInfo

		itemLastIndex := len(returnData.TrackInfo)-1
		itemLastStatus := returnData.TrackInfo[itemLastIndex].StatusDescription
		if itemLastStatus == "delivered at ddu" || itemLastStatus == "Delivered at Sort Facility"{
			returnData.StatusDataNum = 2
		}
		if itemLastStatus == "Your parcel has been delivered to its sender following a return." {
			returnData.TrackInfo[itemLastIndex].SubStatusNum = 711
			returnData.StatusDataNum = 7
		}
		if yunStatus == 3{
			returnData.TrackInfo[itemLastIndex].SubStatusNum = 41
			returnData.StatusDataNum = 4
		}
		if yunStatus == 5{
			returnData.TrackInfo[itemLastIndex].SubStatusNum = 78
			returnData.StatusDataNum = 7
		}
		if yunStatus == 6{
			returnData.TrackInfo[itemLastIndex].SubStatusNum = 61
			returnData.StatusDataNum = 6
		}
		if yunStatus == 7{
			returnData.TrackInfo[itemLastIndex].SubStatusNum = 711
			returnData.StatusDataNum = 7
		}
		//returnDataJson, _ := json.Marshal(returnData)
		//fmt.Println(string(returnDataJson))
		returnData.Function = strings.Join(response.Request.RequestMeta.Function, " ,")
		result.AppendItem(&httpobj.Item{
			ItemStatus: &httpobj.ItemStatus{
				RequestId: response.Request.RequestMeta.RequestId,
			},
			Item: returnData,
		})
	}
	return result, nil
}

//func dealTrackInfoStatusAndLastInfoByFuc(returnData *httpobj.TrackItem) *httpobj.TrackItem {
//	postData := map[string]interface{}{
//		"data": returnData,
//		"function": []string{"yunexpress_api_5_200"},
//		"normal": 1,
//		"lang": "en",
//		"webRoute": []string{"yunexpress_api_on_5"},
//	}
//
//	dl := downloader.NewNetDownloader()
//	req := &httpobj.Request{
//		Url: "https://www.51tracking.com/manage/status_api.php",
//		Method: "POST",
//		Headers: map[string]string{
//			"Content-Type": "application/json",
//		},
//		Json: postData,
//		Timeout: 10,
//	}
//	resp, _ := dl.Fetch(req)
//	fmt.Println(resp.Text)
//	var resultTrackItem httpobj.TrackItem
//	err := json.Unmarshal(resp.Content, &resultTrackItem)
//	if err != nil {
//		fmt.Println(err)
//	}
//	//fmt.Println(resultTrackItem)
//	//returnData = &resultTrackItem
//	return &resultTrackItem
//}