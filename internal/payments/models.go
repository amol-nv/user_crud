package payments

import "time"

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

type Payment struct {
	ID        string
	UserID    string
	Amount    int64
	Currency  string
	Status    PaymentStatus
	CreatedAt time.Time
}
