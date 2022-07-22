package pipeline

import (
	"encoding/json"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	"log"
)

// NativeChanPipeline 原生打印管道
type NativeChanPipeline struct {
	itemChan chan *httpobj.TrackItem
}

func (np *NativeChanPipeline) ProcessItem(item *httpobj.TrackItem) error {
	itemData, _ := json.Marshal(item)
	log.Printf("Native Got item: %v\n", string(itemData))
	return nil
}

func (np *NativeChanPipeline) SubmitItem(item *httpobj.TrackItem) {
	np.itemChan <- item
}

func (np *NativeChanPipeline) Run() {
	//np.itemChan = make(chan *httpobj.Item)
	for item := range np.itemChan{
		err := np.ProcessItem(item)
		if err != nil {
			log.Printf("pipeline error: %s\n", err.Error())
		}
	}
}


//type NativeChanPipelineFactory struct {}
//func (pipeline *NativeChanPipelineFactory) GetPipeline() Pipeline {
//	return &NativeChanPipeline{}
//}

func NewNativeChanPipeline() *NativeChanPipeline {
	return &NativeChanPipeline{
		itemChan: make(chan *httpobj.TrackItem),
	}
}