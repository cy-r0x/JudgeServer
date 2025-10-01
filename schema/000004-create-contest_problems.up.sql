-- +migrate Up
CREATE TABLE IF NOT EXISTS contest_problems (
    contest_id BIGINT REFERENCES contests (id) ON DELETE CASCADE,
    problem_id BIGINT REFERENCES problems (id) ON DELETE CASCADE,
    index INT NOT NULL,
    PRIMARY KEY (contest_id, problem_id)
);

-- Index for quick contest lookups
CREATE INDEX idx_contest_problems_contest ON contest_problems (contest_id, index);