ALTER TABLE subscriptions
    ADD COLUMN provider_subscription_id TEXT;

WITH ranked_active_subscriptions AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id
            ORDER BY created_at DESC, id DESC
        ) AS row_number
    FROM subscriptions
    WHERE status = 'active'
)
UPDATE subscriptions AS subscription
SET
    status = 'cancelled',
    updated_at = now()
FROM ranked_active_subscriptions AS ranked
WHERE subscription.id = ranked.id
  AND ranked.row_number > 1;

WITH ranked_pending_subscriptions AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY user_id
            ORDER BY created_at DESC, id DESC
        ) AS row_number
    FROM subscriptions
    WHERE status = 'pending'
)
UPDATE subscriptions AS subscription
SET
    status = 'cancelled',
    updated_at = now()
FROM ranked_pending_subscriptions AS ranked
WHERE subscription.id = ranked.id
  AND ranked.row_number > 1;

CREATE UNIQUE INDEX subscriptions_one_active_per_user_idx
    ON subscriptions (user_id)
    WHERE status = 'active';

CREATE UNIQUE INDEX subscriptions_one_pending_per_user_idx
    ON subscriptions (user_id)
    WHERE status = 'pending';

CREATE UNIQUE INDEX subscriptions_provider_subscription_id_idx
    ON subscriptions (provider_subscription_id)
    WHERE provider_subscription_id IS NOT NULL;