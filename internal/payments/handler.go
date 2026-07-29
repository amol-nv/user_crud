package payments

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"amol-nv/user_crud/internal/httpapi"
)

type PaymentsHandler struct {
	svc *PaymentsService
}

func NewHandler(svc *PaymentsService) *PaymentsHandler {
	return &PaymentsHandler{svc: svc}
}

func (h *PaymentsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusCreated, resp)
}

func (h *PaymentsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(httpapi.Param(r, "id"))
	if id == "" {
		httpapi.WriteError(w, http.StatusBadRequest, "missing id")
		return
	}

	resp, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}

func (h *PaymentsHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(httpapi.Param(r, "id"))
	if id == "" {
		httpapi.WriteError(w, http.StatusBadRequest, "missing id")
		return
	}

	var req UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	resp, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		writePaymentError(w, err)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}

func (h *PaymentsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(httpapi.Param(r, "id"))
	if id == "" {
		httpapi.WriteError(w, http.StatusBadRequest, "missing id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writePaymentError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writePaymentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrPaymentNotFound):
		httpapi.WriteError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidPayment):
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
	default:
		httpapi.WriteError(w, http.StatusInternalServerError, "internal error")
	}
}
