-- +migrate Up
CREATE TABLE IF NOT EXISTS problems (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    statement TEXT NOT NULL,
    input_statement TEXT NOT NULL,
    output_statement TEXT NOT NULL,
    time_limit FLOAT NOT NULL,
    memory_limit FLOAT NOT NULL,
    checker_type VARCHAR(10) NOT NULL,
    checker_strict_space BOOLEAN NOT NULL,
    checker_precision VARCHAR(10),
    created_by BIGINT REFERENCES users (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_problems_created_by ON problems (created_by);