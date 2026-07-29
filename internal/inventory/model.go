package inventory

import "time"

type ProductInventory struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}
