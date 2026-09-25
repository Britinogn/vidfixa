CREATE TABLE plans (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                   TEXT NOT NULL UNIQUE,
    monthly_download_limit INTEGER NOT NULL,
    price                  NUMERIC(10, 2) NOT NULL DEFAULT 0,
    currency               TEXT NOT NULL DEFAULT 'NGN',
    created_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO plans (name, monthly_download_limit, price, currency) VALUES
    ('free', 10, 0, 'NGN'),
    ('plus', 50, 0, 'NGN'),
    ('pro', 200, 0, 'NGN');

ALTER TABLE subscriptions ADD COLUMN plan_id UUID REFERENCES plans(id);

UPDATE subscriptions s
SET plan_id = p.id
FROM plans p
WHERE s.plan = p.name;

ALTER TABLE subscriptions
    ALTER COLUMN plan_id SET NOT NULL,
    DROP CONSTRAINT subscriptions_plan_check,
    DROP COLUMN plan;