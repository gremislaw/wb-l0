package cache_repository

import (
	"order_service/internal/cache"
	"order_service/internal/models"
)

type Orders interface {
	SetOrder(*models.Order) error
	SetOrders([]*models.Order) error
	GetOrder(string) (*models.Order, error)
	GetOrders() ([]*models.Order, error)
}

type Repository struct {
	Orders
}

func NewRepository(cache *cache.CacheMap) *Repository {
	return &Repository{
		Orders: NewOrdersRepository(cache),
	}
}
