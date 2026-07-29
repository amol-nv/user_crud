package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/product"
	"amol-nv/user_crud/internal/store"
)

func TestProductHandler_CreateGetUpdateDelete(t *testing.T) {
	st := store.NewInMemoryProductStore()
	svc := product.NewProductService(st)
	h := NewProductHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// create
	body := map[string]any{"name": "Widget", "price": 9.99, "description": "ok"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var created product.ProductResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id")
	}

	// get
	w = httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, getReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// update
	newPrice := 10.5
	upd := map[string]any{"price": newPrice}
	ub, _ := json.Marshal(upd)
	w = httptest.NewRecorder()
	updReq := httptest.NewRequest(http.MethodPatch, "/products/"+created.ID, bytes.NewReader(ub))
	mux.ServeHTTP(w, updReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// delete
	w = httptest.NewRecorder()
	delReq := httptest.NewRequest(http.MethodDelete, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, delReq)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}

	// get after delete
	w = httptest.NewRecorder()
	getReq = httptest.NewRequest(http.MethodGet, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, getReq)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestProductHandler_Validation(t *testing.T) {
	st := store.NewInMemoryProductStore()
	svc := product.NewProductService(st)
	h := NewProductHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// invalid create (missing name)
	body := map[string]any{"price": 1.0}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
