package payments

import (
	"context"
	"sync"
	"time"
)

type MemoryPaymentRepository struct {
	mu       sync.RWMutex
	payments map[string]Payment
}

func NewMemoryPaymentRepository() *MemoryPaymentRepository {
	return &MemoryPaymentRepository{payments: make(map[string]Payment)}
}

func (r *MemoryPaymentRepository) Create(ctx context.Context, p Payment) (Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Ensure timestamps are set.
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}

	r.payments[p.ID] = p
	return p, nil
}

func (r *MemoryPaymentRepository) GetByID(ctx context.Context, id string) (Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}
	return p, nil
}

func (r *MemoryPaymentRepository) UpdateStatus(ctx context.Context, id string, status PaymentStatus) (Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrPaymentNotFound
	}
	p.Status = status
	p.UpdatedAt = time.Now()
	r.payments[id] = p
	return p, nil
}
