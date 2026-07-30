package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Store defines the persistence/IO boundary for payment.
// In this repo, "store" is used for external calls as well.
type Store interface {
	PostPayment(ctx context.Context, req PaymentRequest) (PaymentResponse, int, error)
}

type httpStore struct {
	client  *http.Client
	baseURL string
	// endpoint path for POST payment
	postPath string
	// optional headers
	apiKey string
}

// NewHTTPStore creates a Store that calls an external payment provider.
func NewHTTPStore(client *http.Client, baseURL, postPath, apiKey string) Store {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &httpStore{
		client:   client,
		baseURL:  baseURL,
		postPath: postPath,
		apiKey:   apiKey,
	}
}

func (s *httpStore) PostPayment(ctx context.Context, req PaymentRequest) (PaymentResponse, int, error) {
	var out PaymentResponse

	url := fmt.Sprintf("%s%s", s.baseURL, s.postPath)

	b, err := json.Marshal(req)
	if err != nil {
		return out, 0, fmt.Errorf("marshal payment request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return out, 0, fmt.Errorf("create http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))
	}

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return out, 0, fmt.Errorf("post payment: %w", err)
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	body, _ := io.ReadAll(resp.Body)

	if status < 200 || status >= 300 {
		// Best-effort error message.
		if len(body) > 0 {
			return out, status, fmt.Errorf("payment provider error: status=%d body=%s", status, string(body))
		}
		return out, status, fmt.Errorf("payment provider error: status=%d", status)
	}

	if len(body) == 0 {
		return out, status, nil
	}

	if err := json.Unmarshal(body, &out); err != nil {
		return out, status, fmt.Errorf("decode payment response: %w", err)
	}

	return out, status, nil
}
