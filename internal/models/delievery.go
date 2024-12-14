package models

type Delivery struct {
	OrderUID int    `json:"order_uid"`
	Name     string `json:"delivery_name"`
	Phone    string `json:"delivery_phone"`
	Zip      string `json:"delivery_zip"`
	City     string `json:"delivery_city"`
	Address  string `json:"delivery_address"`
	Region   string `json:"delivery_region"`
	Email    string `json:"delivery_email"`
}
