package product

import (
	"context"
	"testing"
)

func TestServiceCRUD(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateProductRequest{Name: "Phone", Description: "X", Price: 10.5})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected id")
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.Name != "Phone" {
		t.Fatalf("expected name")
	}

	items, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item")
	}

	updated, err := svc.Update(ctx, created.ID, UpdateProductRequest{Name: "Phone2", Description: "Y", Price: 20})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if updated.Price != 20 {
		t.Fatalf("expected updated price")
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Fatalf("expected not found")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected not found")
	}
}

func IsNotFound(err error) bool {
	return err == ErrNotFound || (err != nil && (err.Error() == ErrNotFound.Error()))
}
