package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	. "order_service/internal/logger"
	"order_service/internal/models"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func (s *OrderService) CreateTopic(kafkaBrokers []string) {
	conn, err := kafka.Dial("tcp", kafkaBrokers[0])
	if err != nil {
		Logger.Fatal("kafka connection error", zap.Error(err))
	}
	defer conn.Close()

	topic := "orders"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	controller, err := conn.Controller()
	if err != nil {
		Logger.Fatal("controller receiving error", zap.Error(err))
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		Logger.Fatal("controller connection error", zap.Error(err))
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		Logger.Fatal("topic creating error", zap.Error(err))
	}

	Logger.Info("topic successfully created", zap.String("topic", topic))
}

func (s *OrderService) ConsumeMessages(kafkaBrokers []string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  kafkaBrokers,
		GroupID:  "order_service",
		Topic:    "orders",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			Logger.Warn("message reading error", zap.Error(err))
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			Logger.Warn("message unmarshalling error", zap.Error(err))
			continue
		}

		s.Queries.CreateOrder(s.Ctx, order)
		//cache.OrderCache.Store(order.OrderUID, order)
		Logger.Info("message successfully processed", zap.String("orderUID", order.OrderUID))
	}
}
