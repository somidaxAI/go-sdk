package somidax

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebhooksService manages webhook endpoint registration and verification.
type WebhooksService struct{ c *Client }

// WebhookEvent is the type of event delivered to a webhook endpoint.
type WebhookEvent string

const (
	EventOrderCreated   WebhookEvent = "order.created"
	EventOrderUpdated   WebhookEvent = "order.updated"
	EventOrderCancelled WebhookEvent = "order.cancelled"
	EventPaymentSettled WebhookEvent = "payment.settled"
	EventPaymentFailed  WebhookEvent = "payment.failed"
	EventRefundIssued   WebhookEvent = "refund.issued"
	EventRewardEarned   WebhookEvent = "reward.earned"
	EventRewardRedeemed WebhookEvent = "reward.redeemed"
)

// Webhook represents a registered webhook endpoint.
type Webhook struct {
	ID        string         `json:"id"`
	URL       string         `json:"url"`
	Events    []WebhookEvent `json:"events"`
	Active    bool           `json:"active"`
	Secret    string         `json:"secret,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// CreateWebhookRequest registers a new webhook endpoint.
type CreateWebhookRequest struct {
	URL    string         `json:"url"`
	Events []WebhookEvent `json:"events"`
}

// WebhookPayload is the envelope delivered to registered endpoints.
type WebhookPayload struct {
	ID        string       `json:"id"`
	Event     WebhookEvent `json:"event"`
	CreatedAt time.Time    `json:"created_at"`
	Data      any          `json:"data"`
}

// Create registers a webhook endpoint and returns the signing secret.
// Store the secret securely — it is only shown once.
func (s *WebhooksService) Create(ctx context.Context, req CreateWebhookRequest) (*Webhook, error) {
	var out Webhook
	return &out, s.c.do(ctx, "POST", "/webhooks", req, &out)
}

// Get retrieves a registered webhook by ID.
func (s *WebhooksService) Get(ctx context.Context, id string) (*Webhook, error) {
	var out Webhook
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/webhooks/%s", id), nil, &out)
}

// List returns all registered webhook endpoints.
func (s *WebhooksService) List(ctx context.Context) ([]Webhook, error) {
	var out []Webhook
	return out, s.c.do(ctx, "GET", "/webhooks", nil, &out)
}

// Delete removes a webhook endpoint.
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	return s.c.do(ctx, "DELETE", fmt.Sprintf("/webhooks/%s", id), nil, nil)
}

// VerifySignature validates the HMAC-SHA256 signature on an incoming webhook.
// Pass the raw request body (before any JSON decoding) and the signing secret.
//
//	body, _ := io.ReadAll(r.Body)
//	if err := somidax.VerifySignature(r, body, secret); err != nil {
//	    http.Error(w, "invalid signature", http.StatusUnauthorized)
//	    return
//	}
func VerifySignature(r *http.Request, body []byte, secret string) error {
	sig := r.Header.Get("X-Somidax-Signature")
	if sig == "" {
		return fmt.Errorf("somidax: missing X-Somidax-Signature header")
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return fmt.Errorf("somidax: webhook signature mismatch")
	}
	return nil
}

// ReadBody reads the full request body without consuming it for further decoding.
func ReadBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}
