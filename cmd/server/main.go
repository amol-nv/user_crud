package main

import (
	"log"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/users"
)

func main() {
	st := store.NewInMemoryUserStore()
	svc := users.NewService(st)

	r := httpapi.NewRouter(svc)

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
