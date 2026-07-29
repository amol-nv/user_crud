package payments

import (
	"time"

	"amol-nv/user_crud/internal/store"
)

type Service struct {
	store store.PaymentStore
}

func NewService(s store.PaymentStore) *Service {
	return &Service{store: s}
}

type CreatePaymentInput struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type UpdatePaymentInput struct {
	Amount   *float64 `json:"amount"`
	Currency *string  `json:"currency"`
	Status   *string  `json:"status"`
}

func (s *Service) CreatePayment(in CreatePaymentInput) (*Payment, error) {
	if in.Currency == "" {
		return nil, ErrInvalidPayment
	}
	if in.Status == "" {
		in.Status = "created"
	}
	p := &Payment{
		ID:        "", // store will assign
		Amount:    in.Amount,
		Currency:  in.Currency,
		Status:    in.Status,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	created, err := s.store.Create(p)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) GetPayment(id string) (*Payment, error) {
	p, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListPayments() ([]*Payment, error) {
	return s.store.List()
}

func (s *Service) UpdatePayment(id string, in UpdatePaymentInput) (*Payment, error) {
	p, err := s.store.GetByID(id)
	if err != nil {
		return nil, err
	}
	if in.Amount != nil {
		p.Amount = *in.Amount
	}
	if in.Currency != nil {
		if *in.Currency == "" {
			return nil, ErrInvalidPayment
		}
		p.Currency = *in.Currency
	}
	if in.Status != nil {
		if *in.Status == "" {
			return nil, ErrInvalidPayment
		}
		p.Status = *in.Status
	}
	p.UpdatedAt = time.Now().UTC()
	updated, err := s.store.Update(p)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *Service) DeletePayment(id string) error {
	return s.store.Delete(id)
}
