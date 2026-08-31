package somidax

import (
	"context"
	"fmt"
	"time"
)

// PaymentsService handles Somidax Pay settlement endpoints.
type PaymentsService struct{ c *Client }

// PaymentStatus represents the lifecycle of a payment.
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSettled   PaymentStatus = "settled"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentMethod specifies the accepted payment rail.
type PaymentMethod string

const (
	PaymentMethodFiat  PaymentMethod = "fiat"
	PaymentMethodSMDX  PaymentMethod = "smdx"
	PaymentMethodCrypto PaymentMethod = "crypto"
)

// Payment represents a Somidax Pay settlement record.
type Payment struct {
	ID              string        `json:"id"`
	OrderID         string        `json:"order_id"`
	MerchantID      string        `json:"merchant_id"`
	Amount          Money         `json:"amount"`
	Fee             Money         `json:"fee"`
	NetAmount       Money         `json:"net_amount"`
	Method          PaymentMethod `json:"method"`
	Status          PaymentStatus `json:"status"`
	SettledAt       *time.Time    `json:"settled_at,omitempty"`
	TxHash          string        `json:"tx_hash,omitempty"`
	FailureReason   string        `json:"failure_reason,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
}

// CreatePaymentRequest initiates a new payment settlement.
type CreatePaymentRequest struct {
	OrderID  string        `json:"order_id"`
	Amount   int64         `json:"amount"`   // minor units
	Currency string        `json:"currency"` // ISO 4217
	Method   PaymentMethod `json:"method"`
	// ReturnURL is used for redirect-based payment flows.
	ReturnURL string `json:"return_url,omitempty"`
}

// RefundRequest issues a full or partial refund on a settled payment.
type RefundRequest struct {
	Amount int64  `json:"amount,omitempty"` // omit for full refund
	Reason string `json:"reason,omitempty"`
}

// Refund represents a refund record.
type Refund struct {
	ID        string        `json:"id"`
	PaymentID string        `json:"payment_id"`
	Amount    Money         `json:"amount"`
	Reason    string        `json:"reason,omitempty"`
	Status    PaymentStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

// Create initiates a payment and returns the settlement record.
// Sub-2-second settlement for $SMDX payments; fiat settles async.
func (s *PaymentsService) Create(ctx context.Context, req CreatePaymentRequest) (*Payment, error) {
	var out Payment
	return &out, s.c.do(ctx, "POST", "/payments", req, &out)
}

// Get retrieves a payment by ID.
func (s *PaymentsService) Get(ctx context.Context, id string) (*Payment, error) {
	var out Payment
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/payments/%s", id), nil, &out)
}

// List returns paginated payments for the authenticated merchant.
func (s *PaymentsService) List(ctx context.Context, params ListParams) (*ListResponse[Payment], error) {
	var out ListResponse[Payment]
	return &out, s.c.do(ctx, "GET", "/payments", params, &out)
}

// Refund issues a refund on a settled payment.
// Pass an Amount of 0 to refund the full payment amount.
func (s *PaymentsService) Refund(ctx context.Context, paymentID string, req RefundRequest) (*Refund, error) {
	var out Refund
	return &out, s.c.do(ctx, "POST", fmt.Sprintf("/payments/%s/refund", paymentID), req, &out)
}
