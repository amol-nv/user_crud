package product

import (
	"context"

	"amol-nv/user_crud/internal/store"
)

type Service struct {
	store store.ProductStore
}

func NewProductService(st store.ProductStore) *Service {
	return &Service{store: st}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	if req.Name == "" {
		return Product{}, ErrInvalid
	}
	if req.Price < 0 {
		return Product{}, ErrInvalid
	}

	p := Product{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
	}
	return s.store.Create(ctx, p)
}

func (s *Service) GetByID(ctx context.Context, id string) (Product, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	return s.store.List(ctx)
}

func (s *Service) Update(ctx context.Context, id string, req UpdateProductRequest) (Product, error) {
	current, err := s.store.GetByID(ctx, id)
	if err != nil {
		return Product{}, err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return Product{}, ErrInvalid
		}
		current.Name = *req.Name
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return Product{}, ErrInvalid
		}
		current.Price = *req.Price
	}
	if req.Description != nil {
		current.Description = *req.Description
	}

	return s.store.Update(ctx, id, current)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}
