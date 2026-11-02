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
