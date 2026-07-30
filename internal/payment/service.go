package payment

import (
	"context"
	"fmt"
)

// Service is the View/Controller layer for payment.
// It validates input and delegates to Store.
type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

// PostPayment validates and posts a payment.
func (s *Service) PostPayment(ctx context.Context, req PaymentRequest) (PaymentResponse, error) {
	if err := req.Validate(); err != nil {
		return PaymentResponse{}, fmt.Errorf("validate payment request: %w", err)
	}

	resp, _, err := s.store.PostPayment(ctx, req)
	if err != nil {
		return PaymentResponse{}, err
	}
	return resp, nil
}
