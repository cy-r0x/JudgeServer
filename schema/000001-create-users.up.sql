-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    fullname VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (
        role IN ('user', 'setter', 'admin')
    ),
    allowed_contest INT,
    room_no VARCHAR(50),
    pc_no INT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_username_per_contest UNIQUE (username, allowed_contest),
    CONSTRAINT unique_email_per_contest UNIQUE (email, allowed_contest)
);

-- Indexes
CREATE INDEX idx_users_role ON users (role);