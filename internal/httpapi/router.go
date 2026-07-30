package httpapi

import (
	"net/http"

	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/payment"
	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/users"
)

// Router wires HTTP handlers.
type Router struct {
	mux *http.ServeMux
}

func NewRouter(invSvc *inventory.Service, userSvc *users.Service, paySvc *payment.Service) *Router {
	m := http.NewServeMux()
	mux := &Router{mux: m}

	invHandler := NewInventoryHandler(invSvc)
	userHandler := NewUserHandler(userSvc)
	payHandler := NewPaymentHandler(paySvc)

	m.HandleFunc("/inventory", invHandler.GetInventory)
	m.HandleFunc("/users", userHandler.PostUser)
	m.HandleFunc("/payments", payHandler.PostPayment)

	return mux
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
