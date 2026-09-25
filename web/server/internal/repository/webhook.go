package repository

import (
	"context"

	"gitlab.com/britinogn/vidfixa/internal/db"
)

/*
MarkWebhookProcessed is the idempotency guard from workme: attempts
to insert the event_id, and the UNIQUE constraint on that column does
the actual work. A conflict means this event was already handled —
the caller should treat that as "skip, don't reprocess" rather than
an error.
*/
func MarkWebhookProcessed(ctx context.Context, eventID string) (isNew bool, err error) {
	query := `
		INSERT INTO processed_webhooks (event_id)
		VALUES ($1)
		ON CONFLICT (event_id) DO NOTHING
	`
	tag, err := db.Pool.Exec(ctx, query, eventID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

/*
UnmarkWebhookProcessed removes the idempotency marker only after a
webhook processing failure.

HandleWebhook first claims an event with MarkWebhookProcessed so two
simultaneous deliveries cannot both process it. If a later database
or Bachs operation fails, removing the marker allows Bachs's retry to
attempt the event again.
*/
func UnmarkWebhookProcessed(ctx context.Context, eventID string) error {
	query := `
		DELETE FROM processed_webhooks
		WHERE event_id = $1
	`

	_, err := db.Pool.Exec(ctx, query, eventID)
	return err
}
