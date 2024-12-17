package db_repository

import (
	"context"
	"database/sql"
	"order_service/internal/models"
)

type Orders interface {
	SetOrder(*models.Order) error
	GetOrders() ([]*models.Order, error)
	GetOrder(int64) (*models.Order, error)
	GetPayment(int64) (*models.Payment, error)
	GetItems(int64) ([]*models.Item, error)
	GetDelivery(int64) (*models.Delivery, error)
}

type Repository struct {
	Orders
}

func NewRepository(db *sql.DB, ctx context.Context) *Repository {
	queries := NewQueries(db)
	return &Repository{
		Orders: NewOrdersRepository(queries, ctx),
	}
}
