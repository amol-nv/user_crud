package payments

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"amol-nv/user_crud/internal/httpapi"
	"amol-nv/user_crud/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type paymentTestEnv struct {
	r   chi.Router
	app *httptest.Server
	s   *store.MemoryPaymentStore
}

func setupPaymentHandler(t *testing.T) *paymentTestEnv {
	t.Helper()
	mem := store.NewMemoryPaymentStore()
	svc := NewService(mem)
	h := NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	srv := httptest.NewServer(r)
	_ = httpapi.JSONContentType
	return &paymentTestEnv{r: r, app: srv, s: mem}
}

func TestPaymentCRUD(t *testing.T) {
	env := setupPaymentHandler(t)
	defer env.app.Close()

	// Create
	createBody := map[string]any{"amount": 10.5, "currency": "USD", "status": "created"}
	b, _ := json.Marshal(createBody)
	resp, err := http.Post(env.app.URL+"/payments", "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created Payment
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&created))
	require.NotEmpty(t, created.ID)
	require.Equal(t, 10.5, created.Amount)
	require.Equal(t, "USD", created.Currency)

	// Get
	resp2, err := http.Get(env.app.URL + "/payments/" + created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	var got Payment
	require.NoError(t, json.NewDecoder(resp2.Body).Decode(&got))
	require.Equal(t, created.ID, got.ID)

	// Update
	newAmount := 20.0
	updateBody := map[string]any{"amount": newAmount, "status": "paid"}
	ub, _ := json.Marshal(updateBody)
	req, _ := http.NewRequest(http.MethodPut, env.app.URL+"/payments/"+created.ID, bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	resp3, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp3.StatusCode)
	var updated Payment
	require.NoError(t, json.NewDecoder(resp3.Body).Decode(&updated))
	require.Equal(t, 20.0, updated.Amount)
	require.Equal(t, "paid", updated.Status)

	// List
	resp4, err := http.Get(env.app.URL + "/payments")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp4.StatusCode)
	var list []*Payment
	require.NoError(t, json.NewDecoder(resp4.Body).Decode(&list))
	require.NotEmpty(t, list)

	// Delete
	reqDel, _ := http.NewRequest(http.MethodDelete, env.app.URL+"/payments/"+created.ID, nil)
	resp5, err := http.DefaultClient.Do(reqDel)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp5.StatusCode)

	// Get after delete
	resp6, err := http.Get(env.app.URL + "/payments/" + created.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp6.StatusCode)
}
