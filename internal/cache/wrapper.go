package cache

import (
	"context"
	"database/sql"
	"order_service/internal/models"
	"order_service/internal/repository"
	"sync"
)

type CacheWrapper struct {
	Queries *repository.Queries
	Ctx     context.Context
	orders  map[string]*models.Order
	mu      *sync.Mutex
}

func NewCacheWrapper(queries *repository.Queries, ctx context.Context, ordersMap map[string]*models.Order, mu *sync.Mutex) *CacheWrapper {
	return &CacheWrapper{Queries: queries, Ctx: ctx, orders: ordersMap, mu: mu}
}

func CreateCache(db *sql.DB) *CacheWrapper {
	// создание обертки для кэша
	ctx := context.Background()
	queries := repository.New(db)
	orders, err := queries.GetOrders(ctx)
	if err != nil {
		panic(err)
	}
	ordersMap := make(map[string]*models.Order)

	for _, order := range orders {
		ordersMap[order.OrderUID] = &order
	}

	cacheWrapper := NewCacheWrapper(queries, ctx, ordersMap, &sync.Mutex{})

	return cacheWrapper
}
