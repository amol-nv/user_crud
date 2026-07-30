package httpapi

import (
	"net/http"

	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/payments"
	"github.com/amol-nv/user_crud/internal/users"
	"github.com/go-chi/chi/v5"
)

type routerConfig struct {
	userSvc      users.Service
	inventorySvc inventory.Service
	paymentSvc   payments.Service
}

type RouterOption func(*routerConfig)

func WithUserService(s users.Service) RouterOption {
	return func(c *routerConfig) { c.userSvc = s }
}

func WithInventoryService(s inventory.Service) RouterOption {
	return func(c *routerConfig) { c.inventorySvc = s }
}

func WithPaymentService(s payments.Service) RouterOption {
	return func(c *routerConfig) { c.paymentSvc = s }
}

func NewRouter(opts ...RouterOption) http.Handler {
	cfg := &routerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	r := chi.NewRouter()

	if cfg.userSvc != nil {
		uh := NewUserHandler(cfg.userSvc)
		r.Route("/users", func(r chi.Router) {
			r.Post("/", uh.CreateUser)
			r.Get("/{id}", uh.GetUser)
			r.Put("/{id}", uh.UpdateUser)
			r.Delete("/{id}", uh.DeleteUser)
		})
	}

	if cfg.inventorySvc != nil {
		ih := NewInventoryHandler(cfg.inventorySvc)
		r.Route("/inventory", func(r chi.Router) {
			r.Post("/", ih.CreateInventory)
			r.Get("/{id}", ih.GetInventory)
			r.Put("/{id}", ih.UpdateInventory)
			r.Delete("/{id}", ih.DeleteInventory)
		})
	}

	if cfg.paymentSvc != nil {
		ph := NewPaymentHandler(cfg.paymentSvc)
		r.Route("/payments", func(r chi.Router) {
			r.Post("/", ph.CreatePayment)
		})
	}

	return r
}
