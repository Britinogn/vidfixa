package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/internal/model"
)

var ErrSubscriptionNotFound = errors.New("subscription not found")

/*
GetActiveSubscription returns the user's current subscription row.

A subscription only counts as active while its expiry time is still
in the future. An old row with status "active" but an elapsed
expires_at is treated as Free by callers, until it is replaced by
a new successful checkout.
*/
func GetActiveSubscription(ctx context.Context, userID string) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			plan,
			status,
			provider_subscription_id,
			started_at,
			expires_at,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
			AND status = 'active'
			AND expires_at > now()
		ORDER BY created_at DESC
		LIMIT 1
	`

	var s model.Subscription
	err := db.Pool.QueryRow(ctx, query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.Plan,
		&s.Status,
		&s.ProviderSubscriptionID,
		&s.StartedAt,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &s, nil
}

/*
GetPendingSubscription returns the user's one in-progress checkout.

The database migration now enforces at most one pending subscription
per user. This lets the service reject or handle a second checkout
before it can create another provider billing session.
*/
func GetPendingSubscription(ctx context.Context, userID string) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			plan,
			status,
			provider_subscription_id,
			started_at,
			expires_at,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
			AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var s model.Subscription
	err := db.Pool.QueryRow(ctx, query, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.Plan,
		&s.Status,
		&s.ProviderSubscriptionID,
		&s.StartedAt,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &s, nil
}

/*
GetSubscriptionByIDAndUser verifies that a subscription reference
belongs to the user from the Bachs payment event.

Checkout creation sends the local pending subscription ID as Bachs's
Reference. The webhook must still verify ownership before using that
reference to link a payment row.
*/
func GetSubscriptionByIDAndUser(ctx context.Context, subscriptionID, userID string) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			plan,
			status,
			provider_subscription_id,
			started_at,
			expires_at,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
			AND user_id = $2
	`

	var s model.Subscription
	err := db.Pool.QueryRow(ctx, query, subscriptionID, userID).Scan(
		&s.ID,
		&s.UserID,
		&s.Plan,
		&s.Status,
		&s.ProviderSubscriptionID,
		&s.StartedAt,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &s, nil
}

/*
CreatePendingSubscription inserts a new row in "pending" status when
checkout is initiated. Confirmed to "active" later by the webhook,
never by the checkout call itself — matches workme's "never trust
the frontend on payment state" rule.

The database allows one active subscription and one pending
replacement together. The active row stays usable until the provider
confirms payment for the pending replacement.
*/
func CreatePendingSubscription(ctx context.Context, userID, plan string) (*model.Subscription, error) {
	query := `
		INSERT INTO subscriptions (user_id, plan, status)
		VALUES ($1, $2, 'pending')
		RETURNING
			id,
			user_id,
			plan,
			status,
			provider_subscription_id,
			started_at,
			expires_at,
			created_at,
			updated_at
	`

	var s model.Subscription
	err := db.Pool.QueryRow(ctx, query, userID, plan).Scan(
		&s.ID,
		&s.UserID,
		&s.Plan,
		&s.Status,
		&s.ProviderSubscriptionID,
		&s.StartedAt,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

/*
CancelPendingSubscription clears a pending row when Bachs checkout
creation fails. Without this cleanup, the one-pending-subscription
database rule would block the user from trying checkout again.
*/
func CancelPendingSubscription(ctx context.Context, subscriptionID string) error {
	query := `
		UPDATE subscriptions
		SET
			status = 'cancelled',
			updated_at = now()
		WHERE id = $1
			AND status = 'pending'
	`

	_, err := db.Pool.Exec(ctx, query, subscriptionID)
	return err
}

/*
ActivateReplacementSubscription is called after Bachs confirms a
successful subscription payment.

The transaction first cancels any existing active local subscription,
then activates the exact pending row that the webhook matched. Doing
both together is necessary because the database permits only one
active subscription per user.
*/
func ActivateReplacementSubscription(
	ctx context.Context,
	subscriptionID,
	providerSubscriptionID string,
	expiresAt time.Time,
) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userID string
	err = tx.QueryRow(ctx, `
		SELECT user_id
		FROM subscriptions
		WHERE id = $1
			AND status = 'pending'
		FOR UPDATE
	`, subscriptionID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSubscriptionNotFound
		}
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE subscriptions
		SET
			status = 'cancelled',
			updated_at = now()
		WHERE user_id = $1
			AND status = 'active'
	`, userID)
	if err != nil {
		return err
	}

	commandTag, err := tx.Exec(ctx, `
		UPDATE subscriptions
		SET
			status = 'active',
			provider_subscription_id = $2,
			started_at = now(),
			expires_at = $3,
			updated_at = now()
		WHERE id = $1
			AND status = 'pending'
	`, subscriptionID, providerSubscriptionID, expiresAt)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("subscription %s was not activated", subscriptionID)
	}

	return tx.Commit(ctx)
}

/*
UpdateActiveSubscriptionPlan saves the result of a successful Bachs
plan update.

Bachs keeps the same provider subscription ID when a customer moves
between Plus and Pro. It returns the new billing dates, and the local
row must use those dates so plan lookup and usage reset stay correct.
*/

func UpdateActiveSubscriptionPlan(ctx context.Context, subscriptionID, plan string, startedAt, expiresAt time.Time) error {
	query := `
		UPDATE subscriptions
		SET
			plan = $2,
			started_at = $3,
			expires_at = $4,
			updated_at = now()
		WHERE id = $1
			AND status = 'active'
	`
	commandTag, err := db.Pool.Exec(
		ctx, query, subscriptionID, plan, startedAt, expiresAt,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("active subscription %s was not updated", subscriptionID)
	}

	return nil
}

/*
ActivateSubscription is the previous activation function.

It remains temporarily because the current webhook service still
calls it. In the next service/webhook step, we will replace that call
with ActivateReplacementSubscription, then this function can be
commented out rather than deleted.
*/
func ActivateSubscription(ctx context.Context, subscriptionID string, expiresAt time.Time) error {
	query := `
		UPDATE subscriptions
		SET status = 'active', started_at = now(), expires_at = $2, updated_at = now()
		WHERE id = $1
	`

	_, err := db.Pool.Exec(ctx, query, subscriptionID, expiresAt)
	return err
}

/*
UpdateSubscriptionStatus handles every other state transition:
past_due, cancelled, expired.
*/
func UpdateSubscriptionStatus(ctx context.Context, subscriptionID, status string) error {
	query := `UPDATE subscriptions SET status = $2, updated_at = now() WHERE id = $1`

	_, err := db.Pool.Exec(ctx, query, subscriptionID, status)
	return err
}

/*
FindPendingSubscription finds the most recent "pending" subscription
row for a user + plan combination.

It remains temporarily because the current webhook service still
calls it. The updated webhook will use GetPendingSubscription after
the checkout service has guaranteed one pending row per user.
*/
func FindPendingSubscription(ctx context.Context, userID, plan string) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			user_id,
			plan,
			status,
			provider_subscription_id,
			started_at,
			expires_at,
			created_at,
			updated_at
		FROM subscriptions
		WHERE user_id = $1
			AND plan = $2
			AND status = 'pending'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var s model.Subscription
	err := db.Pool.QueryRow(ctx, query, userID, plan).Scan(
		&s.ID,
		&s.UserID,
		&s.Plan,
		&s.Status,
		&s.ProviderSubscriptionID,
		&s.StartedAt,
		&s.ExpiresAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &s, nil
}
