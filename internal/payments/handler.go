package payments

import (
	"net/http"

	"amol-nv/user_crud/internal/httpapi"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createPaymentRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type updatePaymentRequest struct {
	Amount   *float64 `json:"amount"`
	Currency *string  `json:"currency"`
	Status   *string  `json:"status"`
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Post("", h.create)
		r.Get("", h.list)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createPaymentRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.svc.CreatePayment(CreatePaymentInput{Amount: req.Amount, Currency: req.Currency, Status: req.Status})
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusCreated, p)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ps, err := h.svc.ListPayments()
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetPayment(id)
	if err != nil {
		if err == ErrPaymentNotFound {
			httpapi.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpapi.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updatePaymentRequest
	if err := httpapi.DecodeJSON(r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}
	p, err := h.svc.UpdatePayment(id, UpdatePaymentInput{Amount: req.Amount, Currency: req.Currency, Status: req.Status})
	if err != nil {
		if err == ErrPaymentNotFound {
			httpapi.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.DeletePayment(id); err != nil {
		if err == ErrPaymentNotFound {
			httpapi.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpapi.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
