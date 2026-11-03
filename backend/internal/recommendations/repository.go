package recommendations

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("recommendation not found")
)

// Repository manages persistence for personalised problem recommendations.
type Repository interface {
	SaveRecommendations(ctx context.Context, studentID string, recs []ProblemRecommendation) error
	GetActiveRecommendations(ctx context.Context, studentID string) ([]ProblemRecommendation, error)
	GetTopPick(ctx context.Context, studentID string) (*ProblemRecommendation, error)
	GetByCategory(ctx context.Context, studentID, category string) ([]ProblemRecommendation, error)
	GetByID(ctx context.Context, studentID, recommendationID string) (*ProblemRecommendation, error)
	SaveFeedback(ctx context.Context, studentID, recommendationID, rating, notes string) error
	MarkCompleted(ctx context.Context, studentID, problemID string) error
	GetHistory(ctx context.Context, studentID string, limit, offset int) ([]ProblemRecommendation, int, error)
	ClearActive(ctx context.Context, studentID string) error
}

// -------------------------------------------------------------------------
// In-Memory Repository
// -------------------------------------------------------------------------

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]*ProblemRecommendation // key: recommendation ID
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make(map[string]*ProblemRecommendation),
	}
}

func (r *MemoryRepository) SaveRecommendations(ctx context.Context, studentID string, recs []ProblemRecommendation) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	// Mark older active recommendations for this student as expired
	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive {
			item.Status = StatusExpired
			item.UpdatedAt = now
		}
	}

	for i := range recs {
		rec := recs[i]
		if rec.ID == "" {
			rec.ID = uuid.NewString()
		}
		rec.StudentID = studentID
		rec.Status = StatusActive
		rec.CreatedAt = now
		rec.UpdatedAt = now
		r.items[rec.ID] = &rec
	}
	return nil
}

func (r *MemoryRepository) GetActiveRecommendations(ctx context.Context, studentID string) ([]ProblemRecommendation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []ProblemRecommendation
	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive {
			result = append(result, *item)
		}
	}
	return result, nil
}

func (r *MemoryRepository) GetTopPick(ctx context.Context, studentID string) (*ProblemRecommendation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive && item.IsTopPick {
			return item, nil
		}
	}
	// Fallback to highest match score
	var best *ProblemRecommendation
	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive {
			if best == nil || item.MatchScore > best.MatchScore {
				best = item
			}
		}
	}
	if best == nil {
		return nil, ErrNotFound
	}
	return best, nil
}

