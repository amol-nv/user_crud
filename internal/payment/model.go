package payment

// Payment POST request/response models.
//
// This package follows an MVS-style separation used elsewhere in the repo:
// - Model: request/response DTOs
// - Store: HTTP client implementation
// - View/Controller: handler/service that validates and orchestrates

// PaymentRequest is the payload sent to the payment provider.
// Adjust fields to match the provider contract.
type PaymentRequest struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Method   string `json:"method"`
	OrderID  string `json:"order_id"`
}

// PaymentResponse is the payload returned by the payment provider.
type PaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
	Message   string `json:"message,omitempty"`
}
