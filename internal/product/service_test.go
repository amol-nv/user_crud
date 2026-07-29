package product

import (
	"context"
	"testing"

	"amol-nv/user_crud/internal/store"
)

func TestProductService_CRUD(t *testing.T) {
	st := store.NewInMemoryProductStore()
	svc := NewProductService(st)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateProductRequest{Name: "Phone", Price: 10.5, Description: "Nice"})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id")
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.Name != "Phone" {
		t.Fatalf("expected name")
	}

	newName := "Phone Pro"
	newPrice := 12.0
	updated, err := svc.Update(ctx, created.ID, UpdateProductRequest{Name: &newName, Price: &newPrice})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if updated.Name != "Phone Pro" || updated.Price != 12.0 {
		t.Fatalf("update not applied")
	}

	list, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 item")
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}
	_, err = svc.GetByID(ctx, created.ID)
	if err != ErrNotFound {
		t.Fatalf("expected not found, got: %v", err)
	}
}

func TestProductService_NotFound(t *testing.T) {
	st := store.NewInMemoryProductStore()
	svc := NewProductService(st)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, "missing")
	if err != ErrNotFound {
		t.Fatalf("expected not found")
	}

	_, err = svc.Update(ctx, "missing", UpdateProductRequest{})
	if err != ErrNotFound {
		t.Fatalf("expected not found")
	}

	if err := svc.Delete(ctx, "missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
}

func TestProductService_Validation(t *testing.T) {
	st := store.NewInMemoryProductStore()
	svc := NewProductService(st)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateProductRequest{Name: "", Price: 1})
	if err != ErrInvalid {
		t.Fatalf("expected invalid")
	}

	_, err = svc.Create(ctx, CreateProductRequest{Name: "A", Price: -1})
	if err != ErrInvalid {
		t.Fatalf("expected invalid")
	}

	created, _ := svc.Create(ctx, CreateProductRequest{Name: "A", Price: 1})
	badName := ""
	_, err = svc.Update(ctx, created.ID, UpdateProductRequest{Name: &badName})
	if err != ErrInvalid {
		t.Fatalf("expected invalid")
	}
}
