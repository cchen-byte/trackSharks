package constructor

import (
	"errors"
	"github.com/cchen-byte/trackeSharkes/downloader"
	"github.com/cchen-byte/trackeSharkes/middleware"
	"github.com/cchen-byte/trackeSharkes/pipeline"
)

// BaseTrackerLogic 基础追踪逻辑结构
type BaseTrackerLogic struct {
	TrackerDownloader            downloader.Downloader

	TrackerMiddlewaresManager *middleware.Manager
	TrackerDownloaderMiddlewares map[string]middleware.Middlewares
	TrackerPipeline pipeline.Pipeline
}

// GetTrackerDownloader 获取对应的下载器
func (bt *BaseTrackerLogic) GetTrackerDownloader() (downloader.Downloader, error) {
	if bt.TrackerDownloader == nil{
		return nil, errors.New("Downloader is undefined")
	}
	return bt.TrackerDownloader, nil
}

// GetTrackerPipeline 获取对应的管道
func (bt *BaseTrackerLogic) GetTrackerPipeline() pipeline.Pipeline {
	return bt.TrackerPipeline
}

// GetTrackerDownloaderMiddlewares 获取注册的中间件
func (bt *BaseTrackerLogic) GetTrackerDownloaderMiddlewares() map[string]middleware.Middlewares {
	return bt.TrackerDownloaderMiddlewares
}

// GetRequestsDownloaderMiddlewares 获取每个请求需要的中间件
func (bt *BaseTrackerLogic) GetRequestsDownloaderMiddlewares(requestDownloaderMiddleware map[string]int) ([]middleware.Middlewares, error) {
	trackerDownloaderMiddlewares := bt.GetTrackerDownloaderMiddlewares()
	middlewareList, err := bt.TrackerMiddlewaresManager.GetMiddlewares(requestDownloaderMiddleware, trackerDownloaderMiddlewares)
	if err != nil {
		return nil, err
	}
	return middlewareList, err
}