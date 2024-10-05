CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    request_count INT NOT NULL DEFAULT 0,
    rate_limit INT NOT NULL DEFAULT 100,  -- Default rate limit for users
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);