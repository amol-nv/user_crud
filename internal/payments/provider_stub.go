package payments

import "context"

// StubProvider is a minimal provider implementation for this ticket.
// It always succeeds.
type StubProvider struct{}

func (p *StubProvider) Charge(ctx context.Context, pay Payment) (PaymentStatus, error) {
	return PaymentStatusSucceeded, nil
}
