ALTER TABLE usage_counters
    ADD COLUMN user_id UUID;

UPDATE usage_counters
    SET user_id = split_part(identity_key, ':', 2)::uuid
    WHERE identity_key LIKE 'user:%';

DELETE FROM usage_counters WHERE user_id IS NULL;

ALTER TABLE usage_counters
    DROP CONSTRAINT usage_counters_pkey,
    ALTER COLUMN user_id SET NOT NULL,
    DROP COLUMN identity_key,
    ADD PRIMARY KEY (user_id, period),
    ADD CONSTRAINT usage_counters_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE downloads
    DROP CONSTRAINT downloads_identity_check,
    DROP COLUMN ip_address,
    DROP COLUMN anon_id,
    ALTER COLUMN user_id SET NOT NULL;