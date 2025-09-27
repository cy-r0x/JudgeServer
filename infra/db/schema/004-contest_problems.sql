CREATE TABLE contest_problems (
    contest_id BIGINT NOT NULL,
    problem_id BIGINT NOT NULL,
    display_order INTEGER NOT NULL,
    PRIMARY KEY (contest_id, problem_id),
    FOREIGN KEY (contest_id) REFERENCES contests (id) ON DELETE CASCADE,
    FOREIGN KEY (problem_id) REFERENCES problems (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add index for ordering problems within a contest
CREATE INDEX idx_contest_problems_order ON contest_problems (contest_id, display_order);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_contest_problems_modtime
    BEFORE UPDATE ON contest_problems
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();