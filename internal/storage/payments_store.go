package storage

import (
	"context"
	"errors"
	"sync"
	"time"

	"amol-nv/user_crud/internal/payments"
)

type InMemoryPaymentsStore struct {
	mu       sync.RWMutex
	byID     map[string]payments.Payment
	nextID   int64
}

func NewInMemoryPaymentsStore() *InMemoryPaymentsStore {
	return &InMemoryPaymentsStore{
		byID:   make(map[string]payments.Payment),
		nextID: 1,
	}
}

func (s *InMemoryPaymentsStore) Create(ctx context.Context, p payments.Payment) (payments.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	now := time.Now().UTC()
	p.ID = itoa(id)
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	s.byID[p.ID] = p
	return p, nil
}

func (s *InMemoryPaymentsStore) GetByID(ctx context.Context, id string) (payments.Payment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.byID[id]
	if !ok {
		return payments.Payment{}, payments.ErrPaymentNotFound
	}
	return p, nil
}

func (s *InMemoryPaymentsStore) Update(ctx context.Context, id string, p payments.Payment) (payments.Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.byID[id]
	if !ok {
		return payments.Payment{}, payments.ErrPaymentNotFound
	}
	p.ID = id
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now().UTC()
	}
	s.byID[id] = p
	return p, nil
}

func (s *InMemoryPaymentsStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[id]; !ok {
		return payments.ErrPaymentNotFound
	}
	delete(s.byID, id)
	return nil
}

func itoa(v int64) string {
	// small local helper to avoid importing strconv in multiple files
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		d := v % 10
		buf = append(buf, byte('0'+d))
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

var _ payments.PaymentsRepository = (*InMemoryPaymentsStore)(nil)

var _ = errors.New // keep errors import if build tags change
