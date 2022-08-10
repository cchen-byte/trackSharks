package engine

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/pipeline"
	"github.com/cchen-byte/trackeSharkes/scheduler"
)

type ChanEngine struct {
	trackerEngineChan chan *httpobj.ParseResult // tracker 向引擎推送数据管道
	scheduler   scheduler.Scheduler // 任务调度器
	pipeline pipeline.Pipeline
	EngineConstructorChan chan *httpobj.ItemStatus
}

// SubmitRequests 向引擎推送 Request
func (e *ChanEngine) SubmitRequests(request *httpobj.Request){
	e.scheduler.Submit(request)
}

// SubmitParseResult 向引擎推送 ParseResult
func (e *ChanEngine) SubmitParseResult(result *httpobj.ParseResult){
	select {
	case e.trackerEngineChan <- result:
	}
}

func (e *ChanEngine) Run() {
	// e.trackerEngineChan = make(chan *httpobj.ParseResult)
	// 处理解析后的item以及request
	for {
		select {
		// 处理从tracker过来的结果
		case result := <-e.trackerEngineChan:
			// 然后把 Tracker 解析出的 item 提交 Pipeline
			for _, itemData := range result.Items {
				// 无异常
				if !itemData.ItemStatus.IsError{
					go e.pipeline.SubmitItem(itemData.Item)
				}
				go func(itemData *httpobj.Item) {
					e.EngineConstructorChan <- itemData.ItemStatus
				}(itemData)
			}

			// 然后把 Tracker 解析出的 Request 提交 Scheduler
			for _, request := range result.Requests {
				go e.SubmitRequests(request)
			}
		}
	}
}

func NewChanEngine(sc scheduler.Scheduler, p pipeline.Pipeline, ecchan chan *httpobj.ItemStatus) *ChanEngine{
	return &ChanEngine{
		trackerEngineChan: make(chan *httpobj.ParseResult),
		scheduler: sc,
		pipeline: p,
		EngineConstructorChan: ecchan,
	}
}

