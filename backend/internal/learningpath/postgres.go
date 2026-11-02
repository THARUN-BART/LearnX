package learningpath

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository provides PostgreSQL persistence for adaptive learning paths.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetAdaptivePath(ctx context.Context, studentID string) (*AdaptivePath, error) {
	query := `
		SELECT id, student_id, current_topic_id, current_concept_key,
		       recommended_difficulty, overall_readiness_score, learning_stage,
		       created_at, updated_at
		FROM student_adaptive_paths
		WHERE student_id = $1`

	var p AdaptivePath
	err := r.pool.QueryRow(ctx, query, studentID).Scan(
		&p.ID, &p.StudentID, &p.CurrentTopicID, &p.CurrentConceptKey,
		&p.RecommendedDifficulty, &p.OverallReadinessScore, &p.LearningStage,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get adaptive path: %w", err)
	}

	nodes, err := r.GetPathNodes(ctx, studentID)
	if err == nil {
		p.Nodes = nodes
	}
	return &p, nil
}

func (r *PostgresRepository) UpsertAdaptivePath(ctx context.Context, path *AdaptivePath) error {
	query := `
		INSERT INTO student_adaptive_paths (
			student_id, current_topic_id, current_concept_key,
			recommended_difficulty, overall_readiness_score, learning_stage,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (student_id) DO UPDATE SET
			current_topic_id = EXCLUDED.current_topic_id,
			current_concept_key = EXCLUDED.current_concept_key,
			recommended_difficulty = EXCLUDED.recommended_difficulty,
			overall_readiness_score = EXCLUDED.overall_readiness_score,
			learning_stage = EXCLUDED.learning_stage,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		path.StudentID, path.CurrentTopicID, path.CurrentConceptKey,
		path.RecommendedDifficulty, path.OverallReadinessScore, path.LearningStage,
	).Scan(&path.ID, &path.CreatedAt, &path.UpdatedAt)
}

func (r *PostgresRepository) GetPathNodes(ctx context.Context, studentID string) ([]PathNode, error) {
	query := `
		SELECT n.id, n.student_id, n.topic_id, n.topic_name, n.concept_key, n.concept_name,
		       n.node_order, n.state, n.prerequisite_concept_keys, n.readiness_score,
		       n.recommended_difficulty, n.recommended_problem_id, p.title as problem_name,
		       n.updated_at
		FROM student_path_nodes n
		LEFT JOIN problems p ON p.id = n.recommended_problem_id
		WHERE n.student_id = $1
		ORDER BY n.node_order ASC`

	rows, err := r.pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []PathNode
	for rows.Next() {
		var n PathNode
		if err := rows.Scan(
			&n.ID, &n.StudentID, &n.TopicID, &n.TopicName, &n.ConceptKey, &n.ConceptName,
			&n.NodeOrder, &n.State, &n.PrerequisiteKeys, &n.ReadinessScore,
			&n.RecommendedDifficulty, &n.RecommendedProblemID, &n.RecommendedProblemName,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func (r *PostgresRepository) UpsertPathNode(ctx context.Context, n *PathNode) error {
	query := `
		INSERT INTO student_path_nodes (
			student_id, topic_id, topic_name, concept_key, concept_name,
			node_order, state, prerequisite_concept_keys, readiness_score,
			recommended_difficulty, recommended_problem_id, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		ON CONFLICT (student_id, concept_key) DO UPDATE SET
			topic_id = EXCLUDED.topic_id,
			topic_name = EXCLUDED.topic_name,
			concept_name = EXCLUDED.concept_name,
			node_order = EXCLUDED.node_order,
			state = EXCLUDED.state,
			prerequisite_concept_keys = EXCLUDED.prerequisite_concept_keys,
			readiness_score = EXCLUDED.readiness_score,
			recommended_difficulty = EXCLUDED.recommended_difficulty,
			recommended_problem_id = EXCLUDED.recommended_problem_id,
			updated_at = NOW()
		RETURNING id, updated_at`

	return r.pool.QueryRow(ctx, query,
		n.StudentID, n.TopicID, n.TopicName, n.ConceptKey, n.ConceptName,
		n.NodeOrder, n.State, n.PrerequisiteKeys, n.ReadinessScore,
		n.RecommendedDifficulty, n.RecommendedProblemID,
	).Scan(&n.ID, &n.UpdatedAt)
}

func (r *PostgresRepository) BatchUpsertPathNodes(ctx context.Context, studentID string, nodes []PathNode) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO student_path_nodes (
			student_id, topic_id, topic_name, concept_key, concept_name,
			node_order, state, prerequisite_concept_keys, readiness_score,
			recommended_difficulty, recommended_problem_id, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW())
		ON CONFLICT (student_id, concept_key) DO UPDATE SET
			topic_id = EXCLUDED.topic_id,
			topic_name = EXCLUDED.topic_name,
			concept_name = EXCLUDED.concept_name,
			node_order = EXCLUDED.node_order,
			state = EXCLUDED.state,
			prerequisite_concept_keys = EXCLUDED.prerequisite_concept_keys,
			readiness_score = EXCLUDED.readiness_score,
			recommended_difficulty = EXCLUDED.recommended_difficulty,
			recommended_problem_id = EXCLUDED.recommended_problem_id,
			updated_at = NOW()`

	for _, n := range nodes {
		_, err := tx.Exec(ctx, query,
			studentID, n.TopicID, n.TopicName, n.ConceptKey, n.ConceptName,
			n.NodeOrder, n.State, n.PrerequisiteKeys, n.ReadinessScore,
			n.RecommendedDifficulty, n.RecommendedProblemID,
		)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *PostgresRepository) GetActiveRecommendation(ctx context.Context, studentID string) (*LearningRecommendation, error) {
	query := `
		SELECT r.id, r.student_id, r.action_type, r.concept_key,
		       r.topic_id, COALESCE(t.name, '') as topic_name,
		       r.problem_id, COALESCE(p.title, '') as problem_title, COALESCE(p.slug, '') as problem_slug,
		       r.recommended_difficulty, r.headline, r.explanation, r.evidence,
		       r.is_active, r.completed_at, r.created_at
		FROM student_learning_recommendations r
		LEFT JOIN topics t ON t.id = r.topic_id
		LEFT JOIN problems p ON p.id = r.problem_id
		WHERE r.student_id = $1 AND r.is_active = TRUE
		ORDER BY r.created_at DESC
		LIMIT 1`

	var rec LearningRecommendation
	var rawEv []byte
	err := r.pool.QueryRow(ctx, query, studentID).Scan(
		&rec.ID, &rec.StudentID, &rec.ActionType, &rec.ConceptKey,
		&rec.TopicID, &rec.TopicName,
		&rec.ProblemID, &rec.ProblemTitle, &rec.ProblemSlug,
		&rec.RecommendedDifficulty, &rec.Headline, &rec.Explanation, &rawEv,
		&rec.IsActive, &rec.CompletedAt, &rec.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(rawEv, &rec.Evidence)
	if w, ok := rec.Evidence["relevant_weakness"].(string); ok {
		rec.RelevantWeakness = w
	}
	return &rec, nil
}

func (r *PostgresRepository) CreateRecommendation(ctx context.Context, rec *LearningRecommendation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Deactivate existing recommendations for this student
	_, _ = tx.Exec(ctx, `UPDATE student_learning_recommendations SET is_active = FALSE WHERE student_id = $1 AND is_active = TRUE`, rec.StudentID)

	rawEv, _ := json.Marshal(rec.Evidence)
	query := `
		INSERT INTO student_learning_recommendations (
			student_id, action_type, concept_key, topic_id, problem_id,
			recommended_difficulty, headline, explanation, evidence, is_active, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, NOW())
		RETURNING id, created_at`

	err = tx.QueryRow(ctx, query,
		rec.StudentID, rec.ActionType, rec.ConceptKey, rec.TopicID, rec.ProblemID,
		rec.RecommendedDifficulty, rec.Headline, rec.Explanation, rawEv,
	).Scan(&rec.ID, &rec.CreatedAt)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) CompleteRecommendation(ctx context.Context, studentID, recID string) error {
	query := `
		UPDATE student_learning_recommendations
		SET is_active = FALSE, completed_at = NOW()
		WHERE student_id = $1 AND id = $2`
	_, err := r.pool.Exec(ctx, query, studentID, recID)
	return err
}

func (r *PostgresRepository) ListRecommendationHistory(ctx context.Context, studentID string, limit, offset int) ([]LearningRecommendation, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM student_learning_recommendations WHERE student_id = $1`, studentID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT r.id, r.student_id, r.action_type, r.concept_key,
		       r.topic_id, COALESCE(t.name, '') as topic_name,
		       r.problem_id, COALESCE(p.title, '') as problem_title, COALESCE(p.slug, '') as problem_slug,
		       r.recommended_difficulty, r.headline, r.explanation, r.evidence,
		       r.is_active, r.completed_at, r.created_at
		FROM student_learning_recommendations r
		LEFT JOIN topics t ON t.id = r.topic_id
		LEFT JOIN problems p ON p.id = r.problem_id
		WHERE r.student_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []LearningRecommendation
	for rows.Next() {
		var rec LearningRecommendation
		var rawEv []byte
		if err := rows.Scan(
			&rec.ID, &rec.StudentID, &rec.ActionType, &rec.ConceptKey,
			&rec.TopicID, &rec.TopicName,
			&rec.ProblemID, &rec.ProblemTitle, &rec.ProblemSlug,
			&rec.RecommendedDifficulty, &rec.Headline, &rec.Explanation, &rawEv,
			&rec.IsActive, &rec.CompletedAt, &rec.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(rawEv, &rec.Evidence)
		if w, ok := rec.Evidence["relevant_weakness"].(string); ok {
			rec.RelevantWeakness = w
		}
		list = append(list, rec)
	}
	return list, total, nil
}

var _ Repository = (*PostgresRepository)(nil)
