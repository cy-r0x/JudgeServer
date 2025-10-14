-- +migrate Down

DROP INDEX IF EXISTS idx_submissions_contest_user;

DROP INDEX IF EXISTS idx_contest_solves_lookup;

DROP INDEX IF EXISTS idx_submissions_penalty_lookup;

DROP INDEX IF EXISTS idx_submissions_contest_problem_verdict_id;

DROP INDEX IF EXISTS idx_submissions_contest_user_problem;