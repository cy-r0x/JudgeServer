-- +migrate Up

CREATE TABLE IF NOT EXISTS testcases (
    id BIGSERIAL PRIMARY KEY,
    problem_id BIGINT REFERENCES problems (id) ON DELETE CASCADE,
    input TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_sample BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Index
CREATE INDEX idx_testcases_problem ON testcases (problem_id);