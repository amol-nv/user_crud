package main

import (
	"log"
	"net/http"

	"github.com/amol-nv/user_crud/internal/httpapi"
	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/payments"
	"github.com/amol-nv/user_crud/internal/store"
	"github.com/amol-nv/user_crud/internal/users"
)

func main() {
	// Existing repository uses in-memory stores.
	userStore := store.NewInMemoryUserStore()
	inventoryStore := store.NewInMemoryInventoryStore()
	paymentRepo := payments.NewInMemoryPaymentRepository()

	userSvc := users.NewService(userStore)
	inventorySvc := inventory.NewService(inventoryStore)
	paymentSvc := payments.NewService(paymentRepo)

	router := httpapi.NewRouter(
		httpapi.WithUserService(userSvc),
		httpapi.WithInventoryService(inventorySvc),
		httpapi.WithPaymentService(paymentSvc),
	)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
