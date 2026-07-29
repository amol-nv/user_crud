package product

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type HTTPServer struct {
	svc *Service
}

func NewHTTPServer(svc *Service) *HTTPServer {
	return &HTTPServer{svc: svc}
}

func (h *HTTPServer) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/products", h.handleProducts)
	mux.HandleFunc("/products/", h.handleProductByID)
	return mux
}

func (h *HTTPServer) handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		p, err := h.svc.Create(r.Context(), req)
		if err != nil {
			writeProductError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, p)
	case http.MethodGet:
		items, err := h.svc.List(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
			return
		}
		writeJSON(w, http.StatusOK, items)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *HTTPServer) handleProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	if idStr == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := parseID(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := h.svc.GetByID(r.Context(), id)
		if err != nil {
			writeProductError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, p)
	case http.MethodPut:
		var req UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
			return
		}
		p, err := h.svc.Update(r.Context(), id, req)
		if err != nil {
			writeProductError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, p)
	case http.MethodDelete:
		err := h.svc.Delete(r.Context(), id)
		if err != nil {
			writeProductError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func parseID(s string) (ID, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v <= 0 {
		return 0, errors.New("invalid id")
	}
	return ID(v), nil
}

func writeProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, ErrValidation):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "validation error"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
