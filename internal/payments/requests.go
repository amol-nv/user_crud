package payments

import "time"

type CreatePaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}

type UpdatePaymentRequest struct {
	Amount   *int64  `json:"amount"`
	Currency *string `json:"currency"`
	Status   *string `json:"status"`
}

type PaymentResponse struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
