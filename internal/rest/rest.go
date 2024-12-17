package rest

import (
	"net/http"
	"order_service/internal/repository/cache"
)

type RestService struct {
	cache *cache_repository.Repository
}

func NewRestService(cache *cache_repository.Repository) *RestService {
	return &RestService{cache: cache}
}

func CreateRestService(cache *cache_repository.Repository) *RestService {
	// создание REST сервиса
	orderSvc := NewRestService(cache)
	handleRoutes(orderSvc)
	return orderSvc
}

func handleRoutes(orderSvc *RestService) {
	// endpoint для получения order
	http.HandleFunc("/order/", orderSvc.GetOrder)
}
