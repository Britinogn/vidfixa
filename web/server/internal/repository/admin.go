package repository

import (
	"context"
	"time"

	"gitlab.com/britinogn/vidfixa/internal/db"
)

type AdminOverview struct {
	TotalUsers          int `json:"total_users"`
	ActiveSubscriptions int `json:"active_subscriptions"`
	TotalDownloads      int `json:"total_downloads"`
	FailedDownloads     int `json:"failed_downloads"`
	SuccessfulPayments  int `json:"successful_payments"`
}

type AdminUser struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminSubscription struct {
	ID                     string     `json:"id"`
	UserID                 string     `json:"user_id"`
	UserEmail              string     `json:"user_email"`
	Plan                   string     `json:"plan"`
	Status                 string     `json:"status"`
	ProviderSubscriptionID *string    `json:"provider_subscription_id,omitempty"`
	StartedAt              *time.Time `json:"started_at,omitempty"`
	ExpiresAt              *time.Time `json:"expires_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

type AdminPayment struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	UserEmail         string    `json:"user_email"`
	SubscriptionID    *string   `json:"subscription_id,omitempty"`
	Provider          string    `json:"provider"`
	ProviderReference string    `json:"provider_reference"`
	Amount            string    `json:"amount"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type AdminDownload struct {
	ID          string     `json:"id"`
	UserID      *string    `json:"user_id,omitempty"`
	UserEmail   *string    `json:"user_email,omitempty"`
	AnonID      *string    `json:"anon_id,omitempty"`
	URL         string     `json:"url"`
	Platform    string     `json:"platform"`
	Status      string     `json:"status"`
	Error       *string    `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

/*
GetAdminOverview returns small aggregate counters for the first admin
dashboard screen. It intentionally exposes counts only, not private
user information or payment values.
*/
func GetAdminOverview(ctx context.Context) (*AdminOverview, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM users),
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'active' AND expires_at > now()),
			(SELECT COUNT(*) FROM downloads),
			(SELECT COUNT(*) FROM downloads WHERE status = 'failed'),
			(SELECT COUNT(*) FROM payments WHERE status = 'SUCCEEDED')
	`

	var overview AdminOverview
	err := db.Pool.QueryRow(ctx, query).Scan(
		&overview.TotalUsers,
		&overview.ActiveSubscriptions,
		&overview.TotalDownloads,
		&overview.FailedDownloads,
		&overview.SuccessfulPayments,
	)
	if err != nil {
		return nil, err
	}

	return &overview, nil
}

/*
ListAdminUsers returns account information needed for administration.
Password hashes and other authentication secrets are never selected.
*/
func ListAdminUsers(ctx context.Context) ([]AdminUser, error) {
	query := `
		SELECT id, full_name, email, role, created_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]AdminUser, 0)
	for rows.Next() {
		var user AdminUser
		if err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

/*
ListAdminSubscriptions returns subscription lifecycle data together
with the owning user's email for support and reconciliation work.
*/
func ListAdminSubscriptions(ctx context.Context) ([]AdminSubscription, error) {
	query := `
		SELECT
			s.id,
			s.user_id,
			u.email,
			s.plan,
			s.status,
			s.provider_subscription_id,
			s.started_at,
			s.expires_at,
			s.created_at
		FROM subscriptions s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]AdminSubscription, 0)
	for rows.Next() {
		var subscription AdminSubscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.UserEmail,
			&subscription.Plan,
			&subscription.Status,
			&subscription.ProviderSubscriptionID,
			&subscription.StartedAt,
			&subscription.ExpiresAt,
			&subscription.CreatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

/*
ListAdminPayments returns payment bookkeeping data. It includes the
provider reference for support work but does not expose payment method
details because the database does not store them.
*/
func ListAdminPayments(ctx context.Context) ([]AdminPayment, error) {
	query := `
		SELECT
			p.id,
			p.user_id,
			u.email,
			p.subscription_id,
			p.provider,
			p.provider_reference,
			p.amount,
			p.currency,
			p.status,
			p.created_at
		FROM payments p
		JOIN users u ON u.id = p.user_id
		ORDER BY p.created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	payments := make([]AdminPayment, 0)
	for rows.Next() {
		var payment AdminPayment
		if err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.UserEmail,
			&payment.SubscriptionID,
			&payment.Provider,
			&payment.ProviderReference,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.CreatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

/*
ListAdminDownloads returns recent download activity for operational
review. Anonymous records have no user email and retain their anon ID.
*/
func ListAdminDownloads(ctx context.Context) ([]AdminDownload, error) {
	query := `
		SELECT
			d.id,
			d.user_id,
			u.email,
			d.anon_id,
			d.url,
			d.platform,
			d.status,
			d.error,
			d.created_at,
			d.completed_at
		FROM downloads d
		LEFT JOIN users u ON u.id = d.user_id
		ORDER BY d.created_at DESC
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	downloads := make([]AdminDownload, 0)
	for rows.Next() {
		var download AdminDownload
		if err := rows.Scan(
			&download.ID,
			&download.UserID,
			&download.UserEmail,
			&download.AnonID,
			&download.URL,
			&download.Platform,
			&download.Status,
			&download.Error,
			&download.CreatedAt,
			&download.CompletedAt,
		); err != nil {
			return nil, err
		}
		downloads = append(downloads, download)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return downloads, nil
}
