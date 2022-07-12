package constructor

import (
	"github.com/cchen-byte/trackeSharkes/downloader"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
	"github.com/cchen-byte/trackeSharkes/pipeline"
)

// TrackerLogic 追踪逻辑, 由callback链接
type TrackerLogic interface {
	TrackerLogicUtils
	ConstructFirstRequest(httpobj.TrackData) *httpobj.Request 	// 根据数据构造请求
}

//
type TrackerLogicNode struct {
	Logic TrackerLogic
	Next *TrackerLogicNode
}

type TrackerLogicUtils interface {
	GetTrackerDownloader() (downloader.Downloader, error)
	GetTrackerPipeline() pipeline.Pipeline
	GetTrackerDownloaderMiddlewares() map[string]middleware.Middlewares
	GetRequestsDownloaderMiddlewares(requestDownloaderMiddleware map[string]int) ([]middleware.Middlewares, error)
}