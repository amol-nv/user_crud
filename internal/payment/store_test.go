package payment

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"
)

func TestHTTPStore_PostPayment_SendsJSONAndDecodesResponse(t *testing.T) {
	var gotReq PaymentRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/v1/payments" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_ = json.NewDecoder(r.Body).Decode(&gotReq)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(PaymentResponse{PaymentID: "pay_1", Status: "created"})
	}))
	defer srv.Close()

	store := NewHTTPStore(&http.Client{Timeout: 2 * time.Second}, srv.URL, "/v1/payments", "")

	resp, status, err := store.PostPayment(context.Background(), PaymentRequest{Amount: 10, Currency: "USD", Method: "card", OrderID: "ord_1"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if status != 200 {
		t.Fatalf("expected status 200, got %d", status)
	}
	if resp.PaymentID != "pay_1" {
		t.Fatalf("expected payment_id pay_1, got %s", resp.PaymentID)
	}
	if gotReq.OrderID != "ord_1" {
		t.Fatalf("expected order_id ord_1, got %s", gotReq.OrderID)
	}
}
