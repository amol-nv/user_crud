package storage

import (
	"context"
	"sync"
	"time"

	"amol-nv/user_crud/internal/payments"
)

type InMemoryPaymentStore struct {
	mu       sync.RWMutex
	payments map[string]payments.Payment
}

func NewInMemoryPaymentStore() *InMemoryPaymentStore {
	return &InMemoryPaymentStore{payments: make(map[string]payments.Payment)}
}

func (s *InMemoryPaymentStore) Create(ctx context.Context, p payments.Payment) (payments.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Basic validation (service should validate too).
	if p.ID == "" || p.UserID == "" || p.Amount <= 0 || p.Currency == "" {
		return payments.Payment{}, payments.ErrInvalid
	}

	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	s.payments[p.ID] = p
	return p, nil
}

func (s *InMemoryPaymentStore) GetByID(ctx context.Context, id string) (payments.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.payments[id]
	if !ok {
		return payments.Payment{}, payments.ErrNotFound
	}
	return p, nil
}

func (s *InMemoryPaymentStore) Update(ctx context.Context, id string, p payments.Payment) (payments.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" || p.UserID == "" || p.Amount <= 0 || p.Currency == "" || p.Status == "" {
		return payments.Payment{}, payments.ErrInvalid
	}

	existing, ok := s.payments[id]
	if !ok {
		return payments.Payment{}, payments.ErrNotFound
	}

	now := time.Now().UTC()
	existing.UserID = p.UserID
	existing.Amount = p.Amount
	existing.Currency = p.Currency
	existing.Status = p.Status
	existing.UpdatedAt = now

	s.payments[id] = existing
	return existing, nil
}

func (s *InMemoryPaymentStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.payments[id]; !ok {
		return payments.ErrNotFound
	}
	delete(s.payments, id)
	return nil
}
