package payments

type CreatePaymentRequest struct {
	UserID   string `json:"userId"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}
