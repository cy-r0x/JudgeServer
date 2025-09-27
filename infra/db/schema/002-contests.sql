CREATE TABLE contests (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    duration_seconds BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'upcoming',
            'ongoing',
            'ended'
        )
    ),
    created_by BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Trigger function to compute end_time
CREATE FUNCTION set_contest_end_time()
RETURNS trigger AS $$
BEGIN
    NEW.end_time := NEW.start_time + (NEW.duration_seconds * interval '1 second');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for insert and update
CREATE TRIGGER trg_set_contest_end_time
BEFORE INSERT OR UPDATE ON contests
FOR EACH ROW
EXECUTE FUNCTION set_contest_end_time();