package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiter_AllowsFirstTenRequests(t *testing.T) {
	limiter := NewRateLimiter(10, time.Second)
	base := time.Unix(0, 0)
	limiter.SetNowForTest(func() time.Time { return base })

	var called int32
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&called, 1)
		w.WriteHeader(http.StatusOK)
	})

	h := limiter.Middleware(next)

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	if got := atomic.LoadInt32(&called); got != 10 {
		t.Fatalf("expected next called 10 times, got %d", got)
	}
}

func TestRateLimiter_RejectsEleventhRequest(t *testing.T) {
	limiter := NewRateLimiter(10, time.Second)
	base := time.Unix(0, 0)
	limiter.SetNowForTest(func() time.Time { return base })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := limiter.Middleware(next)

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 11th should be rejected.
	req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}
}

func TestRateLimiter_IndependentPerIP(t *testing.T) {
	limiter := NewRateLimiter(10, time.Second)
	base := time.Unix(0, 0)
	limiter.SetNowForTest(func() time.Time { return base })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := limiter.Middleware(next)

	// IP A: send 11 requests -> last should be 429
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
		req.RemoteAddr = "10.0.0.1:1111"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("IP A request %d: expected 200, got %d", i+1, w.Code)
		}
	}
	{
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
		req.RemoteAddr = "10.0.0.1:1111"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusTooManyRequests {
			t.Fatalf("IP A 11th: expected 429, got %d", w.Code)
		}
	}

	// IP B: should still allow first 10 requests.
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "http://example.com/api", nil)
		req.RemoteAddr = "10.0.0.2:2222"
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("IP B request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}
