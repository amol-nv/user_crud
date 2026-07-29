package httpapi

import (
	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/payment"
	"amol-nv/user_crud/internal/storage"
	"amol-nv/user_crud/internal/users"
	"github.com/go-chi/chi/v5"
)

func NewRouter(invSvc *inventory.Service, userSvc *users.Service, paymentSvc *payment.Service) chi.Router {
	r := chi.NewRouter()

	// existing routes
	invHandler := NewInventoryHandler(invSvc)
	invHandler.RegisterRoutes(r)

	userHandler := NewUserHandler(userSvc)
	userHandler.RegisterRoutes(r)

	// payment routes
	paymentHandler := NewPaymentHandler(paymentSvc)
	paymentHandler.RegisterRoutes(r)

	return r
}
