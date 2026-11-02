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
