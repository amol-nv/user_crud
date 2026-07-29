package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"amol-nv/user_crud/internal/product"
)

type Handler struct {
	svc *product.Service
}

func NewHandler(svc *product.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/products", h.handleProductsRoot)
	mux.HandleFunc("/products/", h.handleProductByID)
}

func (h *Handler) handleProductsRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	}
}

func (h *Handler) handleProductByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r, id)
	case http.MethodPut:
		h.handleUpdate(w, r, id)
	case http.MethodDelete:
		h.handleDelete(w, r, id)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req product.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	// Simple ID generation: use name+timestamp-ish via monotonic counter would be better,
	// but keep it deterministic without extra deps.
	id := newID(req.Name)

	resp, err := h.svc.Create(id, req)
	if err != nil {
		if err == product.ErrInvalid {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid product"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request, id string) {
	resp, err := h.svc.GetByID(id)
	if err != nil {
		if err == product.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request, id string) {
	var req product.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	resp, err := h.svc.Update(id, req)
	if err != nil {
		if err == product.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		if err == product.ErrInvalid {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid product"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.svc.Delete(id); err != nil {
		if err == product.ErrNotFound {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
