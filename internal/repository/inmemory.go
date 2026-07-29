package repository

import (
	"sync"
	"time"
)

type InMemoryProductRepository struct {
	mu       sync.RWMutex
	products map[string]Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{products: make(map[string]Product)}
}

func (r *InMemoryProductRepository) Create(p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Overwrite semantics are not desired for create; treat existing as update.
	// For simplicity, we allow overwrite here.
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = p.CreatedAt
	}

	r.products[p.ID] = p
	return p, nil
}

func (r *InMemoryProductRepository) GetByID(id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *InMemoryProductRepository) Update(id string, p Product) (Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}

	// Preserve createdAt; update fields.
	existing.Name = p.Name
	existing.Price = p.Price
	existing.Description = p.Description
	if !p.UpdatedAt.IsZero() {
		existing.UpdatedAt = p.UpdatedAt
	} else {
		existing.UpdatedAt = time.Now().UTC()
	}

	r.products[id] = existing
	return existing, nil
}

func (r *InMemoryProductRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return ErrNotFound
	}
	delete(r.products, id)
	return nil
}
