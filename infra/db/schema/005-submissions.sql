CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    problem_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    contest_id BIGINT,
    language VARCHAR(20) NOT NULL,
    verdict VARCHAR(20) NOT NULL CHECK (
        verdict IN (
            'accepted',
            'wrong_answer',
            'time_limit_exceeded',
            'memory_limit_exceeded',
            'runtime_error',
            'compile_error',
            'pending'
        )
    ),
    time_required FLOAT NOT NULL,
    memory_required FLOAT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (problem_id) REFERENCES problems (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (contest_id) REFERENCES contest (id) ON DELETE SET NULL
);

-- Add indexes for common query patterns
CREATE INDEX idx_submissions_problem_id ON submissions (problem_id);

CREATE INDEX idx_submissions_user_id ON submissions (user_id);

CREATE INDEX idx_submissions_contest_id ON submissions (contest_id);

CREATE INDEX idx_submissions_verdict ON submissions (verdict);

CREATE INDEX idx_submissions_created_at ON submissions (created_at);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_submissions_modtime
    BEFORE UPDATE ON submissions
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();