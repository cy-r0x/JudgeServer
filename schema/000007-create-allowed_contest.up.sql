-- +migrate Up
CREATE TABLE IF NOT EXISTS allowed_contest (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    contest_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contests (id) ON DELETE CASCADE,
    UNIQUE (user_id, contest_id)
);

-- Add indexes for common query patterns
CREATE INDEX idx_allowed_contest_user_id ON allowed_contest (user_id);

CREATE INDEX idx_allowed_contest_contest_id ON allowed_contest (contest_id);