ALTER TABLE subscriptions ADD COLUMN plan TEXT;

UPDATE subscriptions s
SET plan = p.name
FROM plans p
WHERE s.plan_id = p.id;

ALTER TABLE subscriptions
    ALTER COLUMN plan SET NOT NULL,
    ADD CONSTRAINT subscriptions_plan_check CHECK (plan IN ('free', 'plus', 'pro')),
    DROP COLUMN plan_id;

DROP TABLE IF EXISTS plans;