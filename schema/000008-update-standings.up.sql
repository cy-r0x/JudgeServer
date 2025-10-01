-- +migrate Up

CREATE TABLE IF NOT EXISTS contest_solves (
    contest_id BIGINT REFERENCES contests (id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    problem_id BIGINT REFERENCES problems (id) ON DELETE CASCADE,
    solved_at TIMESTAMPTZ NOT NULL,
    penalty INT NOT NULL,
    PRIMARY KEY (
        contest_id,
        user_id,
        problem_id
    )
);

ALTER TABLE contest_standings DROP COLUMN IF EXISTS submitted;

ALTER TABLE contest_standings
ADD COLUMN IF NOT EXISTS solved_count INT NOT NULL DEFAULT 0;