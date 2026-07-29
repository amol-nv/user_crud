package payments

import (
	"context"
	"sync"
	"time"
)

type memoryPaymentRepository struct {
	mu       sync.RWMutex
	payments map[string]Payment
}

func NewMemoryPaymentRepository() PaymentRepository {
	return &memoryPaymentRepository{payments: make(map[string]Payment)}
}

func (r *memoryPaymentRepository) Create(ctx context.Context, p Payment) (Payment, error) {
	_ = ctx
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.payments[p.ID]; exists {
		// For simplicity, treat as invalid argument.
		return Payment{}, ErrInvalidArgument
	}
	r.payments[p.ID] = p
	return p, nil
}

func (r *memoryPaymentRepository) GetByID(ctx context.Context, id string) (Payment, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrNotFound
	}
	return p, nil
}

func (r *memoryPaymentRepository) Update(ctx context.Context, id string, p Payment) (Payment, error) {
	_ = ctx
	now := time.Now().UTC()

	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrNotFound
	}

	p.ID = id
	p.CreatedAt = existing.CreatedAt
	p.UpdatedAt = now

	r.payments[id] = p
	return p, nil
}

func (r *memoryPaymentRepository) Delete(ctx context.Context, id string) (Payment, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.payments[id]
	if !ok {
		return Payment{}, ErrNotFound
	}
	delete(r.payments, id)
	return p, nil
}
