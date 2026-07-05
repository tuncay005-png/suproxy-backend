-- Drop triggers
DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_plans_updated_at ON plans;

-- Drop indexes
DROP INDEX IF EXISTS idx_subscriptions_user_active;
DROP INDEX IF EXISTS idx_subscriptions_user_status;
DROP INDEX IF EXISTS idx_subscriptions_expires_at;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_user_id;

DROP INDEX IF EXISTS idx_plans_is_active;
DROP INDEX IF EXISTS idx_plans_name;

-- Drop tables (order matters due to foreign keys)
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS plans;
