package problems

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Get(ctx context.Context, idOrSlug string) (*model.Problem, error)
	List(ctx context.Context, filter model.ProblemFilter) ([]model.Problem, int, error)
	Related(ctx context.Context, problemID string, limit int) ([]model.Problem, error)
}

type MemoryRepository struct {
	DB *model.Memory
}

func NewMemoryRepository(db *model.Memory) *MemoryRepository {
	return &MemoryRepository{DB: db}
}

func (r *MemoryRepository) Get(ctx context.Context, idOrSlug string) (*model.Problem, error) {
	item, ok := r.DB.GetProblem(idOrSlug)
	if !ok {
		return nil, model.ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) List(ctx context.Context, filter model.ProblemFilter) ([]model.Problem, int, error) {
	items, total := r.DB.ListProblems(filter)
	return items, total, nil
}

func (r *MemoryRepository) Related(ctx context.Context, problemID string, limit int) ([]model.Problem, error) {
	return r.DB.RelatedProblems(problemID, limit), nil
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

const problemSelect = `
	p.id, p.topic_id, p.subtopic_id, p.title, p.slug, p.summary, p.difficulty,
	p.external_platform, p.external_problem_url, p.source_sheet, p.source_problem_id,
	p.problem_order, p.tags, p.is_active, p.created_at, p.updated_at,
	t.name, s.name
`

func scanProblem(row pgx.Row) (*model.Problem, error) {
	var item model.Problem
	if err := row.Scan(
		&item.ID, &item.TopicID, &item.SubtopicID, &item.Title, &item.Slug, &item.Summary, &item.Difficulty,
		&item.ExternalPlatform, &item.ExternalProblemURL, &item.SourceSheet, &item.SourceProblemID,
		&item.ProblemOrder, &item.Tags, &item.IsActive, &item.CreatedAt, &item.UpdatedAt,
		&item.TopicName, &item.SubtopicName,
	); err != nil {
		return nil, err
	}
	if item.Tags == nil {
		item.Tags = []string{}
	}
	return &item, nil
}

func (r *PostgresRepository) Get(ctx context.Context, idOrSlug string) (*model.Problem, error) {
	item, err := scanProblem(r.pool.QueryRow(ctx, `
		SELECT `+problemSelect+`
		FROM problems p
		JOIN topics t ON t.id = p.topic_id
		JOIN subtopics s ON s.id = p.subtopic_id
		WHERE p.id::text = $1 OR p.slug = $1
		LIMIT 1;
	`, idOrSlug))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get problem: %w", err)
	}
	return item, nil
}

func (r *PostgresRepository) List(ctx context.Context, filter model.ProblemFilter) ([]model.Problem, int, error) {
	where := []string{"p.is_active = TRUE"}
	args := []interface{}{}
	add := func(clause string, value interface{}) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(clause, len(args)))
	}

	if filter.TopicID != "" {
		add("p.topic_id = $%d", filter.TopicID)
	}
	if filter.SubtopicID != "" {
		add("p.subtopic_id = $%d", filter.SubtopicID)
	}
	if filter.RoadmapID != "" {
		add("t.roadmap_id = $%d", filter.RoadmapID)
	}
	if filter.Difficulty != "" {
		add("p.difficulty = $%d", filter.Difficulty)
	}
	if filter.Query != "" {
		add("(LOWER(p.title) LIKE LOWER($%d) OR LOWER(p.slug) LIKE LOWER($%d))", "%"+filter.Query+"%")
		args[len(args)-1] = "%" + filter.Query + "%"
		// reuse same arg index for slug by rewriting last clause
		where[len(where)-1] = fmt.Sprintf("(LOWER(p.title) LIKE LOWER($%d) OR LOWER(p.slug) LIKE LOWER($%d))", len(args), len(args))
	}

	progressJoin := ""
	progressSelect := "''"
	if filter.StudentID != "" {
		progressJoin = "LEFT JOIN student_problem_progress sp ON sp.problem_id = p.id AND sp.student_id = $" + fmt.Sprint(len(args)+1)
		args = append(args, filter.StudentID)
		progressSelect = "COALESCE(sp.status, 'Not Started')"
		if filter.Status != "" {
			args = append(args, filter.Status)
			where = append(where, fmt.Sprintf("COALESCE(sp.status, 'Not Started') = $%d", len(args)))
		}
	}

	whereSQL := strings.Join(where, " AND ")
	order := "p.problem_order ASC, p.title ASC"
	if filter.Sort == "title" {
		order = "p.title ASC"
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM problems p
		JOIN topics t ON t.id = p.topic_id
		JOIN subtopics s ON s.id = p.subtopic_id
		%s
		WHERE %s;
	`, progressJoin, whereSQL)

	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count problems: %w", err)
	}

	limit := filter.Limit
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	listQuery := fmt.Sprintf(`
		SELECT %s, %s
		FROM problems p
		JOIN topics t ON t.id = p.topic_id
		JOIN subtopics s ON s.id = p.subtopic_id
		%s
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d;
	`, problemSelect, progressSelect, progressJoin, whereSQL, order, len(args)-1, len(args))

	rows, err := r.pool.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list problems: %w", err)
	}
	defer rows.Close()

	out := make([]model.Problem, 0)
	for rows.Next() {
		var item model.Problem
		if err := rows.Scan(
			&item.ID, &item.TopicID, &item.SubtopicID, &item.Title, &item.Slug, &item.Summary, &item.Difficulty,
			&item.ExternalPlatform, &item.ExternalProblemURL, &item.SourceSheet, &item.SourceProblemID,
			&item.ProblemOrder, &item.Tags, &item.IsActive, &item.CreatedAt, &item.UpdatedAt,
			&item.TopicName, &item.SubtopicName, &item.ProgressStatus,
		); err != nil {
			return nil, 0, err
		}
		if item.Tags == nil {
			item.Tags = []string{}
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *PostgresRepository) Related(ctx context.Context, problemID string, limit int) ([]model.Problem, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.pool.Query(ctx, `
		SELECT `+problemSelect+`
		FROM problems p
		JOIN topics t ON t.id = p.topic_id
		JOIN subtopics s ON s.id = p.subtopic_id
		WHERE p.is_active = TRUE
		  AND p.id <> $1
		  AND p.subtopic_id = (SELECT subtopic_id FROM problems WHERE id::text = $1 OR slug = $1 LIMIT 1)
		ORDER BY p.problem_order
		LIMIT $2;
	`, problemID, limit)
	if err != nil {
		return nil, fmt.Errorf("related problems: %w", err)
	}
	defer rows.Close()
	out := make([]model.Problem, 0)
	for rows.Next() {
		var item model.Problem
		if err := rows.Scan(
			&item.ID, &item.TopicID, &item.SubtopicID, &item.Title, &item.Slug, &item.Summary, &item.Difficulty,
			&item.ExternalPlatform, &item.ExternalProblemURL, &item.SourceSheet, &item.SourceProblemID,
			&item.ProblemOrder, &item.Tags, &item.IsActive, &item.CreatedAt, &item.UpdatedAt,
			&item.TopicName, &item.SubtopicName,
		); err != nil {
			return nil, err
		}
		if item.Tags == nil {
			item.Tags = []string{}
		}
		out = append(out, item)
	}
	return out, rows.Err()
}
