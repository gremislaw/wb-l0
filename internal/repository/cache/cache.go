package cache_repository

import (
	"errors"
	"order_service/internal/cache"
	"order_service/internal/models"
)

type OrdersRepository struct {
	cache *cache.CacheMap
}

func NewOrdersRepository(cache *cache.CacheMap) *OrdersRepository {
	return &OrdersRepository{
		cache: cache,
	}
}

func (c *OrdersRepository) GetOrders() ([]*models.Order, error) {
	return c.cache.GetAll()
}

func (c *OrdersRepository) GetOrder(uid string) (*models.Order, error) {
	return c.cache.Get(uid)
}

func (c *OrdersRepository) SetOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return models.ErrEmptyOrderID
	}

	c.cache.Set(order)

	return nil
}

func (c *OrdersRepository) SetOrders(orders []*models.Order) error {
	for _, order := range orders {
		if order.OrderUID == "" {
			return errors.New("empty order id")
		}
		c.cache.Set(order)
	}

	return nil
}