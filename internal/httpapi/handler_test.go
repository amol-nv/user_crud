package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/product"
	"amol-nv/user_crud/internal/repository"
)

func TestHandlers_CRUD(t *testing.T) {
	repo := repository.NewInMemoryProductRepository()
	svc := product.NewService(repo)
	h := NewHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create
	createBody := map[string]any{"name": "Laptop", "price": 1000, "description": "Dev"}
	b, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(b))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", w.Code)
	}

	var created product.ProductResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if created.ID == "" {
		t.Fatalf("expected id")
	}

	// Get
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	// Update
	updBody := map[string]any{"name": "Laptop2", "price": 1100, "description": "Dev2"}
	b, _ = json.Marshal(updBody)
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPut, "/products/"+created.ID, bytes.NewReader(b))
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	// Delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", w.Code)
	}

	// Get after delete
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/products/"+created.ID, nil)
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}
