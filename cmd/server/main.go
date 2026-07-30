package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/payment"
	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/users"
	"amol-nv/user_crud/internal/httpapi"
)

func main() {
	// Inventory + Users are already wired in this repo.
	invStore := store.NewInventoryStore()
	invSvc := inventory.NewService(invStore)

	userStore := store.NewUserStore()
	userSvc := users.NewService(userStore)

	// Payment provider configuration.
	paymentBaseURL := os.Getenv("PAYMENT_BASE_URL")
	if paymentBaseURL == "" {
		paymentBaseURL = "http://localhost:8081"
	}
	paymentPostPath := os.Getenv("PAYMENT_POST_PATH")
	if paymentPostPath == "" {
		paymentPostPath = "/v1/payments"
	}
	paymentAPIKey := os.Getenv("PAYMENT_API_KEY")

	client := &http.Client{Timeout: 10 * time.Second}
	payStore := payment.NewHTTPStore(client, paymentBaseURL, paymentPostPath, paymentAPIKey)
	paySvc := payment.NewService(payStore)

	router := httpapi.NewRouter(invSvc, userSvc, paySvc)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, router.Handler()); err != nil {
		log.Fatal(err)
	}
}
