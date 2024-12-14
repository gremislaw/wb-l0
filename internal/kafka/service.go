package kafka

import (
	"context"
	"database/sql"
	"order_service/internal/repository"
)

type OrderService struct {
	Queries *repository.Queries
	Ctx     context.Context
}

func NewOrderService(queries *repository.Queries, ctx context.Context) *OrderService {
	return &OrderService{Queries: queries, Ctx: ctx}
}

func CreateKafkaConsumer(db *sql.DB, kafkaBrokers []string) *OrderService {
	// создание REST сервиса
	ctx := context.Background()
	queries := repository.New(db)
	orderSvc := NewOrderService(queries, ctx)
	orderSvc.CreateTopic(kafkaBrokers)
	return orderSvc
}
