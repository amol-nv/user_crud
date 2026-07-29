package storage

import (
	"context"
	"sync"
	"time"

	"amol-nv/user_crud/internal/inventory"
)

type InMemoryInventoryStore struct {
	mu     sync.RWMutex
	nextID int64
	items  map[int64]inventory.ProductInventory
}

func NewInMemoryInventoryStore() *InMemoryInventoryStore {
	return &InMemoryInventoryStore{
		nextID: 1,
		items:  make(map[int64]inventory.ProductInventory),
	}
}

func (s *InMemoryInventoryStore) Create(ctx context.Context, productID int64, quantity int) (inventory.ProductInventory, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	inv := inventory.ProductInventory{
		ID:        id,
		ProductID: productID,
		Quantity:  quantity,
		UpdatedAt: time.Now().UTC(),
	}
	s.items[id] = inv
	return inv, nil
}

func (s *InMemoryInventoryStore) GetByID(ctx context.Context, id int64) (inventory.ProductInventory, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	inv, ok := s.items[id]
	if !ok {
		return inventory.ProductInventory{}, &inventory.NotFoundError{ID: id}
	}
	return inv, nil
}

func (s *InMemoryInventoryStore) Update(ctx context.Context, id int64, productID int64, quantity int) (inventory.ProductInventory, error) {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	inv, ok := s.items[id]
	if !ok {
		return inventory.ProductInventory{}, &inventory.NotFoundError{ID: id}
	}
	inv.ProductID = productID
	inv.Quantity = quantity
	inv.UpdatedAt = time.Now().UTC()
	s.items[id] = inv
	return inv, nil
}

func (s *InMemoryInventoryStore) Delete(ctx context.Context, id int64) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return &inventory.NotFoundError{ID: id}
	}
	delete(s.items, id)
	return nil
}
