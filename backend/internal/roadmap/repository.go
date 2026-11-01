package roadmap

import (
	"context"
	"errors"
	"fmt"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	List(ctx context.Context) ([]model.Roadmap, error)
	Get(ctx context.Context, idOrSlug string) (*model.Roadmap, error)
}

type MemoryRepository struct {
	DB *model.Memory
}

func NewMemoryRepository(db *model.Memory) *MemoryRepository {
	return &MemoryRepository{DB: db}
}

func (r *MemoryRepository) List(ctx context.Context) ([]model.Roadmap, error) {
	return r.DB.ListRoadmaps(), nil
}

func (r *MemoryRepository) Get(ctx context.Context, idOrSlug string) (*model.Roadmap, error) {
	item, ok := r.DB.GetRoadmap(idOrSlug)
	if !ok {
		return nil, model.ErrNotFound
	}
	return item, nil
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) List(ctx context.Context) ([]model.Roadmap, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT r.id, r.name, r.slug, r.description, r.source_name, r.source_url, r.is_active, r.created_at, r.updated_at,
		       (SELECT COUNT(*) FROM topics t WHERE t.roadmap_id = r.id),
		       (SELECT COUNT(*) FROM problems p JOIN topics t ON t.id = p.topic_id WHERE t.roadmap_id = r.id AND p.is_active)
		FROM roadmaps r
		WHERE r.is_active = TRUE
		ORDER BY r.name;
	`)
	if err != nil {
		return nil, fmt.Errorf("list roadmaps: %w", err)
	}
	defer rows.Close()

	out := make([]model.Roadmap, 0)
	for rows.Next() {
		var item model.Roadmap
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug, &item.Description, &item.SourceName, &item.SourceURL, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &item.TopicCount, &item.ProblemCount); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) Get(ctx context.Context, idOrSlug string) (*model.Roadmap, error) {
	var item model.Roadmap
	err := r.pool.QueryRow(ctx, `
		SELECT r.id, r.name, r.slug, r.description, r.source_name, r.source_url, r.is_active, r.created_at, r.updated_at,
		       (SELECT COUNT(*) FROM topics t WHERE t.roadmap_id = r.id),
		       (SELECT COUNT(*) FROM problems p JOIN topics t ON t.id = p.topic_id WHERE t.roadmap_id = r.id AND p.is_active)
		FROM roadmaps r
		WHERE r.id::text = $1 OR r.slug = $1
		LIMIT 1;
	`, idOrSlug).Scan(&item.ID, &item.Name, &item.Slug, &item.Description, &item.SourceName, &item.SourceURL, &item.IsActive, &item.CreatedAt, &item.UpdatedAt, &item.TopicCount, &item.ProblemCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get roadmap: %w", err)
	}
	return &item, nil
}
