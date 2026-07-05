-- Create plans table
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    duration_days INTEGER NOT NULL CHECK (duration_days > 0),
    traffic_limit_gb BIGINT NOT NULL DEFAULT 0 CHECK (traffic_limit_gb >= 0),
    device_limit INTEGER NOT NULL CHECK (device_limit > 0),
    max_sessions INTEGER NOT NULL CHECK (max_sessions > 0),
    price BIGINT NOT NULL CHECK (price >= 0),
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    is_active BOOLEAN NOT NULL DEFAULT true,
    features JSONB DEFAULT '[]'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create subscriptions table
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'suspended', 'cancelled')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    traffic_used_bytes BIGINT NOT NULL DEFAULT 0 CHECK (traffic_used_bytes >= 0),
    traffic_limit_bytes BIGINT NOT NULL DEFAULT 0 CHECK (traffic_limit_bytes >= 0),
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_plans_name ON plans(name);
CREATE INDEX idx_plans_is_active ON plans(is_active);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE INDEX idx_subscriptions_expires_at ON subscriptions(expires_at);
CREATE INDEX idx_subscriptions_user_status ON subscriptions(user_id, status);

-- Prevent users from having multiple active subscriptions
CREATE UNIQUE INDEX idx_subscriptions_user_active ON subscriptions(user_id) 
WHERE status = 'active';

-- Add trigger for updated_at on plans table
CREATE TRIGGER update_plans_updated_at
    BEFORE UPDATE ON plans
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add trigger for updated_at on subscriptions table
CREATE TRIGGER update_subscriptions_updated_at
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE plans IS 'Subscription plans available for users';
COMMENT ON TABLE subscriptions IS 'User subscriptions to plans';

COMMENT ON COLUMN plans.traffic_limit_gb IS '0 means unlimited traffic';
COMMENT ON COLUMN plans.price IS 'Price in smallest currency unit (e.g., cents for USD)';
COMMENT ON COLUMN subscriptions.traffic_limit_bytes IS '0 means unlimited traffic';
COMMENT ON COLUMN subscriptions.auto_renew IS 'Whether subscription should auto-renew on expiration';
