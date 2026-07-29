package payments

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type PaymentHandler struct {
	svc *PaymentService
}

func NewPaymentHandler(svc *PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: svc}
}

func (h *PaymentHandler) Routes(r chi.Router) {
	r.Post("/payments", h.createPayment)
	r.Get("/payments/{id}", h.getPayment)
	r.Put("/payments/{id}", h.updatePayment)
	r.Patch("/payments/{id}", h.updatePayment)
	r.Delete("/payments/{id}", h.deletePayment)
}

func (h *PaymentHandler) createPayment(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	p, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CreatePaymentResponse{Payment: p})
}

func (h *PaymentHandler) getPayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, GetPaymentResponse{Payment: p})
}

func (h *PaymentHandler) updatePayment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	p, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, UpdatePaymentResponse{Payment: p})
}

func (h *PaymentHandler) deletePayment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	p, err := h.svc.Delete(r.Context(), id)
	if err != nil {
		writePaymentError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, DeletePaymentResponse{DeletedID: p.ID, DeletedAt: p.UpdatedAt})
}

func writePaymentError(w http.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case ErrInvalidArgument:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
