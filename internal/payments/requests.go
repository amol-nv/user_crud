package payments

type CreatePaymentRequest struct {
	UserID   string `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type UpdatePaymentRequest struct {
	UserID   string `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
}
