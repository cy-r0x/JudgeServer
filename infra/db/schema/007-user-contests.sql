-- Create a table for user contest registrations
CREATE TABLE user_contests (
    user_id BIGINT NOT NULL,
    contest_id BIGINT NOT NULL,
    registered_by BIGINT NOT NULL, -- admin who registered the user
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (
        status IN (
            'active',
            'withdrawn',
            'disqualified'
        )
    ),
    PRIMARY KEY (user_id, contest_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contest (id) ON DELETE CASCADE,
    FOREIGN KEY (registered_by) REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add index for common lookups
CREATE INDEX idx_user_contests_user ON user_contests (user_id);

CREATE INDEX idx_user_contests_contest ON user_contests (contest_id);

CREATE INDEX idx_user_contests_registered_by ON user_contests (registered_by);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_user_contests_modtime
    BEFORE UPDATE ON user_contests
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();