-- +migrate Up

-- Composite index for submissions filtering by contest and user
CREATE INDEX IF NOT EXISTS idx_submissions_contest_user_problem ON submissions (
    contest_id,
    user_id,
    problem_id,
    submitted_at DESC
);

-- Composite index for submissions by contest, problem, and verdict (for first blood checks)
CREATE INDEX IF NOT EXISTS idx_submissions_contest_problem_verdict_id ON submissions (
    contest_id,
    problem_id,
    verdict,
    id
)
WHERE
    verdict = 'ac';

-- Composite index for penalty calculation
CREATE INDEX IF NOT EXISTS idx_submissions_penalty_lookup ON submissions (
    contest_id,
    user_id,
    problem_id,
    submitted_at,
    verdict
);

-- Index for contest_solves joins
CREATE INDEX IF NOT EXISTS idx_contest_solves_lookup ON contest_solves (
    contest_id,
    user_id,
    problem_id,
    solved_at
);

-- Index for users joined with submissions
CREATE INDEX IF NOT EXISTS idx_submissions_contest_user ON submissions (contest_id, user_id);