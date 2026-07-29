package payments

import (
	"context"
	"strings"
	"time"
)

type PaymentsRepository interface {
	Create(ctx context.Context, p Payment) (Payment, error)
	GetByID(ctx context.Context, id string) (Payment, error)
	Update(ctx context.Context, id string, p Payment) (Payment, error)
	Delete(ctx context.Context, id string) error
}

type PaymentsService struct {
	repo PaymentsRepository
}

func NewService(repo PaymentsRepository) *PaymentsService {
	return &PaymentsService{repo: repo}
}

func (s *PaymentsService) Create(ctx context.Context, req CreatePaymentRequest) (PaymentResponse, error) {
	if err := validateCreate(req); err != nil {
		return PaymentResponse{}, err
	}

	status := PaymentStatus(strings.ToLower(req.Status))
	if status == "" {
		status = PaymentStatusPending
	}

	now := time.Now().UTC()
	p := Payment{
		ID:        "", // repository will assign
		Amount:    req.Amount,
		Currency:  strings.ToUpper(req.Currency),
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(created), nil
}

func (s *PaymentsService) GetByID(ctx context.Context, id string) (PaymentResponse, error) {
	if strings.TrimSpace(id) == "" {
		return PaymentResponse{}, ErrInvalidPayment
	}
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(p), nil
}

func (s *PaymentsService) Update(ctx context.Context, id string, req UpdatePaymentRequest) (PaymentResponse, error) {
	if strings.TrimSpace(id) == "" {
		return PaymentResponse{}, ErrInvalidPayment
	}
	if err := validateUpdate(req); err != nil {
		return PaymentResponse{}, err
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return PaymentResponse{}, err
	}

	updated := current
	if req.Amount != nil {
		updated.Amount = *req.Amount
	}
	if req.Currency != nil {
		updated.Currency = strings.ToUpper(*req.Currency)
	}
	if req.Status != nil {
		updated.Status = PaymentStatus(strings.ToLower(*req.Status))
	}
	updated.UpdatedAt = time.Now().UTC()

	p, err := s.repo.Update(ctx, id, updated)
	if err != nil {
		return PaymentResponse{}, err
	}
	return toPaymentResponse(p), nil
}

func (s *PaymentsService) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidPayment
	}
	return s.repo.Delete(ctx, id)
}

func validateCreate(req CreatePaymentRequest) error {
	if req.Amount <= 0 {
		return ErrInvalidPayment
	}
	if strings.TrimSpace(req.Currency) == "" {
		return ErrInvalidPayment
	}
	// status is optional; if provided, allow any non-empty string
	if strings.TrimSpace(req.Status) != "" {
		if strings.TrimSpace(req.Status) == "" {
			return ErrInvalidPayment
		}
	}
	return nil
}

func validateUpdate(req UpdatePaymentRequest) error {
	if req.Amount != nil && *req.Amount <= 0 {
		return ErrInvalidPayment
	}
	if req.Currency != nil && strings.TrimSpace(*req.Currency) == "" {
		return ErrInvalidPayment
	}
	if req.Status != nil && strings.TrimSpace(*req.Status) == "" {
		return ErrInvalidPayment
	}
	return nil
}

func toPaymentResponse(p Payment) PaymentResponse {
	return PaymentResponse{
		ID:        p.ID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    string(p.Status),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
