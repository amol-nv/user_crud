package payments

import (
	"context"
	"sync"
)

type InMemoryPaymentRepository struct {
	mu       sync.RWMutex
	payments map[string]Payment
}

func NewInMemoryPaymentRepository() *InMemoryPaymentRepository {
	return &InMemoryPaymentRepository{payments: make(map[string]Payment)}
}

func (r *InMemoryPaymentRepository) Create(ctx context.Context, p Payment) (Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.payments[p.ID]; exists {
		// extremely unlikely due to timestamp-based ID; keep deterministic behavior
		return Payment{}, nil
	}
	r.payments[p.ID] = p
	return p, nil
}
