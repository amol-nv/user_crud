package users

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrValidation   = errors.New("validation error")
	ErrInvalidID    = errors.New("invalid id")
)
