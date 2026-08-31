package somidax

import (
	"context"
	"fmt"
	"time"
)

// CatalogService handles product and variant management.
type CatalogService struct{ c *Client }

// Product represents a merchant product listing.
type Product struct {
	ID          string            `json:"id"`
	MerchantID  string            `json:"merchant_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"` // "active" | "draft" | "archived"
	Images      []string          `json:"images"`
	Tags        []string          `json:"tags,omitempty"`
	Variants    []Variant         `json:"variants"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Variant is a specific SKU of a product (e.g. size + colour combination).
type Variant struct {
	ID         string            `json:"id"`
	ProductID  string            `json:"product_id"`
	SKU        string            `json:"sku"`
	Name       string            `json:"name"`
	Price      Money             `json:"price"`
	CompareAt  *Money            `json:"compare_at,omitempty"`
	Stock      int               `json:"stock"`
	Attributes map[string]string `json:"attributes,omitempty"` // e.g. {"size":"L","colour":"Black"}
}

// CreateProductRequest is the payload for creating a product.
type CreateProductRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status,omitempty"` // default: "draft"
	Images      []string          `json:"images,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Variants    []CreateVariant   `json:"variants"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// CreateVariant specifies a variant when creating a product.
type CreateVariant struct {
	SKU        string            `json:"sku"`
	Name       string            `json:"name"`
	Price      int64             `json:"price"`    // minor units
	Currency   string            `json:"currency"` // ISO 4217
	Stock      int               `json:"stock"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// UpdateProductRequest is the payload for updating a product.
type UpdateProductRequest struct {
	Name        *string           `json:"name,omitempty"`
	Description *string           `json:"description,omitempty"`
	Status      *string           `json:"status,omitempty"`
	Images      []string          `json:"images,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Create creates a new product with one or more variants.
func (s *CatalogService) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	var out Product
	return &out, s.c.do(ctx, "POST", "/catalog/products", req, &out)
}

// Get retrieves a single product by ID.
func (s *CatalogService) Get(ctx context.Context, id string) (*Product, error) {
	var out Product
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/catalog/products/%s", id), nil, &out)
}

// List returns a paginated product catalogue for the authenticated merchant.
func (s *CatalogService) List(ctx context.Context, params ListParams) (*ListResponse[Product], error) {
	var out ListResponse[Product]
	return &out, s.c.do(ctx, "GET", "/catalog/products", params, &out)
}

// Update modifies an existing product.
func (s *CatalogService) Update(ctx context.Context, id string, req UpdateProductRequest) (*Product, error) {
	var out Product
	return &out, s.c.do(ctx, "PATCH", fmt.Sprintf("/catalog/products/%s", id), req, &out)
}

// Delete removes a product (sets status to "archived").
func (s *CatalogService) Delete(ctx context.Context, id string) error {
	archived := "archived"
	return s.c.do(ctx, "PATCH", fmt.Sprintf("/catalog/products/%s", id), UpdateProductRequest{Status: &archived}, nil)
}
