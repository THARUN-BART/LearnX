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
