package payments

import (
	"context"
	"testing"

	"amol-nv/user_crud/internal/storage"
)

func TestServiceCRUD(t *testing.T) {
	st := storage.NewInMemoryPaymentStore()
	svc := NewService(st)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreatePaymentRequest{UserID: "u1", Amount: 100, Currency: "USD"})
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
	if got.UserID != "u1" || got.Amount != 100 || got.Currency != "USD" {
		t.Fatalf("unexpected payment")
	}

	updated, err := svc.Update(ctx, created.ID, UpdatePaymentRequest{UserID: "u1", Amount: 200, Currency: "USD", Status: "paid"})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if updated.Amount != 200 || updated.Status != "paid" {
		t.Fatalf("unexpected updated payment")
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err != ErrNotFound {
		t.Fatalf("expected not found, got: %v", err)
	}
}

func TestServiceErrors(t *testing.T) {
	st := storage.NewInMemoryPaymentStore()
	svc := NewService(st)
	ctx := context.Background()

	_, err := svc.GetByID(ctx, "missing")
	if err != ErrNotFound {
		t.Fatalf("expected not found, got: %v", err)
	}

	_, err = svc.Create(ctx, CreatePaymentRequest{UserID: "", Amount: 0, Currency: ""})
	if err != ErrInvalid {
		t.Fatalf("expected invalid, got: %v", err)
	}

	_, err = svc.Update(ctx, "missing", UpdatePaymentRequest{UserID: "u1", Amount: 1, Currency: "USD", Status: "paid"})
	if err != ErrNotFound {
		t.Fatalf("expected not found, got: %v", err)
	}
}
