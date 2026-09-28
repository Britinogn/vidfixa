DROP INDEX IF EXISTS subscriptions_provider_checkout_id_idx;
ALTER TABLE subscriptions
    DROP COLUMN provider_checkout_id,
    DROP COLUMN checkout_url,
    DROP COLUMN checkout_expires_at;
