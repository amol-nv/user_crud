package payments

import (
	"context"
	"testing"
)

type stubRepo struct {
	created Payment
	called  bool
}

func (s *stubRepo) Create(ctx context.Context, p Payment) (Payment, error) {
	s.called = true
	s.created = p
	return p, nil
}

func TestService_CreatePayment_Valid(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	p, err := svc.CreatePayment(context.Background(), CreatePaymentRequest{
		UserID:   "u1",
		Amount:   100,
		Currency: "usd",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !repo.called {
		t.Fatalf("expected repo to be called")
	}
	if p.UserID != "u1" {
		t.Fatalf("expected userId u1, got %s", p.UserID)
	}
	if p.Amount != 100 {
		t.Fatalf("expected amount 100, got %d", p.Amount)
	}
	if p.Currency != "USD" {
		t.Fatalf("expected currency USD, got %s", p.Currency)
	}
	if p.Status != PaymentStatusCompleted {
		t.Fatalf("expected status completed, got %s", p.Status)
	}
}

func TestService_CreatePayment_InvalidAmount(t *testing.T) {
	repo := &stubRepo{}
	svc := NewService(repo)

	_, err := svc.CreatePayment(context.Background(), CreatePaymentRequest{
		UserID:   "u1",
		Amount:   0,
		Currency: "USD",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
