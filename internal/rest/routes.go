package rest

import (
	"context"
	"database/sql"
	"net/http"
	"order_service/internal/repository"
)

func handleRoutes(orderSvc *OrderService) {
	// endpoint для получения order
	http.HandleFunc("/order/", orderSvc.GetOrder)
}

func CreateRestService(db *sql.DB) {
	ctx := context.Background()
	queries := repository.New(db)
	orderSvc := NewOrderService(queries, ctx)
	handleRoutes(orderSvc)
}
