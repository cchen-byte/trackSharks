package async

import (
	"github.com/cchen-byte/trackeSharkes/constructor"
	"github.com/cchen-byte/trackeSharkes/engine"
	"github.com/cchen-byte/trackeSharkes/pipeline"
	"github.com/cchen-byte/trackeSharkes/scheduler"
	"github.com/cchen-byte/trackeSharkes/setting"
	"github.com/cchen-byte/trackeSharkes/tracker"
)

// RunAsync 异步查询
func RunAsync(trackerLogic constructor.TrackerLogic) {
	// 配置scheduler
	schedulerFactory := &scheduler.ChanSchedulerFactory{}
	asyncScheduler := schedulerFactory.GetScheduler()
	go asyncScheduler.Run()

	// 配置pipeline
	pipelineFactory := &pipeline.NativeChanPipelineFactory{}
	asyncPipeline := pipelineFactory.GetPipeline()
	go asyncPipeline.Run()

	// 配置engine
	asyncEngine := &engine.ChanEngine{
		Scheduler:  asyncScheduler,
		Pipeline: asyncPipeline,
	}

	// 配置tracker
	for i:=0; i< setting.Settings.WorkerCount; i++{
		tracker.CreateChanTracker(trackerLogic, asyncScheduler.GetTrackerChan(), asyncEngine, asyncScheduler)
	}

	// 配置数据读取器
	dataCollect := &testDataCollect{}
	go dataCollect.Run(asyncEngine, trackerLogic)

	asyncEngine.Run()
}