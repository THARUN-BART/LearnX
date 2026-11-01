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
