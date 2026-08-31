// Package somidax provides a Go client for the Somidax API.
// Documentation: https://somidax.net/developers
package somidax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL = "https://api.somidax.net/v1"
	defaultTimeout = 30 * time.Second
	sdkVersion     = "1.0.0"
)

// Client is the top-level Somidax API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Orders    *OrdersService
	Catalog   *CatalogService
	Payments  *PaymentsService
	Rewards   *RewardsService
	Analytics *AnalyticsService
	Webhooks  *WebhooksService
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL (useful for sandbox/testing).
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithHTTPClient replaces the default http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithTimeout sets a custom request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// New creates an authenticated Somidax API client.
//
//	client := somidax.New("sk_live_your_api_key")
func New(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	for _, o := range opts {
		o(c)
	}

	c.Orders = &OrdersService{c}
	c.Catalog = &CatalogService{c}
	c.Payments = &PaymentsService{c}
	c.Rewards = &RewardsService{c}
	c.Analytics = &AnalyticsService{c}
	c.Webhooks = &WebhooksService{c}

	return c
}

// do executes an HTTP request and decodes the JSON response into v.
func (c *Client) do(ctx context.Context, method, path string, body, v any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("somidax: marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return fmt.Errorf("somidax: build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "somidax-go/"+sdkVersion)
	req.Header.Set("X-SDK-Version", sdkVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("somidax: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("somidax: decode response: %w", err)
		}
	}
	return nil
}

// APIError represents an error returned by the Somidax API.
type APIError struct {
	StatusCode int
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("somidax: API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

func parseAPIError(resp *http.Response) error {
	apiErr := &APIError{StatusCode: resp.StatusCode}
	_ = json.NewDecoder(resp.Body).Decode(apiErr)
	if apiErr.Message == "" {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}
	return apiErr
}

// ListResponse wraps paginated list endpoints.
type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// ListParams holds common pagination and filter parameters.
type ListParams struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
}
