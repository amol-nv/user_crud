package httpapi

import (
	"net/http"

	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/payments"
	"github.com/amol-nv/user_crud/internal/users"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	r chi.Router
}

func NewRouter(inventoryHandler *inventory.Handler, userHandler *users.Handler, paymentHandler *payments.Handler) http.Handler {
	r := chi.NewRouter()

	if inventoryHandler != nil {
		inventoryHandler.Routes(r)
	}
	if userHandler != nil {
		userHandler.Routes(r)
	}
	if paymentHandler != nil {
		paymentHandler.Routes(r)
	}

	return r
}
