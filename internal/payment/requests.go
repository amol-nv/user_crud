package payment

import "time"

type CreatePaymentRequest struct {
	UserID   string `json:"userId"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

type UpdatePaymentRequest struct {
	UserID   *string `json:"userId"`
	Amount   *int64  `json:"amount"`
	Currency *string `json:"currency"`
	Status   *string `json:"status"`
}

type PaymentResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
