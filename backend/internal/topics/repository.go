package topics

import (
	"context"
	"errors"
	"fmt"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Get(ctx context.Context, id string) (*model.Topic, error)
	ListByRoadmap(ctx context.Context, roadmapID string) ([]model.Topic, error)
	ListSubtopics(ctx context.Context, topicID string) ([]model.Subtopic, error)
	GetSubtopic(ctx context.Context, id string) (*model.Subtopic, error)
}

type MemoryRepository struct {
	DB *model.Memory
}

func NewMemoryRepository(db *model.Memory) *MemoryRepository {
	return &MemoryRepository{DB: db}
}

func (r *MemoryRepository) Get(ctx context.Context, id string) (*model.Topic, error) {
	item, ok := r.DB.GetTopic(id)
	if !ok {
		return nil, model.ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) ListByRoadmap(ctx context.Context, roadmapID string) ([]model.Topic, error) {
	return r.DB.ListTopics(roadmapID), nil
}

func (r *MemoryRepository) ListSubtopics(ctx context.Context, topicID string) ([]model.Subtopic, error) {
	return r.DB.ListSubtopics(topicID), nil
}

func (r *MemoryRepository) GetSubtopic(ctx context.Context, id string) (*model.Subtopic, error) {
	item, ok := r.DB.GetSubtopic(id)
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

func (r *PostgresRepository) Get(ctx context.Context, id string) (*model.Topic, error) {
	var item model.Topic
	err := r.pool.QueryRow(ctx, `
		SELECT t.id, t.roadmap_id, t.name, t.slug, t.description, t.topic_order, t.created_at, t.updated_at,
		       (SELECT COUNT(*) FROM subtopics s WHERE s.topic_id = t.id),
		       (SELECT COUNT(*) FROM problems p WHERE p.topic_id = t.id AND p.is_active)
		FROM topics t
		WHERE t.id = $1;
	`, id).Scan(&item.ID, &item.RoadmapID, &item.Name, &item.Slug, &item.Description, &item.TopicOrder, &item.CreatedAt, &item.UpdatedAt, &item.SubtopicCount, &item.ProblemCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get topic: %w", err)
	}
	return &item, nil
}

func (r *PostgresRepository) ListByRoadmap(ctx context.Context, roadmapID string) ([]model.Topic, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT t.id, t.roadmap_id, t.name, t.slug, t.description, t.topic_order, t.created_at, t.updated_at,
		       (SELECT COUNT(*) FROM subtopics s WHERE s.topic_id = t.id),
		       (SELECT COUNT(*) FROM problems p WHERE p.topic_id = t.id AND p.is_active)
		FROM topics t
		WHERE t.roadmap_id = $1
		ORDER BY t.topic_order;
	`, roadmapID)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}
	defer rows.Close()
	out := make([]model.Topic, 0)
	for rows.Next() {
		var item model.Topic
		if err := rows.Scan(&item.ID, &item.RoadmapID, &item.Name, &item.Slug, &item.Description, &item.TopicOrder, &item.CreatedAt, &item.UpdatedAt, &item.SubtopicCount, &item.ProblemCount); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) ListSubtopics(ctx context.Context, topicID string) ([]model.Subtopic, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.topic_id, s.name, s.slug, s.description, s.subtopic_order, s.created_at, s.updated_at,
		       (SELECT COUNT(*) FROM problems p WHERE p.subtopic_id = s.id AND p.is_active)
		FROM subtopics s
		WHERE s.topic_id = $1
		ORDER BY s.subtopic_order;
	`, topicID)
	if err != nil {
		return nil, fmt.Errorf("list subtopics: %w", err)
	}
	defer rows.Close()
	out := make([]model.Subtopic, 0)
	for rows.Next() {
		var item model.Subtopic
		if err := rows.Scan(&item.ID, &item.TopicID, &item.Name, &item.Slug, &item.Description, &item.SubtopicOrder, &item.CreatedAt, &item.UpdatedAt, &item.ProblemCount); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (r *PostgresRepository) GetSubtopic(ctx context.Context, id string) (*model.Subtopic, error) {
	var item model.Subtopic
	err := r.pool.QueryRow(ctx, `
		SELECT s.id, s.topic_id, s.name, s.slug, s.description, s.subtopic_order, s.created_at, s.updated_at,
		       (SELECT COUNT(*) FROM problems p WHERE p.subtopic_id = s.id AND p.is_active)
		FROM subtopics s
		WHERE s.id = $1;
	`, id).Scan(&item.ID, &item.TopicID, &item.Name, &item.Slug, &item.Description, &item.SubtopicOrder, &item.CreatedAt, &item.UpdatedAt, &item.ProblemCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get subtopic: %w", err)
	}
	return &item, nil
}
