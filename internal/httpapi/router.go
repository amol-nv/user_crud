package httpapi

import (
	"net/http"

	"amol-nv/user_crud/internal/users"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(userSvc *users.Service) *Router {
	r := &Router{mux: http.NewServeMux()}

	uh := users.NewHandler(userSvc)
	uh.RegisterRoutes(r)

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) Handle(method, pattern string, h func(http.ResponseWriter, *http.Request)) {
	// Minimal pattern support for {id} style.
	// We only need /payments/{id} for this ticket.
	if hasParam(pattern) {
		prefix := pattern[:len(pattern)-len("/{id}")]
		r.mux.HandleFunc(prefix, func(w http.ResponseWriter, req *http.Request) {
			if req.Method != method {
				WriteError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
				return
			}
			// Extract last segment as id.
			id := lastPathSegment(req.URL.Path)
			ctx := WithParam(req.Context(), "id", id)
			h(w, req.WithContext(ctx))
		})
		return
	}

	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			WriteError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
			return
		}
		h(w, req)
	})
}

func hasParam(pattern string) bool {
	return len(pattern) >= 5 && pattern[len(pattern)-5:] == "/{id}"
}

func lastPathSegment(path string) string {
	// naive split; good enough for this repo
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			if i+1 < len(path) {
				return path[i+1:]
			}
			return ""
		}
	}
	return path
}
