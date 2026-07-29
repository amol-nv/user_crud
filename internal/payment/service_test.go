package payment

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepo struct {
	create func(ctx context.Context, p Payment) (Payment, error)
	get    func(ctx context.Context, id string) (Payment, error)
	update func(ctx context.Context, id string, p Payment) (Payment, error)
	del    func(ctx context.Context, id string) error
}

func (f *fakeRepo) Create(ctx context.Context, p Payment) (Payment, error) { return f.create(ctx, p) }
func (f *fakeRepo) GetByID(ctx context.Context, id string) (Payment, error) {
	return f.get(ctx, id)
}
func (f *fakeRepo) Update(ctx context.Context, id string, p Payment) (Payment, error) {
	return f.update(ctx, id, p)
}
func (f *fakeRepo) Delete(ctx context.Context, id string) error { return f.del(ctx, id) }

func TestService_Create_Invalid(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	_, err := svc.Create(context.Background(), CreatePaymentRequest{UserID: "u1", Amount: 0, Currency: "USD", Status: "PENDING"})
	if !errors.Is(err, ErrInvalidPayment) {
		t.Fatalf("expected ErrInvalidPayment, got %v", err)
	}
}

func TestService_Create_Valid(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)
	svc.idFn = func() string { return "p1" }

	repo.create = func(ctx context.Context, p Payment) (Payment, error) {
		if p.ID != "p1" {
			t.Fatalf("expected id p1, got %s", p.ID)
		}
		if p.Currency != "USD" {
			t.Fatalf("expected currency USD, got %s", p.Currency)
		}
		if p.Status != "PENDING" {
			t.Fatalf("expected status PENDING, got %s", p.Status)
		}
		if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
			t.Fatalf("expected timestamps")
		}
		return p, nil
	}

	p, err := svc.Create(context.Background(), CreatePaymentRequest{UserID: "u1", Amount: 10, Currency: "usd", Status: "pending"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.ID != "p1" {
		t.Fatalf("expected id p1, got %s", p.ID)
	}
}

func TestService_Update_NotFound(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	repo.get = func(ctx context.Context, id string) (Payment, error) {
		return Payment{}, ErrNotFound
	}

	_, err := svc.Update(context.Background(), "missing", UpdatePaymentRequest{Amount: ptrInt64(20)})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_Update_Valid(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewService(repo)

	old := Payment{ID: "p1", UserID: "u1", Amount: 10, Currency: "USD", Status: "PENDING", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	repo.get = func(ctx context.Context, id string) (Payment, error) { return old, nil }
	repo.update = func(ctx context.Context, id string, p Payment) (Payment, error) {
		if id != "p1" {
			t.Fatalf("expected id p1, got %s", id)
		}
		if p.Amount != 20 {
			t.Fatalf("expected amount 20, got %d", p.Amount)
		}
		if p.Currency != "EUR" {
			t.Fatalf("expected currency EUR, got %s", p.Currency)
		}
		if p.UpdatedAt.IsZero() {
			t.Fatalf("expected updatedAt")
		}
		return p, nil
	}

	p, err := svc.Update(context.Background(), "p1", UpdatePaymentRequest{Amount: ptrInt64(20), Currency: ptrString("eur")})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if p.Amount != 20 {
		t.Fatalf("expected amount 20, got %d", p.Amount)
	}
}

func ptrInt64(v int64) *int64 { return &v }
func ptrString(v string) *string { return &v }
