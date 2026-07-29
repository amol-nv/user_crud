package httpapi

import (
	"amol-nv/user_crud/internal/payments"
	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/users"
	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *users.Handler, paymentHandler *payments.Handler) chi.Router {
	r := chi.NewRouter()

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, 200, map[string]string{"status": "ok"})
	})

	if userHandler != nil {
		userHandler.RegisterRoutes(r)
	}
	if paymentHandler != nil {
		paymentHandler.RegisterRoutes(r)
	}

	return r
}

// Backwards-compatible constructor used by main.
func Router(store *store.MemoryStore) chi.Router {
	userSvc := users.NewService(store.User())
	userHandler := users.NewHandler(userSvc)

	paymentSvc := payments.NewService(store.Payment())
	paymentHandler := payments.NewHandler(paymentSvc)

	return NewRouter(userHandler, paymentHandler)
}
