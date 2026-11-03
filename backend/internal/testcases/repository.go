package testcases

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("test case not found")

type Repository interface {
	Create(ctx context.Context, tc *TestCase) error
	ListByProblem(ctx context.Context, problemID string, includeHidden bool) ([]TestCase, error)
	Get(ctx context.Context, id string) (*TestCase, error)
	Delete(ctx context.Context, id string) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) Create(ctx context.Context, tc *TestCase) error {
	query := `
		INSERT INTO test_cases (id, problem_id, input_data, expected_output, is_hidden, test_order, created_at)
		VALUES (COALESCE(NULLIF($1, '')::uuid, gen_random_uuid()), $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at;
	`
	return r.pool.QueryRow(ctx, query,
		tc.ID, tc.ProblemID, tc.InputData, tc.ExpectedOutput, tc.IsHidden, tc.TestOrder,
	).Scan(&tc.ID, &tc.CreatedAt)
}

func (r *PostgresRepository) ListByProblem(ctx context.Context, problemID string, includeHidden bool) ([]TestCase, error) {
	query := `
		SELECT id, problem_id, input_data, expected_output, is_hidden, test_order, created_at
		FROM test_cases
		WHERE (problem_id::text = $1 OR problem_id = (SELECT id FROM problems WHERE slug = $1 LIMIT 1))
	`
	if !includeHidden {
		query += ` AND is_hidden = false`
	}
	query += ` ORDER BY test_order ASC, created_at ASC;`

	rows, err := r.pool.Query(ctx, query, problemID)
	if err != nil {
		return nil, fmt.Errorf("list test cases: %w", err)
	}
	defer rows.Close()

	list := make([]TestCase, 0)
	for rows.Next() {
		var tc TestCase
		if err := rows.Scan(
			&tc.ID, &tc.ProblemID, &tc.InputData, &tc.ExpectedOutput, &tc.IsHidden, &tc.TestOrder, &tc.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan test case: %w", err)
		}
		list = append(list, tc)
	}
	return list, rows.Err()
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*TestCase, error) {
	query := `
		SELECT id, problem_id, input_data, expected_output, is_hidden, test_order, created_at
		FROM test_cases
		WHERE id::text = $1 LIMIT 1;
	`
	var tc TestCase
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tc.ID, &tc.ProblemID, &tc.InputData, &tc.ExpectedOutput, &tc.IsHidden, &tc.TestOrder, &tc.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get test case: %w", err)
	}
	return &tc, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM test_cases WHERE id::text = $1;`
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete test case: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// MemoryRepository implements Repository in memory for fast unit testing.
type MemoryRepository struct {
	mu        sync.RWMutex
	testCases map[string]*TestCase
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		testCases: make(map[string]*TestCase),
	}
}

func (m *MemoryRepository) Create(ctx context.Context, tc *TestCase) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if tc.ID == "" {
		tc.ID = fmt.Sprintf("tc-%d", len(m.testCases)+1)
	}
	tc.CreatedAt = time.Now()
	m.testCases[tc.ID] = tc
	return nil
}

func (m *MemoryRepository) ListByProblem(ctx context.Context, problemID string, includeHidden bool) ([]TestCase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []TestCase
	for _, tc := range m.testCases {
		if tc.ProblemID == problemID || (tc.ProblemSlug != "" && tc.ProblemSlug == problemID) {
			if !includeHidden && tc.IsHidden {
				continue
			}
			list = append(list, *tc)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].TestOrder < list[j].TestOrder
	})
	return list, nil
}

func (m *MemoryRepository) Get(ctx context.Context, id string) (*TestCase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tc, ok := m.testCases[id]
	if !ok {
		return nil, ErrNotFound
	}
	return tc, nil
}

func (m *MemoryRepository) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.testCases[id]; !ok {
		return ErrNotFound
	}
	delete(m.testCases, id)
	return nil
}
