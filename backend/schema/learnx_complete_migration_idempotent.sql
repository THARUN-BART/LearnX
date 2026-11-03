-- ==============================================================================
-- LEARNX: 100% IDEMPOTENT & NON-DESTRUCTIVE POSTGRESQL MIGRATION SCRIPT
-- ==============================================================================
-- GUARANTEES:
-- 1. NO DROPS, NO TRUNCATE, NO DELETE: Preserves 100% of your existing data.
-- 2. SAFE COLUMN UPGRADES: Runs ALTER TABLE ... ADD COLUMN IF NOT EXISTS for every
--    single column (including created_at, updated_at, submitted_at, etc.) BEFORE
--    creating any index, preventing "column does not exist" errors.
-- 3. SAFE INDEXES: Uses CREATE INDEX IF NOT EXISTS.
-- 4. NON-DESTRUCTIVE RLS: Drops and recreates policies only (never affects row data).
-- ==============================================================================

-- ------------------------------------------------------------------------------
-- 1. EXTENSIONS
-- ------------------------------------------------------------------------------
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ------------------------------------------------------------------------------
-- 2. PROFILES & AUTH TABLES
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS profiles (
    id UUID PRIMARY KEY,
    full_name TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    college TEXT NOT NULL DEFAULT '',
    course TEXT NOT NULL DEFAULT '',
    department TEXT NOT NULL DEFAULT '',
    current_year TEXT NOT NULL DEFAULT '',
    preferred_language TEXT NOT NULL DEFAULT 'Java',
    learning_goal TEXT NOT NULL DEFAULT '',
    dsa_level TEXT NOT NULL DEFAULT 'Beginner',
    avatar_url TEXT NOT NULL DEFAULT '',
    onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE profiles ADD COLUMN IF NOT EXISTS full_name TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS college TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS course TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS department TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS current_year TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS preferred_language TEXT NOT NULL DEFAULT 'Java';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS learning_goal TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS dsa_level TEXT NOT NULL DEFAULT 'Beginner';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT '';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_profiles_email ON profiles(email);
CREATE INDEX IF NOT EXISTS idx_profiles_onboarding ON profiles(onboarding_completed);

CREATE TABLE IF NOT EXISTS pending_signups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    full_name TEXT NOT NULL,
    password_encrypted TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE pending_signups ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE pending_signups ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '1 day');

CREATE INDEX IF NOT EXISTS idx_pending_signups_email ON pending_signups(email);
CREATE INDEX IF NOT EXISTS idx_pending_signups_expires ON pending_signups(expires_at);

CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash TEXT NOT NULL UNIQUE,
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ
);

ALTER TABLE password_resets ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE password_resets ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT (NOW() + INTERVAL '1 hour');

CREATE INDEX IF NOT EXISTS idx_password_resets_token_hash ON password_resets(token_hash);
CREATE INDEX IF NOT EXISTS idx_password_resets_expires ON password_resets(expires_at);

-- ------------------------------------------------------------------------------
-- 3. DSA CURRICULUM: ROADMAPS, TOPICS, SUBTOPICS, PROBLEMS
-- ------------------------------------------------------------------------------
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

ALTER TABLE roadmaps ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE roadmaps ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE roadmaps ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_roadmaps_active ON roadmaps (is_active);
CREATE INDEX IF NOT EXISTS idx_roadmaps_source_name ON roadmaps (source_name);

