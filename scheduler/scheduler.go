package scheduler

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

// Scheduler 任务调度器
type Scheduler interface {
	ReadyNotifier
	Submit(*httpobj.Request) // 提交一个请求对象
	// todo: 模块化
	GetTrackerChan() chan *httpobj.Request	// 返回一个 schedulerTrackerChan
	Run()
}

type ReadyNotifier interface {
	TrackerReady(chan *httpobj.Request)	// 当tracker空闲时向scheduler报告
}


//type GetSchedulerFactory interface {
//	GetScheduler() Scheduler
//}