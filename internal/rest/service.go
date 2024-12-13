package rest

import (
	"context"
	"order_service/internal/repository"
)

type OrderService struct {
	Queries *repository.Queries
	Ctx context.Context
}

func NewOrderService(queries *repository.Queries, ctx context.Context) *OrderService {
	return &OrderService{Queries: queries, Ctx: ctx}
}