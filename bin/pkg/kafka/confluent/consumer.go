package kafka

import (
	"fmt"
	"location-service/bin/pkg/log"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type consumer struct {
	handler  ConsumerHandler
	consumer *kafka.Consumer
	logger   log.Log
}

// NewConsumer is a constructor of kafka consumer
func NewConsumer(config *kafka.ConfigMap, log log.Log) (Consumer, error) {
	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	return &consumer{
		logger:   log,
		consumer: c,
	}, nil
}

func (c *consumer) SetHandler(handler ConsumerHandler) {
	c.handler = handler
}

func (c *consumer) Subscribe(topics ...string) {
	if c.handler == nil {
		joinTopic := strings.Join(topics, ", ")
		msg := fmt.Sprintf("Kafka Consumer Error: Topics: [%s] There is no consumer handler to handle message from incoming event", joinTopic)
		c.logger.Error("", msg, "", "")
		return
	}

	c.consumer.SubscribeTopics(topics, nil)

	go func() {
		for {
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				msg := fmt.Sprintf("Kafka Consumer Error: %v (%v)\n", err, msg)
				c.logger.Error("", msg, "", "")
				continue
			}
			go c.handler.HandleMessage(msg)
			c.consumer.CommitMessage(msg)
		}
	}()

	return
}
