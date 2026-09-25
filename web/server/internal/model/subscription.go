package model

import "time"

// Subscription mirrors the subscriptions table.
type Subscription struct {
	ID                     string     `db:"id" json:"id"`
	UserID                 string     `db:"user_id" json:"user_id"`
	Plan                   string     `db:"plan" json:"plan"` // "free" | "plus" | "pro"
	Status                 string     `db:"status" json:"status"`
	ProviderSubscriptionID *string    `db:"provider_subscription_id" json:"provider_subscription_id,omitempty"`
	StartedAt              *time.Time `db:"started_at" json:"started_at,omitempty"`
	ExpiresAt              *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt              time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at" json:"updated_at"`
}

// Subscription status constants, matching the CHECK constraint.
const (
	SubStatusNone      = "none"
	SubStatusPending   = "pending"
	SubStatusActive    = "active"
	SubStatusPastDue   = "past_due"
	SubStatusCancelled = "cancelled"
	SubStatusExpired   = "expired"
)

// CheckoutRequest is what POST /api/subscription/checkout decodes.
type CheckoutRequest struct {
	Tier string `json:"tier" validate:"required,oneof=plus pro"`
}

// CheckoutResponse is returned to the frontend to redirect the user.
type CheckoutResponse struct {
	CheckoutURL string `json:"checkout_url"`
}
