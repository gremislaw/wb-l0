package kafka

import (
	"errors"
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const flushTimeoute = 5000

var errUnknownType = errors.New("unknown event type")

type (
	Producer struct {
		Producer *kafka.Producer
	}
)

func NewProducer(address []string) (*Producer, error) {

	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new producer: %w", err)
	}

	return &Producer{Producer: p}, nil
}

func (p *Producer) Produce(message []byte, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
		Key:   nil,
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.Producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("error with produce message: %w", err)
	}

	e := <-kafkaChan
	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return fmt.Errorf("error with produce message: %w", ev)
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.Producer.Flush(flushTimeoute)
	p.Producer.Close()
}
