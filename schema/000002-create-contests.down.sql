-- +migrate Down
DROP TRIGGER IF EXISTS trg_set_contest_end_time ON contests;

DROP FUNCTION IF EXISTS set_contest_end_time ();

DROP TABLE contests;