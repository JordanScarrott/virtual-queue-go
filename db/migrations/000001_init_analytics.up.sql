CREATE TABLE IF NOT EXISTS analytics_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    business_id VARCHAR(255),
    user_id VARCHAR(255),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    properties JSONB
);
CREATE INDEX IF NOT EXISTS idx_analytics_type_time ON analytics_events (event_type, timestamp);
