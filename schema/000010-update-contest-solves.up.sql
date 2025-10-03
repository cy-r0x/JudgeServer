-- +migrate Up

ALTER TABLE contest_solves
ADD COLUMN IF NOT EXISTS attempt_count INT NOT NULL DEFAULT 1,
ADD COLUMN IF NOT EXISTS first_blood BOOLEAN NOT NULL DEFAULT false;

-- Index for quick first blood lookups
CREATE INDEX IF NOT EXISTS idx_contest_solves_first_blood ON contest_solves (
    contest_id,
    problem_id,
    first_blood
)
WHERE
    first_blood = true;