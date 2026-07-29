package main

import (
	"log"
	"net/http"

	"github.com/amol-nv/user_crud/internal/httpapi"
	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/amol-nv/user_crud/internal/payments"
	"github.com/amol-nv/user_crud/internal/storage"
	"github.com/amol-nv/user_crud/internal/store"
	"github.com/amol-nv/user_crud/internal/users"
)

func main() {
	invStore := storage.NewInventoryStore()
	invSvc := inventory.NewService(invStore)
	invHandler := inventory.NewHandler(invSvc)

	userStore := store.NewUserMemoryStore()
	userSvc := users.NewService(userStore)
	userHandler := users.NewHandler(userSvc)

	paymentStore := payments.NewMemoryPaymentStore()
	paymentSvc := payments.NewService(paymentStore)
	paymentHandler := payments.NewHandler(paymentSvc)

	r := httpapi.NewRouter(invHandler, userHandler, paymentHandler)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
