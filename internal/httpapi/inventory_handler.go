package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"amol-nv/user_crud/internal/inventory"
)

type InventoryHandler struct {
	service *inventory.Service
}

func NewInventoryHandler(service *inventory.Service) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/inventory", h.handleCollection)
	mux.HandleFunc("/inventory/", h.handleItem)
}

func (h *InventoryHandler) handleCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req inventory.CreateInventoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		inv, err := h.service.Create(r.Context(), req)
		if err != nil {
			writeError(w, httpStatusFromError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, inv)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *InventoryHandler) handleItem(w http.ResponseWriter, r *http.Request) {
	// /inventory/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/inventory/")
	idStr = strings.TrimSpace(idStr)
	if idStr == "" {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		inv, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			writeError(w, httpStatusFromError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, inv)
	case http.MethodPut:
		var req inventory.UpdateInventoryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		inv, err := h.service.Update(r.Context(), id, req)
		if err != nil {
			writeError(w, httpStatusFromError(err), err.Error())
			return
		}
		writeJSON(w, http.StatusOK, inv)
	case http.MethodDelete:
		if err := h.service.Delete(r.Context(), id); err != nil {
			writeError(w, httpStatusFromError(err), err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func httpStatusFromError(err error) int {
	switch err.(type) {
	case *inventory.NotFoundError:
		return http.StatusNotFound
	case *inventory.ValidationError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
