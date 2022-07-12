package tracker

import (
	"fmt"
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/downloader"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
	"github.com/cchen-byte/trackeSharkes/scheduler"
	"log"
)

type ChanTracker struct {
	Downloader downloader.Downloader
}

// fetchWork 根据传入的 request 发起请求， 返回一个 response
func (tracker *ChanTracker) fetchWork(request *httpobj.Request) (*httpobj.Response, error) {
	log.Printf("Fetching %s\n", request.Url)
	resp, err := tracker.Downloader.Fetch(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// parseWork 解析响应
func (tracker *ChanTracker) parseWork(request *httpobj.Request, resp *httpobj.Response) (*httpobj.ParseResult, error) {
	// 根据任务请求中的解析函数解析网页数据
	return request.Callback(resp)
}

func (tracker *ChanTracker) Run(trackerLogicUtils constructor.TrackerLogicUtils, request *httpobj.Request) (*httpobj.ParseResult, error) {
	var downloadMiddlewares []middleware.Middlewares

	//// 当请求未携带下载中间件信息, 直接下载并解析
	//if len(request.DownloadMiddlewares) == 0 {
	//	resp, err := worker.fetchWork(request)
	//	if err != nil {
	//		log.Printf("Fetch error: %s, request.Url: %s\n", err.Error(), request.Url)
	//		return nil, err
	//	}
	//	parseResult, err := worker.parseWork(request, resp)
	//	if err != nil {
	//		log.Printf("Parse error: %s, request.Url: %s\n", err.Error(), request.Url)
	//		return nil, err
	//	}
	//	return parseResult, nil
	//}

	// 获取该请求对应的中间件
	var err error
	downloadMiddlewares, err = trackerLogicUtils.GetRequestsDownloaderMiddlewares(request.DownloadMiddlewares)
	if err != nil {
		log.Printf("Get Middlewares error: %s, request.Url: %s\n", err.Error(), request.Url)
		return nil, err
	}

	// 中间件处理requests
	for _, v := range downloadMiddlewares {
		err := v.ProcessRequest(request)
		if err != nil {
			log.Printf("ProcessRequest error: %s, request.Url: %s\n", err.Error(), request.Url)
			return nil, err
		}
	}

	// 请求下载
	resp, err := tracker.fetchWork(request)
	if err != nil {
		log.Printf("Fetch error: %s, request.Url: %s\n", err.Error(), request.Url)
		return nil, err
	}

	// 中间件处理response
	for _, v := range downloadMiddlewares {
		err := v.ProcessResponse(resp)
		if err != nil {
			log.Printf("ProcessResponse error: %s, request.Url: %s\n", err.Error(), request.Url)
			return nil, err
		}
	}

	// 返回一个parseResult, 若 parseResult 不为空则直接返回 parseResult
	parseResult, err := handleResponse(resp)
	if err != nil {
		log.Printf("HandleResponse error: %s, request.Url: %s\n", err.Error(), request.Url)
		return nil, err
	}
	if parseResult != nil {
		return parseResult, nil
	}

	// 根据任务请求中的解析函数解析网页数据
	parseResult, err = tracker.parseWork(request, resp)
	if err != nil {
		log.Printf("Parse error: %s, request.Url: %s\n", err.Error(), request.Url)
		return nil, err
	}
	return parseResult, nil
}


// handleResponse 内置处理响应, 当前主要是处理重试
func handleResponse(resp *httpobj.Response) (*httpobj.ParseResult, error) {
	parseResult := &httpobj.ParseResult{}
	if resp.Request.IsRetry {
		parseResult.Requests = append(parseResult.Requests, resp.Request)
		return parseResult, nil
	}
	return nil, nil
}


// CreateChanTracker 创建ChanTracker
func CreateChanTracker(trackerLogicUtils constructor.TrackerLogicUtils, schedulerTrackerChan chan *httpobj.Request, engine engine.Engine, ready scheduler.ReadyNotifier) {
	trackerDownloader, err := trackerLogicUtils.GetTrackerDownloader()
	if err != nil{
		panic(fmt.Sprintf("Tracker get downloader err: %s\n", err.Error()))
	}
	// tracker 下载器使用爬虫对应的下载器
	cTracker := &ChanTracker{
		Downloader: trackerDownloader,
	}
	go func(tracker *ChanTracker) {
		for {
			ready.TrackerReady(schedulerTrackerChan)
			// 调度器内无请求则一直阻塞
			request := <-schedulerTrackerChan

			result, err := tracker.Run(trackerLogicUtils, request)
			if err != nil {
				continue
			}
			engine.SubmitParseResult(result)
		}
	}(cTracker)
}