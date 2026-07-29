package store

import (
	"context"

	"github.com/google/uuid"
)

type WishlistStore interface {
	Add(ctx context.Context, userID, productID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]WishlistItem, error)
}

type WishlistItem struct {
	ProductID uuid.UUID
}
