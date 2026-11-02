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
