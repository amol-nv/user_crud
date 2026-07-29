package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/payments"
	"amol-nv/user_crud/internal/users"
	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/store"
)

func main() {
	invRepo := store.NewInventoryStore()
	invSvc := inventory.NewInventoryService(invRepo)
	invHandler := inventory.NewInventoryHandler(invSvc)

	userRepo := store.NewUserStore()
	userSvc := users.NewUserService(userRepo)
	userHandler := users.NewUserHandler(userSvc)

	payRepo := payments.NewMemoryPaymentRepository()
	paySvc := payments.NewPaymentService(payRepo)
	payHandler := payments.NewPaymentHandler(paySvc)

	r := httpapi.NewRouter()
	invHandler.Routes(r)
	userHandler.Routes(r)
	payHandler.Routes(r)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
