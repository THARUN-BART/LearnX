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
