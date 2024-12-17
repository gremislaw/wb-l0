package service

import (
	"order_service/internal/config"
	"order_service/internal/models"
	"order_service/internal/repository/cache"
	"order_service/internal/repository/db"
)

type CacheService struct {
	config    *config.Config
	dbRepo    *db_repository.Repository
	cacheRepo *cache_repository.Repository
}

func NewOrderCacheService(cfg *config.Config, dbRepo *db_repository.Repository, cacheRepo *cache_repository.Repository) *CacheService {
	return &CacheService{
		config:    cfg,
		cacheRepo: cacheRepo,
		dbRepo:    dbRepo,
	}
}

func (c *CacheService) MigrateFromDB() error {
	orders, err := c.dbRepo.Orders.GetOrders()
	if err != nil {
		return err
	}

	err = c.cacheRepo.Orders.SetOrders(orders)
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheService) SetOrder(order *models.Order) error {
	if order.OrderUID == "" {
		return models.ErrEmptyOrderID
	}

	// запись в кэш
	if err := c.cacheRepo.Orders.SetOrder(order); err != nil {
		return err
	}

	// запись в бд
	if err := c.dbRepo.Orders.SetOrder(order); err != nil {
		return err
	}

	return nil
}

func (c *CacheService) GetOrder(id string) (*models.Order, error) {
	return c.cacheRepo.Orders.GetOrder(id)
}