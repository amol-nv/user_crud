package payments

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amol-nv/user_crud/internal/httpapi"
	"github.com/go-chi/chi/v5"
)

func TestPaymentCRUD(t *testing.T) {
	store := NewMemoryPaymentStore()
	svc := NewService(store)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.Routes(r)

	// create
	createBody := []byte(`{"userId":"u1","amount":10.5,"currency":"USD","status":"created"}`)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var created Payment
	if err := httpapi.DecodeJSON(w.Result().Request, &created); err == nil {
		// no-op; DecodeJSON expects request, not response
	}

	// Instead decode from body
	if err := httpapi.DecodeJSONFromBytes(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode created: %v", err)
	}

	// list
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/payments", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// get
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/payments/"+created.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// update
	updateBody := []byte(`{"amount":20.0}`)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/payments/"+created.ID, bytes.NewReader(updateBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/payments/"+created.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d: %s", w.Code, w.Body.String())
	}

	// get after delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/payments/"+created.ID, nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d: %s", w.Code, w.Body.String())
	}
}
