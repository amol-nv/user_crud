package httpapi

import (
	"net/http"

	"github.com/amol-nv/user_crud/internal/payments"
)

type PaymentHandler struct {
	payments payments.Service
}

func NewPaymentHandler(paymentsSvc payments.Service) *PaymentHandler {
	return &PaymentHandler{payments: paymentsSvc}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var req payments.CreatePaymentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.payments.CreatePayment(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}
