package payments

import "errors"

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
)
