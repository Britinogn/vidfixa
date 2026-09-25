package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"gitlab.com/britinogn/vidfixa/internal/db"
)

// var ErrLimitReached = errors.New("monthly download limit reached")
var ErrLimitReached = errors.New("monthly download limit reached. Upgrade your plan to continue.")

func IncrementUsage(ctx context.Context, identityKey, period string, limit int) (int, error) {
	query := `
		INSERT INTO usage_counters (identity_key, period, count)
		VALUES ($1, $2, 1)
		ON CONFLICT (identity_key, period)
		DO UPDATE SET count = usage_counters.count + 1
		WHERE usage_counters.count < $3
		RETURNING count
	`

	var count int
	err := db.Pool.QueryRow(ctx, query, identityKey, period, limit).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrLimitReached
		}
		return 0, err
	}
	return count, nil
}

func GetUsage(ctx context.Context, identityKey, period string) (int, error) {
	query := `SELECT count FROM usage_counters WHERE identity_key = $1 AND period = $2`

	var count int
	err := db.Pool.QueryRow(ctx, query, identityKey, period).Scan(&count)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
