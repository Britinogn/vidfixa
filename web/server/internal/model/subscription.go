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
	// Provider checkout session behind a pending row. Stored so a retry can
	// resume the exact session (same URL) instead of opening a duplicate,
	// and so webhooks can be matched to the exact checkout rather than to
	// "newest pending row with this plan".
	ProviderCheckoutID *string    `db:"provider_checkout_id" json:"provider_checkout_id,omitempty"`
	CheckoutURL        *string    `db:"checkout_url" json:"-"`
	CheckoutExpiresAt  *time.Time `db:"checkout_expires_at" json:"checkout_expires_at,omitempty"`
	CreatedAt          time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at" json:"updated_at"`
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
	// Resumed is true when no new provider session was created — the user
	// is being sent back to their still-open checkout.
	Resumed   bool       `json:"resumed,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// PendingCheckoutResponse describes the user's in-progress checkout, if any,
// so the subscription page can offer "resume" instead of a dead end.
type PendingCheckoutResponse struct {
	Tier        string     `json:"tier"`
	CheckoutURL string     `json:"checkout_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// SubscriptionStatusResponse is GET /api/subscription: current plan plus any
// pending checkout the user could resume.
type SubscriptionStatusResponse struct {
	Plan    string                   `json:"plan"`
	Pending *PendingCheckoutResponse `json:"pending,omitempty"`
}
