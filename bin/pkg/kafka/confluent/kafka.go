package kafka

import (
	"location-service/bin/config"

	k "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// Producer is collection of function of kafka producer
type Producer interface {
	Publish(topic string, message []byte)
}

// Consumer is collection of function of kafka consumer
type Consumer interface {
	SetHandler(handler ConsumerHandler)
	Subscribe(topics ...string)
}

// ConsumerHandler is a collection of function for handling kafka message
type ConsumerHandler interface {
	HandleMessage(message *k.Message)
}

///

type KafkaConfig struct {
	username string
	password string
	address  string
}

var kafkaConfig KafkaConfig

func InitKafkaConfig() {
	kafkaConfig = KafkaConfig{
		address:  config.GetConfig().KafkaUrl,
		username: "",
		password: "",
	}
}

func GetConfig() KafkaConfig {
	return kafkaConfig
}

func (kc KafkaConfig) GetKafkaConfig() *k.ConfigMap {
	kafkaCfg := k.ConfigMap{}

	if kc.username != "" {
		kafkaCfg.SetKey("sasl.mechanism", "PLAIN")
		kafkaCfg.SetKey("sasl.username", kc.username)
		kafkaCfg.SetKey("sasl.password", kc.password)

		//switch securityProtocol {
		//case "sasl_ssl":
		//	kafkaCfg.SetKey("security.protocol", securityProtocol)
		//	kafkaCfg.SetKey("ssl.endpoint.identification.algorithm", "https")
		//	kafkaCfg.SetKey("enable.ssl.certificate.verification", true)
		//	break
		//case "sasl_plaintext":
		//	kafkaCfg.SetKey("security.protocol", securityProtocol)
		//	break
		//default:
		//	kafkaCfg.SetKey("security.protocol", "sasl_plaintext")
		//	break
		//}
	}

	kafkaCfg.SetKey("bootstrap.servers", kc.address)
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
