-- +migrate Down

DROP INDEX IF EXISTS idx_submissions_first_blood;

ALTER TABLE submissions DROP COLUMN IF EXISTS first_blood;