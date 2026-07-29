package users

import (
	"strings"
	"time"

	"amol-nv/user_crud/internal/store"
)

type Service struct {
	st store.UserStore
}

func NewService(st store.UserStore) *Service {
	return &Service{st: st}
}

func (s *Service) Create(req CreateUserRequest) (User, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)
	if name == "" || email == "" {
		return User{}, ErrValidation
	}
	if !strings.Contains(email, "@") {
		return User{}, ErrValidation
	}

	now := time.Now().UTC()
	u := store.User{
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := s.st.Create(u)
	if err != nil {
		return User{}, err
	}
	return toUser(created), nil
}

func (s *Service) Get(id string) (User, error) {
	if strings.TrimSpace(id) == "" {
		return User{}, ErrInvalidID
	}
	u, err := s.st.Get(id)
	if err != nil {
		return User{}, err
	}
	return toUser(u), nil
}

func (s *Service) Update(id string, req UpdateUserRequest) (User, error) {
	if strings.TrimSpace(id) == "" {
		return User{}, ErrInvalidID
	}

	patch := store.UserPatch{}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return User{}, ErrValidation
		}
		patch.Name = &name
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email == "" || !strings.Contains(email, "@") {
			return User{}, ErrValidation
		}
		patch.Email = &email
	}
	if patch.Name == nil && patch.Email == nil {
		return User{}, ErrValidation
	}

	updated, err := s.st.Update(id, patch)
	if err != nil {
		return User{}, err
	}
	return toUser(updated), nil
}

func (s *Service) Delete(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidID
	}
	return s.st.Delete(id)
}

func toUser(u store.User) User {
	return User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
