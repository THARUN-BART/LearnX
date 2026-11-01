package progress

import (
	"context"
	"errors"
	"fmt"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Get(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	ListByStudent(ctx context.Context, studentID string) ([]model.Progress, error)
	MarkAttempted(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	MarkSolved(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	UpdateStatus(ctx context.Context, studentID, problemID, status string, reset bool) (*model.Progress, error)
}

type MemoryRepository struct {
	DB *model.Memory
}

func NewMemoryRepository(db *model.Memory) *MemoryRepository {
	return &MemoryRepository{DB: db}
}

func (r *MemoryRepository) Get(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	item, ok := r.DB.GetProgress(studentID, problemID)
	if !ok {
		return &model.Progress{
			StudentID: studentID,
			ProblemID: problemID,
			Status:    model.StatusNotStarted,
		}, nil
	}
	return item, nil
}

func (r *MemoryRepository) ListByStudent(ctx context.Context, studentID string) ([]model.Progress, error) {
	return r.DB.ListProgress(studentID), nil
}

func (r *MemoryRepository) MarkAttempted(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	return r.DB.MarkAttempted(studentID, problemID)
}

func (r *MemoryRepository) MarkSolved(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	return r.DB.MarkSolved(studentID, problemID)
}

func (r *MemoryRepository) UpdateStatus(ctx context.Context, studentID, problemID, status string, reset bool) (*model.Progress, error) {
	return r.DB.UpdateProgress(studentID, problemID, status, reset)
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const progressSelect = `
	id, student_id, problem_id, status, total_attempts, successful_attempts,
	first_attempt_at, last_attempt_at, solved_at, created_at, updated_at
`

func scanProgress(row pgx.Row) (*model.Progress, error) {
	var item model.Progress
	if err := row.Scan(
		&item.ID, &item.StudentID, &item.ProblemID, &item.Status, &item.TotalAttempts, &item.SuccessfulAttempts,
		&item.FirstAttemptAt, &item.LastAttemptAt, &item.SolvedAt, &item.CreatedAt, &item.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PostgresRepository) Get(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	item, err := scanProgress(r.pool.QueryRow(ctx, `
		SELECT `+progressSelect+`
		FROM student_problem_progress
		WHERE student_id = $1 AND problem_id = (
			SELECT id FROM problems WHERE id::text = $2 OR slug = $2 LIMIT 1
		)
		LIMIT 1;
	`, studentID, problemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.Progress{
				StudentID: studentID,
				ProblemID: problemID,
				Status:    model.StatusNotStarted,
			}, nil
		}
		return nil, fmt.Errorf("get progress: %w", err)
	}
	return item, nil
}

func (r *PostgresRepository) ListByStudent(ctx context.Context, studentID string) ([]model.Progress, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT `+progressSelect+`
		FROM student_problem_progress
		WHERE student_id = $1
		ORDER BY updated_at DESC;
	`, studentID)
	if err != nil {
		return nil, fmt.Errorf("list progress: %w", err)
	}
	defer rows.Close()
	out := make([]model.Progress, 0)
	for rows.Next() {
		var item model.Progress
		if err := rows.Scan(
			&item.ID, &item.StudentID, &item.ProblemID, &item.Status, &item.TotalAttempts, &item.SuccessfulAttempts,
			&item.FirstAttemptAt, &item.LastAttemptAt, &item.SolvedAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) MarkAttempted(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	item, err := scanProgress(r.pool.QueryRow(ctx, `
		INSERT INTO student_problem_progress (
			student_id, problem_id, status, total_attempts, successful_attempts,
			first_attempt_at, last_attempt_at, updated_at
		)
		SELECT $1, p.id, 'Attempted', 1, 0, NOW(), NOW(), NOW()
		FROM problems p
		WHERE p.id::text = $2 OR p.slug = $2
		ON CONFLICT (student_id, problem_id) DO UPDATE SET
			total_attempts = student_problem_progress.total_attempts + 1,
			last_attempt_at = NOW(),
			first_attempt_at = COALESCE(student_problem_progress.first_attempt_at, NOW()),
			status = CASE
				WHEN student_problem_progress.status IN ('Solved', 'Needs Revision') THEN student_problem_progress.status
				ELSE 'Attempted'
			END,
			updated_at = NOW()
		RETURNING `+progressSelect+`;
	`, studentID, problemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("mark attempted: %w", err)
	}
	return item, nil
}

func (r *PostgresRepository) MarkSolved(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	item, err := scanProgress(r.pool.QueryRow(ctx, `
		INSERT INTO student_problem_progress (
			student_id, problem_id, status, total_attempts, successful_attempts,
			first_attempt_at, last_attempt_at, solved_at, updated_at
		)
		SELECT $1, p.id, 'Solved', 1, 1, NOW(), NOW(), NOW(), NOW()
		FROM problems p
		WHERE p.id::text = $2 OR p.slug = $2
		ON CONFLICT (student_id, problem_id) DO UPDATE SET
			status = 'Solved',
			total_attempts = student_problem_progress.total_attempts + 1,
			successful_attempts = student_problem_progress.successful_attempts + 1,
			last_attempt_at = NOW(),
			first_attempt_at = COALESCE(student_problem_progress.first_attempt_at, NOW()),
			solved_at = COALESCE(student_problem_progress.solved_at, NOW()),
			updated_at = NOW()
		RETURNING `+progressSelect+`;
	`, studentID, problemID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("mark solved: %w", err)
	}
	return item, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, studentID, problemID, status string, reset bool) (*model.Progress, error) {
	if reset || status == model.StatusNotStarted {
		item, err := scanProgress(r.pool.QueryRow(ctx, `
			INSERT INTO student_problem_progress (
				student_id, problem_id, status, total_attempts, successful_attempts, updated_at
			)
			SELECT $1, p.id, 'Not Started', 0, 0, NOW()
			FROM problems p
			WHERE p.id::text = $2 OR p.slug = $2
			ON CONFLICT (student_id, problem_id) DO UPDATE SET
				status = 'Not Started',
				total_attempts = 0,
				successful_attempts = 0,
				first_attempt_at = NULL,
				last_attempt_at = NULL,
				solved_at = NULL,
				updated_at = NOW()
			RETURNING `+progressSelect+`;
		`, studentID, problemID))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, model.ErrNotFound
			}
			return nil, fmt.Errorf("reset progress: %w", err)
		}
		return item, nil
	}

	item, err := scanProgress(r.pool.QueryRow(ctx, `
		INSERT INTO student_problem_progress (
			student_id, problem_id, status, total_attempts, successful_attempts,
			first_attempt_at, last_attempt_at, solved_at, updated_at
		)
		SELECT $1, p.id, $3,
			CASE WHEN $3 = 'Solved' THEN 1 ELSE 0 END,
			CASE WHEN $3 = 'Solved' THEN 1 ELSE 0 END,
			CASE WHEN $3 = 'Not Started' THEN NULL ELSE NOW() END,
			CASE WHEN $3 = 'Not Started' THEN NULL ELSE NOW() END,
			CASE WHEN $3 = 'Solved' THEN NOW() ELSE NULL END,
			NOW()
		FROM problems p
		WHERE p.id::text = $2 OR p.slug = $2
		ON CONFLICT (student_id, problem_id) DO UPDATE SET
			status = EXCLUDED.status,
			solved_at = CASE WHEN EXCLUDED.status = 'Solved' THEN COALESCE(student_problem_progress.solved_at, NOW()) ELSE student_problem_progress.solved_at END,
			updated_at = NOW()
		RETURNING `+progressSelect+`;
	`, studentID, problemID, status))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("update progress: %w", err)
	}
	return item, nil
}
