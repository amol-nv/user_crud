package payment

import "errors"

var (
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrInvalidMethod   = errors.New("invalid method")
	ErrInvalidOrderID  = errors.New("invalid order_id")
)

// Validate validates a PaymentRequest.
func (r PaymentRequest) Validate() error {
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if r.Currency == "" {
		return ErrInvalidCurrency
	}
	if r.Method == "" {
		return ErrInvalidMethod
	}
	if r.OrderID == "" {
		return ErrInvalidOrderID
	}
	return nil
}
