package repository

import "time"

type Product struct {
	ID          string
	Name        string
	Price       float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
