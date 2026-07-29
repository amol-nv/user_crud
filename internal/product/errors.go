package product

import "errors"

var (
	ErrNotFound = errors.New("product not found")
	ErrInvalid  = errors.New("invalid product")
)
