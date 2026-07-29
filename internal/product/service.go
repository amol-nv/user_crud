package product

import (
	"time"

	"amol-nv/user_crud/internal/repository"
)

type Repository interface {
	Create(p repository.Product) (repository.Product, error)
	GetByID(id string) (repository.Product, error)
	Update(id string, p repository.Product) (repository.Product, error)
	Delete(id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(id string, req CreateProductRequest) (ProductResponse, error) {
	if err := validateCreate(req); err != nil {
		return ProductResponse{}, ErrInvalid
	}

	now := time.Now().UTC()
	created, err := s.repo.Create(repository.Product{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return ProductResponse{}, err
	}

	return mapToResponse(created), nil
}

func (s *Service) GetByID(id string) (ProductResponse, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return ProductResponse{}, mapRepoErr(err)
	}
	return mapToResponse(p), nil
}

func (s *Service) Update(id string, req UpdateProductRequest) (ProductResponse, error) {
	if err := validateUpdate(req); err != nil {
		return ProductResponse{}, ErrInvalid
	}

	now := time.Now().UTC()
	updated, err := s.repo.Update(id, repository.Product{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		UpdatedAt:   now,
	})
	if err != nil {
		return ProductResponse{}, mapRepoErr(err)
	}
	return mapToResponse(updated), nil
}

func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return mapRepoErr(err)
	}
	return nil
}

func validateCreate(req CreateProductRequest) error {
	if req.Name == "" {
		return ErrInvalid
	}
	if req.Price < 0 {
		return ErrInvalid
	}
	return nil
}

func validateUpdate(req UpdateProductRequest) error {
	if req.Name == "" {
		return ErrInvalid
	}
	if req.Price < 0 {
		return ErrInvalid
	}
	return nil
}

func mapRepoErr(err error) error {
	if err == repository.ErrNotFound {
		return ErrNotFound
	}
	return err
}
