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
-- 014_create_mistake_diagnoses.sql
-- Records intelligent diagnoses of student coding mistakes for failed submissions

CREATE TABLE IF NOT EXISTS mistake_diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    submission_id UUID NOT NULL,
    problem_id UUID NOT NULL,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    related_concept VARCHAR(128) NOT NULL,
    detection_method VARCHAR(64) NOT NULL,
    confidence_level VARCHAR(32) NOT NULL,
    confidence_score NUMERIC(5, 2) NOT NULL DEFAULT 0.0,
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    suggested_correction TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_diagnosis_submission UNIQUE (submission_id),
    CONSTRAINT fk_diag_submission FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_diag_problem FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_diagnoses_student ON mistake_diagnoses (student_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_problem ON mistake_diagnoses (problem_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_category ON mistake_diagnoses (category, subtype);
CREATE INDEX IF NOT EXISTS idx_diagnoses_created ON mistake_diagnoses (created_at DESC);
-- 015_create_student_mistake_patterns.sql
-- Aggregates and tracks recurring student weakness patterns across submissions

CREATE TABLE IF NOT EXISTS student_mistake_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    related_concept VARCHAR(128) NOT NULL,
    occurrence_count INT NOT NULL DEFAULT 1,
    first_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    correction_status VARCHAR(32) NOT NULL DEFAULT 'active',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_mistake_pattern UNIQUE (student_id, category, subtype, related_concept)
);

CREATE INDEX IF NOT EXISTS idx_mistake_patterns_student ON student_mistake_patterns (student_id);
CREATE INDEX IF NOT EXISTS idx_mistake_patterns_category ON student_mistake_patterns (category, subtype);
CREATE INDEX IF NOT EXISTS idx_mistake_patterns_status ON student_mistake_patterns (correction_status);
CREATE INDEX IF NOT EXISTS idx_mistake_patterns_count ON student_mistake_patterns (occurrence_count DESC);
-- 016_create_ai_tutor_tables.sql
-- Conversation history, progressive hint logs, and student feedback for AI Tutor

CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    conversation_type VARCHAR(64) NOT NULL DEFAULT 'general',
    title VARCHAR(255) NOT NULL DEFAULT 'New Conversation',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_conversations_student ON ai_conversations (student_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_problem ON ai_conversations (problem_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_updated ON ai_conversations (updated_at DESC);

CREATE TABLE IF NOT EXISTS ai_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL,
    message TEXT NOT NULL,
    response_type VARCHAR(64) NOT NULL DEFAULT 'text',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    token_usage JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ai_messages_conv ON ai_messages (conversation_id, created_at ASC);

CREATE TABLE IF NOT EXISTS ai_hint_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    hint_level INT NOT NULL,
    hint_content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_problem_hint_level UNIQUE (student_id, problem_id, hint_level)
);

CREATE INDEX IF NOT EXISTS idx_hint_history_student ON ai_hint_history (student_id, problem_id);

CREATE TABLE IF NOT EXISTS mistake_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    diagnosis_id UUID NOT NULL REFERENCES mistake_diagnoses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL,
    feedback_type VARCHAR(32) NOT NULL,
    student_feedback TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_feedback_diagnosis ON mistake_feedback (diagnosis_id);
CREATE INDEX IF NOT EXISTS idx_feedback_student ON mistake_feedback (student_id);
-- 017_rls_mistakes_and_ai.sql
-- Enables RLS and defines policies for mistake engine and AI tutor tables

ALTER TABLE mistake_diagnoses ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_mistake_patterns ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_hint_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE mistake_feedback ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    -- mistake_diagnoses
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'mistake_diagnoses' AND policyname = 'students_select_own_diagnoses') THEN
        CREATE POLICY students_select_own_diagnoses ON mistake_diagnoses
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    -- student_mistake_patterns
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_mistake_patterns' AND policyname = 'students_select_own_patterns') THEN
        CREATE POLICY students_select_own_patterns ON student_mistake_patterns
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    -- ai_conversations
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_conversations' AND policyname = 'students_manage_own_conversations') THEN
        CREATE POLICY students_manage_own_conversations ON ai_conversations
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    -- ai_messages
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_messages' AND policyname = 'students_manage_own_messages') THEN
        CREATE POLICY students_manage_own_messages ON ai_messages
            FOR ALL TO authenticated
            USING (
                conversation_id IN (
                    SELECT id FROM ai_conversations WHERE student_id = auth.uid()
                )
            )
            WITH CHECK (
                conversation_id IN (
                    SELECT id FROM ai_conversations WHERE student_id = auth.uid()
                )
            );
    END IF;

    -- ai_hint_history
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_hint_history' AND policyname = 'students_select_own_hints') THEN
        CREATE POLICY students_select_own_hints ON ai_hint_history
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_hint_history' AND policyname = 'students_insert_own_hints') THEN
        CREATE POLICY students_insert_own_hints ON ai_hint_history
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;

    -- mistake_feedback
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'mistake_feedback' AND policyname = 'students_manage_own_feedback') THEN
        CREATE POLICY students_manage_own_feedback ON mistake_feedback
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;
-- 018_create_student_learning_memory.sql
-- Persistent cognitive profile and learning memory for each student

