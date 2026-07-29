package store

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPatch struct {
	Name  *string
	Email *string
}

type UserStore interface {
	Create(u User) (User, error)
	Get(id string) (User, error)
	Update(id string, patch UserPatch) (User, error)
	Delete(id string) error
}
