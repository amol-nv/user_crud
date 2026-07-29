package inventory

import "fmt"

type NotFoundError struct {
	ID int64
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("inventory item %d not found", e.ID)
}

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string {
	return e.Msg
}
