package repository

import "order_service/internal/models"

type OrderRepository interface {
	Get(id int) (*models.Order, error)
}