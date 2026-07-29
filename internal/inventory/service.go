package inventory

import (
	"context"
	"time"
)

type InventoryStore interface {
	Create(ctx context.Context, productID int64, quantity int) (ProductInventory, error)
	GetByID(ctx context.Context, id int64) (ProductInventory, error)
	Update(ctx context.Context, id int64, productID int64, quantity int) (ProductInventory, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	store InventoryStore
}

func NewService(store InventoryStore) *Service {
	return &Service{store: store}
}

func (s *Service) Create(ctx context.Context, req CreateInventoryRequest) (ProductInventory, error) {
	if req.ProductID <= 0 {
		return ProductInventory{}, &ValidationError{Msg: "product_id must be positive"}
	}
	if req.Quantity < 0 {
		return ProductInventory{}, &ValidationError{Msg: "quantity must be non-negative"}
	}
	return s.store.Create(ctx, req.ProductID, req.Quantity)
}

func (s *Service) GetByID(ctx context.Context, id int64) (ProductInventory, error) {
	if id <= 0 {
		return ProductInventory{}, &ValidationError{Msg: "id must be positive"}
	}
	return s.store.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, req UpdateInventoryRequest) (ProductInventory, error) {
	if id <= 0 {
		return ProductInventory{}, &ValidationError{Msg: "id must be positive"}
	}
	if req.ProductID <= 0 {
		return ProductInventory{}, &ValidationError{Msg: "product_id must be positive"}
	}
	if req.Quantity < 0 {
		return ProductInventory{}, &ValidationError{Msg: "quantity must be non-negative"}
	}
	inv, err := s.store.Update(ctx, id, req.ProductID, req.Quantity)
	if err != nil {
		return ProductInventory{}, err
	}
	// Ensure updated_at is always set even if store doesn't.
	if inv.UpdatedAt.IsZero() {
		inv.UpdatedAt = time.Now().UTC()
	}
	return inv, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return &ValidationError{Msg: "id must be positive"}
	}
	return s.store.Delete(ctx, id)
}
