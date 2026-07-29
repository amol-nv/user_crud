package payments

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestPaymentHandler_CRUD(t *testing.T) {
	repo := NewMemoryPaymentRepository()
	svc := NewPaymentService(repo)
	h := NewPaymentHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	// Create
	createBody := []byte(`{"amount":10.5,"currency":"usd","status":"created"}`)
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(createBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var createResp CreatePaymentResponse
	if err := decodeJSONFromRecorder(w, &createResp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if createResp.Payment.ID == "" {
		t.Fatalf("expected payment id")
	}

	id := createResp.Payment.ID

	// Get
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/payments/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Update
	updateBody := []byte(`{"amount":20,"currency":"USD","status":"updated"}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/payments/"+id, bytes.NewReader(updateBody))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/payments/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Get after delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/payments/"+id, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
