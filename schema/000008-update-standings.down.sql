-- +migrate Down

ALTER TABLE contest_standings DROP COLUMN IF EXISTS solved_count;

ALTER TABLE contest_standings ADD COLUMN submitted TEXT;

DROP TABLE IF EXISTS contest_solves;