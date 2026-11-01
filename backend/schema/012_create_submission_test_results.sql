-- 012_create_submission_test_results.sql
-- LearnX Submission Test Case Execution Results

CREATE TABLE IF NOT EXISTS submission_test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    test_case_id UUID NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL,
    execution_time_ms INTEGER DEFAULT 0,
    memory_used_kb INTEGER DEFAULT 0,
    actual_output TEXT DEFAULT '',
    error_message TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_submission_test_results_submission_id ON submission_test_results(submission_id);
CREATE INDEX IF NOT EXISTS idx_submission_test_results_test_case_id ON submission_test_results(test_case_id);
