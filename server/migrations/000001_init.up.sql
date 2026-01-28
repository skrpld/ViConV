CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET timezone = 'UTC';

CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    refresh_token TEXT,
    refresh_token_expiry_time TIMESTAMP WITH TIME ZONE
);