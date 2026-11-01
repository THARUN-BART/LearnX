-- Complete Database Schema for LearnX Authentication, Student Profiles, and DSA Roadmaps

-- 001_create_profiles.sql
-- Student Profile schema with full onboarding fields

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY,
    full_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    college TEXT NOT NULL DEFAULT '',
    course TEXT NOT NULL DEFAULT '',
    department TEXT NOT NULL DEFAULT '',
    current_year TEXT NOT NULL DEFAULT '',
    preferred_language TEXT NOT NULL DEFAULT '',
    learning_goal TEXT NOT NULL DEFAULT '',
    dsa_level TEXT NOT NULL DEFAULT '',
    avatar_url TEXT NOT NULL DEFAULT '',
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);
CREATE INDEX IF NOT EXISTS idx_profiles_onboarding ON profiles(onboarding_completed);

-- 002_create_pending_signups.sql
-- Pending signups waiting for OTP verification

CREATE TABLE IF NOT EXISTS pending_signups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    password_encrypted TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pending_signups_email ON pending_signups(email);
CREATE INDEX IF NOT EXISTS idx_pending_signups_expires ON pending_signups(expires_at);

-- 003_create_password_resets.sql
-- Single-use password reset authorization tokens

CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_password_resets_token_hash ON password_resets(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_resets_expires ON password_resets(expires_at);

-- 004_create_roadmaps.sql
-- Shared DSA roadmap catalogue (readable by authenticated students; writes are admin/service-only)

CREATE TABLE IF NOT EXISTS roadmaps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    source_name TEXT NOT NULL DEFAULT '',
    source_url TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_roadmaps_active ON roadmaps (is_active);
CREATE INDEX IF NOT EXISTS idx_roadmaps_source_name ON roadmaps (source_name);

-- 005_create_topics.sql

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_id UUID NOT NULL REFERENCES roadmaps(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    topic_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_topics_roadmap_slug UNIQUE (roadmap_id, slug),
    CONSTRAINT uq_topics_roadmap_order UNIQUE (roadmap_id, topic_order),
    CONSTRAINT chk_topics_order_positive CHECK (topic_order > 0)
);

CREATE INDEX IF NOT EXISTS idx_topics_roadmap_id ON topics (roadmap_id);
CREATE INDEX IF NOT EXISTS idx_topics_order ON topics (roadmap_id, topic_order);

-- 006_create_subtopics.sql

CREATE TABLE IF NOT EXISTS subtopics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    subtopic_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_subtopics_topic_slug UNIQUE (topic_id, slug),
    CONSTRAINT uq_subtopics_topic_order UNIQUE (topic_id, subtopic_order),
    CONSTRAINT chk_subtopics_order_positive CHECK (subtopic_order > 0)
);

CREATE INDEX IF NOT EXISTS idx_subtopics_topic_id ON subtopics (topic_id);
CREATE INDEX IF NOT EXISTS idx_subtopics_order ON subtopics (topic_id, subtopic_order);

-- 007_create_problems.sql
-- Central problem catalogue. Stores metadata and canonical external URLs only (no editorials/solutions).

CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    subtopic_id UUID NOT NULL REFERENCES subtopics(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    summary TEXT NOT NULL DEFAULT '',
    difficulty TEXT NOT NULL,
    external_platform TEXT NOT NULL DEFAULT '',
    external_problem_url TEXT NOT NULL,
    source_sheet TEXT NOT NULL DEFAULT '',
    source_problem_id TEXT NOT NULL DEFAULT '',
    problem_order INTEGER NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_problems_difficulty CHECK (difficulty IN ('Easy', 'Medium', 'Hard', 'Unspecified')),
    CONSTRAINT chk_problems_order_positive CHECK (problem_order > 0),
    CONSTRAINT uq_problems_source UNIQUE (source_sheet, source_problem_id),
    CONSTRAINT uq_problems_subtopic_order UNIQUE (subtopic_id, problem_order)
);

CREATE INDEX IF NOT EXISTS idx_problems_topic_id ON problems (topic_id);
CREATE INDEX IF NOT EXISTS idx_problems_subtopic_id ON problems (subtopic_id);
CREATE INDEX IF NOT EXISTS idx_problems_difficulty ON problems (difficulty);
CREATE INDEX IF NOT EXISTS idx_problems_order ON problems (problem_order);
CREATE INDEX IF NOT EXISTS idx_problems_title_lower ON problems (LOWER(title));
CREATE INDEX IF NOT EXISTS idx_problems_source_sheet ON problems (source_sheet);
CREATE INDEX IF NOT EXISTS idx_problems_active ON problems (is_active);
CREATE INDEX IF NOT EXISTS idx_problems_external_url ON problems (external_problem_url);

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

