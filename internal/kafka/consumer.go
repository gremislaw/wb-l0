package kafka

import (
	"encoding/json"
	"fmt"
	"order_service/internal/config"
	. "order_service/internal/logger"
	"order_service/internal/models"
	"order_service/internal/service"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer *kafka.Consumer
	service  *service.Service
	stop     bool
}

func NewConsumer(cfg *config.Config, service *service.Service) (*Consumer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.KAFKA_BOOTSTRAP_SERVERS, ","),
		"group.id":          cfg.KAFKA_CONSUMER_GROUP,
	}

	c, err := kafka.NewConsumer(conf)
	if err != nil {
		return nil, fmt.Errorf("error with new consumer: %w", err)
	}

	if err := c.Subscribe(cfg.KAFKA_TOPIC, nil); err != nil {
		return nil, fmt.Errorf("error with consumer subscribe: %w", err)
	}

	return &Consumer{
		consumer: c,
		service: service,
	}, nil
}

func (c *Consumer) Start() {
	for {

		if c.stop {
			break
		}

		// чтение сообщений из брокера
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			Logger.Error("error with read message", zap.Any("error", err))
		}

		if kafkaMsg == nil {
			continue
		}

		if err := c.ConsumeMessage(kafkaMsg.Value); err != nil {
			switch err {

			case models.ErrJSONUnmarshal: // случай несоответствия модели orders
				Logger.Warn("unexpected message content")
			default:
				Logger.Error("error with handle kafka message", zap.Any("error", err))
			}
			continue
		}

		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			Logger.Error("error with offset kafka message", zap.Any("error", err))
			continue
		}
	}
}

func (c *Consumer) ConsumeMessage(message []byte) error {
	var order models.Order

	if err := json.Unmarshal(message, &order); err != nil {
		return models.ErrJSONUnmarshal
	}

	// запись в кэш и бд
	if err := c.service.SetOrder(&order); err != nil {
		return err
	}

	return nil
}

func (c *Consumer) Stop() error {
	c.stop = true

	if _, err := c.consumer.Commit(); err != nil {
		return fmt.Errorf("error with commit consumer: %w", err)
	}

	Logger.Info("commitet offset")

	return c.consumer.Close()
}