func (r *MemoryRepository) GetByCategory(ctx context.Context, studentID, category string) ([]ProblemRecommendation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []ProblemRecommendation
	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive && item.Category == category {
			result = append(result, *item)
		}
	}
	return result, nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, studentID, recommendationID string) (*ProblemRecommendation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[recommendationID]
	if !ok || item.StudentID != studentID {
		return nil, ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) SaveFeedback(ctx context.Context, studentID, recommendationID, rating, notes string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[recommendationID]
	if !ok || item.StudentID != studentID {
		return ErrNotFound
	}

	item.FeedbackRating = &rating
	if notes != "" {
		item.FeedbackNotes = &notes
	}
	item.UpdatedAt = time.Now()
	return nil
}

func (r *MemoryRepository) MarkCompleted(ctx context.Context, studentID, problemID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, item := range r.items {
		if item.StudentID == studentID && item.ProblemID == problemID && item.Status == StatusActive {
			item.Status = StatusCompleted
			item.CompletedAt = &now
			item.UpdatedAt = now
		}
	}
	return nil
}

func (r *MemoryRepository) GetHistory(ctx context.Context, studentID string, limit, offset int) ([]ProblemRecommendation, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []ProblemRecommendation
	for _, item := range r.items {
		if item.StudentID == studentID {
			all = append(all, *item)
		}
	}

	total := len(all)
	if offset >= total {
		return []ProblemRecommendation{}, total, nil
	}
	end := offset + limit
	if limit <= 0 || end > total {
		end = total
	}
	return all[offset:end], total, nil
}

func (r *MemoryRepository) ClearActive(ctx context.Context, studentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, item := range r.items {
		if item.StudentID == studentID && item.Status == StatusActive {
			item.Status = StatusExpired
			item.UpdatedAt = now
		}
	}
	return nil
}

// -------------------------------------------------------------------------
// PostgreSQL Repository
// -------------------------------------------------------------------------

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) SaveRecommendations(ctx context.Context, studentID string, recs []ProblemRecommendation) error {
	if len(recs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Expire existing active recommendations
	_, err = tx.Exec(ctx, `
		UPDATE student_problem_recommendations
		SET status = 'expired', updated_at = NOW()
		WHERE student_id = $1 AND status = 'active'
	`, studentID)
	if err != nil {
		return fmt.Errorf("failed to expire previous recommendations: %w", err)
	}

	// Insert fresh recommendations batch
	for i := range recs {
		rec := recs[i]
		if rec.ID == "" {
			rec.ID = uuid.NewString()
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO student_problem_recommendations (
				id, student_id, problem_id, category, recommended_difficulty,
				primary_concept, related_concepts, match_score, headline,
				explanation, learning_objective, status, is_top_pick,
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'active', $12, NOW(), NOW()
			)
		`,
			rec.ID, studentID, rec.ProblemID, rec.Category, rec.RecommendedDifficulty,
			rec.PrimaryConcept, rec.RelatedConcepts, rec.MatchScore, rec.Headline,
			rec.Explanation, rec.LearningObjective, rec.IsTopPick,
		)
		if err != nil {
			return fmt.Errorf("failed inserting recommendation: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetActiveRecommendations(ctx context.Context, studentID string) ([]ProblemRecommendation, error) {
	query := `
		SELECT
			r.id, r.student_id, r.problem_id, r.category, r.recommended_difficulty,
			r.primary_concept, r.related_concepts, r.match_score, r.headline,
			r.explanation, r.learning_objective, r.status, r.is_top_pick,
			r.feedback_rating, r.feedback_notes, r.completed_at, r.created_at, r.updated_at,
			p.title, p.slug, p.difficulty, p.summary, p.external_problem_url,
			COALESCE(t.name, ''), COALESCE(st.name, ''), COALESCE(p.tags, '{}'),
			COALESCE(sp.status, 'Not Started')
		FROM student_problem_recommendations r
		JOIN problems p ON p.id = r.problem_id
		LEFT JOIN topics t ON t.id = p.topic_id
		LEFT JOIN subtopics st ON st.id = p.subtopic_id
		LEFT JOIN student_problem_progress sp ON sp.student_id = r.student_id AND sp.problem_id = p.id
		WHERE r.student_id = $1 AND r.status = 'active'
		ORDER BY r.is_top_pick DESC, r.match_score DESC
	`
	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed querying recommendations: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

func (r *PostgresRepository) GetTopPick(ctx context.Context, studentID string) (*ProblemRecommendation, error) {
	query := `
		SELECT
			r.id, r.student_id, r.problem_id, r.category, r.recommended_difficulty,
			r.primary_concept, r.related_concepts, r.match_score, r.headline,
			r.explanation, r.learning_objective, r.status, r.is_top_pick,
			r.feedback_rating, r.feedback_notes, r.completed_at, r.created_at, r.updated_at,
			p.title, p.slug, p.difficulty, p.summary, p.external_problem_url,
			COALESCE(t.name, ''), COALESCE(st.name, ''), COALESCE(p.tags, '{}'),
			COALESCE(sp.status, 'Not Started')
		FROM student_problem_recommendations r
		JOIN problems p ON p.id = r.problem_id
		LEFT JOIN topics t ON t.id = p.topic_id
		LEFT JOIN subtopics st ON st.id = p.subtopic_id
		LEFT JOIN student_problem_progress sp ON sp.student_id = r.student_id AND sp.problem_id = p.id
		WHERE r.student_id = $1 AND r.status = 'active'
		ORDER BY r.is_top_pick DESC, r.match_score DESC
		LIMIT 1
	`
	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed querying top pick: %w", err)
	}
	defer rows.Close()

	recs, err := r.scanRecommendations(rows)
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, ErrNotFound
	}
	return &recs[0], nil
}

func (r *PostgresRepository) GetByCategory(ctx context.Context, studentID, category string) ([]ProblemRecommendation, error) {
	query := `
		SELECT
			r.id, r.student_id, r.problem_id, r.category, r.recommended_difficulty,
			r.primary_concept, r.related_concepts, r.match_score, r.headline,
			r.explanation, r.learning_objective, r.status, r.is_top_pick,
			r.feedback_rating, r.feedback_notes, r.completed_at, r.created_at, r.updated_at,
			p.title, p.slug, p.difficulty, p.summary, p.external_problem_url,
			COALESCE(t.name, ''), COALESCE(st.name, ''), COALESCE(p.tags, '{}'),
			COALESCE(sp.status, 'Not Started')
		FROM student_problem_recommendations r
		JOIN problems p ON p.id = r.problem_id
		LEFT JOIN topics t ON t.id = p.topic_id
		LEFT JOIN subtopics st ON st.id = p.subtopic_id
		LEFT JOIN student_problem_progress sp ON sp.student_id = r.student_id AND sp.problem_id = p.id
		WHERE r.student_id = $1 AND r.category = $2 AND r.status = 'active'
		ORDER BY r.match_score DESC
	`
	rows, err := r.pool.Query(ctx, query, studentID, category)
	if err != nil {
		return nil, fmt.Errorf("failed querying category recommendations: %w", err)
	}
	defer rows.Close()

	return r.scanRecommendations(rows)
}

func (r *PostgresRepository) GetByID(ctx context.Context, studentID, recommendationID string) (*ProblemRecommendation, error) {
	query := `
		SELECT
			r.id, r.student_id, r.problem_id, r.category, r.recommended_difficulty,
			r.primary_concept, r.related_concepts, r.match_score, r.headline,
			r.explanation, r.learning_objective, r.status, r.is_top_pick,
			r.feedback_rating, r.feedback_notes, r.completed_at, r.created_at, r.updated_at,
			p.title, p.slug, p.difficulty, p.summary, p.external_problem_url,
			COALESCE(t.name, ''), COALESCE(st.name, ''), COALESCE(p.tags, '{}'),
			COALESCE(sp.status, 'Not Started')
		FROM student_problem_recommendations r
		JOIN problems p ON p.id = r.problem_id
		LEFT JOIN topics t ON t.id = p.topic_id
		LEFT JOIN subtopics st ON st.id = p.subtopic_id
		LEFT JOIN student_problem_progress sp ON sp.student_id = r.student_id AND sp.problem_id = p.id
		WHERE r.student_id = $1 AND r.id = $2
		LIMIT 1
	`
	rows, err := r.pool.Query(ctx, query, studentID, recommendationID)
	if err != nil {
		return nil, fmt.Errorf("failed querying recommendation by id: %w", err)
	}
	defer rows.Close()

	recs, err := r.scanRecommendations(rows)
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, ErrNotFound
	}
	return &recs[0], nil
}

func (r *PostgresRepository) SaveFeedback(ctx context.Context, studentID, recommendationID, rating, notes string) error {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE student_problem_recommendations
		SET feedback_rating = $3, feedback_notes = $4, updated_at = NOW()
		WHERE student_id = $1 AND id = $2
	`, studentID, recommendationID, rating, notes)
	if err != nil {
		return fmt.Errorf("failed updating feedback: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepository) MarkCompleted(ctx context.Context, studentID, problemID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE student_problem_recommendations
		SET status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE student_id = $1 AND problem_id = $2 AND status = 'active'
	`, studentID, problemID)
	if err != nil {
		return fmt.Errorf("failed completing recommendation: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetHistory(ctx context.Context, studentID string, limit, offset int) ([]ProblemRecommendation, int, error) {
	if limit <= 0 {
		limit = 20
	}
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM student_problem_recommendations WHERE student_id = $1
	`, studentID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting history: %w", err)
	}

	query := `
		SELECT
			r.id, r.student_id, r.problem_id, r.category, r.recommended_difficulty,
			r.primary_concept, r.related_concepts, r.match_score, r.headline,
			r.explanation, r.learning_objective, r.status, r.is_top_pick,
			r.feedback_rating, r.feedback_notes, r.completed_at, r.created_at, r.updated_at,
			p.title, p.slug, p.difficulty, p.summary, p.external_problem_url,
			COALESCE(t.name, ''), COALESCE(st.name, ''), COALESCE(p.tags, '{}'),
			COALESCE(sp.status, 'Not Started')
		FROM student_problem_recommendations r
		JOIN problems p ON p.id = r.problem_id
		LEFT JOIN topics t ON t.id = p.topic_id
		LEFT JOIN subtopics st ON st.id = p.subtopic_id
		LEFT JOIN student_problem_progress sp ON sp.student_id = r.student_id AND sp.problem_id = p.id
		WHERE r.student_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed querying history: %w", err)
	}
	defer rows.Close()

	recs, err := r.scanRecommendations(rows)
	if err != nil {
		return nil, 0, err
	}
	return recs, total, nil
}

func (r *PostgresRepository) ClearActive(ctx context.Context, studentID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE student_problem_recommendations
		SET status = 'expired', updated_at = NOW()
		WHERE student_id = $1 AND status = 'active'
	`, studentID)
	return err
}

func (r *PostgresRepository) scanRecommendations(rows pgx.Rows) ([]ProblemRecommendation, error) {
	var result []ProblemRecommendation
	for rows.Next() {
		var rec ProblemRecommendation
		err := rows.Scan(
			&rec.ID, &rec.StudentID, &rec.ProblemID, &rec.Category, &rec.RecommendedDifficulty,
			&rec.PrimaryConcept, &rec.RelatedConcepts, &rec.MatchScore, &rec.Headline,
			&rec.Explanation, &rec.LearningObjective, &rec.Status, &rec.IsTopPick,
			&rec.FeedbackRating, &rec.FeedbackNotes, &rec.CompletedAt, &rec.CreatedAt, &rec.UpdatedAt,
			&rec.ProblemTitle, &rec.ProblemSlug, &rec.ProblemDifficulty, &rec.ProblemSummary, &rec.ExternalProblemURL,
			&rec.TopicName, &rec.SubtopicName, &rec.Tags,
			&rec.AttemptStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("failed scanning recommendation: %w", err)
		}
		result = append(result, rec)
	}
	return result, nil
}

