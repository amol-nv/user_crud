package wishlist

import "github.com/google/uuid"

type WishlistItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
}
