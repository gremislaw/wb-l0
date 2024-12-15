package kafka

import (
	"context"
	"database/sql"
	"order_service/internal/cache"
	"order_service/internal/repository"
)

type ConsumerWrapper struct {
	Queries      *repository.Queries
	Ctx          context.Context
	KafkaBrokers []string
	Cache        *cache.CacheWrapper
}

func NewConsumer(queries *repository.Queries, ctx context.Context, kafkaBrokers []string, cache *cache.CacheWrapper) *ConsumerWrapper {
	return &ConsumerWrapper{Queries: queries, Ctx: ctx, KafkaBrokers: kafkaBrokers, Cache: cache}
}

func CreateConsumerWrapper(db *sql.DB, kafkaBrokers []string, cache *cache.CacheWrapper) *ConsumerWrapper {
	// создание обертки для consumer
	ctx := context.Background()
	queries := repository.New(db)
	orderSvc := NewConsumer(queries, ctx, kafkaBrokers, cache)
	orderSvc.CreateTopic()
	return orderSvc
}
