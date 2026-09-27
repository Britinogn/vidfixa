-- Safe to roll back: usage_refunded_at is only an idempotency marker.
-- Rolling back restores the old double-refund risk, but it never deletes
-- usage_counters data, so no credits are lost either way.
ALTER TABLE downloads
    DROP COLUMN usage_refunded_at;
