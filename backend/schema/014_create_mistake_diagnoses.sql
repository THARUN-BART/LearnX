-- 014_create_mistake_diagnoses.sql
-- Records intelligent diagnoses of student coding mistakes for failed submissions

CREATE TABLE IF NOT EXISTS mistake_diagnoses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL,
    submission_id UUID NOT NULL,
    problem_id UUID NOT NULL,
    category VARCHAR(64) NOT NULL,
    subtype VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    related_concept VARCHAR(128) NOT NULL,
    detection_method VARCHAR(64) NOT NULL,
    confidence_level VARCHAR(32) NOT NULL,
    confidence_score NUMERIC(5, 2) NOT NULL DEFAULT 0.0,
    evidence JSONB NOT NULL DEFAULT '{}'::jsonb,
    suggested_correction TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_diagnosis_submission UNIQUE (submission_id),
    CONSTRAINT fk_diag_submission FOREIGN KEY (submission_id) REFERENCES submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_diag_problem FOREIGN KEY (problem_id) REFERENCES problems(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_diagnoses_student ON mistake_diagnoses (student_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_problem ON mistake_diagnoses (problem_id);
CREATE INDEX IF NOT EXISTS idx_diagnoses_category ON mistake_diagnoses (category, subtype);
CREATE INDEX IF NOT EXISTS idx_diagnoses_created ON mistake_diagnoses (created_at DESC);