CREATE TABLE IF NOT EXISTS student_concept_memory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    subtopic_id UUID REFERENCES subtopics(id) ON DELETE SET NULL,
    concept_key VARCHAR(128) NOT NULL,
    learning_state VARCHAR(32) NOT NULL DEFAULT 'discovered',
    estimated_mastery NUMERIC(5, 2) NOT NULL DEFAULT 0.0,
    confidence_level VARCHAR(32) NOT NULL DEFAULT 'low',
    total_attempts INTEGER NOT NULL DEFAULT 0,
    independent_solves INTEGER NOT NULL DEFAULT 0,
    assisted_solves INTEGER NOT NULL DEFAULT 0,
    failed_attempts INTEGER NOT NULL DEFAULT 0,
    hint_requests INTEGER NOT NULL DEFAULT 0,
    mistake_count INTEGER NOT NULL DEFAULT 0,
    corrected_mistakes INTEGER NOT NULL DEFAULT 0,
    first_practiced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_practiced_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_successful_at TIMESTAMPTZ,
    next_review_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_concept UNIQUE (student_id, concept_key)
);

CREATE INDEX IF NOT EXISTS idx_concept_memory_student ON student_concept_memory (student_id);
CREATE INDEX IF NOT EXISTS idx_concept_memory_topic ON student_concept_memory (student_id, topic_id);
CREATE INDEX IF NOT EXISTS idx_concept_memory_state ON student_concept_memory (student_id, learning_state);
CREATE INDEX IF NOT EXISTS idx_concept_memory_updated ON student_concept_memory (student_id, updated_at DESC);

-- Append-only record of important learning activities
CREATE TABLE IF NOT EXISTS student_learning_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    event_type VARCHAR(64) NOT NULL,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    diagnosis_id UUID REFERENCES mistake_diagnoses(id) ON DELETE SET NULL,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    subtopic_id UUID REFERENCES subtopics(id) ON DELETE SET NULL,
    concept_key VARCHAR(128),
    event_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_learning_events_student ON student_learning_events (student_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_learning_events_problem ON student_learning_events (problem_id);
CREATE INDEX IF NOT EXISTS idx_learning_events_concept ON student_learning_events (student_id, concept_key);
CREATE INDEX IF NOT EXISTS idx_learning_events_type ON student_learning_events (student_id, event_type);

-- Recurring mistake memory tracking active vs resolved weakness patterns
CREATE TABLE IF NOT EXISTS student_mistake_memory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    concept_key VARCHAR(128) NOT NULL,
    occurrence_count INTEGER NOT NULL DEFAULT 1,
    corrected_count INTEGER NOT NULL DEFAULT 0,
    first_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_corrected_at TIMESTAMPTZ,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    supporting_diagnosis_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_mistake_memory UNIQUE (student_id, category, subtype, concept_key)
);

CREATE INDEX IF NOT EXISTS idx_mistake_memory_student ON student_mistake_memory (student_id);
CREATE INDEX IF NOT EXISTS idx_mistake_memory_status ON student_mistake_memory (student_id, status);
CREATE INDEX IF NOT EXISTS idx_mistake_memory_concept ON student_mistake_memory (student_id, concept_key);

-- Historical learning snapshots for efficient trend charting
CREATE TABLE IF NOT EXISTS student_learning_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    snapshot_date DATE NOT NULL,
    total_concepts_practised INTEGER NOT NULL DEFAULT 0,
    independent_solves INTEGER NOT NULL DEFAULT 0,
    assisted_solves INTEGER NOT NULL DEFAULT 0,
    mistake_count INTEGER NOT NULL DEFAULT 0,
    corrected_mistake_count INTEGER NOT NULL DEFAULT 0,
    active_learning_days INTEGER NOT NULL DEFAULT 1,
    summary JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_snapshot_date UNIQUE (student_id, snapshot_date)
);

