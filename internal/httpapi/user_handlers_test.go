package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/store"
)

type createReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type errResp struct {
	Error string `json:"error"`
}

func TestUserCRUD(t *testing.T) {
	st := store.NewInMemoryUserStore()
	router := NewRouter(st)

	// Create
	cr := createReq{Name: "Alice", Email: "alice@example.com"}
	body, _ := json.Marshal(cr)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created userResp
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected id > 0")
	}
	if created.Name != cr.Name || created.Email != cr.Email {
		t.Fatalf("unexpected created fields: %+v", created)
	}

	// Get
	w = httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/users/"+itoa(created.ID), nil)
	router.ServeHTTP(w, getReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var got userResp
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID != created.ID || got.Name != cr.Name || got.Email != cr.Email {
		t.Fatalf("unexpected get fields: %+v", got)
	}

	// Update
	ur := createReq{Name: "Alice Updated", Email: "alice2@example.com"}
	body, _ = json.Marshal(ur)
	w = httptest.NewRecorder()
	updReq := httptest.NewRequest(http.MethodPut, "/users/"+itoa(created.ID), bytes.NewReader(body))
	router.ServeHTTP(w, updReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var updated userResp
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if updated.ID != created.ID || updated.Name != ur.Name || updated.Email != ur.Email {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}

	// Delete
	w = httptest.NewRecorder()
	delReq := httptest.NewRequest(http.MethodDelete, "/users/"+itoa(created.ID), nil)
	router.ServeHTTP(w, delReq)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	// Get after delete => 404
	w = httptest.NewRecorder()
	getReq = httptest.NewRequest(http.MethodGet, "/users/"+itoa(created.ID), nil)
	router.ServeHTTP(w, getReq)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
	var er errResp
	_ = json.Unmarshal(w.Body.Bytes(), &er)
	if er.Error == "" {
		t.Fatalf("expected error message")
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
