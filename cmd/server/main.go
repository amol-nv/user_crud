package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/product"
	"amol-nv/user_crud/internal/repository"
)

func main() {
	repo := repository.NewInMemoryProductRepository()
	service := product.NewService(repo)
	h := httpapi.NewHandler(service)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
