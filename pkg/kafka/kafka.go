package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog/log"
	//"log"
)

const (
	Int32Max = 2147483647 - 1000
)

func DoInitProducer() *kafka.Producer {
	//cfg := conf.Setting.Kafka
	cfg := kafkaConfig

	log.Info().Msg("init kafka producer, it may take a few seconds to init the connection")
	//common arguments
	// TODO:这里需要权衡一下，如果linger.ms=0代表有消息就立马发送;batch.size是存够了一定大小才发送，
	// 这样如果太大可能造成一直阻塞提高延迟，但一定程度能提高吞吞吐率！
	// 发送条件是只取linger.ms 和 batch.size 其中一个，只要达到了里面发送！

	// Local Queu full 正解
	// https://www.notion.so/webhook-kafka-queue-full-385d224a51514b47acfb30763f495d6f#eb6324d68828449cb5a7b9921072d2b0
	// 本质就是没确认消息。但是很奇怪，ACK=0不需要确认，为什么阿里云还要，有点神奇。。。
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":           "true",
		"message.max.bytes":             1000000,
		"batch.size":                    563840,
		"linger.ms":                     100,
		"sticky.partitioning.linger.ms": 1000,
		"retries":                       Int32Max,
		"retry.backoff.ms":              1000,
		"acks":                          "1",
	}
	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)

	switch cfg.SecurityProtocol {
	case "PLAINTEXT":
		kafkaconf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		kafkaconf.SetKey("security.protocol", "sasl_ssl")
		kafkaconf.SetKey("ssl.ca.location", "../pkg/kafka/ca-cert.pem")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	case "SASL_PLAINTEXT":
		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	default:
		panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
	}

	producer, err := kafka.NewProducer(kafkaconf)
	if err != nil {
		panic(err)
	}
	log.Info().Msg("init kafka producer success")
	return producer
}

func doInitConsumer(cfg Config) *kafka.Consumer {
	log.Info().Msg("init kafka consumer, it may take a few seconds to init the connection")
	//common arguments
	var kafkaconf = &kafka.ConfigMap{
		"api.version.request":   "true",
		"enable.auto.commit":    "false",
		"auto.offset.reset":     "latest",
		"heartbeat.interval.ms": 3000,
		"session.timeout.ms":    30000,
		"max.poll.interval.ms":  120000,
		// "max.poll.records":          1000, //There is no such property, messages are fetched in batches from the broker and queued locally, but the application may only receive one message at the time.
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000}
	kafkaconf.SetKey("bootstrap.servers", cfg.BootstrapServers)
	kafkaconf.SetKey("group.id", cfg.GroupId)

	switch cfg.SecurityProtocol {
	case "PLAINTEXT":
		kafkaconf.SetKey("security.protocol", "plaintext")
	case "SASL_SSL":
		kafkaconf.SetKey("security.protocol", "sasl_ssl")
		kafkaconf.SetKey("ssl.ca.location", "ca-cert.pem")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)
	case "SASL_PLAINTEXT":
		kafkaconf.SetKey("security.protocol", "sasl_plaintext")
		kafkaconf.SetKey("sasl.username", cfg.SaslUsername)
		kafkaconf.SetKey("sasl.password", cfg.SaslPassword)
		kafkaconf.SetKey("sasl.mechanism", cfg.SaslMechanism)

	default:
		panic(kafka.NewError(kafka.ErrUnknownProtocol, "unknown protocol", true))
	}

	consumer, err := kafka.NewConsumer(kafkaconf)
	if err != nil {
		panic(err)
	}
	fmt.Print("init kafka consumer success\n")
	return consumer
}

func Produce(KafkaProducer *kafka.Producer, msgBody string) error {

	// Produce messages to topic (asynchronously)
	var msg *kafka.Message = nil

	msg = &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaConfig.ProducerTopic, Partition: kafka.PartitionAny},
		Value:          []byte(msgBody),
	}

	err := KafkaProducer.Produce(msg, nil)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("KafkaProducer Produce Error: %s", err.Error()))
		return err
	}
	// Wait for message deliveries before shutting down
	KafkaProducer.Flush(1)
	return nil
}

func Consume(ConsumeChan chan string) {
	consumer := doInitConsumer(kafkaConfig)
	err := consumer.SubscribeTopics([]string{kafkaConfig.ConsumerTopic}, nil)
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Consumer SubscribeTopics Error: %s", err.Error()))
	}
	log.Printf("Consumer SubscribeTopics success: %s", kafkaConfig.ConsumerTopic)
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Consumer Close Error: %s", err.Error()))
		}
	}(consumer)
	//fmt.Printf("get %v", consumer.GetMetadata())
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			// 接受消息进行处理，并针对消息进行消费确认，这里用同步方式
			// handleErr := ConsumerT.HandleKafkaMessage(string(msg.Value))
			log.Info().Msg(fmt.Sprintf("Consumer Get Message: %s", msg))
			ConsumeChan <- string(msg.Value)
			_, err := consumer.CommitMessage(msg)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("Consumer CommitMessage Error: %s", err.Error()))
			}

		} else {
			// The client will automatically try to recover from all errors.
			log.Info().Msg(fmt.Sprintf("Consumer error: %v \n msg: (%v)", err, msg))
		}
	}
}
