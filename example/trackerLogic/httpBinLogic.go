package trackerLogic

import (
	exampleMiddleware "github.com/cchen-byte/trackeSharkes/example/middleware"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
)

func init(){
	var hbUrlMiddleware = new(exampleMiddleware.HBUrlMiddleware)
	// 注册中间件
	middleware.TrackerMiddlewaresManager.MiddlewaresMap.Store("hbUrlMiddleware", hbUrlMiddleware)
}

type HttpBinLogic1 struct {
}

func (trackerLogic *HttpBinLogic1) ConstructFirstRequest (httpobj.TrackData) *httpobj.Request {
	return &httpobj.Request{
		Url:      "https://httpbin.org",
		DownloadMiddlewares: map[string]int{
			"hbUrlMiddleware": 101,
		},
		Callback: Parse,
	}
}

// Parse 解析返回的响应
func Parse(response *httpobj.Response) (*httpobj.ParseResult, error) {
	respJsonDom, _ := response.GetJsonDom()

	tracId := respJsonDom.XpathOne("headers/X-Amzn-Trace-Id").InnerText()
	itemData := &httpobj.Item{
		"TracId": tracId,
	}

	result := httpobj.NewParseResult()
	result.AppendItem(itemData)
	return result, nil
}

func NewLogic() *HttpBinLogic1 {
	httpBinLogic := &HttpBinLogic1{}
	return httpBinLogic
}