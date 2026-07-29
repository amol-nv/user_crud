package httpapi

import (
	"net/http"

	"amol-nv/user_crud/internal/users"
)

func NewRouter(svc *users.Service) http.Handler {
	mux := http.NewServeMux()

	// Collection
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h := users.NewHandler(svc)
			h.CreateUser(w, r)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})

	// Resource
	mux.HandleFunc("/api/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/api/v1/users/"):]
		if id == "" {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}

		h := users.NewHandler(svc)
		switch r.Method {
		case http.MethodGet:
			h.GetUser(w, r, id)
		case http.MethodPut:
			h.UpdateUser(w, r, id)
		case http.MethodDelete:
			h.DeleteUser(w, r, id)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		}
	})

	return withJSONMiddleware(mux)
}
