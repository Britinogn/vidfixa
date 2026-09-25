DROP INDEX IF EXISTS subscriptions_provider_subscription_id_idx;
DROP INDEX IF EXISTS subscriptions_one_pending_per_user_idx;
DROP INDEX IF EXISTS subscriptions_one_active_per_user_idx;

ALTER TABLE subscriptions
    DROP COLUMN provider_subscription_id;