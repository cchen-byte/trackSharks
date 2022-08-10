package async

//import (
//	"encoding/json"
//	"errors"
//	"github.com/cchen-byte/trackeSharkes/httpobj"
//	"github.com/cchen-byte/trackeSharkes/pkg/kafka"
//	"log"
//	"time"
//)
//
//type kafkaDataCollect struct {}
//
//// singleRead
//func (kc *kafkaDataCollect) singleRead(collectDataChan chan *httpobj.TrackData){
//	// 创建管道并开启 kafka 消费者
//	consumeChan := make(chan string)
//	go kafka.Consume(consumeChan)
//
//	// 读取 kafka 消息
//	kafkaMessage := <- consumeChan
//
//	// 解析结构化 kafka 消息
//	var kafkaTrackData httpobj.TrackData
//	var kafkaBaseTrackData *httpobj.BaseTrackData
//	err := json.Unmarshal([]byte(kafkaMessage), &kafkaBaseTrackData)
//	if err != nil {
//		log.Printf("HandleKafkaMessage Json Unmarshal Error: %s", err)
//	}
//	kafkaTrackData = append(kafkaTrackData, kafkaBaseTrackData)
//
//	// 推送出去
//	collectDataChan <- &kafkaTrackData
//}
//
//
//// UnBlockRead 无阻塞读取channel消息
//func UnBlockRead(consumeChan chan string) (string, error) {
//	select {
//	case kafkaMessage := <- consumeChan:
//		return kafkaMessage, nil
//	case <-time.After(time.Millisecond):
//		return "", errors.New("timeout")
//	}
//}
//
//// BatchRead 读取kafka
//// 按照指定 batchSize 分组返回
//// 若指定 timeOut 时间内没有凑成一组拥有 batchSize 个的分组, 则返回当前所有元素
//func (kc *kafkaDataCollect) BatchRead(collectDataChan chan *httpobj.TrackData) {
//	var batchSize int = 1
//	var timeOut time.Duration = time.Second * 10
//
//	var wait time.Time   // 等待时间
//	var kafkaMsgQ httpobj.TrackData
//	isFirst := true // 是否为新分组的第一个
//
//	consumeChan := make(chan string)
//	go kafka.Consume(consumeChan)
//	for{
//		var sliceNum int
//		kafkaMessage, err := UnBlockRead(consumeChan)
//		if err == nil{
//			var kafkaTrackData *httpobj.BaseTrackData
//			err = json.Unmarshal([]byte(kafkaMessage), &kafkaTrackData)
//			if err != nil {
//				log.Printf("HandleKafkaMessage Json Unmarshal Error: %s", err)
//			}
//			kafkaMsgQ = append(kafkaMsgQ, kafkaTrackData)
//
//			// 是否达到批量要求
//			isBatch := len(kafkaMsgQ) >= batchSize
//			if isBatch {
//				sliceNum = batchSize
//			}
//			// 若是分组的第一个则开始计时
//			if isFirst == true{
//				wait = time.Now()
//				isFirst = false
//			}
//		}
//
//		// 是否达到超时限制
//		isTimeOut := time.Since(wait) > timeOut
//		if isTimeOut {
//			sliceNum = len(kafkaMsgQ)
//		}
//
//		if sliceNum > 0{
//			collectDataChan <- &kafkaMsgQ
//			kafkaMsgQ = httpobj.TrackData{}
//			wait = time.Now()
//			isFirst = true
//		}
//	}
//}
//
//func (kc *kafkaDataCollect) Run(collectDataChan chan *httpobj.TrackData){
//	kc.BatchRead(collectDataChan)
//}