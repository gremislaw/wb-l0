package models

import "time"

type Payment struct {
	Transaction  string    `json:"payment_transaction"`
	RequestID    string    `json:"payment_request_id"`
	Currency     string    `json:"payment_currency"`
	Provider     string    `json:"payment_provider"`
	Amount       float64   `json:"payment_amount"`
	PaymentDt    time.Time `json:"payment_payment_dt"`
	Bank         string    `json:"payment_bank"`
	DeliveryCost float64   `json:"payment_delivery_cost"`
	GoodsTotal   float64   `json:"payment_goods_total"`
	CustomFee    float64   `json:"payment_custom_fee"`
}
