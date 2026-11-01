package submissions

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSubmissionNotFound = errors.New("submission not found")
	ErrUnauthorized       = errors.New("unauthorized: submission belongs to another student")
)

type Repository interface {
	Create(ctx context.Context, sub *Submission) error
	Update(ctx context.Context, sub *Submission) error
	SaveTestResults(ctx context.Context, results []SubmissionTestResult) error
	GetByID(ctx context.Context, id string) (*Submission, error)
	GetTestResultsBySubmissionID(ctx context.Context, submissionID string) ([]SubmissionTestResult, error)
	ListByStudent(ctx context.Context, studentID, problemID, status string, limit, offset int) ([]Submission, int, error)
	GetLatestByStudentAndProblem(ctx context.Context, studentID, problemID string) (*Submission, error)
	GetPendingSubmissions(ctx context.Context, limit int) ([]Submission, error)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const submissionColumns = `
	id, student_id, problem_id, language, source_code, status,
	execution_time_ms, memory_used_kb, total_test_cases, passed_test_cases,
	failed_test_cases, error_message, execution_provider, provider_submission_id,
	submitted_at, completed_at
`

func scanSubmission(row pgx.Row) (*Submission, error) {
	var s Submission
	err := row.Scan(
		&s.ID, &s.StudentID, &s.ProblemID, &s.Language, &s.SourceCode, &s.Status,
		&s.ExecutionTimeMs, &s.MemoryUsedKB, &s.TotalTestCases, &s.PassedTestCases,
		&s.FailedTestCases, &s.ErrorMessage, &s.ExecutionProvider, &s.ProviderSubmissionID,
		&s.SubmittedAt, &s.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSubmissionNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *PostgresRepository) Create(ctx context.Context, sub *Submission) error {
	query := `
		INSERT INTO submissions (
			id, student_id, problem_id, language, source_code, status,
			execution_time_ms, memory_used_kb, total_test_cases, passed_test_cases,
			failed_test_cases, error_message, execution_provider, provider_submission_id,
			submitted_at
		) VALUES (
			COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()),
			$2,
			(SELECT id FROM problems WHERE id::text = $3 OR slug = $3 LIMIT 1),
			$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW()
		)
		RETURNING id, problem_id, submitted_at;
	`
	return r.pool.QueryRow(ctx, query,
		sub.ID, sub.StudentID, sub.ProblemID, sub.Language, sub.SourceCode, sub.Status,
		sub.ExecutionTimeMs, sub.MemoryUsedKB, sub.TotalTestCases, sub.PassedTestCases,
		sub.FailedTestCases, sub.ErrorMessage, sub.ExecutionProvider, sub.ProviderSubmissionID,
	).Scan(&sub.ID, &sub.ProblemID, &sub.SubmittedAt)
}

func (r *PostgresRepository) Update(ctx context.Context, sub *Submission) error {
	query := `
		UPDATE submissions SET
			status = $2,
			execution_time_ms = $3,
			memory_used_kb = $4,
			total_test_cases = $5,
			passed_test_cases = $6,
			failed_test_cases = $7,
			error_message = $8,
			provider_submission_id = $9,
			completed_at = $10
		WHERE id::text = $1;
	`
	tag, err := r.pool.Exec(ctx, query,
		sub.ID, sub.Status, sub.ExecutionTimeMs, sub.MemoryUsedKB, sub.TotalTestCases,
		sub.PassedTestCases, sub.FailedTestCases, sub.ErrorMessage, sub.ProviderSubmissionID, sub.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("update submission: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrSubmissionNotFound
	}
	return nil
}

func (r *PostgresRepository) SaveTestResults(ctx context.Context, results []SubmissionTestResult) error {
	if len(results) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, res := range results {
		batch.Queue(`
			INSERT INTO submission_test_results (
				id, submission_id, test_case_id, status, execution_time_ms,
				memory_used_kb, actual_output, error_message, created_at
			) VALUES (
				gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, NOW()
			);
		`, res.SubmissionID, res.TestCaseID, res.Status, res.ExecutionTimeMs, res.MemoryUsedKB, res.ActualOutput, res.ErrorMessage)
	}

	br := r.pool.SendBatch(ctx, batch)
	defer br.Close()
	for i := 0; i < len(results); i++ {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("save submission test result: %w", err)
		}
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*Submission, error) {
	query := fmt.Sprintf(`SELECT %s FROM submissions WHERE id::text = $1 LIMIT 1;`, submissionColumns)
	return scanSubmission(r.pool.QueryRow(ctx, query, id))
}

func (r *PostgresRepository) GetTestResultsBySubmissionID(ctx context.Context, submissionID string) ([]SubmissionTestResult, error) {
	query := `
		SELECT str.id, str.submission_id, str.test_case_id, str.status,
		       str.execution_time_ms, str.memory_used_kb, str.actual_output, str.error_message,
		       str.created_at, tc.is_hidden
		FROM submission_test_results str
		JOIN test_cases tc ON tc.id = str.test_case_id
		WHERE str.submission_id::text = $1
		ORDER BY tc.test_order ASC, str.created_at ASC;
	`
	rows, err := r.pool.Query(ctx, query, submissionID)
	if err != nil {
		return nil, fmt.Errorf("get test results: %w", err)
	}
	defer rows.Close()

	var list []SubmissionTestResult
	for rows.Next() {
		var res SubmissionTestResult
		if err := rows.Scan(
			&res.ID, &res.SubmissionID, &res.TestCaseID, &res.Status,
			&res.ExecutionTimeMs, &res.MemoryUsedKB, &res.ActualOutput, &res.ErrorMessage,
			&res.CreatedAt, &res.IsHidden,
		); err != nil {
			return nil, fmt.Errorf("scan test result: %w", err)
		}
		list = append(list, res)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, studentID, problemID, status string, limit, offset int) ([]Submission, int, error) {
	whereClauses := []string{"student_id = $1"}
	args := []interface{}{studentID}
	argIdx := 2

	if problemID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(problem_id::text = $%d OR problem_id = (SELECT id FROM problems WHERE slug = $%d LIMIT 1))", argIdx, argIdx))
		args = append(args, problemID)
		argIdx++
	}
	if status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, status)
		argIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM submissions WHERE %s;", whereSQL)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count submissions: %w", err)
	}

	query := fmt.Sprintf(`
		SELECT %s FROM submissions
		WHERE %s
		ORDER BY submitted_at DESC
		LIMIT $%d OFFSET $%d;
	`, submissionColumns, whereSQL, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list submissions: %w", err)
	}
	defer rows.Close()

	list := make([]Submission, 0)
	for rows.Next() {
		sub, err := scanSubmission(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan list submission: %w", err)
		}
		list = append(list, *sub)
	}
	return list, total, rows.Err()
}

