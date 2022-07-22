package pipeline

import "github.com/cchen-byte/trackeSharkes/httpobj"

type Pipeline interface {
	ProcessItem(item *httpobj.TrackItem) error // 持久化item回调函数
	SubmitItem(item *httpobj.TrackItem)
	Run()
}

//type GetPipelineFactory interface {
//	GetPipeline() Pipeline
//}
