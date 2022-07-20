package yunexpress

import (
	"fmt"
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

func api2ConstructFirstRequest(trackData httpobj.TrackData) *httpobj.Request{
	return &httpobj.Request{
		Url: fmt.Sprintf("http://yunapi.yunexpress.com/LMS.API/api/WayBill/GetOrder?number=%s", trackData[0].TrackNumber),
		Headers: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
			"Authorization": "Basic QzYyMjQ5Jjd4aTFSRXJPWEg0PQ==",
		},
		MetaData: &httpobj.MetaData{
			"trackData": trackData,
		},
		Callback: parseApi2TrackData,
	}
}

func parseApi2TrackData(response *httpobj.Response) (*httpobj.ParseResult, error) {
	result := httpobj.NewParseResult()
	returnData := httpobj.NewTrackItem()
	returnData.Function += fmt.Sprintf("web_api_2_%d", response.StatusCode)

	// 异常处理

	// 正常逻辑
	//respJsonDom, _ := response.GetJsonDom()


	return result, nil
}