-- 009_rls_dsa.sql
-- Row Level Security for the DSA catalogue and student-owned progress.
-- Policies are applied only when the Supabase `auth` schema is present.

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'auth') THEN
        RAISE NOTICE 'Skipping DSA RLS policies: auth schema not present (local/dev PostgreSQL)';
        RETURN;
    END IF;

    ALTER TABLE roadmaps ENABLE ROW LEVEL SECURITY;
    ALTER TABLE topics ENABLE ROW LEVEL SECURITY;
    ALTER TABLE subtopics ENABLE ROW LEVEL SECURITY;
    ALTER TABLE problems ENABLE ROW LEVEL SECURITY;
    ALTER TABLE student_problem_progress ENABLE ROW LEVEL SECURITY;
    ALTER TABLE student_learning_events ENABLE ROW LEVEL SECURITY;

    -- Catalogue is readable by signed-in students. Mutations stay with service_role (bypasses RLS).
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'roadmaps' AND policyname = 'roadmaps_select_authenticated') THEN
        CREATE POLICY roadmaps_select_authenticated ON roadmaps
            FOR SELECT TO authenticated
            USING (is_active = TRUE);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'topics' AND policyname = 'topics_select_authenticated') THEN
        CREATE POLICY topics_select_authenticated ON topics
            FOR SELECT TO authenticated
            USING (
                EXISTS (SELECT 1 FROM roadmaps r WHERE r.id = topics.roadmap_id AND r.is_active = TRUE)
            );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'subtopics' AND policyname = 'subtopics_select_authenticated') THEN
        CREATE POLICY subtopics_select_authenticated ON subtopics
            FOR SELECT TO authenticated
            USING (
                EXISTS (
                    SELECT 1 FROM topics t
                    JOIN roadmaps r ON r.id = t.roadmap_id
                    WHERE t.id = subtopics.topic_id AND r.is_active = TRUE
                )
            );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'problems' AND policyname = 'problems_select_authenticated') THEN
        CREATE POLICY problems_select_authenticated ON problems
            FOR SELECT TO authenticated
            USING (is_active = TRUE);
    END IF;

    -- Students may only read/write their own progress rows.
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_select_own') THEN
        CREATE POLICY progress_select_own ON student_problem_progress
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_insert_own') THEN
        CREATE POLICY progress_insert_own ON student_problem_progress
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_update_own') THEN
        CREATE POLICY progress_update_own ON student_problem_progress
            FOR UPDATE TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_delete_own') THEN
        CREATE POLICY progress_delete_own ON student_problem_progress
            FOR DELETE TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_events' AND policyname = 'learning_events_select_own') THEN
        CREATE POLICY learning_events_select_own ON student_learning_events
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_events' AND policyname = 'learning_events_insert_own') THEN
        CREATE POLICY learning_events_insert_own ON student_learning_events
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;


-- ============================================================================
-- Problem Submission & Code Assessment Tables (010, 011, 012, 013)
-- ============================================================================

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

CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input_data TEXT NOT NULL DEFAULT '',
    expected_output TEXT NOT NULL DEFAULT '',
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    test_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id ON test_cases(problem_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_order ON test_cases(problem_id, test_order ASC);
CREATE INDEX IF NOT EXISTS idx_test_cases_is_hidden ON test_cases(problem_id, is_hidden);

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

ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE test_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE submission_test_results ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'submissions' AND policyname = 'students_select_own_submissions') THEN
        CREATE POLICY students_select_own_submissions ON submissions
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'submissions' AND policyname = 'students_insert_own_submissions') THEN
        CREATE POLICY students_insert_own_submissions ON submissions
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'test_cases' AND policyname = 'authenticated_select_sample_test_cases') THEN
        CREATE POLICY authenticated_select_sample_test_cases ON test_cases
            FOR SELECT TO authenticated
            USING (is_hidden = false);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'submission_test_results' AND policyname = 'students_select_own_test_results') THEN
        CREATE POLICY students_select_own_test_results ON submission_test_results
            FOR SELECT TO authenticated
            USING (
                submission_id IN (
                    SELECT id FROM submissions WHERE student_id = auth.uid()
                )
            );
    END IF;
END $$;
