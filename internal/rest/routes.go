package rest

import (
	"net/http"
	"database/sql"
)

func handleRoutes(db *sql.DB) {
	// Обертка для репозитория БД
	productSvc := NewOrderService(db).(*OrderService)
	// endpoint для получения order
	http.HandleFunc("/order/", productSvc.GetOrder)
}

func CreateRestService(db *sql.DB) {
	handleRoutes(db)
}