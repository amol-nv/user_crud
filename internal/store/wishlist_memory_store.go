package store

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type WishlistMemoryStore struct {
	mu sync.RWMutex
	// userID -> set(productID)
	items map[uuid.UUID]map[uuid.UUID]struct{}
}

func NewWishlistMemoryStore() *WishlistMemoryStore {
	return &WishlistMemoryStore{items: make(map[uuid.UUID]map[uuid.UUID]struct{})}
}

func (s *WishlistMemoryStore) Add(ctx context.Context, userID, productID uuid.UUID) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	set, ok := s.items[userID]
	if !ok {
		set = make(map[uuid.UUID]struct{})
		s.items[userID] = set
	}
	set[productID] = struct{}{}
	return nil
}

func (s *WishlistMemoryStore) List(ctx context.Context, userID uuid.UUID) ([]WishlistItem, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()
	set, ok := s.items[userID]
	if !ok {
		return []WishlistItem{}, nil
	}
	out := make([]WishlistItem, 0, len(set))
	for pid := range set {
		out = append(out, WishlistItem{ProductID: pid})
	}
	return out, nil
}
