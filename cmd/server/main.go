package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/payments"
	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/storage"
	"amol-nv/user_crud/internal/users"
)

func main() {
	userStore := store.NewInMemoryUserStore()
	userSvc := users.NewService(userStore)

	paymentStore := storage.NewInMemoryPaymentStore()
	paymentSvc := payments.NewService(paymentStore)

	r := httpapi.NewRouter(userSvc)
	payments.NewHandler(paymentSvc).RegisterRoutes(r)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
