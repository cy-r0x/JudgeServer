-- +migrate Down

-- Drop new tables
DROP TABLE IF EXISTS contest_problem_stats;

DROP TABLE IF EXISTS contest_user_problems;

-- Remove added columns from contest_standings
ALTER TABLE contest_standings
DROP COLUMN IF EXISTS wrong_attempts,
DROP COLUMN IF EXISTS last_solved_at;