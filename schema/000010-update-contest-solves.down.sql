-- +migrate Down

DROP INDEX IF EXISTS idx_contest_solves_first_blood;

ALTER TABLE contest_solves
DROP COLUMN IF EXISTS first_blood,
DROP COLUMN IF EXISTS attempt_count;