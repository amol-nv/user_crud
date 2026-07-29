package store

import (
	"context"
	"sync"

	"amol-nv/user_crud/internal/product"
)

type ProductStore interface {
	Create(ctx context.Context, p product.Product) (product.Product, error)
	GetByID(ctx context.Context, id string) (product.Product, error)
	List(ctx context.Context) ([]product.Product, error)
	Update(ctx context.Context, id string, p product.Product) (product.Product, error)
	Delete(ctx context.Context, id string) error
}

type InMemoryProductStore struct {
	mu       sync.RWMutex
	items    map[string]product.Product
	idSeq    int
}

func NewInMemoryProductStore() *InMemoryProductStore {
	return &InMemoryProductStore{items: make(map[string]product.Product)}
}

func (s *InMemoryProductStore) Create(ctx context.Context, p product.Product) (product.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idSeq++
	p.ID = itoa(s.idSeq)
	s.items[p.ID] = p
	return p, nil
}

func (s *InMemoryProductStore) GetByID(ctx context.Context, id string) (product.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.items[id]
	if !ok {
		return product.Product{}, product.ErrNotFound
	}
	return p, nil
}

func (s *InMemoryProductStore) List(ctx context.Context) ([]product.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]product.Product, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out, nil
}

func (s *InMemoryProductStore) Update(ctx context.Context, id string, p product.Product) (product.Product, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return product.Product{}, product.ErrNotFound
	}
	p.ID = id
	s.items[id] = p
	return p, nil
}

func (s *InMemoryProductStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return product.ErrNotFound
	}
	delete(s.items, id)
	return nil
}

func itoa(i int) string {
	// small helper to avoid importing strconv in multiple files
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	buf := make([]byte, 0, 12)
	for i > 0 {
		buf = append(buf, byte('0'+(i%10)))
		i /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for l, r := 0, len(buf)-1; l < r; l, r = l+1, r-1 {
		buf[l], buf[r] = buf[r], buf[l]
	}
	return string(buf)
}
