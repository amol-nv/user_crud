package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/storage"
)

func TestInventoryHandler_CRUD(t *testing.T) {
	store := storage.NewInMemoryInventoryStore()
	svc := inventory.NewService(store)
	h := NewInventoryHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create
	createBody := map[string]any{"product_id": 1, "quantity": 3}
	b, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewReader(b))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var created inventory.ProductInventory
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	// Get
	w = httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/inventory/"+itoa(created.ID), nil)
	mux.ServeHTTP(w, getReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Update
	updateBody := map[string]any{"product_id": 1, "quantity": 4}
	ub, _ := json.Marshal(updateBody)
	w = httptest.NewRecorder()
	updReq := httptest.NewRequest(http.MethodPut, "/inventory/"+itoa(created.ID), bytes.NewReader(ub))
	mux.ServeHTTP(w, updReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// Delete
	w = httptest.NewRecorder()
	delReq := httptest.NewRequest(http.MethodDelete, "/inventory/"+itoa(created.ID), nil)
	mux.ServeHTTP(w, delReq)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	// Get after delete
	w = httptest.NewRecorder()
	getReq = httptest.NewRequest(http.MethodGet, "/inventory/"+itoa(created.ID), nil)
	mux.ServeHTTP(w, getReq)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}

func TestInventoryHandler_InvalidJSON(t *testing.T) {
	store := storage.NewInMemoryInventoryStore()
	svc := inventory.NewService(store)
	h := NewInventoryHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/inventory", bytes.NewBufferString("{bad json"))
	req = req.WithContext(context.Background())
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func itoa(v int64) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		d := v % 10
		buf = append(buf, byte('0'+d))
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
