package wishlist

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrProductNotFound = errors.New("product not found")
)