func (r *PostgresRepository) GetLatestByStudentAndProblem(ctx context.Context, studentID, problemID string) (*Submission, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM submissions
		WHERE student_id = $1 AND (problem_id::text = $2 OR problem_id = (SELECT id FROM problems WHERE slug = $2 LIMIT 1))
		ORDER BY submitted_at DESC
		LIMIT 1;
	`, submissionColumns)
	return scanSubmission(r.pool.QueryRow(ctx, query, studentID, problemID))
}

func (r *PostgresRepository) GetPendingSubmissions(ctx context.Context, limit int) ([]Submission, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM submissions
		WHERE status IN ('Queued', 'Processing')
		ORDER BY submitted_at ASC
		LIMIT $1;
	`, submissionColumns)
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get pending submissions: %w", err)
	}
	defer rows.Close()

	list := make([]Submission, 0)
	for rows.Next() {
		sub, err := scanSubmission(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *sub)
	}
	return list, rows.Err()
}

// MemoryRepository provides thread-safe in-memory storage for unit & integration testing.
type MemoryRepository struct {
	mu          sync.RWMutex
	submissions map[string]*Submission
	results     map[string][]SubmissionTestResult
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		submissions: make(map[string]*Submission),
		results:     make(map[string][]SubmissionTestResult),
	}
}

func (m *MemoryRepository) Create(ctx context.Context, sub *Submission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sub.ID == "" {
		sub.ID = fmt.Sprintf("sub-%d", len(m.submissions)+1)
	}
	sub.SubmittedAt = time.Now()
	m.submissions[sub.ID] = sub
	return nil
}

func (m *MemoryRepository) Update(ctx context.Context, sub *Submission) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.submissions[sub.ID]; !ok {
		return ErrSubmissionNotFound
	}
	m.submissions[sub.ID] = sub
	return nil
}

func (m *MemoryRepository) SaveTestResults(ctx context.Context, results []SubmissionTestResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, res := range results {
		m.results[res.SubmissionID] = append(m.results[res.SubmissionID], res)
	}
	return nil
}

func (m *MemoryRepository) GetByID(ctx context.Context, id string) (*Submission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sub, ok := m.submissions[id]
	if !ok {
		return nil, ErrSubmissionNotFound
	}
	// Return copy
	cpy := *sub
	return &cpy, nil
}

func (m *MemoryRepository) GetTestResultsBySubmissionID(ctx context.Context, submissionID string) ([]SubmissionTestResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.results[submissionID], nil
}

func (m *MemoryRepository) ListByStudent(ctx context.Context, studentID, problemID, status string, limit, offset int) ([]Submission, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var filtered []Submission
	for _, s := range m.submissions {
		if s.StudentID == studentID {
			if problemID != "" && s.ProblemID != problemID {
				continue
			}
			if status != "" && s.Status != status {
				continue
			}
			filtered = append(filtered, *s)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].SubmittedAt.After(filtered[j].SubmittedAt)
	})
	total := len(filtered)
	if offset >= total {
		return []Submission{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (m *MemoryRepository) GetLatestByStudentAndProblem(ctx context.Context, studentID, problemID string) (*Submission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var latest *Submission
	for _, s := range m.submissions {
		if s.StudentID == studentID && s.ProblemID == problemID {
			if latest == nil || s.SubmittedAt.After(latest.SubmittedAt) {
				latest = s
			}
		}
	}
	if latest == nil {
		return nil, ErrSubmissionNotFound
	}
	cpy := *latest
	return &cpy, nil
}

func (m *MemoryRepository) GetPendingSubmissions(ctx context.Context, limit int) ([]Submission, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var pending []Submission
	for _, s := range m.submissions {
		if s.Status == "Queued" || s.Status == "Processing" {
			pending = append(pending, *s)
		}
	}
	return pending, nil
}
