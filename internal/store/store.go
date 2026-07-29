package store

import "amol-nv/user_crud/internal/users"

type Store interface {
	UserStore
	PaymentStore
}

type UserStore interface {
		Create(u *users.User) (*users.User, error)
		GetByID(id string) (*users.User, error)
		List() ([]*users.User, error)
		Update(u *users.User) (*users.User, error)
		Delete(id string) error
}

type PaymentStore interface {
		Create(p *Payment) (*Payment, error)
		GetByID(id string) (*Payment, error)
		List() ([]*Payment, error)
		Update(p *Payment) (*Payment, error)
		Delete(id string) error
}

// Payment is the persistence model used by the in-memory store.
// The payments service maps this to its own model.
// Keeping it in store avoids circular imports.
type Payment struct {
	ID        string
	Amount    float64
	Currency  string
	Status    string
	CreatedAt interface{}
	UpdatedAt interface{}
}
