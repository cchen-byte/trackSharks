package setting

import (
	"time"
)

type Options struct {
	// 请求间隔
	RateLimit time.Duration
	// 请求重试次数
	RetryTime int

	// 下载器类型
	Downloader string
	// 下载器数量
	WorkerCount int

	// 调度器类型
	Scheduler string
	// 调度器内部队列类型
	SchedulerQueue string
	// 管道并发处理数量
	PipelineCount int
}

// DefaultSettingOptions 返回一个默认的全局设置
func DefaultSettingOptions() *Options {
	return &Options{
		WorkerCount: 16,
		RetryTime: 3,
		RateLimit: 0,
		Scheduler: "ChanScheduler",
		SchedulerQueue: "LIFO",
		Downloader: "NetDownloader",
		PipelineCount: 1,
	}
}

// Settings 全局设置
var Settings = DefaultSettingOptions()

// UpdateSettings 更新全局设置
func UpdateSettings(opt *Options) {
	if opt.WorkerCount != 0{
		Settings.WorkerCount = opt.WorkerCount
	}

	if opt.RetryTime != 0{
		Settings.RetryTime = opt.RetryTime
	}

	var emptyTimeDuration time.Duration
	if opt.RateLimit != emptyTimeDuration{
		Settings.RateLimit = opt.RateLimit
	}

	if len(opt.Scheduler) != 0{
		Settings.Scheduler = opt.Scheduler
	}

	if len(opt.SchedulerQueue) != 0{
		Settings.SchedulerQueue = opt.SchedulerQueue
	}

	if len(opt.Downloader) != 0{
		Settings.Downloader = opt.Downloader
	}

	if opt.PipelineCount != 0{
		Settings.PipelineCount = opt.PipelineCount
	}
}