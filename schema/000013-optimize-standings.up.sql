-- +migrate Up

-- 1. Enhance contest_standings with additional tracking columns
ALTER TABLE contest_standings
ADD COLUMN IF NOT EXISTS wrong_attempts INT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS last_solved_at TIMESTAMPTZ;

-- 2. Create contest_user_problems table for denormalized user-problem interactions
CREATE TABLE IF NOT EXISTS contest_user_problems (
    contest_id BIGINT NOT NULL REFERENCES contests (id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL REFERENCES problems (id) ON DELETE CASCADE,
    problem_index INT NOT NULL,
    is_solved BOOLEAN NOT NULL DEFAULT FALSE,
    solved_at TIMESTAMPTZ,
    penalty INT NOT NULL DEFAULT 0,
    attempt_count INT NOT NULL DEFAULT 0,
    first_blood BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (
        contest_id,
        user_id,
        problem_id
    )
);

-- Index for efficient user-contest lookups
CREATE INDEX idx_cup_user_contest ON contest_user_problems (
    contest_id,
    user_id,
    problem_index ASC
);

-- Index for problem-wide queries (e.g., finding first blood)
CREATE INDEX idx_cup_problem ON contest_user_problems (
    contest_id,
    problem_id,
    solved_at ASC
)
WHERE
    is_solved = TRUE;

-- 3. Create contest_problem_stats table for materialized problem statistics
CREATE TABLE IF NOT EXISTS contest_problem_stats (
    contest_id BIGINT NOT NULL REFERENCES contests (id) ON DELETE CASCADE,
    problem_id BIGINT NOT NULL REFERENCES problems (id) ON DELETE CASCADE,
    problem_index INT NOT NULL,
    solved_count INT NOT NULL DEFAULT 0,
    attempted_users INT NOT NULL DEFAULT 0,
    PRIMARY KEY (contest_id, problem_id)
);

-- Index for contest-wide problem stats
CREATE INDEX idx_cps_contest ON contest_problem_stats (contest_id, problem_index ASC);

-- 4. Migrate existing data from contest_solves to contest_user_problems
INSERT INTO
    contest_user_problems (
        contest_id,
        user_id,
        problem_id,
        problem_index,
        is_solved,
        solved_at,
        penalty,
        attempt_count,
        first_blood
    )
SELECT cs.contest_id, cs.user_id, cs.problem_id, cp.index, TRUE, cs.solved_at, cs.penalty, cs.attempt_count, cs.first_blood
FROM
    contest_solves cs
    JOIN contest_problems cp ON cp.contest_id = cs.contest_id
    AND cp.problem_id = cs.problem_id ON CONFLICT (
        contest_id,
        user_id,
        problem_id
    ) DO NOTHING;

-- 5. Populate contest_user_problems with unsolved attempts
INSERT INTO
    contest_user_problems (
        contest_id,
        user_id,
        problem_id,
        problem_index,
        is_solved,
        attempt_count
    )
SELECT s.contest_id, s.user_id, s.problem_id, cp.index, FALSE, COUNT(*)
FROM
    submissions s
    JOIN contest_problems cp ON cp.contest_id = s.contest_id
    AND cp.problem_id = s.problem_id
WHERE
    s.contest_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1
        FROM contest_user_problems cup
        WHERE
            cup.contest_id = s.contest_id
            AND cup.user_id = s.user_id
            AND cup.problem_id = s.problem_id
    )
GROUP BY
    s.contest_id,
    s.user_id,
    s.problem_id,
    cp.index ON CONFLICT (
        contest_id,
        user_id,
        problem_id
    ) DO NOTHING;

-- 6. Populate contest_problem_stats
INSERT INTO
    contest_problem_stats (
        contest_id,
        problem_id,
        problem_index,
        solved_count,
        attempted_users
    )
SELECT cp.contest_id, cp.problem_id, cp.index, COUNT(
        DISTINCT CASE
            WHEN cup.is_solved THEN cup.user_id
        END
    ), COUNT(DISTINCT cup.user_id)
FROM
    contest_problems cp
    LEFT JOIN contest_user_problems cup ON cup.contest_id = cp.contest_id
    AND cup.problem_id = cp.problem_id
GROUP BY
    cp.contest_id,
    cp.problem_id,
    cp.index ON CONFLICT (contest_id, problem_id) DO
UPDATE
SET
    solved_count = EXCLUDED.solved_count,
    attempted_users = EXCLUDED.attempted_users;

-- 7. Update contest_standings with wrong_attempts and last_solved_at
UPDATE contest_standings cs
SET
    wrong_attempts = (
        SELECT COUNT(*)
        FROM submissions s
        WHERE
            s.contest_id = cs.contest_id
            AND s.user_id = cs.user_id
            AND LOWER(s.verdict) IN ('wa', 'tle', 're', 'mle')
    ),
    last_solved_at = (
        SELECT MAX(solved_at)
        FROM contest_user_problems cup
        WHERE
            cup.contest_id = cs.contest_id
            AND cup.user_id = cs.user_id
            AND cup.is_solved = TRUE
    );