package payments

import "time"

type CreatePaymentRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type CreatePaymentResponse struct {
	Payment Payment `json:"payment"`
}

type GetPaymentResponse struct {
	Payment Payment `json:"payment"`
}

type UpdatePaymentRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type UpdatePaymentResponse struct {
	Payment Payment `json:"payment"`
}

type DeletePaymentResponse struct {
	DeletedID string    `json:"deletedId"`
	DeletedAt time.Time `json:"deletedAt"`
}
