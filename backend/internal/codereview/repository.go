package codereview

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrReviewNotFound = errors.New("code review not found")

type Repository interface {
	SaveReview(ctx context.Context, review *CodeReviewResult) error
	GetReview(ctx context.Context, reviewID, studentID string) (*CodeReviewResult, error)
	ListStudentReviews(ctx context.Context, studentID string, limit, offset int) ([]ReviewHistoryItem, error)
}

// ----------------- Postgres Repository -----------------

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) SaveReview(ctx context.Context, rev *CodeReviewResult) error {
	if rev.ID == "" {
		rev.ID = uuid.NewString()
	}
	rev.CreatedAt = time.Now()

	dataJSON, err := json.Marshal(rev)
	if err != nil {
		dataJSON = []byte("{}")
	}

	var probID *string
	if rev.ProblemID != "" {
		probID = &rev.ProblemID
	}
	var subID *string
	if rev.SubmissionID != "" {
		subID = &rev.SubmissionID
	}

	query := `
		INSERT INTO code_reviews (
			id, student_id, problem_id, submission_id, language, source_code,
			overall_score, readability_score, efficiency_score, maintainability_score,
			review_data, optimized_code, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err = r.pool.Exec(ctx, query,
		rev.ID, rev.StudentID, probID, subID, rev.Language, rev.SourceCode,
		rev.Scores.Overall, rev.Scores.Readability, rev.Scores.Efficiency, rev.Scores.Maintainability,
		dataJSON, rev.OptimizedCode, rev.CreatedAt,
	)
	return err
}

func (r *postgresRepository) GetReview(ctx context.Context, reviewID, studentID string) (*CodeReviewResult, error) {
	query := `
		SELECT review_data
		FROM code_reviews
		WHERE id = $1 AND student_id = $2
	`
	row := r.pool.QueryRow(ctx, query, reviewID, studentID)

	var dataBytes []byte
	if err := row.Scan(&dataBytes); err != nil {
		return nil, ErrReviewNotFound
	}

	var rev CodeReviewResult
	if err := json.Unmarshal(dataBytes, &rev); err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *postgresRepository) ListStudentReviews(ctx context.Context, studentID string, limit, offset int) ([]ReviewHistoryItem, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT cr.id, COALESCE(cr.problem_id::text, ''), COALESCE(p.title, 'General Code Review'),
		       cr.language, cr.overall_score, cr.created_at
		FROM code_reviews cr
		LEFT JOIN problems p ON cr.problem_id = p.id
		WHERE cr.student_id = $1
		ORDER BY cr.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ReviewHistoryItem
	for rows.Next() {
		var item ReviewHistoryItem
		if err := rows.Scan(&item.ID, &item.ProblemID, &item.ProblemTitle, &item.Language, &item.OverallScore, &item.CreatedAt); err == nil {
			list = append(list, item)
		}
	}
	return list, nil
}

// ----------------- Memory Repository (fallback) -----------------

type memoryRepository struct {
	mu      sync.RWMutex
	reviews map[string]CodeReviewResult
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		reviews: make(map[string]CodeReviewResult),
	}
}

func (m *memoryRepository) SaveReview(_ context.Context, rev *CodeReviewResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if rev.ID == "" {
		rev.ID = uuid.NewString()
	}
	rev.CreatedAt = time.Now()
	m.reviews[rev.ID] = *rev
	return nil
}

func (m *memoryRepository) GetReview(_ context.Context, reviewID, studentID string) (*CodeReviewResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rev, ok := m.reviews[reviewID]
	if !ok || rev.StudentID != studentID {
		return nil, ErrReviewNotFound
	}
	return &rev, nil
}

func (m *memoryRepository) ListStudentReviews(_ context.Context, studentID string, limit, offset int) ([]ReviewHistoryItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var items []ReviewHistoryItem
	for _, r := range m.reviews {
		if r.StudentID == studentID {
			items = append(items, ReviewHistoryItem{
				ID:           r.ID,
				ProblemID:    r.ProblemID,
				ProblemTitle: r.ProblemTitle,
				Language:     r.Language,
				OverallScore: r.Scores.Overall,
				CreatedAt:    r.CreatedAt,
			})
		}
	}

	if offset >= len(items) {
		return []ReviewHistoryItem{}, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], nil
}

