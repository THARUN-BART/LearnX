-- 010_create_submissions.sql
-- LearnX Problem Submission System

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language VARCHAR(32) NOT NULL,
    source_code TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'Queued'
        CHECK (status IN ('Queued', 'Processing', 'Accepted', 'Wrong Answer', 'Compilation Error', 'Runtime Error', 'Time Limit Exceeded', 'Memory Limit Exceeded', 'Internal Error')),
    execution_time_ms INTEGER DEFAULT 0,
    memory_used_kb INTEGER DEFAULT 0,
    total_test_cases INTEGER DEFAULT 0,
    passed_test_cases INTEGER DEFAULT 0,
    failed_test_cases INTEGER DEFAULT 0,
    error_message TEXT DEFAULT '',
    execution_provider VARCHAR(64) DEFAULT 'judge0',
    provider_submission_id TEXT DEFAULT '',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_submissions_student_id ON submissions(student_id);
CREATE INDEX IF NOT EXISTS idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_submitted_at ON submissions(submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_submissions_student_problem ON submissions(student_id, problem_id);
