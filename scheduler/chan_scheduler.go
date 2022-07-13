package scheduler

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"github.com/cchen-byte/trackeSharkes/setting"
	"github.com/cchen-byte/trackeSharkes/utils"
	"time"
)

// ChanScheduler 使用channel作为内置队列的调度器
type ChanScheduler struct {
	requestChan          chan *httpobj.Request      // 请求对象队列
	trackerSchedulerChan chan chan *httpobj.Request // tracker channel
}

func (scheduler *ChanScheduler) Submit(request *httpobj.Request) {
	scheduler.requestChan <- request
}

func (scheduler *ChanScheduler) GetTrackerChan() chan *httpobj.Request {
	return make(chan *httpobj.Request) // 无缓存通道
}

// TrackerReady Tracker准备会向管道内发送
func (scheduler *ChanScheduler) TrackerReady(w chan *httpobj.Request) {
	scheduler.trackerSchedulerChan <- w
}

func (scheduler *ChanScheduler) Run() {
	//scheduler.requestChan = make(chan *httpobj.Request)               // 请求channel
	//scheduler.trackerSchedulerChan = make(chan chan *httpobj.Request) // tracker channel

	// 创建请求队列和工作队列
	var workerQ []chan *httpobj.Request
	RequestQ := utils.NewQueue(setting.Settings.SchedulerQueue)

	// 请求时间间隔
	RateLimitTimeDuration := setting.Settings.RateLimit
	start := time.Now().Add(time.Hour * - 1)
	go func() {
		for {
			//fmt.Println(len(workerQ), RequestQ.GetSize())
			var activeWorker chan *httpobj.Request
			var activeRequest *httpobj.Request

			var EngineSchedulerChan chan *httpobj.Request

			// 当本地队列存量小于16个才进行写入
			if RequestQ.GetSize() < 16 {
				EngineSchedulerChan = scheduler.requestChan
			}

			// requestChan, workerQ同时有数据时
			if len(workerQ) > 0 && RequestQ.GetSize() > 0 && time.Since(start) >= RateLimitTimeDuration {
				activeWorker = workerQ[0]
				activeRequest = RequestQ.Top().Val.(*httpobj.Request)
			}

			select {
			// 速率限制, 不然会造成生产速度远大于消耗
			case r := <-EngineSchedulerChan: // 当 EngineSchedulerChan 为 nil 会阻塞
				priority := 0
				if r.Priority > 0 {
					priority = r.Priority
				}
				// 推入队列
				RequestQ.Push(utils.Node{
					Val: r,
					Point: priority,
				})
			case w := <-scheduler.trackerSchedulerChan: // 当 workerChan 收到数据
				workerQ = append(workerQ, w)
			case activeWorker <- activeRequest: // 给任务队列分配任务, 当 activeWorker 为 nil 会阻塞
				workerQ = workerQ[1:]
				RequestQ.Pop()
				start = time.Now()
			case <- time.Tick(time.Millisecond*100):
			}
		}
	}()
}

//type ChanSchedulerFactory struct {}
//func (scheduler *ChanSchedulerFactory) GetScheduler() Scheduler {
//	return &ChanScheduler{}
//}

func NewChanScheduler() *ChanScheduler {
	return &ChanScheduler{
		requestChan: make(chan *httpobj.Request),
		trackerSchedulerChan: make(chan chan *httpobj.Request),
	}
}