package payment

import "errors"

var (
	ErrInvalidPayment = errors.New("invalid payment")
	ErrNotFound        = errors.New("payment not found")
)
