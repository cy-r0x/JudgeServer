-- +migrate Up
CREATE TABLE IF NOT EXISTS filepath (
    id BIGSERIAL PRIMARY KEY,
    contest_id BIGINT NOT NULL REFERENCES contests (id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL
);