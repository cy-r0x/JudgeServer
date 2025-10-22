-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (
        role IN ('user', 'setter', 'admin')
    ),
    allowed_contest BIGINT,
    clan VARCHAR(255),
    room_no VARCHAR(50),
    pc_no VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_username_per_contest UNIQUE (username, allowed_contest)
);

-- Indexes
CREATE INDEX idx_users_role ON users (role);

INSERT INTO
    users (
        full_name,
        username,
        password,
        role
    )
VALUES (
        'admin',
        'admin',
        '$2a$12$Ncde3vjx7AbBXwyDlzgN5ue8PKgD1XexbvWdityKLbQHsHJAi1jKG',
        'admin'
    ) ON CONFLICT (username) DO NOTHING;