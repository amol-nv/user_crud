package payment

import (
	"context"
	"strings"
	"time"

	"amol-nv/user_crud/internal/store"
)

type PaymentRepository interface {
	Create(ctx context.Context, p Payment) (Payment, error)
	GetByID(ctx context.Context, id string) (Payment, error)
	Update(ctx context.Context, id string, p Payment) (Payment, error)
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo PaymentRepository
	idFn func() string
}

func NewService(repo PaymentRepository) *Service {
	return &Service{repo: repo, idFn: store.NewID}
}

func (s *Service) Create(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	if strings.TrimSpace(req.UserID) == "" || req.Amount <= 0 || strings.TrimSpace(req.Currency) == "" || strings.TrimSpace(req.Status) == "" {
		return Payment{}, ErrInvalidPayment
	}

	now := time.Now().UTC()
	p := Payment{
		ID:        s.idFn(),
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  strings.ToUpper(strings.TrimSpace(req.Currency)),
		Status:    strings.ToUpper(strings.TrimSpace(req.Status)),
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.Create(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id string) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidPayment
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id string, req UpdatePaymentRequest) (Payment, error) {
	if strings.TrimSpace(id) == "" {
		return Payment{}, ErrInvalidPayment
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Payment{}, err
	}

	if req.UserID != nil {
		if strings.TrimSpace(*req.UserID) == "" {
			return Payment{}, ErrInvalidPayment
		}
		current.UserID = *req.UserID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return Payment{}, ErrInvalidPayment
		}
		current.Amount = *req.Amount
	}
	if req.Currency != nil {
		if strings.TrimSpace(*req.Currency) == "" {
			return Payment{}, ErrInvalidPayment
		}
		current.Currency = strings.ToUpper(strings.TrimSpace(*req.Currency))
	}
	if req.Status != nil {
		if strings.TrimSpace(*req.Status) == "" {
			return Payment{}, ErrInvalidPayment
		}
		current.Status = strings.ToUpper(strings.TrimSpace(*req.Status))
	}

	current.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, id, current)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidPayment
	}
	return s.repo.Delete(ctx, id)
}
