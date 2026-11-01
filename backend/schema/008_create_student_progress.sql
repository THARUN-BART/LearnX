-- 008_create_student_progress.sql
-- Per-student progress is isolated by student_id. One student can never mutate another student's rows.

CREATE TABLE IF NOT EXISTS student_problem_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'Not Started',
    total_attempts INTEGER NOT NULL DEFAULT 0,
    successful_attempts INTEGER NOT NULL DEFAULT 0,
    first_attempt_at TIMESTAMPTZ,
    last_attempt_at TIMESTAMPTZ,
    solved_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_problem UNIQUE (student_id, problem_id),
    CONSTRAINT chk_progress_status CHECK (status IN ('Not Started', 'Attempted', 'Solved', 'Needs Revision')),
    CONSTRAINT chk_progress_attempts CHECK (total_attempts >= 0 AND successful_attempts >= 0 AND successful_attempts <= total_attempts)
);

CREATE INDEX IF NOT EXISTS idx_progress_student ON student_problem_progress (student_id);
CREATE INDEX IF NOT EXISTS idx_progress_problem ON student_problem_progress (problem_id);
CREATE INDEX IF NOT EXISTS idx_progress_student_status ON student_problem_progress (student_id, status);
CREATE INDEX IF NOT EXISTS idx_progress_last_attempt ON student_problem_progress (student_id, last_attempt_at DESC);

-- Future Adaptive Learning Engine hook: append-only events (submissions, mistakes, hints, revisions).
-- No AI behaviour is implemented in this migration.
CREATE TABLE IF NOT EXISTS student_learning_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_learning_event_type CHECK (event_type IN (
        'submission',
        'mistake',
        'hint_request',
        'revision_scheduled',
        'concept_mastery',
        'recommendation'
    ))
);

CREATE INDEX IF NOT EXISTS idx_learning_events_student ON student_learning_events (student_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_learning_events_problem ON student_learning_events (problem_id);
CREATE INDEX IF NOT EXISTS idx_learning_events_type ON student_learning_events (event_type);