CREATE INDEX IF NOT EXISTS idx_learning_snapshots_student ON student_learning_snapshots (student_id, snapshot_date DESC);

-- Enable RLS on all 4 learning memory tables
ALTER TABLE student_concept_memory ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_mistake_memory ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_snapshots ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_concept_memory' AND policyname = 'students_manage_own_concept_memory') THEN
        CREATE POLICY students_manage_own_concept_memory ON student_concept_memory
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_events' AND policyname = 'students_select_own_learning_events') THEN
        CREATE POLICY students_select_own_learning_events ON student_learning_events
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_mistake_memory' AND policyname = 'students_manage_own_mistake_memory') THEN
        CREATE POLICY students_manage_own_mistake_memory ON student_mistake_memory
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_snapshots' AND policyname = 'students_select_own_snapshots') THEN
        CREATE POLICY students_select_own_snapshots ON student_learning_snapshots
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;
END $$;
-- 019_create_adaptive_learning_path.sql
-- Persistent Adaptive Learning Path, Node Progression States, and Recommendation History

CREATE TABLE IF NOT EXISTS student_adaptive_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    current_topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    current_concept_key VARCHAR(128) NOT NULL DEFAULT 'arrays',
    recommended_difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy',
    overall_readiness_score NUMERIC(5, 2) NOT NULL DEFAULT 50.0,
    learning_stage VARCHAR(64) NOT NULL DEFAULT 'Fundamentals & Foundation',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_adaptive_path UNIQUE (student_id)
);

CREATE INDEX IF NOT EXISTS idx_adaptive_path_student ON student_adaptive_paths (student_id);

-- Computed state for each node along the curriculum DAG
CREATE TABLE IF NOT EXISTS student_path_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    topic_name VARCHAR(128) NOT NULL,
    concept_key VARCHAR(128) NOT NULL,
    concept_name VARCHAR(128) NOT NULL,
    node_order INT NOT NULL DEFAULT 1,
    state VARCHAR(32) NOT NULL DEFAULT 'not_started',
    prerequisite_concept_keys TEXT[] NOT NULL DEFAULT '{}',
    readiness_score NUMERIC(5, 2) NOT NULL DEFAULT 0.0,
    recommended_difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy',
    recommended_problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_path_node UNIQUE (student_id, concept_key),
    CONSTRAINT chk_node_state CHECK (state IN (
        'not_started',
        'ready_to_learn',
        'in_progress',
        'practicing',
        'needs_revision',
        'mastered',
        'locked_by_prerequisite'
    ))
);

CREATE INDEX IF NOT EXISTS idx_path_nodes_student ON student_path_nodes (student_id);
CREATE INDEX IF NOT EXISTS idx_path_nodes_state ON student_path_nodes (student_id, state);
CREATE INDEX IF NOT EXISTS idx_path_nodes_order ON student_path_nodes (student_id, node_order);

-- Append-only recommendation history answering "What should this student learn next, and why?"
CREATE TABLE IF NOT EXISTS student_learning_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    action_type VARCHAR(64) NOT NULL,
    concept_key VARCHAR(128) NOT NULL,
    topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    recommended_difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy',
    headline TEXT NOT NULL,
    explanation TEXT NOT NULL,
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_recommendations_student ON student_learning_recommendations (student_id, is_active);
CREATE INDEX IF NOT EXISTS idx_recommendations_created ON student_learning_recommendations (student_id, created_at DESC);

-- Enable Row Level Security (RLS)
ALTER TABLE student_adaptive_paths ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_path_nodes ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_recommendations ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_adaptive_paths' AND policyname = 'students_manage_own_adaptive_path') THEN
        CREATE POLICY students_manage_own_adaptive_path ON student_adaptive_paths
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_path_nodes' AND policyname = 'students_manage_own_path_nodes') THEN
        CREATE POLICY students_manage_own_path_nodes ON student_path_nodes
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_recommendations' AND policyname = 'students_manage_own_recommendations') THEN
        CREATE POLICY students_manage_own_recommendations ON student_learning_recommendations
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;


-- 020_create_problem_recommendations.sql
-- Personalised Problem Recommendation Engine: Dynamic problem suggestions, explanations, feedback loop

