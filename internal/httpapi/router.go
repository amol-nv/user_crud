package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	r *mux.Router
}

func NewRouter() *Router {
	return &Router{r: mux.NewRouter()}
}

func (rt *Router) Handle(method, path string, handler http.HandlerFunc) {
	rt.r.HandleFunc(path, handler).Methods(method)
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.r.ServeHTTP(w, r)
}
