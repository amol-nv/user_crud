package payments

import "errors"

var (
	ErrNotFound = errors.New("payment not found")
	ErrInvalid  = errors.New("invalid payment")
)
