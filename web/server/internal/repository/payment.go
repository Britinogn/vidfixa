package repository

import (
	"context"

	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/internal/model"
)

func CreatePayment(ctx context.Context, userID string, subscriptionID *string, provider, providerRef, amount, currency, status string) (*model.Payment, error) {
	query := `
		INSERT INTO payments (user_id, subscription_id, provider, provider_reference, amount, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, user_id, subscription_id, provider, provider_reference, amount, currency, status, created_at
	`

	var p model.Payment
	err := db.Pool.QueryRow(ctx, query, userID, subscriptionID, provider, providerRef, amount, currency, status).Scan(
		&p.ID, &p.UserID, &p.SubscriptionID, &p.Provider, &p.ProviderReference, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
