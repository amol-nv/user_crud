package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"amol-nv/user_crud/internal/models"
	"amol-nv/user_crud/internal/store"
)

type userCreateUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func RegisterUserRoutes(mux *http.ServeMux, st store.UserStore) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateUser(w, r, st)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user id"})
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r, st, id)
		case http.MethodPut:
			handleUpdateUser(w, r, st, id)
		case http.MethodDelete:
			handleDeleteUser(w, r, st, id)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		}
	})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, st store.UserStore) {
	var req userCreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := validateUserReq(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	u, err := st.CreateUser(context.Background(), models.User{Name: req.Name, Email: req.Email})
	if err != nil {
		if err == store.ErrConflict {
			writeJSON(w, http.StatusConflict, errorResponse{Error: "email already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusCreated, u)
}

func handleGetUser(w http.ResponseWriter, r *http.Request, st store.UserStore, id int64) {
	u, err := st.GetUserByID(context.Background(), id)
	if err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, st store.UserStore, id int64) {
	var req userCreateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	if err := validateUserReq(req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	u, err := st.UpdateUser(context.Background(), id, models.User{Name: req.Name, Email: req.Email})
	if err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "user not found"})
			return
		}
		if err == store.ErrConflict {
			writeJSON(w, http.StatusConflict, errorResponse{Error: "email already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}

	writeJSON(w, http.StatusOK, u)
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, st store.UserStore, id int64) {
	if err := st.DeleteUser(context.Background(), id); err != nil {
		if err == store.ErrNotFound {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateUserReq(req userCreateUpdateRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return &validationError{msg: "name is required"}
	}
	if strings.TrimSpace(req.Email) == "" {
		return &validationError{msg: "email is required"}
	}
	// Minimal email validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return &validationError{msg: "email is invalid"}
	}
	return nil
}

type validationError struct{ msg string }

func (e *validationError) Error() string { return e.msg }

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}

// Ensure time is referenced in this package for go vet in some setups.
var _ = time.Time{}
