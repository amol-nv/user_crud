package wishlist

import (
	"net/http"

	"github.com/amol-nv/user_crud/internal/httpapi"
	"github.com/amol-nv/user_crud/internal/inventory"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type addWishlistRequest struct {
	ProductID uuid.UUID `json:"productId"`
}

type addWishlistResponse struct {
	Saved bool `json:"saved"`
}

func (h *Handler) RegisterRoutes(r *httpapi.Router) {
	// POST /api/wishlist/items
	r.Handle(http.MethodPost, "/api/wishlist/items", h.add)
	// GET /api/wishlist/items
	r.Handle(http.MethodGet, "/api/wishlist/items", h.list)
}

func (h *Handler) add(w http.ResponseWriter, req *http.Request) {
	userID, ok := httpapi.UserIDFromRequest(req)
	if !ok {
		httpapi.WriteError(w, http.StatusUnauthorized, httpapi.ErrUnauthorized)
		return
	}

	var body addWishlistRequest
	if err := httpapi.DecodeJSON(req, &body); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.svc.Add(req.Context(), userID, body.ProductID); err != nil {
		if err == ErrProductNotFound {
			httpapi.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpapi.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Idempotent add: always return saved=true
	httpapi.WriteJSON(w, http.StatusOK, addWishlistResponse{Saved: true})
}

func (h *Handler) list(w http.ResponseWriter, req *http.Request) {
	userID, ok := httpapi.UserIDFromRequest(req)
	if !ok {
		httpapi.WriteError(w, http.StatusUnauthorized, httpapi.ErrUnauthorized)
		return
	}

	products, err := h.svc.List(req.Context(), userID)
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Reuse inventory.Product JSON shape
	var resp = struct {
		Items []inventory.Product `json:"items"`
	}{Items: products}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}
