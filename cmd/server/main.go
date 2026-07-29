package main

import (
	"log"
	"net/http"
	"os"

	"amol-nv/user_crud/internal/handlers"
	"amol-nv/user_crud/internal/product"
	"amol-nv/user_crud/internal/store"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	st := store.NewInMemoryProductStore()
	svc := product.NewProductService(st)
	h := handlers.NewProductHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