CREATE TABLE IF NOT EXISTS student_problem_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    category VARCHAR(64) NOT NULL DEFAULT 'recommended_for_you',
    recommended_difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy',
    primary_concept VARCHAR(128) NOT NULL,
    related_concepts TEXT[] NOT NULL DEFAULT '{}',
    match_score NUMERIC(5, 2) NOT NULL DEFAULT 50.0,
    headline TEXT NOT NULL,
    explanation TEXT NOT NULL,
    learning_objective TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    is_top_pick BOOLEAN NOT NULL DEFAULT FALSE,
    feedback_rating VARCHAR(32),
    feedback_notes TEXT,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rec_status CHECK (status IN ('active', 'completed', 'skipped', 'expired')),
    CONSTRAINT chk_feedback_rating CHECK (feedback_rating IS NULL OR feedback_rating IN ('helpful', 'not_relevant', 'too_easy', 'too_difficult'))
);

CREATE INDEX IF NOT EXISTS idx_problem_rec_student ON student_problem_recommendations (student_id, status);
CREATE INDEX IF NOT EXISTS idx_problem_rec_category ON student_problem_recommendations (student_id, category);
CREATE INDEX IF NOT EXISTS idx_problem_rec_problem ON student_problem_recommendations (student_id, problem_id);
CREATE INDEX IF NOT EXISTS idx_problem_rec_top_pick ON student_problem_recommendations (student_id, is_top_pick);
CREATE INDEX IF NOT EXISTS idx_problem_rec_created ON student_problem_recommendations (student_id, created_at DESC);

-- Enable Row Level Security (RLS)
ALTER TABLE student_problem_recommendations ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_recommendations' AND policyname = 'students_manage_own_problem_recommendations') THEN
        CREATE POLICY students_manage_own_problem_recommendations ON student_problem_recommendations
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;


-- Migration: 021_create_visualizer_and_code_reviews.sql
-- Description: Schema for Algorithm Visualizer Sessions and AI Code Reviews

-- 1. Visualizer Sessions Table
CREATE TABLE IF NOT EXISTS visualizer_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    algorithm_id VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    input_config JSONB NOT NULL DEFAULT '{}'::jsonb,
    steps_completed INT NOT NULL DEFAULT 0,
    total_steps INT NOT NULL DEFAULT 0,
    time_spent_seconds INT NOT NULL DEFAULT 0,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_student_id ON visualizer_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_algorithm_id ON visualizer_sessions(algorithm_id);
CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_updated_at ON visualizer_sessions(updated_at DESC);

-- 2. Code Reviews Table
CREATE TABLE IF NOT EXISTS code_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    language VARCHAR(50) NOT NULL,
    source_code TEXT NOT NULL,
    overall_score INT NOT NULL CHECK (overall_score >= 0 AND overall_score <= 100),
    readability_score INT NOT NULL CHECK (readability_score >= 0 AND readability_score <= 100),
    efficiency_score INT NOT NULL CHECK (efficiency_score >= 0 AND efficiency_score <= 100),
    maintainability_score INT NOT NULL CHECK (maintainability_score >= 0 AND maintainability_score <= 100),
    review_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    optimized_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_code_reviews_student_id ON code_reviews(student_id);
CREATE INDEX IF NOT EXISTS idx_code_reviews_problem_id ON code_reviews(problem_id);
CREATE INDEX IF NOT EXISTS idx_code_reviews_created_at ON code_reviews(created_at DESC);

-- 3. Row-Level Security (RLS)
ALTER TABLE visualizer_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE code_reviews ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS "Students can view their own visualizer sessions" ON visualizer_sessions;
CREATE POLICY "Students can view their own visualizer sessions"
    ON visualizer_sessions FOR SELECT
    USING (auth.uid() = student_id);

DROP POLICY IF EXISTS "Students can insert their own visualizer sessions" ON visualizer_sessions;
CREATE POLICY "Students can insert their own visualizer sessions"
    ON visualizer_sessions FOR INSERT
    WITH CHECK (auth.uid() = student_id);

DROP POLICY IF EXISTS "Students can update their own visualizer sessions" ON visualizer_sessions;
CREATE POLICY "Students can update their own visualizer sessions"
    ON visualizer_sessions FOR UPDATE
    USING (auth.uid() = student_id);

DROP POLICY IF EXISTS "Students can view their own code reviews" ON code_reviews;
CREATE POLICY "Students can view their own code reviews"
    ON code_reviews FOR SELECT
    USING (auth.uid() = student_id);

DROP POLICY IF EXISTS "Students can insert their own code reviews" ON code_reviews;
CREATE POLICY "Students can insert their own code reviews"
    ON code_reviews FOR INSERT
    WITH CHECK (auth.uid() = student_id);
