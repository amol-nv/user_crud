package httpapi

import (
	"net/http"

	"amol-nv/user_crud/internal/payments"
	"amol-nv/user_crud/internal/users"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	r *chi.Mux
}

func NewRouter(userSvc *users.Service, paymentsSvc *payments.PaymentsService) http.Handler {
	r := chi.NewRouter()

	uh := users.NewHandler(userSvc)
	r.Post("/users", uh.Create)
	r.Get("/users/{id}", uh.GetByID)
	r.Put("/users/{id}", uh.Update)
	r.Delete("/users/{id}", uh.Delete)

	if paymentsSvc != nil {
		ph := payments.NewHandler(paymentsSvc)
		r.Post("/payments", ph.Create)
		r.Get("/payments/{id}", ph.GetByID)
		r.Put("/payments/{id}", ph.Update)
		r.Delete("/payments/{id}", ph.Delete)
	}

	return r
}

func Param(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}
