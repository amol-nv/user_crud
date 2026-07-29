package product

import (
	"context"
	"sync"
	"time"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	nextID ID
	data   map[ID]Product
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		nextID: 1,
		data:   make(map[ID]Product),
	}
}

func (r *MemoryRepository) Create(ctx context.Context, p Product) (Product, error) {
	_ = ctx
	now := time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.nextID++
	p.CreatedAt = now
	p.UpdatedAt = now

	r.data[p.ID] = p
	return p, nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id ID) (Product, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.data[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Product, error) {
	_ = ctx
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Product, 0, len(r.data))
	for _, p := range r.data {
		out = append(out, p)
	}
	return out, nil
}

func (r *MemoryRepository) Update(ctx context.Context, id ID, p Product) (Product, error) {
	_ = ctx
	now := time.Now().UTC()
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.data[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	existing.Name = p.Name
	existing.Description = p.Description
	existing.Price = p.Price
	existing.UpdatedAt = now

	r.data[id] = existing
	return existing, nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id ID) error {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return ErrNotFound
	}
	delete(r.data, id)
	return nil
}
