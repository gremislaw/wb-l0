package repository

import (
	"context"
	"order_service/internal/models"
	"strconv"
)

// шаблон запроса получения orders по id
const getOrders = `-- name: GetOrder :many
SELECT id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders
`

func (q *Queries) GetOrders(ctx context.Context) ([]models.Order, error) {
	rows, err := q.db.QueryContext(ctx, getOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order
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

		id, err := strconv.Atoi(i.OrderUID)
		if err != nil {
			return nil, err
		}

		delivery, err := q.GetDelivery(ctx, id)
		if err != nil {
			return nil, err
		}
		i.Delivery = delivery

		payment, err := q.GetPayment(ctx, id)
		if err != nil {
			return nil, err
		}
		i.Payment = payment

		items, err := q.GetItems(ctx, id)
		if err != nil {
			return nil, err
		}
		i.Items = items
		orders = append(orders, i)
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
SELECT id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders WHERE id = $1
`

// функция обращения к бд для получения order
func (q *Queries) GetOrder(ctx context.Context, id int) (models.Order, error) {
	row := q.db.QueryRowContext(ctx, getOrder, id)
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
		return i, err
	}

	delivery, err := q.GetDelivery(ctx, id)
	if err != nil {
		return i, err
	}
	i.Delivery = delivery

	payment, err := q.GetPayment(ctx, id)
	if err != nil {
		return i, err
	}
	i.Payment = payment

	items, err := q.GetItems(ctx, id)
	if err != nil {
		return i, err
	}
	i.Items = items

	return i, nil
}

// шаблон запроса получения delivery по id
const getDelivery = `-- name: GetDelivery :one
SELECT name, phone, zip, city, address, region, email
FROM deliveries WHERE order_uid = $1
`

// функция обращения к бд для получения delivery
func (q *Queries) GetDelivery(ctx context.Context, id int) (models.Delivery, error) {
	row := q.db.QueryRowContext(ctx, getDelivery, id)
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
	return i, err
}

// шаблон запроса получения payment по id
const getPayment = `-- name: GetPayment :one
SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
FROM payments WHERE order_uid = $1
`

// функция обращения к бд для получения payment
func (q *Queries) GetPayment(ctx context.Context, id int) (models.Payment, error) {
	row := q.db.QueryRowContext(ctx, getPayment, id)
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
	return i, err
}

// шаблон запроса получения items по id
const getItems = `-- name: GetItem :many
SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM items WHERE order_uid = $1
`

// функция обращения к бд для получения items
func (q *Queries) GetItems(ctx context.Context, id int) ([]models.Item, error) {
	rows, err := q.db.QueryContext(ctx, getItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
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
		items = append(items, i)
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
const createOrder = `-- name: GetOrder :one
INSERT INTO orders (id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

func (q *Queries) CreateOrder(ctx context.Context, order models.Order) error {
	err := q.CreateDelivery(ctx, order.Delivery)
	if err != nil {
		return err
	}

	err = q.CreatePayment(ctx, order.Payment)
	if err != nil {
		return err
	}

	err = q.CreateItems(ctx, order.Items)
	if err != nil {
		return err
	}

	_, err = q.db.ExecContext(ctx, createOrder,
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
const createDelivery = `-- name: GetDelivery :one
INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

func (q *Queries) CreateDelivery(ctx context.Context, delivery models.Delivery) error {
	_, err := q.db.ExecContext(ctx, createDelivery,
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
const createPayment = `-- name: GetPayment :one
INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

func (q *Queries) CreatePayment(ctx context.Context, payment models.Payment) error {
	_, err := q.db.ExecContext(ctx, createPayment,
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
const createItems = `-- name: GetItem :many
INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

func (q *Queries) CreateItems(ctx context.Context, items []models.Item) error {
	for _, item := range items {
		_, err := q.db.ExecContext(ctx, createItems,
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
