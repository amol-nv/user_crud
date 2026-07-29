package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/inventory"
	"amol-nv/user_crud/internal/storage"
)

func main() {
	store := storage.NewInMemoryInventoryStore()
	service := inventory.NewService(store)
	h := httpapi.NewInventoryHandler(service)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
