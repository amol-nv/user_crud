package users

import (
	"encoding/json"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/store"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := h.svc.Create(req)
	if err != nil {
		httpapi.WriteErrorFromService(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, u)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	u, err := h.svc.Get(id)
	if err != nil {
		httpapi.WriteErrorFromService(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, err := h.svc.Update(id, req)
	if err != nil {
		httpapi.WriteErrorFromService(w, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, u)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.svc.Delete(id); err != nil {
		httpapi.WriteErrorFromService(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Ensure store errors map correctly.
var _ = store.ErrNotFound
