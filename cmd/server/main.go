package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/payment"
	"amol-nv/user_crud/internal/storage"
	"amol-nv/user_crud/internal/users"
	"amol-nv/user_crud/internal/httpapi"
)

func main() {
	invStore := storage.NewInventoryStore()
	invSvc := inventory.NewService(invStore)

	userStore := storage.NewUserStore()
	userSvc := users.NewService(userStore)

	paymentStore := storage.NewPaymentStore()
	paymentSvc := payment.NewService(paymentStore)

	r := httpapi.NewRouter(invSvc, userSvc, paymentSvc)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
