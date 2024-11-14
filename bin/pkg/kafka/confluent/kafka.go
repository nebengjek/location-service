package kafka

import (
	"location-service/bin/config"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Producer interface {
	Publish(topic string, message []byte)
}

type Consumer interface {
	SetHandler(handler ConsumerHandler)
	Subscribe(topics ...string)
}

type ConsumerHandler interface {
	HandleMessage(message *k.Message)
}

type KafkaConfig struct {
	Username      string
	Password      string
	Address       string
	SaslMechanism string
}

var kafkaConfig KafkaConfig

func InitKafkaConfig() {

	kafkaConfig = KafkaConfig{
		Address:       config.GetConfig().KafkaUrl,
		Username:      config.GetConfig().KafkaUsername,
		Password:      config.GetConfig().KafkaPassword,
		SaslMechanism: "PLAIN",
	}
}

func GetConfig() KafkaConfig {
	return kafkaConfig
}

func (kc KafkaConfig) GetKafkaConfig() *k.ConfigMap {
	kafkaCfg := k.ConfigMap{}

	if kc.Username != "" {
		kafkaCfg["sasl.mechanism"] = kc.SaslMechanism
		kafkaCfg["sasl.username"] = kc.Username
		kafkaCfg["sasl.password"] = kc.Password
		kafkaCfg["security.protocol"] = "sasl_ssl"
	}
	kafkaCfg.SetKey("bootstrap.servers", kc.Address)
	kafkaCfg.SetKey("group.id", config.GetConfig().AppName)
	kafkaCfg.SetKey("retry.backoff.ms", 500)
	kafkaCfg.SetKey("socket.max.fails", 10)
	kafkaCfg.SetKey("reconnect.backoff.ms", 200)
	kafkaCfg.SetKey("reconnect.backoff.max.ms", 5000)
	kafkaCfg.SetKey("request.timeout.ms", 5000)
	kafkaCfg.SetKey("partition.assignment.strategy", "roundrobin")
	kafkaCfg.SetKey("auto.offset.reset", "earliest")

	return &kafkaCfg
}
