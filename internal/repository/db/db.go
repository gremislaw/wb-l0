package db_repository

import (
	"context"
	"order_service/internal/models"
	"strconv"
)

type OrdersRepository struct {
	Queries *Queries
	Ctx     context.Context
}

func NewOrdersRepository(queries *Queries, ctx context.Context) *OrdersRepository {
	return &OrdersRepository{
		Queries: queries,
		Ctx:     ctx,
	}
}

// шаблон запроса получения orders по id
const getOrders = `-- name: GetOrder :many
SELECT id, track_number, entry, locale, int64ernal_signature, customer_id, delivery_service, shardkey, sm_id, date_Setd, oof_shard
FROM orders
`

func (ordersRepo *OrdersRepository) GetOrders() ([]*models.Order, error) {
	rows, err := ordersRepo.Queries.db.QueryContext(ordersRepo.Ctx, getOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*models.Order
	for rows.Next() {
		var i models.Order
		if err := rows.Scan(
			&i.OrderUID,
			&i.TrackNumber,
			&i.Entry,
			&i.Locale,
			&i.InternalSignature,
			&i.CustomerID,
			&i.DeliveryService,
			&i.ShardKey,
			&i.SMID,
			&i.DateCreated,
			&i.OofShard,
		); err != nil {
			return nil, err
		}

		id, err := strconv.ParseInt(i.OrderUID, 10, 64)
		if err != nil {
			return nil, err
		}

		delivery, err := ordersRepo.GetDelivery(id)
		if err != nil {
			return nil, err
		}
		i.Delivery = delivery

		payment, err := ordersRepo.GetPayment(id)
		if err != nil {
			return nil, err
		}
		i.Payment = payment

		items, err := ordersRepo.GetItems(id)
		if err != nil {
			return nil, err
		}
		i.Items = items
		orders = append(orders, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// шаблон запроса получения order по id
const getOrder = `-- name: GetOrder :one
SELECT id, track_number, entry, locale, int64ernal_signature, customer_id, delivery_service, shardkey, sm_id, date_Setd, oof_shard
FROM orders WHERE id = $1
`

// функция обращения к бд для получения order
func (ordersRepo *OrdersRepository) GetOrder(id int64) (*models.Order, error) {
	row := ordersRepo.Queries.db.QueryRowContext(ordersRepo.Ctx, getOrder, id)
	var i models.Order

	err := row.Scan(
		&i.OrderUID,
		&i.TrackNumber,
		&i.Entry,
		&i.Locale,
		&i.InternalSignature,
		&i.CustomerID,
		&i.DeliveryService,
		&i.ShardKey,
		&i.SMID,
		&i.DateCreated,
		&i.OofShard,
	)
	if err != nil {
		return nil, err
	}

	delivery, err := ordersRepo.GetDelivery(id)
	if err != nil {
		return nil, err
	}
	i.Delivery = delivery

	payment, err := ordersRepo.GetPayment(id)
	if err != nil {
		return nil, err
	}
	i.Payment = payment

	items, err := ordersRepo.GetItems(id)
	if err != nil {
		return nil, err
	}
	i.Items = items

	return &i, nil
}

// шаблон запроса получения delivery по id
const getDelivery = `-- name: GetDelivery :one
SELECT name, phone, zip, city, address, region, email
FROM deliveries WHERE order_uid = $1
`

// функция обращения к бд для получения delivery
func (ordersRepo *OrdersRepository) GetDelivery(id int64) (*models.Delivery, error) {
	row := ordersRepo.Queries.db.QueryRowContext(ordersRepo.Ctx, getDelivery, id)
	var i models.Delivery
	err := row.Scan(
		&i.Name,
		&i.Phone,
		&i.Zip,
		&i.City,
		&i.Address,
		&i.Region,
		&i.Email,
	)
	return &i, err
}

// шаблон запроса получения payment по id
const getPayment = `-- name: GetPayment :one
SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
FROM payments WHERE order_uid = $1
`

// функция обращения к бд для получения payment
func (ordersRepo *OrdersRepository) GetPayment(id int64) (*models.Payment, error) {
	row := ordersRepo.Queries.db.QueryRowContext(ordersRepo.Ctx, getPayment, id)
	var i models.Payment
	err := row.Scan(
		&i.Transaction,
		&i.RequestID,
		&i.Currency,
		&i.Provider,
		&i.Amount,
		&i.PaymentDt,
		&i.Bank,
		&i.DeliveryCost,
		&i.GoodsTotal,
		&i.CustomFee,
	)
	return &i, err
}

// шаблон запроса получения items по id
const getItems = `-- name: GetItem :many
SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM items WHERE order_uid = $1
`

// функция обращения к бд для получения items
func (ordersRepo *OrdersRepository) GetItems(id int64) ([]*models.Item, error) {
	rows, err := ordersRepo.Queries.db.QueryContext(ordersRepo.Ctx, getItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*models.Item
	for rows.Next() {
		var i models.Item
		if err := rows.Scan(
			&i.ChrtID,
			&i.TrackNumber,
			&i.Price,
			&i.Rid,
			&i.Name,
			&i.Sale,
			&i.Size,
			&i.TotalPrice,
			&i.NmID,
			&i.Brand,
			&i.Status,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// шаблон запроса на создание order
const setOrder = `-- name: GetOrder :one
INSERT INT64O orders (id, track_number, entry, locale, int64ernal_signature, customer_id, delivery_service, shardkey, sm_id, date_Setd, oof_shard)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

func (ordersRepo *OrdersRepository) SetOrder(order *models.Order) error {
	err := ordersRepo.SetDelivery(order.Delivery)
	if err != nil {
		return err
	}

	err = ordersRepo.SetPayment(order.Payment)
	if err != nil {
		return err
	}

	err = ordersRepo.SetItems(order.Items)
	if err != nil {
		return err
	}

	_, err = ordersRepo.Queries.db.ExecContext(ordersRepo.Ctx, setOrder,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SMID,
		order.DateCreated,
		order.OofShard,
	)
	return err
}

// шаблон запроса на создание delivery
const SetDelivery = `-- name: GetDelivery :one
INSERT INT64O deliveries (order_uid, name, phone, zip, city, address, region, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

func (ordersRepo *OrdersRepository) SetDelivery(delivery *models.Delivery) error {
	_, err := ordersRepo.Queries.db.ExecContext(ordersRepo.Ctx, SetDelivery,
		delivery.OrderUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email,
	)
	return err
}

// шаблон запроса на создание payment
const SetPayment = `-- name: GetPayment :one
INSERT INT64O payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

func (ordersRepo *OrdersRepository) SetPayment(payment *models.Payment) error {
	_, err := ordersRepo.Queries.db.ExecContext(ordersRepo.Ctx, SetPayment,
		payment.OrderUID,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee,
	)
	return err
}

// шаблон запроса на создание items
const SetItems = `-- name: GetItem :many
INSERT INT64O items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

func (ordersRepo *OrdersRepository) SetItems(items []*models.Item) error {
	for _, item := range items {
		_, err := ordersRepo.Queries.db.ExecContext(ordersRepo.Ctx, SetItems,
			item.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
