package main

import (
	"log"
	"net/http"
	"os"

	"amol-nv/user_crud/internal/product"
)

func main() {
	repo := product.NewMemoryRepository()
	svc := product.NewService(repo)
	httpServer := product.NewHTTPServer(svc)

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	log.Printf("product service listening on %s", addr)
	if err := http.ListenAndServe(addr, httpServer.Routes()); err != nil {
		log.Fatal(err)
	}
}
