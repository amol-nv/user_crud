package payments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, p Payment) (Payment, error)
	GetByID(ctx context.Context, id string) (Payment, error)
	UpdateStatus(ctx context.Context, id string, status PaymentStatus) (Payment, error)
}

type Provider interface {
	// Charge attempts to charge the payment.
	// For this ticket we keep it minimal and synchronous.
	Charge(ctx context.Context, p Payment) (PaymentStatus, error)
}

type Service interface {
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (CreatePaymentResponse, error)
	GetPayment(ctx context.Context, id string) (Payment, error)
}

type service struct {
	repo      PaymentRepository
	provider  Provider
	nowFn     func() time.Time
	idGen     func() string
	minAmount int64
}

func NewService(repo PaymentRepository, provider Provider) Service {
	return &service{
		repo:      repo,
		provider:  provider,
		nowFn:     time.Now,
		idGen:     func() string { return uuid.NewString() },
		minAmount: 1,
	}
}

func (s *service) CreatePayment(ctx context.Context, req CreatePaymentRequest) (CreatePaymentResponse, error) {
	if strings.TrimSpace(req.UserID) == "" {
		return CreatePaymentResponse{}, fmt.Errorf("%w: user_id is required", ErrInvalidAmount)
	}
	if strings.TrimSpace(req.OrderID) == "" {
		return CreatePaymentResponse{}, fmt.Errorf("%w: order_id is required", ErrInvalidAmount)
	}
	if req.Amount < s.minAmount {
		return CreatePaymentResponse{}, ErrInvalidAmount
	}
	cur := strings.ToUpper(strings.TrimSpace(req.Currency))
	if cur == "" {
		return CreatePaymentResponse{}, ErrInvalidCurrency
	}

	now := s.nowFn()
	p := Payment{
		ID:        s.idGen(),
		UserID:    req.UserID,
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Currency:  cur,
		Status:    PaymentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repo.Create(ctx, p)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	status, err := s.provider.Charge(ctx, created)
	if err != nil {
		// Keep payment record as pending/failed depending on provider error.
		_, _ = s.repo.UpdateStatus(ctx, created.ID, PaymentStatusFailed)
		return CreatePaymentResponse{}, err
	}

	updated, err := s.repo.UpdateStatus(ctx, created.ID, status)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	return CreatePaymentResponse{Payment: updated}, nil
}

func (s *service) GetPayment(ctx context.Context, id string) (Payment, error) {
	return s.repo.GetByID(ctx, id)
}
