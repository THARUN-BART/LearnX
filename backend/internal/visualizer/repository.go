package visualizer

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("visualizer session not found")

type Repository interface {
	CreateSession(ctx context.Context, session *VisualizerSession) error
	GetSession(ctx context.Context, sessionID, studentID string) (*VisualizerSession, error)
	ListSessions(ctx context.Context, studentID string, limit, offset int) ([]VisualizerSession, error)
	UpdateSession(ctx context.Context, sessionID, studentID string, stepsCompleted, timeSpent int, completed bool) error
}

// ----------------- Postgres Repository -----------------

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) CreateSession(ctx context.Context, s *VisualizerSession) error {
	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	cfgJSON, err := json.Marshal(s.InputConfig)
	if err != nil {
		cfgJSON = []byte("{}")
	}

	query := `
		INSERT INTO visualizer_sessions (
			id, student_id, algorithm_id, category, input_config,
			steps_completed, total_steps, time_spent_seconds, completed, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = r.pool.Exec(ctx, query,
		s.ID, s.StudentID, s.AlgorithmID, string(s.Category), cfgJSON,
		s.StepsCompleted, s.TotalSteps, s.TimeSpentSeconds, s.Completed, s.CreatedAt, s.UpdatedAt,
	)
	return err
}

func (r *postgresRepository) GetSession(ctx context.Context, sessionID, studentID string) (*VisualizerSession, error) {
	query := `
		SELECT id, student_id, algorithm_id, category, input_config,
		       steps_completed, total_steps, time_spent_seconds, completed, created_at, updated_at
		FROM visualizer_sessions
		WHERE id = $1 AND student_id = $2
	`
	row := r.pool.QueryRow(ctx, query, sessionID, studentID)

	var s VisualizerSession
	var cfgBytes []byte
	var cat string

	err := row.Scan(
		&s.ID, &s.StudentID, &s.AlgorithmID, &cat, &cfgBytes,
		&s.StepsCompleted, &s.TotalSteps, &s.TimeSpentSeconds, &s.Completed, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, ErrSessionNotFound
	}
	s.Category = AlgorithmCategory(cat)
	_ = json.Unmarshal(cfgBytes, &s.InputConfig)

	return &s, nil
}

func (r *postgresRepository) ListSessions(ctx context.Context, studentID string, limit, offset int) ([]VisualizerSession, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, student_id, algorithm_id, category, input_config,
		       steps_completed, total_steps, time_spent_seconds, completed, created_at, updated_at
		FROM visualizer_sessions
		WHERE student_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []VisualizerSession
	for rows.Next() {
		var s VisualizerSession
		var cfgBytes []byte
		var cat string

		if err := rows.Scan(
			&s.ID, &s.StudentID, &s.AlgorithmID, &cat, &cfgBytes,
			&s.StepsCompleted, &s.TotalSteps, &s.TimeSpentSeconds, &s.Completed, &s.CreatedAt, &s.UpdatedAt,
		); err == nil {
			s.Category = AlgorithmCategory(cat)
			_ = json.Unmarshal(cfgBytes, &s.InputConfig)
			sessions = append(sessions, s)
		}
	}
	return sessions, nil
}

func (r *postgresRepository) UpdateSession(ctx context.Context, sessionID, studentID string, stepsCompleted, timeSpent int, completed bool) error {
	query := `
		UPDATE visualizer_sessions
		SET steps_completed = $1, time_spent_seconds = time_spent_seconds + $2, completed = $3, updated_at = NOW()
		WHERE id = $4 AND student_id = $5
	`
	tag, err := r.pool.Exec(ctx, query, stepsCompleted, timeSpent, completed, sessionID, studentID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSessionNotFound
	}
	return nil
}

// ----------------- Memory Repository (fallback) -----------------

type memoryRepository struct {
	mu       sync.RWMutex
	sessions map[string]VisualizerSession
}

func NewMemoryRepository() Repository {
	return &memoryRepository{
		sessions: make(map[string]VisualizerSession),
	}
}

func (m *memoryRepository) CreateSession(_ context.Context, s *VisualizerSession) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = uuid.NewString()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	m.sessions[s.ID] = *s
	return nil
}

func (m *memoryRepository) GetSession(_ context.Context, sessionID, studentID string) (*VisualizerSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s, ok := m.sessions[sessionID]
	if !ok || s.StudentID != studentID {
		return nil, ErrSessionNotFound
	}
	return &s, nil
}

func (m *memoryRepository) ListSessions(_ context.Context, studentID string, limit, offset int) ([]VisualizerSession, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var res []VisualizerSession
	for _, s := range m.sessions {
		if s.StudentID == studentID {
			res = append(res, s)
		}
	}
	// Simple slice paging
	if offset >= len(res) {
		return []VisualizerSession{}, nil
	}
	end := offset + limit
	if end > len(res) {
		end = len(res)
	}
	return res[offset:end], nil
}

func (m *memoryRepository) UpdateSession(_ context.Context, sessionID, studentID string, stepsCompleted, timeSpent int, completed bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[sessionID]
	if !ok || s.StudentID != studentID {
		return ErrSessionNotFound
	}
	s.StepsCompleted = stepsCompleted
	s.TimeSpentSeconds += timeSpent
	if completed {
		s.Completed = true
	}
	s.UpdatedAt = time.Now()
	m.sessions[sessionID] = s
	return nil
}

