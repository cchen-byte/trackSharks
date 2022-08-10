package pipeline

import (
	"encoding/json"
	"github.com/cchen-byte/trackeSharkes/httpobj"
	exampleKafka "github.com/cchen-byte/trackeSharkes/pkg/kafka"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type KafkaPipeline struct {
	KafkaProducer *kafka.Producer
	itemChan chan *httpobj.TrackItem
}

func (kp *KafkaPipeline) ProcessItem(item *httpobj.TrackItem) error {
	//log.Printf("KAFKA Got item: #: %v\n", item)
	msg, _ := json.Marshal(item)
	log.Printf("KAFKA Got item: #: %s\n", string(msg))
	//return nil
	return exampleKafka.Produce(kp.KafkaProducer, string(msg))
}

func (kp *KafkaPipeline) SubmitItem(item *httpobj.TrackItem) {
	kp.itemChan <- item
}

func (kp *KafkaPipeline) Run() {
	for item := range kp.itemChan{
		err := kp.ProcessItem(item)
		if err != nil {
			log.Printf("pipeline error: %s\n", err.Error())
		}
	}
}

func NewKafkaPipeline() *KafkaPipeline {
	kafkaProducer := exampleKafka.DoInitProducer()
	kafkaPipeline := &KafkaPipeline{
		KafkaProducer: kafkaProducer,
		itemChan: make(chan *httpobj.TrackItem),
	}

	// ack
	go func() {
		for e := range kafkaProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return kafkaPipeline
}