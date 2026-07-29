package store

import (
	"amol-nv/user_crud/internal/storage"
)

func NewInMemoryUserStore() *storage.InMemoryUserStore {
	return storage.NewInMemoryUserStore()
}

func NewInMemoryPaymentsStore() *storage.InMemoryPaymentsStore {
	return storage.NewInMemoryPaymentsStore()
}
