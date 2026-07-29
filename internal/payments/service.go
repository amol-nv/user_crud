package payments

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	if err := validateCreate(req); err != nil {
		return Payment{}, err
	}

	id := newPaymentID()
	p := Payment{
		ID:       id,
		Amount:   req.Amount,
		Currency: strings.ToUpper(strings.TrimSpace(req.Currency)),
		Status:   strings.TrimSpace(req.Status),
	}
	return s.repo.Create(ctx, p)
}

func (s *PaymentService) GetByID(ctx context.Context, id string) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidArgument
	}
	return s.repo.GetByID(ctx, id)
}

func (s *PaymentService) Update(ctx context.Context, id string, req UpdatePaymentRequest) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidArgument
	}
	if err := validateUpdate(req); err != nil {
		return Payment{}, err
	}

	p := Payment{
		ID:       id,
		Amount:   req.Amount,
		Currency: strings.ToUpper(strings.TrimSpace(req.Currency)),
		Status:   strings.TrimSpace(req.Status),
	}
	return s.repo.Update(ctx, id, p)
}

func (s *PaymentService) Delete(ctx context.Context, id string) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidArgument
	}
	return s.repo.Delete(ctx, id)
}

func validateCreate(req CreatePaymentRequest) error {
	if req.Amount <= 0 {
		return ErrInvalidArgument
	}
	if strings.TrimSpace(req.Currency) == "" {
		return ErrInvalidArgument
	}
	if strings.TrimSpace(req.Status) == "" {
		return ErrInvalidArgument
	}
	return nil
}

func validateUpdate(req UpdatePaymentRequest) error {
	if req.Amount <= 0 {
		return ErrInvalidArgument
	}
	if strings.TrimSpace(req.Currency) == "" {
		return ErrInvalidArgument
	}
	if strings.TrimSpace(req.Status) == "" {
		return ErrInvalidArgument
	}
	return nil
}

func newPaymentID() string {
	// Simple deterministic-ish ID for local development.
	return fmt.Sprintf("pay_%d", time.Now().UTC().UnixNano())
}
