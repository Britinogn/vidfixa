package model

import "time"

// Payment mirrors the payments table.
type Payment struct {
	ID                string    `db:"id" json:"id"`
	UserID            string    `db:"user_id" json:"user_id"`
	SubscriptionID    *string   `db:"subscription_id" json:"subscription_id,omitempty"`
	Provider          string    `db:"provider" json:"provider"`
	ProviderReference string    `db:"provider_reference" json:"provider_reference"`
	Amount            string    `db:"amount" json:"amount"` // decimal string, matches Bachs's own format
	Currency          string    `db:"currency" json:"currency"`
	Status            string    `db:"status" json:"status"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}
