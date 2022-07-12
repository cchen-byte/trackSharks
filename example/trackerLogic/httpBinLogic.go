package trackerLogic

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/downloader"
	exampleMiddlewares "github.com/cchen-byte/trackeSharkes/example/middleware"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
)

type HttpBinLogic struct {
	constructor.BaseTrackerLogic
}

func (trackerLogic *HttpBinLogic) ConstructFirstRequest (httpobj.TrackData) *httpobj.Request {
	return &httpobj.Request{
		Url:      "https://httpbin.org",
		DownloadMiddlewares: map[string]int{
			"middleware.HBUrlMiddleware": 101,
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

func NewLogic() *HttpBinLogic {
	downloaderFactory := &downloader.NetDownloaderFactory{}
	httpBinLogic := &HttpBinLogic{}
	httpBinLogic.TrackerDownloader = downloaderFactory.GetDownloader()
	httpBinLogic.TrackerDownloaderMiddlewares = map[string]middleware.Middlewares{
		"HBUrlMiddleware": &exampleMiddlewares.HBUrlMiddleware{},
	}
	httpBinLogic.TrackerMiddlewaresManager = &middleware.Manager{}
	return httpBinLogic
}