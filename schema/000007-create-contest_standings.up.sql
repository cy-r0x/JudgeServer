-- +migrate Up

CREATE TABLE IF NOT EXISTS contest_standings (
    contest_id BIGINT REFERENCES contests (id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    score INT NOT NULL DEFAULT 0,
    penalty INT NOT NULL DEFAULT 0,
    PRIMARY KEY (contest_id, user_id)
);

-- Index for leaderboard lookups
CREATE INDEX idx_standings_contest ON contest_standings (
    contest_id,
    score DESC,
    penalty ASC
);