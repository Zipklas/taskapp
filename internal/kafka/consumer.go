package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer   sarama.Consumer
	topic      string
	handlers   []func(message []byte)
	handlersMu sync.RWMutex
}

func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (c *Consumer) AddHandler(handler func(message []byte)) {
	c.handlersMu.Lock()
	defer c.handlersMu.Unlock()
	c.handlers = append(c.handlers, handler)
}

func (c *Consumer) Start() {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				c.handlersMu.RLock()
				for _, handler := range c.handlers {
					handler(msg.Value)
				}
				c.handlersMu.RUnlock()

			case err := <-partitionConsumer.Errors():
				log.Printf("Kafka consumer error: %v", err)
			}
		}
	}()
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
