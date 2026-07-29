package payments

import "time"

type CreatePaymentRequest struct {
	UserID   string  `json:"userId"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type UpdatePaymentRequest struct {
	UserID   *string  `json:"userId,omitempty"`
	Amount   *float64 `json:"amount,omitempty"`
	Currency *string  `json:"currency,omitempty"`
	Status   *string  `json:"status,omitempty"`
}

type PaymentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
