package pipeline

import (
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"log"
)

// NativeChanPipeline 原生打印管道
type NativeChanPipeline struct {
	itemChan chan *httpobj.Item
}

func (np *NativeChanPipeline) ProcessItem(item *httpobj.Item) error {
	log.Printf("Native Got item: #: %v\n", item)
	return nil
}

func (np *NativeChanPipeline) SubmitItem(item *httpobj.Item) {
	np.itemChan <- item
}

func (np *NativeChanPipeline) Run() {
	np.itemChan = make(chan *httpobj.Item)
	for item := range np.itemChan{
		err := np.ProcessItem(item)
		if err != nil {
			log.Printf("pipeline error: %s\n", err.Error())
		}
	}
}


type NativeChanPipelineFactory struct {}
func (pipeline *NativeChanPipelineFactory) GetPipeline() Pipeline {
	return &NativeChanPipeline{}
}
