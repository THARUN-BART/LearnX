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

