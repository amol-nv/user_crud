package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/store"
	"github.com/amol-nv/user_crud/internal/users"
	"github.com/google/uuid"
)

func TestWishlistAddAndList(t *testing.T) {
	invStore := store.NewInventoryMemoryStore()
	invSvc := inventory.NewService(invStore)

	userStore := store.NewUserMemoryStore()
	userSvc := users.NewService(userStore)

	wishStore := store.NewWishlistMemoryStore()
	wishSvc := wishlist.NewService(wishStore, invSvc)
	wishHandler := wishlist.NewHandler(wishSvc)

	r := NewRouter()
	// register existing user/inventory routes
	userHandler := users.NewHandler(userSvc)
	userHandler.RegisterRoutes(r)
	invHandler := inventory.NewHandler(invSvc)
	invHandler.RegisterRoutes(r)
	wishHandler.RegisterRoutes(r)

	// create user
	userID := uuid.New()
	_, _ = userStore.Create(t.Context(), userID, "a@b.com", "pass")

	// create product
	productID := uuid.New()
	_, _ = invStore.Create(t.Context(), productID, "p1", 10)

	// auth token: repo uses httpapi.UserIDFromRequest; tests should set header accordingly.
	// We'll set X-User-ID which is what UserIDFromRequest expects.
	addBody := map[string]any{"productId": productID}
	b, _ := json.Marshal(addBody)

	req := httptest.NewRequest(http.MethodPost, "/api/wishlist/items", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID.String())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	// list
	req2 := httptest.NewRequest(http.MethodGet, "/api/wishlist/items", nil)
	req2.Header.Set("X-User-ID", userID.String())
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w2.Code, w2.Body.String())
	}
}
