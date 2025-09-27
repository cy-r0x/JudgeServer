CREATE TABLE submissions_data (
    id BIGSERIAL PRIMARY KEY,
    submission_id BIGINT NOT NULL,
    submitted_code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    FOREIGN KEY (submission_id) REFERENCES submissions (id) ON DELETE CASCADE
);

-- Add index for faster lookups by submission_id
CREATE INDEX idx_submissions_data_submission_id ON submissions_data (submission_id);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_submissions_data_modtime
    BEFORE UPDATE ON submissions_data
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();