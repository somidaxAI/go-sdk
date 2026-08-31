package somidax

import (
	"context"
	"time"
)

// AnalyticsService provides merchant performance metrics.
type AnalyticsService struct{ c *Client }

// Period defines the time window for analytics queries.
type Period string

const (
	PeriodDay   Period = "day"
	PeriodWeek  Period = "week"
	PeriodMonth Period = "month"
	PeriodYear  Period = "year"
)

// OverviewMetrics contains high-level merchant KPIs.
type OverviewMetrics struct {
	Period         Period    `json:"period"`
	From           time.Time `json:"from"`
	To             time.Time `json:"to"`
	TotalRevenue   Money     `json:"total_revenue"`
	TotalOrders    int       `json:"total_orders"`
	AOV            Money     `json:"aov"` // average order value
	UniqueCustomers int      `json:"unique_customers"`
	NewCustomers   int       `json:"new_customers"`
	RepeatRate     float64   `json:"repeat_rate"`     // 0–1
	SMDXDistributed float64  `json:"smdx_distributed"` // tokens rewarded
	SMDXRedeemed   float64   `json:"smdx_redeemed"`   // tokens spent
}

// TopProduct summarises performance for a single product.
type TopProduct struct {
	ProductID  string  `json:"product_id"`
	Name       string  `json:"name"`
	Units      int     `json:"units"`
	Revenue    Money   `json:"revenue"`
	ReturnRate float64 `json:"return_rate"`
}

// AnalyticsParams controls the query window for analytics endpoints.
type AnalyticsParams struct {
	Period Period    `json:"period,omitempty"`
	From   time.Time `json:"from,omitempty"`
	To     time.Time `json:"to,omitempty"`
}

// Overview returns aggregated KPIs for the given time window.
func (s *AnalyticsService) Overview(ctx context.Context, params AnalyticsParams) (*OverviewMetrics, error) {
	var out OverviewMetrics
	return &out, s.c.do(ctx, "GET", "/analytics/overview", params, &out)
}

// TopProducts returns the best-performing products ranked by revenue.
func (s *AnalyticsService) TopProducts(ctx context.Context, params AnalyticsParams) ([]TopProduct, error) {
	var out []TopProduct
	return out, s.c.do(ctx, "GET", "/analytics/top-products", params, &out)
}
