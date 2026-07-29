package payments

import (
	"sync"
	"time"
)

type Store interface {
	Create(payment Payment) (Payment, error)
	GetByID(id string) (Payment, error)
	List() ([]Payment, error)
	Update(id string, payment Payment) (Payment, error)
	Delete(id string) error
}

type Service struct {
	store Store
	mu    sync.Mutex
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Create(req CreatePaymentRequest) (Payment, error) {
	if req.UserID == "" || req.Currency == "" || req.Status == "" {
		return Payment{}, ErrInvalidPayment
	}
	if req.Amount <= 0 {
		return Payment{}, ErrInvalidPayment
	}

	now := time.Now().UTC()
	p := Payment{
		ID:        newID(),
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    req.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.store.Create(p)
}

func (s *Service) GetByID(id string) (Payment, error) {
	if id == "" {
		return Payment{}, ErrInvalidPayment
	}
	return s.store.GetByID(id)
}

func (s *Service) List() ([]Payment, error) {
	return s.store.List()
}

func (s *Service) Update(id string, req UpdatePaymentRequest) (Payment, error) {
	if id == "" {
		return Payment{}, ErrInvalidPayment
	}

	existing, err := s.store.GetByID(id)
	if err != nil {
		return Payment{}, err
	}

	if req.UserID != nil {
		if *req.UserID == "" {
			return Payment{}, ErrInvalidPayment
		}
		existing.UserID = *req.UserID
	}
	if req.Amount != nil {
		if *req.Amount <= 0 {
			return Payment{}, ErrInvalidPayment
		}
		existing.Amount = *req.Amount
	}
	if req.Currency != nil {
		if *req.Currency == "" {
			return Payment{}, ErrInvalidPayment
		}
		existing.Currency = *req.Currency
	}
	if req.Status != nil {
		if *req.Status == "" {
			return Payment{}, ErrInvalidPayment
		}
		existing.Status = *req.Status
	}

	existing.UpdatedAt = time.Now().UTC()
	return s.store.Update(id, existing)
}

func (s *Service) Delete(id string) error {
	if id == "" {
		return ErrInvalidPayment
	}
	return s.store.Delete(id)
}
