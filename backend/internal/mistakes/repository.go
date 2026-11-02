package mistakes

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDiagnosisNotFound = errors.New("diagnosis not found")
	ErrPatternNotFound   = errors.New("mistake pattern not found")
)

type Repository interface {
	SaveDiagnosis(ctx context.Context, diagnosis *Diagnosis) error
	GetDiagnosisBySubmissionID(ctx context.Context, submissionID string) (*Diagnosis, error)
	GetDiagnosisByID(ctx context.Context, id string) (*Diagnosis, error)
	ListDiagnosesByStudent(ctx context.Context, studentID string, limit, offset int) ([]Diagnosis, error)
	UpsertPattern(ctx context.Context, pattern *StudentMistakePattern) error
	GetActivePatterns(ctx context.Context, studentID string) ([]StudentMistakePattern, error)
	MarkPatternCorrected(ctx context.Context, studentID, category, subtype, relatedConcept string) error
	SaveFeedback(ctx context.Context, feedback *MistakeFeedback) error
	GetFeedbackByDiagnosis(ctx context.Context, diagnosisID string) ([]MistakeFeedback, error)
}

// PostgresRepository provides SQL implementation backed by Supabase pgxpool.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) SaveDiagnosis(ctx context.Context, diag *Diagnosis) error {
	if diag.ID == "" {
		diag.ID = crypto.GenerateUUID()
	}
	if diag.CreatedAt.IsZero() {
		diag.CreatedAt = time.Now()
	}
	evidenceJSON, err := json.Marshal(diag.Evidence)
	if err != nil {
		evidenceJSON = []byte("{}")
	}

	query := `
		INSERT INTO mistake_diagnoses (
			id, student_id, submission_id, problem_id, category, subtype,
			description, related_concept, detection_method, confidence_level,
			confidence_score, evidence, suggested_correction, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (submission_id) DO UPDATE SET
			category = EXCLUDED.category,
			subtype = EXCLUDED.subtype,
			description = EXCLUDED.description,
			related_concept = EXCLUDED.related_concept,
			detection_method = EXCLUDED.detection_method,
			confidence_level = EXCLUDED.confidence_level,
			confidence_score = EXCLUDED.confidence_score,
			evidence = EXCLUDED.evidence,
			suggested_correction = EXCLUDED.suggested_correction;
	`
	_, err = r.pool.Exec(ctx, query,
		diag.ID, diag.StudentID, diag.SubmissionID, diag.ProblemID, diag.Category, diag.Subtype,
		diag.Description, diag.RelatedConcept, diag.DetectionMethod, diag.ConfidenceLevel,
		diag.ConfidenceScore, evidenceJSON, diag.SuggestedCorrection, diag.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetDiagnosisBySubmissionID(ctx context.Context, submissionID string) (*Diagnosis, error) {
	query := `
		SELECT id, student_id, submission_id, problem_id, category, subtype,
		       description, related_concept, detection_method, confidence_level,
		       confidence_score, evidence, suggested_correction, created_at
		FROM mistake_diagnoses
		WHERE submission_id = $1
		LIMIT 1;
	`
	var diag Diagnosis
	var evidenceBytes []byte
	err := r.pool.QueryRow(ctx, query, submissionID).Scan(
		&diag.ID, &diag.StudentID, &diag.SubmissionID, &diag.ProblemID, &diag.Category, &diag.Subtype,
		&diag.Description, &diag.RelatedConcept, &diag.DetectionMethod, &diag.ConfidenceLevel,
		&diag.ConfidenceScore, &evidenceBytes, &diag.SuggestedCorrection, &diag.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDiagnosisNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(evidenceBytes, &diag.Evidence)
	return &diag, nil
}

func (r *PostgresRepository) GetDiagnosisByID(ctx context.Context, id string) (*Diagnosis, error) {
	query := `
		SELECT id, student_id, submission_id, problem_id, category, subtype,
		       description, related_concept, detection_method, confidence_level,
		       confidence_score, evidence, suggested_correction, created_at
		FROM mistake_diagnoses
		WHERE id = $1
		LIMIT 1;
	`
	var diag Diagnosis
	var evidenceBytes []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&diag.ID, &diag.StudentID, &diag.SubmissionID, &diag.ProblemID, &diag.Category, &diag.Subtype,
		&diag.Description, &diag.RelatedConcept, &diag.DetectionMethod, &diag.ConfidenceLevel,
		&diag.ConfidenceScore, &evidenceBytes, &diag.SuggestedCorrection, &diag.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDiagnosisNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(evidenceBytes, &diag.Evidence)
	return &diag, nil
}

func (r *PostgresRepository) ListDiagnosesByStudent(ctx context.Context, studentID string, limit, offset int) ([]Diagnosis, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, student_id, submission_id, problem_id, category, subtype,
		       description, related_concept, detection_method, confidence_level,
		       confidence_score, evidence, suggested_correction, created_at
		FROM mistake_diagnoses
		WHERE student_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Diagnosis
	for rows.Next() {
		var diag Diagnosis
		var evidenceBytes []byte
		if err := rows.Scan(
			&diag.ID, &diag.StudentID, &diag.SubmissionID, &diag.ProblemID, &diag.Category, &diag.Subtype,
			&diag.Description, &diag.RelatedConcept, &diag.DetectionMethod, &diag.ConfidenceLevel,
			&diag.ConfidenceScore, &evidenceBytes, &diag.SuggestedCorrection, &diag.CreatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(evidenceBytes, &diag.Evidence)
		result = append(result, diag)
	}
	return result, rows.Err()
}

func (r *PostgresRepository) UpsertPattern(ctx context.Context, p *StudentMistakePattern) error {
	now := time.Now()
	if p.ID == "" {
		p.ID = crypto.GenerateUUID()
	}
	p.UpdatedAt = now

	query := `
		INSERT INTO student_mistake_patterns (
			id, student_id, category, subtype, related_concept, occurrence_count,
			first_occurrence_at, last_occurrence_at, correction_status, updated_at
		) VALUES ($1, $2, $3, $4, $5, 1, $6, $6, 'active', $6)
		ON CONFLICT (student_id, category, subtype, related_concept) DO UPDATE SET
			occurrence_count = student_mistake_patterns.occurrence_count + 1,
			last_occurrence_at = EXCLUDED.last_occurrence_at,
			correction_status = 'active',
			updated_at = EXCLUDED.updated_at;
	`
	_, err := r.pool.Exec(ctx, query,
		p.ID, p.StudentID, p.Category, p.Subtype, p.RelatedConcept, now,
	)
	return err
}

func (r *PostgresRepository) GetActivePatterns(ctx context.Context, studentID string) ([]StudentMistakePattern, error) {
	query := `
		SELECT id, student_id, category, subtype, related_concept, occurrence_count,
		       first_occurrence_at, last_occurrence_at, correction_status, updated_at
		FROM student_mistake_patterns
		WHERE student_id = $1 AND correction_status = 'active'
		ORDER BY occurrence_count DESC, updated_at DESC;
	`
	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patterns []StudentMistakePattern
	for rows.Next() {
		var p StudentMistakePattern
		if err := rows.Scan(
			&p.ID, &p.StudentID, &p.Category, &p.Subtype, &p.RelatedConcept, &p.OccurrenceCount,
			&p.FirstOccurrenceAt, &p.LastOccurrenceAt, &p.CorrectionStatus, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
	}
	return patterns, rows.Err()
}

func (r *PostgresRepository) MarkPatternCorrected(ctx context.Context, studentID, category, subtype, relatedConcept string) error {
	query := `
		UPDATE student_mistake_patterns
		SET correction_status = 'corrected', updated_at = NOW()
		WHERE student_id = $1 AND (category = $2 OR $2 = '') AND (subtype = $3 OR $3 = '');
	`
	_, err := r.pool.Exec(ctx, query, studentID, category, subtype)
	return err
}

func (r *PostgresRepository) SaveFeedback(ctx context.Context, feedback *MistakeFeedback) error {
	if feedback.ID == "" {
		feedback.ID = crypto.GenerateUUID()
	}
	if feedback.CreatedAt.IsZero() {
		feedback.CreatedAt = time.Now()
	}
	query := `
		INSERT INTO mistake_feedback (id, diagnosis_id, student_id, feedback_type, student_feedback, created_at)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := r.pool.Exec(ctx, query,
		feedback.ID, feedback.DiagnosisID, feedback.StudentID, feedback.FeedbackType, feedback.StudentFeedback, feedback.CreatedAt,
	)
	return err
}

func (r *PostgresRepository) GetFeedbackByDiagnosis(ctx context.Context, diagnosisID string) ([]MistakeFeedback, error) {
	query := `
		SELECT id, diagnosis_id, student_id, feedback_type, student_feedback, created_at
		FROM mistake_feedback
		WHERE diagnosis_id = $1
		ORDER BY created_at DESC;
	`
	rows, err := r.pool.Query(ctx, query, diagnosisID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedbackList []MistakeFeedback
	for rows.Next() {
		var fb MistakeFeedback
		if err := rows.Scan(&fb.ID, &fb.DiagnosisID, &fb.StudentID, &fb.FeedbackType, &fb.StudentFeedback, &fb.CreatedAt); err != nil {
			return nil, err
		}
		feedbackList = append(feedbackList, fb)
	}
	return feedbackList, rows.Err()
}

// MemoryRepository is a thread-safe in-memory store for unit tests.
type MemoryRepository struct {
	mu           sync.RWMutex
	diagnoses    map[string]*Diagnosis
	bySubmission map[string]*Diagnosis
	patterns     map[string]*StudentMistakePattern
	feedback     map[string][]*MistakeFeedback
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		diagnoses:    make(map[string]*Diagnosis),
		bySubmission: make(map[string]*Diagnosis),
		patterns:     make(map[string]*StudentMistakePattern),
		feedback:     make(map[string][]*MistakeFeedback),
	}
}

func (m *MemoryRepository) SaveDiagnosis(ctx context.Context, diag *Diagnosis) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if diag.ID == "" {
		diag.ID = crypto.GenerateUUID()
	}
	if diag.CreatedAt.IsZero() {
		diag.CreatedAt = time.Now()
	}
	m.diagnoses[diag.ID] = diag
	m.bySubmission[diag.SubmissionID] = diag
	return nil
}

func (m *MemoryRepository) GetDiagnosisBySubmissionID(ctx context.Context, submissionID string) (*Diagnosis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	d, ok := m.bySubmission[submissionID]
	if !ok {
		return nil, ErrDiagnosisNotFound
	}
	return d, nil
}

func (m *MemoryRepository) GetDiagnosisByID(ctx context.Context, id string) (*Diagnosis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	d, ok := m.diagnoses[id]
	if !ok {
		return nil, ErrDiagnosisNotFound
	}
	return d, nil
}

func (m *MemoryRepository) ListDiagnosesByStudent(ctx context.Context, studentID string, limit, offset int) ([]Diagnosis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var res []Diagnosis
	for _, d := range m.diagnoses {
		if d.StudentID == studentID {
			res = append(res, *d)
		}
	}
	// Sort by newest first
	for i := 0; i < len(res)-1; i++ {
		for j := i + 1; j < len(res); j++ {
			if res[j].CreatedAt.After(res[i].CreatedAt) {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	if offset > len(res) {
		return []Diagnosis{}, nil
	}
	res = res[offset:]
	if limit > 0 && len(res) > limit {
		res = res[:limit]
	}
	return res, nil
}

func (m *MemoryRepository) UpsertPattern(ctx context.Context, p *StudentMistakePattern) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := p.StudentID + ":" + p.Category + ":" + p.Subtype + ":" + p.RelatedConcept
	existing, ok := m.patterns[key]
	now := time.Now()
	if ok {
		existing.OccurrenceCount++
		existing.LastOccurrenceAt = now
		existing.CorrectionStatus = "active"
		existing.UpdatedAt = now
	} else {
		if p.ID == "" {
			p.ID = crypto.GenerateUUID()
		}
		p.OccurrenceCount = 1
		p.FirstOccurrenceAt = now
		p.LastOccurrenceAt = now
		p.CorrectionStatus = "active"
		p.UpdatedAt = now
		m.patterns[key] = p
	}
	return nil
}

func (m *MemoryRepository) GetActivePatterns(ctx context.Context, studentID string) ([]StudentMistakePattern, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active []StudentMistakePattern
	for _, p := range m.patterns {
		if p.StudentID == studentID && p.CorrectionStatus == "active" {
			active = append(active, *p)
		}
	}
	return active, nil
}

func (m *MemoryRepository) MarkPatternCorrected(ctx context.Context, studentID, category, subtype, relatedConcept string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.patterns {
		if p.StudentID == studentID {
			if (category == "" || p.Category == category) && (subtype == "" || p.Subtype == subtype) {
				p.CorrectionStatus = "corrected"
				p.UpdatedAt = time.Now()
			}
		}
	}
	return nil
}

func (m *MemoryRepository) SaveFeedback(ctx context.Context, feedback *MistakeFeedback) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if feedback.ID == "" {
		feedback.ID = crypto.GenerateUUID()
	}
	if feedback.CreatedAt.IsZero() {
		feedback.CreatedAt = time.Now()
	}
	m.feedback[feedback.DiagnosisID] = append(m.feedback[feedback.DiagnosisID], feedback)
	return nil
}

func (m *MemoryRepository) GetFeedbackByDiagnosis(ctx context.Context, diagnosisID string) ([]MistakeFeedback, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fb, ok := m.feedback[diagnosisID]
	if !ok {
		return []MistakeFeedback{}, nil
	}
	res := make([]MistakeFeedback, len(fb))
	for i, f := range fb {
		res[i] = *f
	}
	return res, nil
}
