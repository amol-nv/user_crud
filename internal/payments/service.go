package payments

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type PaymentRepository interface {
	Create(ctx context.Context, p Payment) (Payment, error)
}

type Service interface {
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (Payment, error)
}

type service struct {
	repo PaymentRepository
}

func NewService(repo PaymentRepository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePayment(ctx context.Context, req CreatePaymentRequest) (Payment, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return Payment{}, ErrInvalidUserID
	}
	if req.Amount <= 0 {
		return Payment{}, ErrInvalidAmount
	}
	cur := strings.ToUpper(strings.TrimSpace(req.Currency))
	if cur == "" {
		return Payment{}, ErrInvalidCurrency
	}

	p := Payment{
		ID:        fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  cur,
		Status:    PaymentStatusCompleted,
		CreatedAt: time.Now().UTC(),
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return Payment{}, err
	}
	return created, nil
}
