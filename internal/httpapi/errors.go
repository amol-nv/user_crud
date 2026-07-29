package httpapi

import (
	"errors"
	"net/http"

	"amol-nv/user_crud/internal/store"
	"amol-nv/user_crud/internal/users"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	writeJSON(w, status, v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}

func WriteErrorFromService(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, users.ErrNotFound), errors.Is(err, store.ErrNotFound):
		WriteError(w, http.StatusNotFound, "not found")
	case errors.Is(err, users.ErrValidation):
		WriteError(w, http.StatusBadRequest, "validation error")
	case errors.Is(err, users.ErrInvalidID):
		WriteError(w, http.StatusBadRequest, "invalid id")
	default:
		WriteError(w, http.StatusInternalServerError, "internal error")
	}
}
