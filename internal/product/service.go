package product

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	if err := validateCreate(req); err != nil {
		return Product{}, err
	}

	p := Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	return s.repo.Create(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id ID) (Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}

func (s *Service) Update(ctx context.Context, id ID, req UpdateProductRequest) (Product, error) {
	if err := validateUpdate(req); err != nil {
		return Product{}, err
	}

	p := Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
	return s.repo.Update(ctx, id, p)
}

func (s *Service) Delete(ctx context.Context, id ID) error {
	return s.repo.Delete(ctx, id)
}
