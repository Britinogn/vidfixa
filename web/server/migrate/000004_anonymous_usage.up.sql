-- downloads: user_id becomes optional, anon_id added for anonymous jobs
ALTER TABLE downloads
    ALTER COLUMN user_id DROP NOT NULL,
    ADD COLUMN anon_id TEXT,
    ADD COLUMN ip_address TEXT,
    ADD CONSTRAINT downloads_identity_check
        CHECK (user_id IS NOT NULL OR anon_id IS NOT NULL);

-- usage_counters: single identity_key replaces user_id
ALTER TABLE usage_counters
    ADD COLUMN identity_key TEXT;

UPDATE usage_counters SET identity_key = 'user:' || user_id::text;

ALTER TABLE usage_counters
    DROP CONSTRAINT usage_counters_pkey,
    ALTER COLUMN identity_key SET NOT NULL,
    DROP COLUMN user_id,
    ADD PRIMARY KEY (identity_key, period);