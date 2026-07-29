package payments

import "context"

type Store interface {
	Create(ctx context.Context, p Payment) (Payment, error)
	GetByID(ctx context.Context, id string) (Payment, error)
	Update(ctx context.Context, id string, p Payment) (Payment, error)
	Delete(ctx context.Context, id string) error
}
