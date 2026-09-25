-- 000006: reconcile legacy active that blocks webhook.go:306
-- Active rows before 000005 have provider_subscription_id IS NULL
-- webhook handleSubscriptionCreated fails with "has no provider subscription ID"
-- and pending upgrade/downgrade stays pending forever (plus→pro AND pro→plus)
-- This cancels ONLY legacy actives where a pending replacement already exists
-- (so we don't touch standalone actives without a pending)
-- IMPORTANT: cancel the matching Bachs recurring sub in Dashboard FIRST, else double-bill risk

WITH pending_users AS (
  SELECT DISTINCT user_id FROM subscriptions WHERE status = 'pending'
)
UPDATE subscriptions AS s
SET status = 'cancelled',
    updated_at = now()
FROM pending_users p
WHERE s.user_id = p.user_id
  AND s.status = 'active'
  AND s.provider_subscription_id IS NULL;