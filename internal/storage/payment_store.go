package storage

import (
	"context"
	"sync"

	"amol-nv/user_crud/internal/payment"
)

type PaymentStore struct {
	mu    sync.RWMutex
	data  map[string]payment.Payment
}

func NewPaymentStore() *PaymentStore {
	return &PaymentStore{data: make(map[string]payment.Payment)}
}

func (s *PaymentStore) Create(ctx context.Context, p payment.Payment) (payment.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[p.ID]; ok {
		// treat as invalid for simplicity
		return payment.Payment{}, payment.ErrInvalidPayment
	}
	s.data[p.ID] = p
	return p, nil
}

func (s *PaymentStore) GetByID(ctx context.Context, id string) (payment.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.data[id]
	if !ok {
		return payment.Payment{}, payment.ErrNotFound
	}
	return p, nil
}

func (s *PaymentStore) Update(ctx context.Context, id string, p payment.Payment) (payment.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return payment.Payment{}, payment.ErrNotFound
	}
	p.ID = id
	s.data[id] = p
	return p, nil
}

func (s *PaymentStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return payment.ErrNotFound
	}
	delete(s.data, id)
	return nil
}
