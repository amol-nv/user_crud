package product

import "context"

type Repository interface {
	Create(ctx context.Context, p Product) (Product, error)
	GetByID(ctx context.Context, id ID) (Product, error)
	List(ctx context.Context) ([]Product, error)
	Update(ctx context.Context, id ID, p Product) (Product, error)
	Delete(ctx context.Context, id ID) error
}
