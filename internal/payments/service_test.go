package payments

import (
	"context"
	"testing"
)

func TestPaymentService_CRUD(t *testing.T) {
	repo := NewMemoryPaymentRepository()
	svc := NewPaymentService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreatePaymentRequest{Amount: 10.5, Currency: "usd", Status: "created"})
	if err != nil {
		t.Fatalf("expected create success, got err: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id")
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected get success, got err: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected same id")
	}

	updated, err := svc.Update(ctx, created.ID, UpdatePaymentRequest{Amount: 20, Currency: "USD", Status: "updated"})
	if err != nil {
		t.Fatalf("expected update success, got err: %v", err)
	}
	if updated.Amount != 20 {
		t.Fatalf("expected updated amount")
	}

	_, err = svc.GetByID(ctx, "missing")
	if err != ErrNotFound {
		t.Fatalf("expected not found, got: %v", err)
	}

	_, err = svc.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected delete success, got err: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err != ErrNotFound {
		t.Fatalf("expected not found after delete, got: %v", err)
	}
}

func TestPaymentService_Validation(t *testing.T) {
	repo := NewMemoryPaymentRepository()
	svc := NewPaymentService(repo)
	ctx := context.Background()

	_, err := svc.Create(ctx, CreatePaymentRequest{Amount: 0, Currency: "USD", Status: "created"})
	if err != ErrInvalidArgument {
		t.Fatalf("expected invalid argument, got: %v", err)
	}

	_, err = svc.Update(ctx, "pay_1", UpdatePaymentRequest{Amount: 1, Currency: "", Status: "updated"})
	if err != ErrInvalidArgument {
		t.Fatalf("expected invalid argument, got: %v", err)
	}
}
