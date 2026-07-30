package payments

import (
	"context"
	"testing"
)

type fakeProvider struct {
	status PaymentStatus
	err    error
}

func (f *fakeProvider) Charge(ctx context.Context, p Payment) (PaymentStatus, error) {
	return f.status, f.err
}

type fakeRepo struct {
	created Payment
	gotID   string
}

func (r *fakeRepo) Create(ctx context.Context, p Payment) (Payment, error) {
	r.created = p
	return p, nil
}

func (r *fakeRepo) GetByID(ctx context.Context, id string) (Payment, error) {
	r.gotID = id
	return r.created, nil
}

func (r *fakeRepo) UpdateStatus(ctx context.Context, id string, status PaymentStatus) (Payment, error) {
	r.created.ID = id
	r.created.Status = status
	return r.created, nil
}

func TestService_CreatePayment_Succeeds(t *testing.T) {
	repo := &fakeRepo{}
	provider := &fakeProvider{status: PaymentStatusSucceeded}
	svc := NewService(repo, provider)

	resp, err := svc.CreatePayment(context.Background(), CreatePaymentRequest{
		UserID:   "u1",
		OrderID:  "o1",
		Amount:   100,
		Currency: "usd",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.Payment.Status != PaymentStatusSucceeded {
		t.Fatalf("expected succeeded, got %s", resp.Payment.Status)
	}
	if resp.Payment.Currency != "USD" {
		t.Fatalf("expected currency USD, got %s", resp.Payment.Currency)
	}
}

func TestService_CreatePayment_InvalidAmount(t *testing.T) {
	repo := &fakeRepo{}
	provider := &fakeProvider{status: PaymentStatusSucceeded}
	svc := NewService(repo, provider)

	_, err := svc.CreatePayment(context.Background(), CreatePaymentRequest{
		UserID:   "u1",
		OrderID:  "o1",
		Amount:   0,
		Currency: "USD",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
}
