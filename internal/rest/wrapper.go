package rest

import (
	"context"
	"database/sql"
	"net/http"
	"order_service/internal/cache"
	"order_service/internal/repository"
)

type OrderService struct {
	Queries *repository.Queries
	Ctx     context.Context
	Cache *cache.CacheWrapper
}

func NewOrderService(queries *repository.Queries, ctx context.Context, cache *cache.CacheWrapper) *OrderService {
	return &OrderService{Queries: queries, Ctx: ctx}
}

func CreateRestService(db *sql.DB, cache *cache.CacheWrapper) *OrderService {
	// создание REST сервиса
	ctx := context.Background()
	queries := repository.New(db)
	orderSvc := NewOrderService(queries, ctx, cache)
	handleRoutes(orderSvc)
	return orderSvc
}

func handleRoutes(orderSvc *OrderService) {
	// endpoint для получения order
	http.HandleFunc("/order/", orderSvc.GetOrder)
}