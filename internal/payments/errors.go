package payments

import "errors"

var (
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrPaymentNotFound = errors.New("payment not found")
)
