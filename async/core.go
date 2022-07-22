package async

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/downloader"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/middleware"
	"github.com/cchen-byte/trackeSharkes/pipeline"
	"github.com/cchen-byte/trackeSharkes/scheduler"
	"github.com/cchen-byte/trackeSharkes/setting"
	"github.com/cchen-byte/trackeSharkes/tracker"
)

// RunAsync 异步查询
func RunAsync() {
	// 配置scheduler
	asyncScheduler := scheduler.NewChanScheduler()
	go asyncScheduler.Run()

	// 配置pipeline
	asyncPipeline := pipeline.NewNativeChanPipeline()
	go asyncPipeline.Run()

	//
	engineConstructorChan := make(chan *httpobj.ItemStatus)

	// 配置engine
	asyncEngine := engine.NewChanEngine(asyncScheduler, asyncPipeline, engineConstructorChan)
	go asyncEngine.Run()

	// 配置downloader
	trackerDownloader := downloader.NewNetDownloader()

	// 配置tracker
	for i:=0; i< setting.Settings.WorkerCount; i++{
		tracker.CreateChanTracker(trackerDownloader, middleware.TrackerMiddlewaresManager, asyncScheduler.GetTrackerChan(), asyncEngine, asyncScheduler)
	}

	collectDataChan := make(chan *httpobj.TrackData)
	// 配置数据读取器
	dataCollect := &testDataCollect{}
	go dataCollect.Run(collectDataChan)

	// 线路管理器
	cm := constructor.NewConstructorManager(engineConstructorChan)
	go cm.Run(asyncEngine)

	for {
		select {
		case collectData := <- collectDataChan:
			cm.SubmitTrackData(collectData)
		}
	}

}