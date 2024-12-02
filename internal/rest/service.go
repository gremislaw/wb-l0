package rest

import (
	"database/sql"
	"order_service/internal/models"
	"order_service/internal/repository"
)

type OrderService struct {
	DB *sql.DB
}

func NewOrderService(db *sql.DB) repository.OrderRepository {
	return &OrderService{DB: db}
}

func (s *OrderService) Get(id int) (*models.Order, error) {
	return nil, nil
}