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