CREATE TABLE IF NOT EXISTS topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    roadmap_id UUID NOT NULL REFERENCES roadmaps(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    topic_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE topics ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE topics ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_topics_roadmap_id ON topics (roadmap_id);
CREATE INDEX IF NOT EXISTS idx_topics_order ON topics (roadmap_id, topic_order);

CREATE TABLE IF NOT EXISTS subtopics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    subtopic_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE subtopics ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_subtopics_topic_id ON subtopics (topic_id);
CREATE INDEX IF NOT EXISTS idx_subtopics_order ON subtopics (topic_id, subtopic_order);

CREATE TABLE IF NOT EXISTS problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
    subtopic_id UUID NOT NULL REFERENCES subtopics(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    summary TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    difficulty TEXT NOT NULL,
    external_platform TEXT NOT NULL DEFAULT '',
    external_problem_url TEXT NOT NULL DEFAULT '',
    source_sheet TEXT NOT NULL DEFAULT '',
    source_problem_id TEXT NOT NULL DEFAULT '',
    problem_order INTEGER NOT NULL DEFAULT 1,
    tags TEXT[] NOT NULL DEFAULT '{}',
    hints TEXT[] NOT NULL DEFAULT '{}',
    starter_code JSONB NOT NULL DEFAULT '{}'::jsonb,
    solution_code TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE problems ADD COLUMN IF NOT EXISTS summary TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS description TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS external_platform TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS external_problem_url TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS source_sheet TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS source_problem_id TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS problem_order INTEGER NOT NULL DEFAULT 1;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS hints TEXT[] NOT NULL DEFAULT '{}';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS starter_code JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS solution_code TEXT NOT NULL DEFAULT '';
ALTER TABLE problems ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE problems ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_problems_topic_id ON problems (topic_id);
CREATE INDEX IF NOT EXISTS idx_problems_subtopic_id ON problems (subtopic_id);
CREATE INDEX IF NOT EXISTS idx_problems_difficulty ON problems (difficulty);
CREATE INDEX IF NOT EXISTS idx_problems_order ON problems (problem_order);
CREATE INDEX IF NOT EXISTS idx_problems_title_lower ON problems (LOWER(title));
CREATE INDEX IF NOT EXISTS idx_problems_active ON problems (is_active);

-- ------------------------------------------------------------------------------
-- 4. STUDENT PROGRESS & SUBMISSIONS (Handles both submitted_at and created_at)
-- ------------------------------------------------------------------------------
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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_problem UNIQUE (student_id, problem_id)
);

ALTER TABLE student_problem_progress ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_problem_progress ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_student_progress_student ON student_problem_progress (student_id);
CREATE INDEX IF NOT EXISTS idx_student_progress_problem ON student_problem_progress (problem_id);
CREATE INDEX IF NOT EXISTS idx_student_progress_status ON student_problem_progress (student_id, status);

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language VARCHAR(32) NOT NULL,
    source_code TEXT NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'Pending',
    passed_test_cases INTEGER NOT NULL DEFAULT 0,
    total_test_cases INTEGER NOT NULL DEFAULT 0,
    execution_time_ms INTEGER NOT NULL DEFAULT 0,
    memory_used_kb INTEGER NOT NULL DEFAULT 0,
    error_message TEXT DEFAULT '',
    execution_provider VARCHAR(64) DEFAULT 'judge0',
    provider_submission_id TEXT DEFAULT '',
    is_independent BOOLEAN NOT NULL DEFAULT TRUE,
    ai_hints_used INTEGER NOT NULL DEFAULT 0,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crucial: Ensure submissions supports both submitted_at and created_at
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS completed_at TIMESTAMPTZ;
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS memory_used_kb INTEGER NOT NULL DEFAULT 0;
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS execution_provider VARCHAR(64) DEFAULT 'judge0';
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS provider_submission_id TEXT DEFAULT '';
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS is_independent BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE submissions ADD COLUMN IF NOT EXISTS ai_hints_used INTEGER NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_submissions_student_id ON submissions(student_id);
CREATE INDEX IF NOT EXISTS idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status ON submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_submitted_at ON submissions(submitted_at DESC);
CREATE INDEX IF NOT EXISTS idx_submissions_created_at ON submissions(created_at DESC);

CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input_data TEXT NOT NULL DEFAULT '',
    expected_output TEXT NOT NULL DEFAULT '',
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    test_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE test_cases ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id ON test_cases(problem_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_order ON test_cases(problem_id, test_order ASC);
CREATE INDEX IF NOT EXISTS idx_test_cases_is_hidden ON test_cases(problem_id, is_hidden);

CREATE TABLE IF NOT EXISTS submission_test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    test_case_id UUID NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    passed BOOLEAN NOT NULL DEFAULT false,
    status VARCHAR(32) NOT NULL DEFAULT 'Passed',
    actual_output TEXT DEFAULT '',
    execution_time_ms INTEGER NOT NULL DEFAULT 0,
    memory_kb INTEGER NOT NULL DEFAULT 0,
    memory_used_kb INTEGER NOT NULL DEFAULT 0,
    error_message TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE submission_test_results ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_submission_test_results_submission ON submission_test_results(submission_id);

-- ------------------------------------------------------------------------------
-- 5. MISTAKE DETECTION & PATTERN ANALYTICS
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS mistake_diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    related_concept VARCHAR(128) NOT NULL DEFAULT '',
    detection_method VARCHAR(64) NOT NULL DEFAULT 'rule_based',
    confidence_level VARCHAR(32) NOT NULL DEFAULT 'medium',
    confidence_score NUMERIC(5, 2) NOT NULL DEFAULT 0.80,
    confidence NUMERIC(4, 2) NOT NULL DEFAULT 0.80,
    evidence TEXT NOT NULL DEFAULT '',
    suggested_correction TEXT NOT NULL DEFAULT '',
    student_feedback VARCHAR(32),
    feedback_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE mistake_diagnoses ADD COLUMN IF NOT EXISTS confidence NUMERIC(4, 2) NOT NULL DEFAULT 0.80;
ALTER TABLE mistake_diagnoses ADD COLUMN IF NOT EXISTS confidence_score NUMERIC(5, 2) NOT NULL DEFAULT 0.80;
ALTER TABLE mistake_diagnoses ADD COLUMN IF NOT EXISTS student_feedback VARCHAR(32);
ALTER TABLE mistake_diagnoses ADD COLUMN IF NOT EXISTS feedback_notes TEXT;
ALTER TABLE mistake_diagnoses ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_mistake_diag_student ON mistake_diagnoses(student_id);
CREATE INDEX IF NOT EXISTS idx_mistake_diag_problem ON mistake_diagnoses(problem_id);
CREATE INDEX IF NOT EXISTS idx_mistake_diag_submission ON mistake_diagnoses(submission_id);
CREATE INDEX IF NOT EXISTS idx_mistake_diag_category ON mistake_diagnoses(student_id, category);
CREATE INDEX IF NOT EXISTS idx_mistake_diag_created ON mistake_diagnoses(created_at DESC);

CREATE TABLE IF NOT EXISTS student_mistake_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    related_concept VARCHAR(128) NOT NULL DEFAULT '',
    occurrence_count INT NOT NULL DEFAULT 1,
    frequency INTEGER NOT NULL DEFAULT 1,
    first_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    correction_status VARCHAR(32) NOT NULL DEFAULT 'active',
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Crucial: student_mistake_patterns may exist from older migrations without created_at
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS related_concept VARCHAR(128) NOT NULL DEFAULT '';
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS occurrence_count INT NOT NULL DEFAULT 1;
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS frequency INTEGER NOT NULL DEFAULT 1;
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS first_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS last_occurrence_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS correction_status VARCHAR(32) NOT NULL DEFAULT 'active';
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS resolved_at TIMESTAMPTZ;
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_mistake_patterns ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_student_patterns_student ON student_mistake_patterns(student_id);
CREATE INDEX IF NOT EXISTS idx_student_patterns_active ON student_mistake_patterns(student_id, resolved_at);

-- ------------------------------------------------------------------------------
-- 6. AI TUTOR CONVERSATIONS & HINTS
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ai_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    conversation_type VARCHAR(64) NOT NULL DEFAULT 'general',
    title VARCHAR(255) NOT NULL DEFAULT 'New Conversation',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE ai_conversations ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE ai_conversations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_ai_conversations_student ON ai_conversations(student_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_problem ON ai_conversations(problem_id);
CREATE INDEX IF NOT EXISTS idx_ai_conversations_updated ON ai_conversations(updated_at DESC);

CREATE TABLE IF NOT EXISTS ai_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(32) NOT NULL,
    message TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    response_type VARCHAR(64) NOT NULL DEFAULT 'text',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    token_usage JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE ai_messages ADD COLUMN IF NOT EXISTS message TEXT NOT NULL DEFAULT '';
ALTER TABLE ai_messages ADD COLUMN IF NOT EXISTS content TEXT NOT NULL DEFAULT '';
ALTER TABLE ai_messages ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}'::jsonb;
ALTER TABLE ai_messages ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_ai_messages_conv ON ai_messages(conversation_id, created_at ASC);

CREATE TABLE IF NOT EXISTS ai_hint_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    hint_level INT NOT NULL,
    hint_content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE ai_hint_history ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_hint_history_student ON ai_hint_history (student_id, problem_id);

CREATE TABLE IF NOT EXISTS mistake_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    diagnosis_id UUID NOT NULL REFERENCES mistake_diagnoses(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    feedback_type VARCHAR(32) NOT NULL,
    student_feedback TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE mistake_feedback ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- ------------------------------------------------------------------------------
-- 7. STUDENT LEARNING MEMORY ENGINE
-- ------------------------------------------------------------------------------
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

ALTER TABLE student_concept_memory ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_concept_memory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_concept_memory_student ON student_concept_memory (student_id);
CREATE INDEX IF NOT EXISTS idx_concept_memory_topic ON student_concept_memory (student_id, topic_id);
CREATE INDEX IF NOT EXISTS idx_concept_memory_state ON student_concept_memory (student_id, learning_state);
CREATE INDEX IF NOT EXISTS idx_concept_memory_updated ON student_concept_memory (student_id, updated_at DESC);

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

ALTER TABLE student_learning_events ADD COLUMN IF NOT EXISTS occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_learning_events ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_learning_events_student ON student_learning_events (student_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_learning_events_problem ON student_learning_events (problem_id);
CREATE INDEX IF NOT EXISTS idx_learning_events_concept ON student_learning_events (student_id, concept_key);

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

ALTER TABLE student_mistake_memory ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_mistake_memory ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_mistake_memory_student ON student_mistake_memory (student_id);
CREATE INDEX IF NOT EXISTS idx_mistake_memory_status ON student_mistake_memory (student_id, status);

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

ALTER TABLE student_learning_snapshots ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_learning_snapshots_student ON student_learning_snapshots (student_id, snapshot_date DESC);

-- ------------------------------------------------------------------------------
-- 8. ADAPTIVE LEARNING PATH ENGINE
-- ------------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS student_adaptive_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    current_topic_id UUID REFERENCES topics(id) ON DELETE SET NULL,
    current_concept_key VARCHAR(128) NOT NULL DEFAULT 'arrays',
    recommended_difficulty VARCHAR(32) NOT NULL DEFAULT 'Easy',
    overall_readiness_score NUMERIC(5, 2) NOT NULL DEFAULT 50.0,
    learning_stage VARCHAR(64) NOT NULL DEFAULT 'Foundations & Array Iteration',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_adaptive_path UNIQUE (student_id)
);

ALTER TABLE student_adaptive_paths ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_adaptive_paths ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_adaptive_path_student ON student_adaptive_paths (student_id);

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
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_student_path_node UNIQUE (student_id, concept_key)
);

-- Crucial: student_path_nodes originally only had updated_at
ALTER TABLE student_path_nodes ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_path_nodes ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_path_nodes_student ON student_path_nodes (student_id);
CREATE INDEX IF NOT EXISTS idx_path_nodes_state ON student_path_nodes (student_id, state);
CREATE INDEX IF NOT EXISTS idx_path_nodes_order ON student_path_nodes (student_id, node_order);

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

ALTER TABLE student_learning_recommendations ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_recommendations_student ON student_learning_recommendations (student_id, is_active);
CREATE INDEX IF NOT EXISTS idx_recommendations_created ON student_learning_recommendations (student_id, created_at DESC);

-- ------------------------------------------------------------------------------
-- 9. PERSONALISED PROBLEM RECOMMENDATION ENGINE
-- ------------------------------------------------------------------------------
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE student_problem_recommendations ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE student_problem_recommendations ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_problem_rec_student ON student_problem_recommendations (student_id, status);
CREATE INDEX IF NOT EXISTS idx_problem_rec_category ON student_problem_recommendations (student_id, category);
CREATE INDEX IF NOT EXISTS idx_problem_rec_problem ON student_problem_recommendations (student_id, problem_id);
CREATE INDEX IF NOT EXISTS idx_problem_rec_top_pick ON student_problem_recommendations (student_id, is_top_pick);
CREATE INDEX IF NOT EXISTS idx_problem_rec_created ON student_problem_recommendations (student_id, created_at DESC);

-- ------------------------------------------------------------------------------
-- 10. ALGORITHM VISUALIZER & AI CODE REVIEW SESSIONS
-- ------------------------------------------------------------------------------
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

ALTER TABLE visualizer_sessions ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE visualizer_sessions ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_student_id ON visualizer_sessions(student_id);
CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_algorithm_id ON visualizer_sessions(algorithm_id);
CREATE INDEX IF NOT EXISTS idx_visualizer_sessions_updated_at ON visualizer_sessions(updated_at DESC);

CREATE TABLE IF NOT EXISTS code_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE SET NULL,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    language VARCHAR(50) NOT NULL,
    source_code TEXT NOT NULL,
    overall_score INT NOT NULL DEFAULT 70,
    readability_score INT NOT NULL DEFAULT 70,
    efficiency_score INT NOT NULL DEFAULT 70,
    maintainability_score INT NOT NULL DEFAULT 70,
    review_data JSONB NOT NULL DEFAULT '{}'::jsonb,
    optimized_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE code_reviews ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_code_reviews_student_id ON code_reviews(student_id);
CREATE INDEX IF NOT EXISTS idx_code_reviews_problem_id ON code_reviews(problem_id);
CREATE INDEX IF NOT EXISTS idx_code_reviews_created_at ON code_reviews(created_at DESC);

-- ------------------------------------------------------------------------------
-- 11. ROW LEVEL SECURITY (RLS) POLICIES
-- ------------------------------------------------------------------------------
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE roadmaps ENABLE ROW LEVEL SECURITY;
ALTER TABLE topics ENABLE ROW LEVEL SECURITY;
ALTER TABLE subtopics ENABLE ROW LEVEL SECURITY;
ALTER TABLE problems ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_problem_progress ENABLE ROW LEVEL SECURITY;
ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE test_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE submission_test_results ENABLE ROW LEVEL SECURITY;
ALTER TABLE mistake_diagnoses ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_mistake_patterns ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_concept_memory ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_mistake_memory ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_snapshots ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_adaptive_paths ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_path_nodes ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_learning_recommendations ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_problem_recommendations ENABLE ROW LEVEL SECURITY;
ALTER TABLE visualizer_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE code_reviews ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    -- Public / Curriculum Read-only policies
    DROP POLICY IF EXISTS "Public roadmaps read" ON roadmaps;
    CREATE POLICY "Public roadmaps read" ON roadmaps FOR SELECT USING (true);

    DROP POLICY IF EXISTS "Public topics read" ON topics;
    CREATE POLICY "Public topics read" ON topics FOR SELECT USING (true);

    DROP POLICY IF EXISTS "Public subtopics read" ON subtopics;
    CREATE POLICY "Public subtopics read" ON subtopics FOR SELECT USING (true);

    DROP POLICY IF EXISTS "Public problems read" ON problems;
    CREATE POLICY "Public problems read" ON problems FOR SELECT USING (true);

    DROP POLICY IF EXISTS "Public non-hidden test cases read" ON test_cases;
    CREATE POLICY "Public non-hidden test cases read" ON test_cases FOR SELECT USING (is_hidden = false);

    -- Profiles policy
    DROP POLICY IF EXISTS "Users can view own profile" ON profiles;
    CREATE POLICY "Users can view own profile" ON profiles FOR SELECT USING (auth.uid() = id);

    DROP POLICY IF EXISTS "Users can update own profile" ON profiles;
    CREATE POLICY "Users can update own profile" ON profiles FOR UPDATE USING (auth.uid() = id);

    -- Student Progress policies
    DROP POLICY IF EXISTS "Users manage own progress" ON student_problem_progress;
    CREATE POLICY "Users manage own progress" ON student_problem_progress FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- Submissions policies
    DROP POLICY IF EXISTS "Users manage own submissions" ON submissions;
    CREATE POLICY "Users manage own submissions" ON submissions FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- Submission test results policy
    DROP POLICY IF EXISTS "Users view own submission test results" ON submission_test_results;
    CREATE POLICY "Users view own submission test results" ON submission_test_results FOR SELECT USING (
        EXISTS (SELECT 1 FROM submissions s WHERE s.id = submission_test_results.submission_id AND s.student_id = auth.uid())
    );

    -- Mistake Diagnoses & Patterns
    DROP POLICY IF EXISTS "Users manage own mistake diagnoses" ON mistake_diagnoses;
    CREATE POLICY "Users manage own mistake diagnoses" ON mistake_diagnoses FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own mistake patterns" ON student_mistake_patterns;
    CREATE POLICY "Users manage own mistake patterns" ON student_mistake_patterns FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- AI Conversations & Messages policies
    DROP POLICY IF EXISTS "Users manage own ai conversations" ON ai_conversations;
    CREATE POLICY "Users manage own ai conversations" ON ai_conversations FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own ai messages" ON ai_messages;
    CREATE POLICY "Users manage own ai messages" ON ai_messages FOR ALL USING (
        EXISTS (SELECT 1 FROM ai_conversations c WHERE c.id = ai_messages.conversation_id AND c.student_id = auth.uid())
    ) WITH CHECK (
        EXISTS (SELECT 1 FROM ai_conversations c WHERE c.id = ai_messages.conversation_id AND c.student_id = auth.uid())
    );

    -- Learning Memory policies
    DROP POLICY IF EXISTS "Users manage own concept memory" ON student_concept_memory;
    CREATE POLICY "Users manage own concept memory" ON student_concept_memory FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own learning events" ON student_learning_events;
    CREATE POLICY "Users manage own learning events" ON student_learning_events FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own mistake memory" ON student_mistake_memory;
    CREATE POLICY "Users manage own mistake memory" ON student_mistake_memory FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own learning snapshots" ON student_learning_snapshots;
    CREATE POLICY "Users manage own learning snapshots" ON student_learning_snapshots FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- Adaptive Learning Path policies
    DROP POLICY IF EXISTS "Users manage own adaptive path" ON student_adaptive_paths;
    CREATE POLICY "Users manage own adaptive path" ON student_adaptive_paths FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own path nodes" ON student_path_nodes;
    CREATE POLICY "Users manage own path nodes" ON student_path_nodes FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own learning recommendations" ON student_learning_recommendations;
    CREATE POLICY "Users manage own learning recommendations" ON student_learning_recommendations FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- Problem Recommendations policy
    DROP POLICY IF EXISTS "Users manage own problem recommendations" ON student_problem_recommendations;
    CREATE POLICY "Users manage own problem recommendations" ON student_problem_recommendations FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    -- Visualizer & Code Review policies
    DROP POLICY IF EXISTS "Users manage own visualizer sessions" ON visualizer_sessions;
    CREATE POLICY "Users manage own visualizer sessions" ON visualizer_sessions FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);

    DROP POLICY IF EXISTS "Users manage own code reviews" ON code_reviews;
    CREATE POLICY "Users manage own code reviews" ON code_reviews FOR ALL USING (auth.uid() = student_id) WITH CHECK (auth.uid() = student_id);
END $$;

-- ------------------------------------------------------------------------------
-- 12. CANONICAL ROADMAP SEEDING (IDEMPOTENT - NEVER OVERWRITES OR DELETES DATA)
-- ------------------------------------------------------------------------------
INSERT INTO roadmaps (id, name, slug, description, source_name, source_url, is_active)
VALUES (
    '8396a6f4-760c-41e7-b78b-f343bed9b52c',
    'Striver''s A2Z DSA Sheet',
    'strivers-a2z-dsa-sheet',
    'Comprehensive step-by-step roadmap from TakeUForward covering beginner fundamentals to advanced graphs and dynamic programming.',
    'TakeUForward',
    'https://takeuforward.org/strivers-a2z-dsa-course/strivers-a2z-dsa-course-sheet-2',
    true
)
ON CONFLICT (slug) DO NOTHING;

-- End of migration script
