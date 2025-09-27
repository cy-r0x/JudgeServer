CREATE TABLE problems (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    setter_id BIGINT NOT NULL,
    description TEXT NOT NULL,
    input_description TEXT NOT NULL,
    output_description TEXT NOT NULL,
    sample_testcase TEXT NOT NULL,
    regular_testcase TEXT NOT NULL,
    time_limit FLOAT NOT NULL,
    memory_limit FLOAT NOT NULL,
    FOREIGN KEY (setter_id) REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add indexes for common lookups
CREATE INDEX idx_problems_title ON problems (title);

CREATE INDEX idx_problems_setter ON problems (setter_id);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_problems_modtime
    BEFORE UPDATE ON problems
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();