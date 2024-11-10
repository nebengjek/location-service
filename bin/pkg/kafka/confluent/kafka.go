package kafka

import (
	"encoding/base64"
	"os"

	"location-service/bin/config"
	"location-service/bin/pkg/log"

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
	CA            string
	SaslMechanism string
}

var kafkaConfig KafkaConfig

func InitKafkaConfig() {
	encodedCA := config.GetConfig().KafkaCaCert

	decodedCA, err := base64.StdEncoding.DecodeString(encodedCA)
	if err != nil {
		log.GetLogger().Error("decoded-ca", "Failed to decode Kafka CA certificate: %s", "", err.Error())
	}

	caFilePath := "/tmp/kafka-ca-cert.pem"
	err = os.WriteFile(caFilePath, decodedCA, 0644)
	if err != nil {
		log.GetLogger().Error("write-ca", "Failed to write decoded Kafka CA certificate to file: %s", "", err.Error())
	}
	kafkaConfig = KafkaConfig{
		Address:       config.GetConfig().KafkaUrl,
		Username:      config.GetConfig().KafkaSaslUsername,
		Password:      config.GetConfig().KafkaSaslPassword,
		CA:            caFilePath,
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
		kafkaCfg["ssl.ca.location"] = kc.CA
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
