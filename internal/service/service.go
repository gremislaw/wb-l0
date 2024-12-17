package service

import (
	"order_service/internal/config"
	"order_service/internal/models"
	"order_service/internal/repository/cache"
	"order_service/internal/repository/db"
)

type Orders interface {
	GetOrder(string) (*models.Order, error)
	SetOrder(*models.Order) error
	MigrateFromDB() error
}

type Service struct {
	Orders
}

type Deployments struct {
	Config    *config.Config
	DbRepo    *db_repository.Repository
	CacheRepo *cache_repository.Repository
}

func NewOrderService(d Deployments) *Service {
	return &Service{
		Orders: NewOrderCacheService(d.Config, d.DbRepo, d.CacheRepo),
	}
}