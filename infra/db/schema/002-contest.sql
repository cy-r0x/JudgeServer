CREATE TABLE contest (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ GENERATED ALWAYS AS (
        start_time + (
            duration_seconds * interval '1 second'
        )
    ) STORED,
    duration_seconds BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'upcoming',
            'ongoing',
            'ended'
        )
    ),
    created_by BIGINT NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Add index for faster status-based lookups
CREATE INDEX idx_contest_status ON contest (status);

CREATE INDEX idx_contest_created_by ON contest (created_by);

CREATE INDEX idx_contest_time ON contest (start_time, end_time);

-- Create a trigger for updating the updated_at timestamp
CREATE TRIGGER update_contest_modtime
    BEFORE UPDATE ON contest
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_column();