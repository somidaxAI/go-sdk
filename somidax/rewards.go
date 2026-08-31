package somidax

import (
	"context"
	"fmt"
	"time"
)

// RewardsService handles $SMDX token loyalty operations.
type RewardsService struct{ c *Client }

// StakingTier represents a customer's $SMDX staking level.
type StakingTier string

const (
	TierSpark StakingTier = "spark" // 0–999 SMDX
	TierCore  StakingTier = "core"  // 1,000–9,999 SMDX
	TierSurge StakingTier = "surge" // 10,000–49,999 SMDX
	TierApex  StakingTier = "apex"  // 50,000+ SMDX
)

// RewardAccount holds a customer's $SMDX token balance and tier.
type RewardAccount struct {
	CustomerID  string      `json:"customer_id"`
	Balance     float64     `json:"balance"`      // SMDX tokens
	Staked      float64     `json:"staked"`       // locked for tier
	Tier        StakingTier `json:"tier"`
	TotalEarned float64     `json:"total_earned"` // lifetime earned
	TotalSpent  float64     `json:"total_spent"`  // lifetime redeemed
	UpdatedAt   time.Time   `json:"updated_at"`
}

// RewardTransaction records an earn or redeem event.
type RewardTransaction struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Type        string    `json:"type"`    // "earn" | "redeem" | "stake" | "unstake" | "burn"
	Amount      float64   `json:"amount"`  // SMDX tokens (positive = credit, negative = debit)
	OrderID     string    `json:"order_id,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// EarnRequest records tokens earned from a purchase.
type EarnRequest struct {
	CustomerID string  `json:"customer_id"`
	OrderID    string  `json:"order_id"`
	Amount     float64 `json:"amount"` // SMDX tokens to credit
}

// RedeemRequest burns reward tokens against a purchase.
type RedeemRequest struct {
	CustomerID string  `json:"customer_id"`
	OrderID    string  `json:"order_id,omitempty"`
	Amount     float64 `json:"amount"` // SMDX tokens to debit
}

// GetAccount retrieves a customer's reward balance and tier.
func (s *RewardsService) GetAccount(ctx context.Context, customerID string) (*RewardAccount, error) {
	var out RewardAccount
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/rewards/accounts/%s", customerID), nil, &out)
}

// Earn credits $SMDX tokens to a customer after a qualifying purchase.
func (s *RewardsService) Earn(ctx context.Context, req EarnRequest) (*RewardTransaction, error) {
	var out RewardTransaction
	return &out, s.c.do(ctx, "POST", "/rewards/earn", req, &out)
}

// Redeem debits $SMDX tokens from a customer's balance.
func (s *RewardsService) Redeem(ctx context.Context, req RedeemRequest) (*RewardTransaction, error) {
	var out RewardTransaction
	return &out, s.c.do(ctx, "POST", "/rewards/redeem", req, &out)
}

// ListTransactions returns the reward transaction history for a customer.
func (s *RewardsService) ListTransactions(ctx context.Context, customerID string, params ListParams) (*ListResponse[RewardTransaction], error) {
	var out ListResponse[RewardTransaction]
	return &out, s.c.do(ctx, "GET", fmt.Sprintf("/rewards/accounts/%s/transactions", customerID), params, &out)
}
