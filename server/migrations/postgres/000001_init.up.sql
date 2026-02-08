CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET timezone = 'UTC';

CREATE OR REPLACE FUNCTION update_modified_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    refresh_token TEXT,
    refresh_token_expiry_time TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS posts (
    post_id UUID PRIMARY KEY  DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    idempotency_key TEXT UNIQUE NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT fk_posts_user
                                 FOREIGN KEY (user_id)
                                 REFERENCES users(user_id)
                                 ON DELETE CASCADE
);

CREATE TRIGGER update_posts_modtime
    BEFORE UPDATE ON posts
    FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

CREATE OR REPLACE FUNCTION haversine_distance(
    lat1 DOUBLE PRECISION,
    lon1 DOUBLE PRECISION,
    lat2 DOUBLE PRECISION,
    lon2 DOUBLE PRECISION
)
    RETURNS DOUBLE PRECISION AS $$
DECLARE
    r DOUBLE PRECISION := 6371;
    p DOUBLE PRECISION := PI() / 180;
    a DOUBLE PRECISION;
BEGIN
    a := 0.5 - COS((lat2 - lat1) * p) / 2
        + COS(lat1 * p) * COS(lat2 * p) *
          (1 - COS((lon2 - lon1) * p)) / 2;

    RETURN 2 * r * ASIN(SQRT(a));
END;
$$ LANGUAGE plpgsql IMMUTABLE;