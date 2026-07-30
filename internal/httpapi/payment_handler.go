package httpapi

import (
	"encoding/json"
	"net/http"

	"amol-nv/user_crud/internal/payment"
)

// PaymentHandler exposes HTTP endpoints for payment.
type PaymentHandler struct {
	service *payment.Service
}

func NewPaymentHandler(service *payment.Service) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// PostPayment handles POST /payments.
func (h *PaymentHandler) PostPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req payment.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	resp, err := h.service.PostPayment(r.Context(), req)
	if err != nil {
		// Keep it simple: validation errors -> 400, provider errors -> 502.
		// If the error message contains validation markers, treat as 400.
		// (Repo has its own error mapping; we keep this scoped.)
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
