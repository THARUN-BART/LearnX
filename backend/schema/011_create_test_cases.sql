-- 011_create_test_cases.sql
-- LearnX Problem Test Cases

CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input_data TEXT NOT NULL DEFAULT '',
    expected_output TEXT NOT NULL DEFAULT '',
    is_hidden BOOLEAN NOT NULL DEFAULT false,
    test_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id ON test_cases(problem_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_order ON test_cases(problem_id, test_order ASC);
CREATE INDEX IF NOT EXISTS idx_test_cases_is_hidden ON test_cases(problem_id, is_hidden);

-- Seed initial test cases for common A2Z problems if they exist
DO $$
DECLARE
    v_prob_id UUID;
BEGIN
    -- Seed test cases for Two Sum if present
    SELECT id INTO v_prob_id FROM problems WHERE slug = 'two-sum' LIMIT 1;
    IF v_prob_id IS NOT NULL AND NOT EXISTS (SELECT 1 FROM test_cases WHERE problem_id = v_prob_id) THEN
        INSERT INTO test_cases (problem_id, input_data, expected_output, is_hidden, test_order) VALUES
        (v_prob_id, '4\n2 7 11 15\n9', '0 1', false, 1),
        (v_prob_id, '3\n3 2 4\n6', '1 2', false, 2),
        (v_prob_id, '2\n3 3\n6', '0 1', true, 3),
        (v_prob_id, '5\n-1 -2 -3 -4 -5\n-8', '2 4', true, 4);
    END IF;

    -- Seed test cases for Binary Search if present
    SELECT id INTO v_prob_id FROM problems WHERE slug = 'binary-search' LIMIT 1;
    IF v_prob_id IS NOT NULL AND NOT EXISTS (SELECT 1 FROM test_cases WHERE problem_id = v_prob_id) THEN
        INSERT INTO test_cases (problem_id, input_data, expected_output, is_hidden, test_order) VALUES
        (v_prob_id, '6\n-1 0 3 5 9 12\n9', '4', false, 1),
        (v_prob_id, '6\n-1 0 3 5 9 12\n2', '-1', false, 2),
        (v_prob_id, '1\n5\n5', '0', true, 3);
    END IF;
END $$;
