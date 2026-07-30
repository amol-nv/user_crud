package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/payment"
)

type stubStore struct {
	lastReq payment.PaymentRequest
	resp    payment.PaymentResponse
	err     error
}

func (s *stubStore) PostPayment(ctx context.Context, req payment.PaymentRequest) (payment.PaymentResponse, int, error) {
	s.lastReq = req
	return s.resp, 200, s.err
}

func TestPaymentHandler_PostPayment_OK(t *testing.T) {
	st := &stubStore{
		resp: payment.PaymentResponse{PaymentID: "p1", Status: "created"},
	}
	svc := payment.NewService(st)
	h := NewPaymentHandler(svc)

	body := payment.PaymentRequest{Amount: 100, Currency: "USD", Method: "card", OrderID: "o1"}
	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(b))
	w := httptest.NewRecorder()

	h.PostPayment(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if st.lastReq.OrderID != "o1" {
		t.Fatalf("expected order_id o1, got %s", st.lastReq.OrderID)
	}
}

func TestPaymentHandler_PostPayment_BadJSON(t *testing.T) {
	st := &stubStore{}
	svc := payment.NewService(st)
	h := NewPaymentHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewBufferString("{bad json"))
	w := httptest.NewRecorder()

	h.PostPayment(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}
