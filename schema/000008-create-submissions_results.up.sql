-- +migrate Up

CREATE TABLE IF NOT EXISTS submission_results (
    id BIGSERIAL PRIMARY KEY,
    submission_id BIGINT REFERENCES submissions (id) ON DELETE CASCADE,
    test_case_id BIGINT REFERENCES testcases (id) ON DELETE CASCADE,
    verdict VARCHAR(30) NOT NULL,
    execution_time_ms INT,
    memory_used_kb INT
);

-- Indexes
CREATE INDEX idx_results_submission ON submission_results (submission_id);

CREATE INDEX idx_results_test_case ON submission_results (test_case_id);