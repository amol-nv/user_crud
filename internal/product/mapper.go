package product

import (
	"time"

	"amol-nv/user_crud/internal/repository"
)

func mapToResponse(p repository.Product) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Price:       p.Price,
		Description: p.Description,
		CreatedAt:   formatTime(p.CreatedAt),
		UpdatedAt:   formatTime(p.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}
