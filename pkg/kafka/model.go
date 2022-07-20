package kafka

type Config struct {
	ProducerTopic    string `yaml:"producer_topic"`
	ConsumerTopic    string `yaml:"consumer_topic"`
	//TopicForApi      string `yaml:"topicForApi"`
	GroupId          string `yaml:"groupId"`
	BootstrapServers string `yaml:"bootstrapServers"`
	SecurityProtocol string `yaml:"securityProtocol"`
	SslCaLocation    string `yaml:"sslCaLocation"`
	SaslMechanism    string `yaml:"saslMechanism"`
	SaslUsername     string `yaml:"saslUsername"`
	SaslPassword     string `yaml:"saslPassword"`
}

// TODO: 账号密码等由全局配置定义, 分区等由具体爬虫定义
var kafkaConfig = Config{
	ProducerTopic: "data_clean",
	//ConsumerTopic: "spider_test",
	ConsumerTopic: "yunexpress",
	GroupId: "spider_test_consumer",
	BootstrapServers: "alikafka-post-cn-zvp2nv48l009-1.alikafka.aliyuncs.com:9093,alikafka-post-cn-zvp2nv48l009-2.alikafka.aliyuncs.com:9093,alikafka-post-cn-zvp2nv48l009-3.alikafka.aliyuncs.com:9093",
	SecurityProtocol: "SASL_SSL",
	//SecurityProtocol: "SASL_PLAINTEXT",
	SslCaLocation: "",
	SaslMechanism: "PLAIN",
	SaslUsername: "alikafka_post-cn-zvp2nv48l009",
	SaslPassword: "tQBTHPCcPXjxxLwU60pOQW3dTDY9n9rx",
}

