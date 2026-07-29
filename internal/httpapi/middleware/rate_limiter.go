package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter is a simple in-memory fixed-window rate limiter.
// It limits requests per client IP.
//
// Default behavior:
// - window: 1 second
// - limit: 10 requests per window
// - status code: 429 Too Many Requests
//
// This implementation is concurrency-safe.
// It uses a per-IP state stored in a sync.Map.
//
// Note: This is intended for a single-process deployment.

type RateLimiter struct {
	limit  int
	window time.Duration
	// now is injectable for tests.
	now func() time.Time
	// clients maps ip -> *clientState
	clients sync.Map
}

type clientState struct {
	mu       sync.Mutex
	windowID int64
	count    int
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
		now:    time.Now,
	}
}

// SetNowForTest allows tests to inject a controllable time source.
func (r *RateLimiter) SetNowForTest(now func() time.Time) {
	r.now = now
}

func (r *RateLimiter) windowID(t time.Time) int64 {
	// Fixed window based on UnixNano / window.
	return t.UnixNano() / r.window.Nanoseconds()
}

func clientIPFromRequest(req *http.Request) string {
	// Prefer X-Forwarded-For only if present.
	// If multiple IPs are present, the first is the client.
	if xff := strings.TrimSpace(req.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// RemoteAddr is typically "ip:port".
	if host, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		if host != "" {
			return host
		}
	}
	// Fallback: use RemoteAddr as-is.
	return strings.TrimSpace(req.RemoteAddr)
}

// Middleware returns an http middleware that enforces the rate limit.
func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ip := clientIPFromRequest(req)
		if ip == "" {
			ip = "unknown"
		}

		stateAny, _ := r.clients.LoadOrStore(ip, &clientState{})
		state := stateAny.(*clientState)

		now := r.now()
		wid := r.windowID(now)

		state.mu.Lock()
		defer state.mu.Unlock()

		if state.windowID != wid {
			state.windowID = wid
			state.count = 0
		}

		if state.count >= r.limit {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		state.count++
		next.ServeHTTP(w, req)
	})
}
