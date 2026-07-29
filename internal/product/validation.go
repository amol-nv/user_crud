package product

import "strings"

func validateCreate(req CreateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrValidation
	}
	if req.Price < 0 {
		return ErrValidation
	}
	return nil
}

func validateUpdate(req UpdateProductRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrValidation
	}
	if req.Price < 0 {
		return ErrValidation
	}
	return nil
}
