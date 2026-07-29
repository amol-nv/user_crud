package payments

import "errors"

var (
	ErrNotFound        = errors.New("payment not found")
	ErrInvalidArgument = errors.New("invalid payment argument")
)
