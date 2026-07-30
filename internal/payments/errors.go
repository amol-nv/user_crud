package payments

import "errors"

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidUserID = errors.New("invalid user id")
)
