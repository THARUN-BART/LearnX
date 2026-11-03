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

