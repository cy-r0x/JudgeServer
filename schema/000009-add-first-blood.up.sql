-- +migrate Up

ALTER TABLE submissions
ADD COLUMN IF NOT EXISTS first_blood BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_submissions_first_blood ON submissions (
    problem_id,
    contest_id,
    first_blood
)
WHERE
    first_blood = true;