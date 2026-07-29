package payments

import (
	"encoding/json"
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *httpapi.Router) {
	r.Handle(http.MethodPost, "/payments", h.create)
	r.Handle(http.MethodGet, "/payments/{id}", h.getByID)
	r.Handle(http.MethodPut, "/payments/{id}", h.update)
	r.Handle(http.MethodDelete, "/payments/{id}", h.delete)
}

func (h *Handler) create(w http.ResponseWriter, req *http.Request) {
	var body CreatePaymentRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, httpapi.ErrBadRequest)
		return
	}

	resp, err := h.svc.Create(req.Context(), body)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		if err == ErrInvalid {
			status = http.StatusBadRequest
		}
		httpapi.WriteError(w, status, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) getByID(w http.ResponseWriter, req *http.Request) {
	id := httpapi.Param(req, "id")
	resp, err := h.svc.GetByID(req.Context(), id)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		if err == ErrInvalid {
			status = http.StatusBadRequest
		}
		httpapi.WriteError(w, status, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) update(w http.ResponseWriter, req *http.Request) {
	id := httpapi.Param(req, "id")

	var body UpdatePaymentRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, httpapi.ErrBadRequest)
		return
	}

	resp, err := h.svc.Update(req.Context(), id, body)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		if err == ErrInvalid {
			status = http.StatusBadRequest
		}
		httpapi.WriteError(w, status, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) delete(w http.ResponseWriter, req *http.Request) {
	id := httpapi.Param(req, "id")
	if err := h.svc.Delete(req.Context(), id); err != nil {
		status := http.StatusBadRequest
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		if err == ErrInvalid {
			status = http.StatusBadRequest
		}
		httpapi.WriteError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
