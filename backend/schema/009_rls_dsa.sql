-- 009_rls_dsa.sql
-- Row Level Security for the DSA catalogue and student-owned progress.
-- Policies are applied only when the Supabase `auth` schema is present.

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'auth') THEN
        RAISE NOTICE 'Skipping DSA RLS policies: auth schema not present (local/dev PostgreSQL)';
        RETURN;
    END IF;

    ALTER TABLE roadmaps ENABLE ROW LEVEL SECURITY;
    ALTER TABLE topics ENABLE ROW LEVEL SECURITY;
    ALTER TABLE subtopics ENABLE ROW LEVEL SECURITY;
    ALTER TABLE problems ENABLE ROW LEVEL SECURITY;
    ALTER TABLE student_problem_progress ENABLE ROW LEVEL SECURITY;
    ALTER TABLE student_learning_events ENABLE ROW LEVEL SECURITY;

    -- Catalogue is readable by signed-in students. Mutations stay with service_role (bypasses RLS).
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'roadmaps' AND policyname = 'roadmaps_select_authenticated') THEN
        CREATE POLICY roadmaps_select_authenticated ON roadmaps
            FOR SELECT TO authenticated
            USING (is_active = TRUE);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'topics' AND policyname = 'topics_select_authenticated') THEN
        CREATE POLICY topics_select_authenticated ON topics
            FOR SELECT TO authenticated
            USING (
                EXISTS (SELECT 1 FROM roadmaps r WHERE r.id = topics.roadmap_id AND r.is_active = TRUE)
            );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'subtopics' AND policyname = 'subtopics_select_authenticated') THEN
        CREATE POLICY subtopics_select_authenticated ON subtopics
            FOR SELECT TO authenticated
            USING (
                EXISTS (
                    SELECT 1 FROM topics t
                    JOIN roadmaps r ON r.id = t.roadmap_id
                    WHERE t.id = subtopics.topic_id AND r.is_active = TRUE
                )
            );
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'problems' AND policyname = 'problems_select_authenticated') THEN
        CREATE POLICY problems_select_authenticated ON problems
            FOR SELECT TO authenticated
            USING (is_active = TRUE);
    END IF;

    -- Students may only read/write their own progress rows.
    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_select_own') THEN
        CREATE POLICY progress_select_own ON student_problem_progress
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_insert_own') THEN
        CREATE POLICY progress_insert_own ON student_problem_progress
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_update_own') THEN
        CREATE POLICY progress_update_own ON student_problem_progress
            FOR UPDATE TO authenticated
            USING (student_id = auth.uid())
            WITH CHECK (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_problem_progress' AND policyname = 'progress_delete_own') THEN
        CREATE POLICY progress_delete_own ON student_problem_progress
            FOR DELETE TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_events' AND policyname = 'learning_events_select_own') THEN
        CREATE POLICY learning_events_select_own ON student_learning_events
            FOR SELECT TO authenticated
            USING (student_id = auth.uid());
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_policies WHERE tablename = 'student_learning_events' AND policyname = 'learning_events_insert_own') THEN
        CREATE POLICY learning_events_insert_own ON student_learning_events
            FOR INSERT TO authenticated
            WITH CHECK (student_id = auth.uid());
    END IF;
END $$;
