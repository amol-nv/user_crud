package store

import (
	"context"
	"errors"
	"sync"
	"time"

	"amol-nv/user_crud/internal/models"
)

type UserStore interface {
	CreateUser(ctx context.Context, u models.User) (models.User, error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
	UpdateUser(ctx context.Context, id int64, u models.User) (models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

var (
	ErrNotFound = errors.New("user not found")
	ErrConflict = errors.New("user already exists")
)

type InMemoryUserStore struct {
	mu     sync.RWMutex
	byID   map[int64]models.User
	byEmail map[string]int64
	nextID int64
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		byID:    make(map[int64]models.User),
		byEmail: make(map[string]int64),
		nextID:  1,
	}
}

func (s *InMemoryUserStore) CreateUser(ctx context.Context, u models.User) (models.User, error) {
	_ = ctx
	now := time.Now().UTC()

	s.mu.Lock()
	defer s.mu.Unlock()

	if u.Email != "" {
		if existingID, ok := s.byEmail[u.Email]; ok {
			return models.User{}, ErrConflict
			_ = existingID
		}
	}

	u.ID = s.nextID
	s.nextID++
	u.CreatedAt = now
	u.UpdatedAt = now

	s.byID[u.ID] = u
	if u.Email != "" {
		s.byEmail[u.Email] = u.ID
	}

	return u, nil
}

func (s *InMemoryUserStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.byID[id]
	if !ok {
		return models.User{}, ErrNotFound
	}
	return u, nil
}

func (s *InMemoryUserStore) UpdateUser(ctx context.Context, id int64, u models.User) (models.User, error) {
	_ = ctx
	now := time.Now().UTC()

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.byID[id]
	if !ok {
		return models.User{}, ErrNotFound
	}

	// Email uniqueness check
	if u.Email != "" && u.Email != existing.Email {
		if existingID, ok := s.byEmail[u.Email]; ok && existingID != id {
			return models.User{}, ErrConflict
		}
	}

	// Update fields
	existing.Name = u.Name
	if u.Email != "" {
		existing.Email = u.Email
	}
	existing.UpdatedAt = now

	// Update email index if changed
	if u.Email != "" && u.Email != s.byID[id].Email {
		delete(s.byEmail, s.byID[id].Email)
		s.byEmail[u.Email] = id
	}

	s.byID[id] = existing
	return existing, nil
}

func (s *InMemoryUserStore) DeleteUser(ctx context.Context, id int64) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.byID[id]
	if !ok {
		return ErrNotFound
	}

	delete(s.byID, id)
	if u.Email != "" {
		delete(s.byEmail, u.Email)
	}
	return nil
}
