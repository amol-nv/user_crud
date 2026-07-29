package httpapi

import (
	"context"
	"errors"
	"net/http"

	"amol-nv/user_crud/internal/payment"
	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	svc *payment.Service
}

func NewPaymentHandler(svc *payment.Service) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) RegisterRoutes(r chi.Router) {
	r.Route("/payments", func(r chi.Router) {
		r.Post("", h.create)
		r.Get("/{id}", h.getByID)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *PaymentHandler) create(w http.ResponseWriter, r *http.Request) {
	var req payment.CreatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.svc.Create(context.Background(), req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, payment.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusCreated, paymentResponse(p))
}

func (h *PaymentHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetByID(context.Background(), id)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, payment.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusOK, paymentResponse(p))
}

func (h *PaymentHandler) update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req payment.UpdatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.svc.Update(context.Background(), id, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, payment.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	writeJSON(w, http.StatusOK, paymentResponse(p))
}

func (h *PaymentHandler) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.svc.Delete(context.Background(), id); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, payment.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func paymentResponse(p payment.Payment) payment.PaymentResponse {
	return payment.PaymentResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Amount:    p.Amount,
		Currency:  p.Currency,
		Status:    p.Status,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
