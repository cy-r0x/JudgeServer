-- +migrate Up
CREATE TABLE IF NOT EXISTS problems (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    statement TEXT NOT NULL, -- can later switch to JSONB if needed
    time_limit_ms INT NOT NULL,
    memory_limit_mb INT NOT NULL,
    created_by BIGINT REFERENCES users (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_problems_created_by ON problems (created_by);