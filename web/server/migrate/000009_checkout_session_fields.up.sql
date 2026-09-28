ALTER TABLE subscriptions
    ADD COLUMN provider_checkout_id TEXT,
    ADD COLUMN checkout_url TEXT,
    ADD COLUMN checkout_expires_at TIMESTAMPTZ;

-- Partial index: many rows will have NULL here (active rows, legacy
-- pending), and those must be allowed to coexist. Only real provider
-- checkout IDs have to be unique, so a stuck or retried checkout can
-- always be traced back to its exact provider session.
CREATE UNIQUE INDEX subscriptions_provider_checkout_id_idx
    ON subscriptions (provider_checkout_id)
    WHERE provider_checkout_id IS NOT NULL;
