package payments

import "time"

type Payment struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
