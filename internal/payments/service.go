package payments

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, req CreatePaymentRequest) (PaymentResponse, error) {
	if err := validateCreate(req); err != nil {
		return PaymentResponse{}, err
	}

	now := time.Now().UTC()
	p := Payment{
		ID:        newID(),
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    "created",
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.store.Create(ctx, p)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(created), nil
}

func (s *Service) GetByID(ctx context.Context, id string) (PaymentResponse, error) {
	if id == "" {
		return PaymentResponse{}, ErrInvalid
	}

	p, err := s.store.GetByID(ctx, id)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(p), nil
}

func (s *Service) Update(ctx context.Context, id string, req UpdatePaymentRequest) (PaymentResponse, error) {
	if id == "" {
		return PaymentResponse{}, ErrInvalid
	}
	if err := validateUpdate(req); err != nil {
		return PaymentResponse{}, err
	}

	p := Payment{
		ID:       id,
		UserID:   req.UserID,
		Amount:   req.Amount,
		Currency: req.Currency,
		Status:   req.Status,
	}

	updated, err := s.store.Update(ctx, id, p)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(updated), nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalid
	}
	return s.store.Delete(ctx, id)
}

func validateCreate(req CreatePaymentRequest) error {
	if req.UserID == "" {
		return ErrInvalid
	}
	if req.Amount <= 0 {
		return ErrInvalid
	}
	if req.Currency == "" {
		return ErrInvalid
	}
	return nil
}

func validateUpdate(req UpdatePaymentRequest) error {
	if req.UserID == "" {
		return ErrInvalid
	}
	if req.Amount <= 0 {
		return ErrInvalid
	}
	if req.Currency == "" {
		return ErrInvalid
	}
	if req.Status == "" {
		return ErrInvalid
	}
	return nil
}

func toPaymentResponse(p Payment) PaymentResponse {
	return PaymentResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func newID() string {
	// Simple deterministic-ish ID without external deps.
	return fmt.Sprintf("pay_%d", time.Now().UnixNano())
}
