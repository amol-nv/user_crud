package payments

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID        string
	UserID    string
	OrderID   string
	Amount    int64
	Currency  string
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePaymentRequest struct {
	UserID   string
	OrderID  string
	Amount   int64
	Currency string
}

type CreatePaymentResponse struct {
	Payment Payment
}
