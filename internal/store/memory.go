package store

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type InMemoryUserStore struct {
	mu     sync.RWMutex
	nextID int
	byID   map[string]User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{byID: make(map[string]User)}
}

func (s *InMemoryUserStore) Create(u User) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := fmt.Sprintf("%d", s.nextID)
	u.ID = id
	s.byID[id] = u
	return u, nil
}

func (s *InMemoryUserStore) Get(id string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.byID[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (s *InMemoryUserStore) Update(id string, patch UserPatch) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.byID[id]
	if !ok {
		return User{}, ErrNotFound
	}

	if patch.Name != nil {
		u.Name = *patch.Name
	}
	if patch.Email != nil {
		u.Email = *patch.Email
	}
	u.UpdatedAt = time.Now().UTC()

	s.byID[id] = u
	return u, nil
}

func (s *InMemoryUserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byID[id]; !ok {
		return ErrNotFound
	}
	delete(s.byID, id)
	return nil
}
