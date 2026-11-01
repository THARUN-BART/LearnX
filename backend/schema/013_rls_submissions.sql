-- 013_rls_submissions.sql
-- LearnX Submission & Test Case Row Level Security Policies

ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE test_cases ENABLE ROW LEVEL SECURITY;
ALTER TABLE submission_test_results ENABLE ROW LEVEL SECURITY;

-- Submissions policies
DROP POLICY IF EXISTS "students_select_own_submissions" ON submissions;
CREATE POLICY "students_select_own_submissions" ON submissions
    FOR SELECT TO authenticated
    USING (student_id = auth.uid());

DROP POLICY IF EXISTS "students_insert_own_submissions" ON submissions;
CREATE POLICY "students_insert_own_submissions" ON submissions
    FOR INSERT TO authenticated
    WITH CHECK (student_id = auth.uid());

-- Test cases policies: public/students can ONLY read sample (non-hidden) test cases
DROP POLICY IF EXISTS "authenticated_select_sample_test_cases" ON test_cases;
CREATE POLICY "authenticated_select_sample_test_cases" ON test_cases
    FOR SELECT TO authenticated
    USING (is_hidden = false);

-- Submission test results: students can ONLY read results of their own submissions
DROP POLICY IF EXISTS "students_select_own_test_results" ON submission_test_results;
CREATE POLICY "students_select_own_test_results" ON submission_test_results
    FOR SELECT TO authenticated
    USING (
        submission_id IN (
            SELECT id FROM submissions WHERE student_id = auth.uid()
        )
    );
