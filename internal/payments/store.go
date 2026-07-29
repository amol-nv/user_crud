package payments

import (
	"sync"
	"time"
)

type MemoryPaymentStore struct {
	mu       sync.RWMutex
	payments map[string]Payment
}

func NewMemoryPaymentStore() *MemoryPaymentStore {
	return &MemoryPaymentStore{payments: make(map[string]Payment)}
}

func (m *MemoryPaymentStore) Create(payment Payment) (Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.payments[payment.ID] = payment
	return payment, nil
}

func (m *MemoryPaymentStore) GetByID(id string) (Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}
	return p, nil
}

func (m *MemoryPaymentStore) List() ([]Payment, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Payment, 0, len(m.payments))
	for _, p := range m.payments {
		// ensure timestamps are non-zero for older entries
		if p.CreatedAt.IsZero() {
			p.CreatedAt = time.Now().UTC()
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = p.CreatedAt
		}
		out = append(out, p)
	}
	return out, nil
}

func (m *MemoryPaymentStore) Update(id string, payment Payment) (Payment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.payments[id]; !ok {
		return Payment{}, ErrPaymentNotFound
	}
	payment.ID = id
	m.payments[id] = payment
	return payment, nil
}

func (m *MemoryPaymentStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.payments[id]; !ok {
		return ErrPaymentNotFound
	}
	delete(m.payments, id)
	return nil
}
