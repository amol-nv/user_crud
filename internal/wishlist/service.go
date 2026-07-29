package wishlist

import (
	"context"
	"github.com/google/uuid"
	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/store"
)

type Service struct {
	store store.WishlistStore
	inv   inventory.Service
}

func NewService(w store.WishlistStore, inv inventory.Service) *Service {
	return &Service{store: w, inv: inv}
}

// Add adds a product to a user's wishlist. It is idempotent.
func (s *Service) Add(ctx context.Context, userID, productID uuid.UUID) error {
	// Validate product exists
	_, err := s.inv.GetByID(ctx, productID)
	if err != nil {
		return ErrProductNotFound
	}
	return s.store.Add(ctx, userID, productID)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]inventory.Product, error) {
	items, err := s.store.List(ctx, userID)
	if err != nil {
		return nil, err
	}
	products := make([]inventory.Product, 0, len(items))
	for _, it := range items {
		p, err := s.inv.GetByID(ctx, it.ProductID)
		if err != nil {
			// If product was deleted, skip it.
			continue
		}
		products = append(products, p)
	}
	return products, nil
}
