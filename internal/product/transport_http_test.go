package product

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPProductCRUD(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewService(repo)
	h := NewHTTPServer(svc)

	srv := httptest.NewServer(h.Routes())
	defer srv.Close()

	// Create
	createBody := CreateProductRequest{Name: "Laptop", Description: "Z", Price: 99.9}
	b, _ := json.Marshal(createBody)
	resp, err := http.Post(srv.URL+"/products", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("post err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var created Product
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode err: %v", err)
	}

	// Get
	resp, err = http.Get(srv.URL + "/products/" + itoa(int64(created.ID)))
	if err != nil {
		t.Fatalf("get err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Update
	updateBody := UpdateProductRequest{Name: "Laptop2", Description: "Z2", Price: 120}
	ub, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/products/"+itoa(int64(created.ID)), bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("put err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Delete
	req, _ = http.NewRequest(http.MethodDelete, srv.URL+"/products/"+itoa(int64(created.ID)), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("delete err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
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
