package somidax

import (
	"context"
	"fmt"
	"time"
)

// OrdersService handles order management endpoints.
type OrdersService struct{ c *Client }

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// Order represents a Somidax order.
type Order struct {
	ID          string      `json:"id"`
	MerchantID  string      `json:"merchant_id"`
	CustomerID  string      `json:"customer_id"`
	Status      OrderStatus `json:"status"`
	Items       []OrderItem `json:"items"`
	Subtotal    Money       `json:"subtotal"`
	Tax         Money       `json:"tax"`
	Shipping    Money       `json:"shipping"`
	Total       Money       `json:"total"`
	Currency    string      `json:"currency"`
	RewardEarned *RewardAmount `json:"reward_earned,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// OrderItem is a line item within an order.
type OrderItem struct {
	ProductID  string  `json:"product_id"`
	VariantID  string  `json:"variant_id,omitempty"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku,omitempty"`
	Quantity   int     `json:"quantity"`
	UnitPrice  Money   `json:"unit_price"`
	TotalPrice Money   `json:"total_price"`
}

// Money represents a monetary value with currency.
type Money struct {
	Amount   int64  `json:"amount"`   // in minor units (e.g. pence, cents)
	Currency string `json:"currency"` // ISO 4217
}

// RewardAmount represents $SMDX tokens earned on an order.
type RewardAmount struct {
	SMDX    float64 `json:"smdx"`
	USDValue float64 `json:"usd_value"`
}

// CreateOrderRequest is the payload to create a new order.
type CreateOrderRequest struct {
	CustomerID string            `json:"customer_id"`
	Items      []CreateOrderItem `json:"items"`
	Currency   string            `json:"currency"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// CreateOrderItem specifies a product line item when creating an order.
type CreateOrderItem struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id,omitempty"`
	Quantity  int    `json:"quantity"`
}

// UpdateOrderRequest is the payload to update order status or metadata.
type UpdateOrderRequest struct {
	Status   *OrderStatus      `json:"status,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ListOrdersParams extends ListParams with order-specific filters.
type ListOrdersParams struct {
	ListParams
	Status     OrderStatus `json:"status,omitempty"`
	CustomerID string      `json:"customer_id,omitempty"`
}

// Create places a new order.
func (s *OrdersService) Create(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	var out Order
	return &out, s.c.do(ctx, "POST", "/orders", req, &out)
}

// Get retrieves a single order by ID.
func (s *OrdersService) Get(ctx context.Context, id string) (*Order, error) {
	var out Order
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/orders/%s", id), nil, &out)
}

// List returns a paginated list of orders.
func (s *OrdersService) List(ctx context.Context, params ListOrdersParams) (*ListResponse[Order], error) {
	var out ListResponse[Order]
	return &out, s.c.do(ctx, "GET", "/orders", params, &out)
}

// Update changes the status or metadata of an existing order.
func (s *OrdersService) Update(ctx context.Context, id string, req UpdateOrderRequest) (*Order, error) {
	var out Order
	return &out, s.c.do(ctx, "PATCH", fmt.Sprintf("/orders/%s", id), req, &out)
}

// Cancel cancels an order. Returns the updated order.
func (s *OrdersService) Cancel(ctx context.Context, id string) (*Order, error) {
	status := OrderStatusCancelled
	return s.Update(ctx, id, UpdateOrderRequest{Status: &status})
}
