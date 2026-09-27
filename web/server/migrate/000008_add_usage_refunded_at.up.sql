-- Marks that a download's reserved usage slot was already refunded.
--
-- downloads.usage_period is written when a download is created, because
-- usage is reserved up front (service/download.go) and released again if
-- the job never completes. This column is the idempotency guard for that
-- release: a row is only ever refunded once, no matter how many times a
-- terminal path fires (worker failure, queue-full, cancel, restart recovery).
--
-- NULL means "still holding its reservation". A timestamp means the slot
-- was already returned to usage_counters and must not be returned again.
ALTER TABLE downloads
    ADD COLUMN usage_refunded_at TIMESTAMPTZ;
