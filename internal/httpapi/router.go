package httpapi

import (
	"net/http"

	"amol-nv/user_crud/internal/store"
)

func NewRouter(st store.UserStore) http.Handler {
	mux := http.NewServeMux()
	RegisterUserRoutes(mux, st)
	return mux
}
