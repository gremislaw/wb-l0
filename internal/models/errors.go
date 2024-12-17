package models

import "errors"

var (
	ErrJSONUnmarshal = errors.New("unexpected content of the JSON message")
	ErrEmptyOrderID = errors.New("empty order id")
)
