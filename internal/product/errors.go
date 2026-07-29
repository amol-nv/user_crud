package product

import "errors"

var (
	ErrNotFound      = errors.New("product not found")
	ErrValidation   = errors.New("validation error")
	ErrInvalidInput = errors.New("invalid input")
)
