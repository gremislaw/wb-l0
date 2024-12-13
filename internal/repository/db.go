package repository

import (
	"context"
	"order_service/internal/models"
)

const getOrder = `-- name: GetOrder :one
SELECT id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
FROM orders WHERE id = $1
`

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


const getDelivery = `-- name: GetDelivery :one
SELECT name, phone, zip, city, address, region, email
FROM deliveries WHERE order_uid = $1
`

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


const getPayment = `-- name: GetPayment :one
SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
FROM payments WHERE order_uid = $1
`

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


const getItems = `-- name: GetItem :many
SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
FROM items WHERE order_uid = $1
`

func (q *Queries) GetItems(ctx context.Context, id int) ([]models.Item, error) {
	rows, err := q.db.QueryContext(ctx, getItems, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.Item
	for rows.Next() {
		var i  models.Item
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