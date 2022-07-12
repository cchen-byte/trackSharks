package pipeline

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
)

type Pipeline interface {
	ProcessItem(item *httpobj.Item) error // 持久化item回调函数
	SubmitItem(item *httpobj.Item)
	Run()
}

type GetPipelineFactory interface {
	GetPipeline() Pipeline
}
