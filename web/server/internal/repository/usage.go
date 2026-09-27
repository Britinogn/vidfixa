package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"gitlab.com/britinogn/vidfixa/internal/db"
	"gitlab.com/britinogn/vidfixa/pkg/utils"
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

func DecrementUsage(ctx context.Context, identityKey, period string) error {
	query := `
		UPDATE usage_counters
		SET count = GREATEST(count - 1, 0)
		WHERE identity_key = $1 AND period = $2
	`
	_, err := db.Pool.Exec(ctx, query, identityKey, period)
	return err
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

/*
identityKeyForDownload mirrors service.Identity.key(): exactly one of
user_id/anon_id is set on a downloads row, enforced by the
downloads_identity_check constraint.
*/
func identityKeyForDownload(userID, anonID *string) string {
	if userID != nil {
		return utils.IdentityKeyForUser(*userID)
	}
	if anonID != nil {
		return utils.IdentityKeyForAnon(*anonID)
	}
	return ""
}

/*
decrementUsageTx releases one reserved slot inside an existing transaction.
Kept separate from DecrementUsage so a refund and its guard marker commit
together — a partially applied refund is worse than none, because the
marker would then block the retry that would have fixed it.
*/
func decrementUsageTx(ctx context.Context, tx pgx.Tx, identityKey, period string) error {
	_, err := tx.Exec(ctx, `
		UPDATE usage_counters
		SET count = GREATEST(count - 1, 0)
		WHERE identity_key = $1 AND period = $2
	`, identityKey, period)
	return err
}

/*
RefundFailedDownload is the single entry point for returning a reserved
usage slot when a download ends in any non-completed state (failed,
cancelled, or abandoned by a restart).

It replaces the hand-rolled refund that used to live in main.go's
processJob closure. Three things it fixes:

 1. Completeness — the old version only ran on worker failure, so
    queue-full and INSERT-failure silently ate a credit.
 2. Idempotency — the guarded UPDATE matches only rows that are still
    queued/processing and whose usage_refunded_at IS NULL, so a second
    call affects zero rows and refunds nothing. The status guard is
    what stops a stray call from demoting a download that already
    completed, which would hand back a credit for a real download.
 3. Atomicity — the marker and the counter update share one transaction,
    so a crash can't leave a download marked refunded while the slot is
    still held.

refunded reports whether this call was the one that released the slot.
A false return with a nil error means the download was already handled,
which callers should treat as success, not failure.

Rows whose usage_period is NULL predate migration 000007, so there is no
period key to credit back. They still get marked terminal, and we never
guess a period — that would corrupt someone else's counter.
*/
func RefundFailedDownload(ctx context.Context, id, status string, errMsg *string) (refunded bool, err error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userID, anonID, usagePeriod *string

	err = tx.QueryRow(ctx, `
		UPDATE downloads
		SET status = $2,
			error = $3,
			completed_at = now(),
			usage_refunded_at = now()
		WHERE id = $1
			AND usage_refunded_at IS NULL
			AND status IN ('queued', 'processing')
		RETURNING user_id, anon_id, usage_period
	`, id, status, errMsg).Scan(&userID, &anonID, &usagePeriod)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Already refunded, already terminal, or completed.
			return false, nil
		}
		return false, err
	}

	if usagePeriod != nil {
		identityKey := identityKeyForDownload(userID, anonID)
		if identityKey != "" {
			if err := decrementUsageTx(ctx, tx, identityKey, *usagePeriod); err != nil {
				return false, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

/*
RecoverStaleDownloads handles the durability gap: the jobs channel lives
only in memory, so a kill or crash leaves rows stuck at queued/processing
with their usage slots reserved forever.

Called once at startup with a cutoff older than the 30-minute job timeout
in worker/pool.go, so anything still non-terminal that long is definitionally
dead — a live job would have been cancelled by its context and written a
terminal status of its own.

Returns the number of downloads reconciled.
*/
func RecoverStaleDownloads(ctx context.Context, staleBefore time.Time) (int, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id
		FROM downloads
		WHERE status IN ('queued', 'processing')
			AND created_at < $1
	`, staleBefore)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	// Release the connection before the refund loop opens its own
	// transactions — rows is already drained, so the pool is free.
	rows.Close()

	recovered := 0
	for _, id := range ids {
		errMsg := "server restarted before this download finished"
		refunded, err := RefundFailedDownload(ctx, id, "failed", &errMsg)
		if err != nil {
			return recovered, err
		}
		if refunded {
			recovered++
		}
	}

	return recovered, nil
}
