-- +migrate Up

CREATE TABLE IF NOT EXISTS contests (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    duration_seconds BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'upcoming',
            'ongoing',
            'ended'
        )
    ),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_contests_status ON contests (status);

CREATE INDEX idx_contests_start_time ON contests (start_time);