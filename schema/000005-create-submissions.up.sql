-- +migrate Up

CREATE TABLE IF NOT EXISTS submissions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    problem_id BIGINT REFERENCES problems (id) ON DELETE CASCADE,
    contest_id BIGINT REFERENCES contests (id) ON DELETE SET NULL,
    language VARCHAR(30) NOT NULL,
    source_code TEXT NOT NULL,
    verdict VARCHAR(30) DEFAULT 'Pending',
    execution_time INT,
    memory_used INT,
    submitted_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_submissions_user ON submissions (user_id);

CREATE INDEX idx_submissions_problem ON submissions (problem_id);

CREATE INDEX idx_submissions_contest ON submissions (contest_id);

CREATE INDEX idx_submissions_verdict ON submissions (verdict);

CREATE INDEX idx_submissions_submitted_at ON submissions (submitted_at DESC);