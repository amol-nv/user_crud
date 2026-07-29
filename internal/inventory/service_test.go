package inventory

import (
	"context"
	"testing"

	"amol-nv/user_crud/internal/storage"
)

func TestService_CRUD(t *testing.T) {
	store := storage.NewInMemoryInventoryStore()
	svc := NewService(store)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreateInventoryRequest{ProductID: 10, Quantity: 5})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected id")
	}
	if created.ProductID != 10 || created.Quantity != 5 {
		t.Fatalf("unexpected created: %+v", created)
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.Quantity != 5 {
		t.Fatalf("expected quantity 5")
	}

	updated, err := svc.Update(ctx, created.ID, UpdateInventoryRequest{ProductID: 10, Quantity: 7})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if updated.Quantity != 7 {
		t.Fatalf("expected quantity 7")
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Fatalf("expected not found")
	}
	if _, ok := err.(*NotFoundError); !ok {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
}

func TestService_Validation(t *testing.T) {
	store := storage.NewInMemoryInventoryStore()
	svc := NewService(store)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreateInventoryRequest{ProductID: 0, Quantity: 1})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	_, err = svc.Create(ctx, CreateInventoryRequest{ProductID: 1, Quantity: -1})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if _, ok := err.(*ValidationError); !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}
