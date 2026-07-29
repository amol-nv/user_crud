package payments

import (
	"net/http"

	"github.com/amol-nv/user_crud/internal/httpapi"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Post("", h.create)
		r.Get("", h.list)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *Handler) create(w http.ResponseWriter, req *http.Request) {
	var body CreatePaymentRequest
	if err := httpapi.DecodeJSON(req, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payment, err := h.service.Create(body)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		httpapi.WriteError(w, status, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, payment)
}

func (h *Handler) list(w http.ResponseWriter, req *http.Request) {
	payments, err := h.service.List()
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, payments)
}

func (h *Handler) getByID(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	payment, err := h.service.GetByID(id)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		httpapi.WriteError(w, status, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, payment)
}

func (h *Handler) update(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var body UpdatePaymentRequest
	if err := httpapi.DecodeJSON(req, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}

	payment, err := h.service.Update(id, body)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		httpapi.WriteError(w, status, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, payment)
}

func (h *Handler) delete(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if err := h.service.Delete(id); err != nil {
		status := http.StatusBadRequest
		if err == ErrPaymentNotFound {
			status = http.StatusNotFound
		}
		httpapi.WriteError(w, status, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
