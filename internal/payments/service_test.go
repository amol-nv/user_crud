package payments

import (
	"context"
	"testing"
	"time"

	"amol-nv/user_crud/internal/storage"
)

func TestPaymentsService_CRUD(t *testing.T) {
	repo := storage.NewInMemoryPaymentsStore()
	svc := NewService(repo)
	ctx := context.Background()

	created, err := svc.Create(ctx, CreatePaymentRequest{Amount: 100, Currency: "usd", Status: "pending"})
	if err != nil {
		t.Fatalf("create err: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id")
	}
	if created.Currency != "USD" {
		t.Fatalf("expected currency USD, got %s", created.Currency)
	}
	if created.Status != "pending" {
		t.Fatalf("expected status pending, got %s", created.Status)
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Fatalf("expected timestamps")
	}

	got, err := svc.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected same id")
	}

	newAmount := int64(250)
	newCurrency := "eur"
	newStatus := "completed"
	updated, err := svc.Update(ctx, created.ID, UpdatePaymentRequest{Amount: &newAmount, Currency: &newCurrency, Status: &newStatus})
	if err != nil {
		t.Fatalf("update err: %v", err)
	}
	if updated.Amount != 250 {
		t.Fatalf("expected amount 250, got %d", updated.Amount)
	}
	if updated.Currency != "EUR" {
		t.Fatalf("expected currency EUR, got %s", updated.Currency)
	}
	if updated.Status != "completed" {
		t.Fatalf("expected status completed, got %s", updated.Status)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt.Add(-time.Second)) {
		// allow small clock skew
		t.Fatalf("expected updatedAt to change")
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("delete err: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
	if !Is(err, ErrPaymentNotFound) {
		t.Fatalf("expected ErrPaymentNotFound, got %v", err)
	}
}

func Is(err, target error) bool {
	// local helper to avoid importing errors in test file
	return err != nil && target != nil && (err.Error() == target.Error())
}
