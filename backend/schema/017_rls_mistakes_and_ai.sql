-- 017_rls_mistakes_and_ai.sql
-- Enables RLS and defines policies for mistake engine and AI tutor tables

ALTER TABLE mistake_diagnoses ENABLE ROW LEVEL SECURITY;
ALTER TABLE student_mistake_patterns ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_hint_history ENABLE ROW LEVEL SECURITY;
ALTER TABLE mistake_feedback ENABLE ROW LEVEL SECURITY;

DO $$
BEGIN
    -- mistake_diagnoses
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'mistake_diagnoses' AND policyname = 'students_select_own_diagnoses') THEN
        CREATE POLICY students_select_own_diagnoses ON mistake_diagnoses
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    -- student_mistake_patterns
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_mistake_patterns' AND policyname = 'students_select_own_patterns') THEN
        CREATE POLICY students_select_own_patterns ON student_mistake_patterns
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    -- ai_conversations
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_conversations' AND policyname = 'students_manage_own_conversations') THEN
        CREATE POLICY students_manage_own_conversations ON ai_conversations
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    -- ai_messages
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_messages' AND policyname = 'students_manage_own_messages') THEN
        CREATE POLICY students_manage_own_messages ON ai_messages
            FOR ALL TO authenticated
            USING (
                conversation_id IN (
                    SELECT id FROM ai_conversations WHERE student_id = auth.uid()
                )
            )
            WITH CHECK (
                conversation_id IN (
                    SELECT id FROM ai_conversations WHERE student_id = auth.uid()
                )
            );
    END IF;

    -- ai_hint_history
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_hint_history' AND policyname = 'students_select_own_hints') THEN
        CREATE POLICY students_select_own_hints ON ai_hint_history
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'ai_hint_history' AND policyname = 'students_insert_own_hints') THEN
        CREATE POLICY students_insert_own_hints ON ai_hint_history
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;

    -- mistake_feedback
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'mistake_feedback' AND policyname = 'students_manage_own_feedback') THEN
        CREATE POLICY students_manage_own_feedback ON mistake_feedback
            FOR ALL TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;
