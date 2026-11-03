package learningmemory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository provides PostgreSQL persistence using pgxpool.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetConceptMemory(ctx context.Context, studentID, conceptKey string) (*ConceptMemory, error) {
	query := `
		SELECT
			cm.id, cm.student_id, cm.topic_id, cm.subtopic_id, cm.concept_key,
			COALESCE(t.name, '') as topic_name,
			cm.learning_state, cm.estimated_mastery, cm.confidence_level,
			cm.total_attempts, cm.independent_solves, cm.assisted_solves,
			cm.failed_attempts, cm.hint_requests, cm.mistake_count, cm.corrected_mistakes,
			cm.first_practiced_at, cm.last_practiced_at, cm.last_successful_at, cm.next_review_at,
			cm.created_at, cm.updated_at
		FROM student_concept_memory cm
		LEFT JOIN topics t ON t.id = cm.topic_id
		WHERE cm.student_id = $1 AND cm.concept_key = $2`

	var cm ConceptMemory
	err := r.pool.QueryRow(ctx, query, studentID, conceptKey).Scan(
		&cm.ID, &cm.StudentID, &cm.TopicID, &cm.SubtopicID, &cm.ConceptKey,
		&cm.TopicName,
		&cm.LearningState, &cm.EstimatedMastery, &cm.ConfidenceLevel,
		&cm.TotalAttempts, &cm.IndependentSolves, &cm.AssistedSolves,
		&cm.FailedAttempts, &cm.HintRequests, &cm.MistakeCount, &cm.CorrectedMistakes,
		&cm.FirstPracticedAt, &cm.LastPracticedAt, &cm.LastSuccessfulAt, &cm.NextReviewAt,
		&cm.CreatedAt, &cm.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to query concept memory: %w", err)
	}
	cm.ConceptName = formatConceptName(cm.ConceptKey)
	return &cm, nil
}

func (r *PostgresRepository) UpsertConceptMemory(ctx context.Context, cm *ConceptMemory) error {
	query := `
		INSERT INTO student_concept_memory (
			student_id, topic_id, subtopic_id, concept_key,
			learning_state, estimated_mastery, confidence_level,
			total_attempts, independent_solves, assisted_solves,
			failed_attempts, hint_requests, mistake_count, corrected_mistakes,
			first_practiced_at, last_practiced_at, last_successful_at, next_review_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW()
		)
		ON CONFLICT (student_id, concept_key) DO UPDATE SET
			topic_id = COALESCE(EXCLUDED.topic_id, student_concept_memory.topic_id),
			subtopic_id = COALESCE(EXCLUDED.subtopic_id, student_concept_memory.subtopic_id),
			learning_state = EXCLUDED.learning_state,
			estimated_mastery = EXCLUDED.estimated_mastery,
			confidence_level = EXCLUDED.confidence_level,
			total_attempts = EXCLUDED.total_attempts,
			independent_solves = EXCLUDED.independent_solves,
			assisted_solves = EXCLUDED.assisted_solves,
			failed_attempts = EXCLUDED.failed_attempts,
			hint_requests = EXCLUDED.hint_requests,
			mistake_count = EXCLUDED.mistake_count,
			corrected_mistakes = EXCLUDED.corrected_mistakes,
			last_practiced_at = EXCLUDED.last_practiced_at,
			last_successful_at = COALESCE(EXCLUDED.last_successful_at, student_concept_memory.last_successful_at),
			next_review_at = EXCLUDED.next_review_at,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		cm.StudentID, cm.TopicID, cm.SubtopicID, cm.ConceptKey,
		cm.LearningState, cm.EstimatedMastery, cm.ConfidenceLevel,
		cm.TotalAttempts, cm.IndependentSolves, cm.AssistedSolves,
		cm.FailedAttempts, cm.HintRequests, cm.MistakeCount, cm.CorrectedMistakes,
		cm.FirstPracticedAt, cm.LastPracticedAt, cm.LastSuccessfulAt, cm.NextReviewAt,
	).Scan(&cm.ID, &cm.CreatedAt, &cm.UpdatedAt)
}

func (r *PostgresRepository) ListConceptMemories(ctx context.Context, studentID string, filter ConceptFilter) ([]ConceptMemory, int, error) {
	baseQuery := `
		FROM student_concept_memory cm
		LEFT JOIN topics t ON t.id = cm.topic_id
		WHERE cm.student_id = $1`
	args := []interface{}{studentID}
	idx := 2

	if filter.TopicID != "" {
		baseQuery += fmt.Sprintf(" AND cm.topic_id = $%d", idx)
		args = append(args, filter.TopicID)
		idx++
	}
	if filter.LearningState != "" {
		baseQuery += fmt.Sprintf(" AND cm.learning_state = $%d", idx)
		args = append(args, filter.LearningState)
		idx++
	}
	if filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND (cm.concept_key ILIKE $%d OR t.name ILIKE $%d)", idx, idx)
		args = append(args, "%"+filter.Search+"%")
		idx++
	}

	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count concepts: %w", err)
	}

	selectQuery := `
		SELECT
			cm.id, cm.student_id, cm.topic_id, cm.subtopic_id, cm.concept_key,
			COALESCE(t.name, '') as topic_name,
			cm.learning_state, cm.estimated_mastery, cm.confidence_level,
			cm.total_attempts, cm.independent_solves, cm.assisted_solves,
			cm.failed_attempts, cm.hint_requests, cm.mistake_count, cm.corrected_mistakes,
			cm.first_practiced_at, cm.last_practiced_at, cm.last_successful_at, cm.next_review_at,
			cm.created_at, cm.updated_at ` + baseQuery + ` ORDER BY cm.updated_at DESC`

	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list concepts: %w", err)
	}
	defer rows.Close()

	var list []ConceptMemory
	for rows.Next() {
		var cm ConceptMemory
		if err := rows.Scan(
			&cm.ID, &cm.StudentID, &cm.TopicID, &cm.SubtopicID, &cm.ConceptKey,
			&cm.TopicName,
			&cm.LearningState, &cm.EstimatedMastery, &cm.ConfidenceLevel,
			&cm.TotalAttempts, &cm.IndependentSolves, &cm.AssistedSolves,
			&cm.FailedAttempts, &cm.HintRequests, &cm.MistakeCount, &cm.CorrectedMistakes,
			&cm.FirstPracticedAt, &cm.LastPracticedAt, &cm.LastSuccessfulAt, &cm.NextReviewAt,
			&cm.CreatedAt, &cm.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		cm.ConceptName = formatConceptName(cm.ConceptKey)
		list = append(list, cm)
	}
	return list, total, nil
}

func (r *PostgresRepository) GetStrengths(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error) {
	if limit <= 0 {
		limit = 5
	}
	query := `
		SELECT
			cm.id, cm.student_id, cm.topic_id, cm.subtopic_id, cm.concept_key,
			COALESCE(t.name, '') as topic_name,
			cm.learning_state, cm.estimated_mastery, cm.confidence_level,
			cm.total_attempts, cm.independent_solves, cm.assisted_solves,
			cm.failed_attempts, cm.hint_requests, cm.mistake_count, cm.corrected_mistakes,
			cm.first_practiced_at, cm.last_practiced_at, cm.last_successful_at, cm.next_review_at,
			cm.created_at, cm.updated_at
		FROM student_concept_memory cm
		LEFT JOIN topics t ON t.id = cm.topic_id
		WHERE cm.student_id = $1 AND cm.independent_solves > 0
		ORDER BY cm.independent_solves DESC, cm.estimated_mastery DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ConceptMemory
	for rows.Next() {
		var cm ConceptMemory
		if err := rows.Scan(
			&cm.ID, &cm.StudentID, &cm.TopicID, &cm.SubtopicID, &cm.ConceptKey,
			&cm.TopicName,
			&cm.LearningState, &cm.EstimatedMastery, &cm.ConfidenceLevel,
			&cm.TotalAttempts, &cm.IndependentSolves, &cm.AssistedSolves,
			&cm.FailedAttempts, &cm.HintRequests, &cm.MistakeCount, &cm.CorrectedMistakes,
			&cm.FirstPracticedAt, &cm.LastPracticedAt, &cm.LastSuccessfulAt, &cm.NextReviewAt,
			&cm.CreatedAt, &cm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cm.ConceptName = formatConceptName(cm.ConceptKey)
		list = append(list, cm)
	}
	return list, nil
}

func (r *PostgresRepository) GetWeaknesses(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error) {
	if limit <= 0 {
		limit = 5
	}
	query := `
		SELECT
			cm.id, cm.student_id, cm.topic_id, cm.subtopic_id, cm.concept_key,
			COALESCE(t.name, '') as topic_name,
			cm.learning_state, cm.estimated_mastery, cm.confidence_level,
			cm.total_attempts, cm.independent_solves, cm.assisted_solves,
			cm.failed_attempts, cm.hint_requests, cm.mistake_count, cm.corrected_mistakes,
			cm.first_practiced_at, cm.last_practiced_at, cm.last_successful_at, cm.next_review_at,
			cm.created_at, cm.updated_at
		FROM student_concept_memory cm
		LEFT JOIN topics t ON t.id = cm.topic_id
		WHERE cm.student_id = $1 AND (cm.mistake_count > cm.corrected_mistakes OR cm.learning_state = 'struggling')
		ORDER BY (cm.mistake_count - cm.corrected_mistakes) DESC, cm.failed_attempts DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, studentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []ConceptMemory
	for rows.Next() {
		var cm ConceptMemory
		if err := rows.Scan(
			&cm.ID, &cm.StudentID, &cm.TopicID, &cm.SubtopicID, &cm.ConceptKey,
			&cm.TopicName,
			&cm.LearningState, &cm.EstimatedMastery, &cm.ConfidenceLevel,
			&cm.TotalAttempts, &cm.IndependentSolves, &cm.AssistedSolves,
			&cm.FailedAttempts, &cm.HintRequests, &cm.MistakeCount, &cm.CorrectedMistakes,
			&cm.FirstPracticedAt, &cm.LastPracticedAt, &cm.LastSuccessfulAt, &cm.NextReviewAt,
			&cm.CreatedAt, &cm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cm.ConceptName = formatConceptName(cm.ConceptKey)
		list = append(list, cm)
	}
	return list, nil
}

func (r *PostgresRepository) RecordEvent(ctx context.Context, event *LearningEvent) error {
	dataBytes, _ := json.Marshal(event.EventData)
	query := `
		INSERT INTO student_learning_events (
			student_id, event_type, problem_id, submission_id, diagnosis_id,
			topic_id, subtopic_id, concept_key, event_data, occurred_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		event.StudentID, event.EventType, event.ProblemID, event.SubmissionID, event.DiagnosisID,
		event.TopicID, event.SubtopicID, event.ConceptKey, dataBytes, event.OccurredAt,
	).Scan(&event.ID, &event.CreatedAt)
}

func (r *PostgresRepository) ListEvents(ctx context.Context, studentID string, filter EventFilter) ([]LearningEvent, int, error) {
	baseQuery := `FROM student_learning_events WHERE student_id = $1`
	args := []interface{}{studentID}
	idx := 2

	if filter.EventType != "" {
		baseQuery += fmt.Sprintf(" AND event_type = $%d", idx)
		args = append(args, filter.EventType)
		idx++
	}
	if filter.ConceptKey != "" {
		baseQuery += fmt.Sprintf(" AND concept_key = $%d", idx)
		args = append(args, filter.ConceptKey)
		idx++
	}
	if filter.ProblemID != "" {
		baseQuery += fmt.Sprintf(" AND problem_id = $%d", idx)
		args = append(args, filter.ProblemID)
		idx++
	}
	if filter.StartDate != nil {
		baseQuery += fmt.Sprintf(" AND occurred_at >= $%d", idx)
		args = append(args, *filter.StartDate)
		idx++
	}
	if filter.EndDate != nil {
		baseQuery += fmt.Sprintf(" AND occurred_at <= $%d", idx)
		args = append(args, *filter.EndDate)
		idx++
	}

	var total int
	if err := r.pool.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	selectQuery := `
		SELECT id, student_id, event_type, problem_id, submission_id, diagnosis_id,
		       topic_id, subtopic_id, concept_key, event_data, occurred_at, created_at ` +
		baseQuery + ` ORDER BY occurred_at DESC`

	if filter.Limit > 0 {
		selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
		args = append(args, filter.Limit, filter.Offset)
	}

	rows, err := r.pool.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []LearningEvent
	for rows.Next() {
		var e LearningEvent
		var rawData []byte
		if err := rows.Scan(
			&e.ID, &e.StudentID, &e.EventType, &e.ProblemID, &e.SubmissionID, &e.DiagnosisID,
			&e.TopicID, &e.SubtopicID, &e.ConceptKey, &rawData, &e.OccurredAt, &e.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(rawData, &e.EventData)
		list = append(list, e)
	}
	return list, total, nil
}

func (r *PostgresRepository) HasRecentHint(ctx context.Context, studentID, problemID string, before time.Time) (bool, error) {
	threshold := before.Add(-2 * time.Hour)
	query := `
		SELECT EXISTS (
			SELECT 1 FROM student_learning_events
			WHERE student_id = $1 AND problem_id = $2 AND event_type = $3
			  AND occurred_at >= $4 AND occurred_at <= $5
		)`
	var exists bool
	err := r.pool.QueryRow(ctx, query, studentID, problemID, EventHintRequested, threshold, before.Add(time.Minute)).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetMistakeMemory(ctx context.Context, studentID, category, subtype, conceptKey string) (*MistakeMemory, error) {
	query := `
		SELECT id, student_id, category, subtype, concept_key, occurrence_count,
		       corrected_count, first_occurrence_at, last_occurrence_at, last_corrected_at,
		       status, supporting_diagnosis_ids, created_at, updated_at
		FROM student_mistake_memory
		WHERE student_id = $1 AND category = $2 AND subtype = $3 AND concept_key = $4`

	var mm MistakeMemory
	var rawIDs []byte
	err := r.pool.QueryRow(ctx, query, studentID, category, subtype, conceptKey).Scan(
		&mm.ID, &mm.StudentID, &mm.Category, &mm.Subtype, &mm.ConceptKey,
		&mm.OccurrenceCount, &mm.CorrectedCount, &mm.FirstOccurrenceAt,
		&mm.LastOccurrenceAt, &mm.LastCorrectedAt, &mm.Status,
		&rawIDs, &mm.CreatedAt, &mm.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	_ = json.Unmarshal(rawIDs, &mm.SupportingDiagnosisIDs)
	return &mm, nil
}

func (r *PostgresRepository) UpsertMistakeMemory(ctx context.Context, mm *MistakeMemory) error {
	rawIDs, _ := json.Marshal(mm.SupportingDiagnosisIDs)
	query := `
		INSERT INTO student_mistake_memory (
			student_id, category, subtype, concept_key, occurrence_count,
			corrected_count, first_occurrence_at, last_occurrence_at, last_corrected_at,
			status, supporting_diagnosis_ids, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
		)
		ON CONFLICT (student_id, category, subtype, concept_key) DO UPDATE SET
			occurrence_count = EXCLUDED.occurrence_count,
			corrected_count = EXCLUDED.corrected_count,
			last_occurrence_at = EXCLUDED.last_occurrence_at,
			last_corrected_at = COALESCE(EXCLUDED.last_corrected_at, student_mistake_memory.last_corrected_at),
			status = EXCLUDED.status,
			supporting_diagnosis_ids = EXCLUDED.supporting_diagnosis_ids,
			updated_at = NOW()
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		mm.StudentID, mm.Category, mm.Subtype, mm.ConceptKey, mm.OccurrenceCount,
		mm.CorrectedCount, mm.FirstOccurrenceAt, mm.LastOccurrenceAt, mm.LastCorrectedAt,
		mm.Status, rawIDs,
	).Scan(&mm.ID, &mm.CreatedAt, &mm.UpdatedAt)
}

func (r *PostgresRepository) ListMistakeMemories(ctx context.Context, studentID string, status string) ([]MistakeMemory, error) {
	query := `
		SELECT id, student_id, category, subtype, concept_key, occurrence_count,
		       corrected_count, first_occurrence_at, last_occurrence_at, last_corrected_at,
		       status, supporting_diagnosis_ids, created_at, updated_at
		FROM student_mistake_memory
		WHERE student_id = $1`
	args := []interface{}{studentID}
	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}
	query += " ORDER BY occurrence_count DESC, updated_at DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []MistakeMemory
	for rows.Next() {
		var mm MistakeMemory
		var rawIDs []byte
		if err := rows.Scan(
			&mm.ID, &mm.StudentID, &mm.Category, &mm.Subtype, &mm.ConceptKey,
			&mm.OccurrenceCount, &mm.CorrectedCount, &mm.FirstOccurrenceAt,
			&mm.LastOccurrenceAt, &mm.LastCorrectedAt, &mm.Status,
			&rawIDs, &mm.CreatedAt, &mm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(rawIDs, &mm.SupportingDiagnosisIDs)
		list = append(list, mm)
	}
	return list, nil
}

func (r *PostgresRepository) UpsertSnapshot(ctx context.Context, snap *LearningSnapshot) error {
	rawSummary, _ := json.Marshal(snap.Summary)
	query := `
		INSERT INTO student_learning_snapshots (
			student_id, snapshot_date, total_concepts_practised, independent_solves,
			assisted_solves, mistake_count, corrected_mistake_count, active_learning_days,
			summary, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		ON CONFLICT (student_id, snapshot_date) DO UPDATE SET
			total_concepts_practised = EXCLUDED.total_concepts_practised,
			independent_solves = EXCLUDED.independent_solves,
			assisted_solves = EXCLUDED.assisted_solves,
			mistake_count = EXCLUDED.mistake_count,
			corrected_mistake_count = EXCLUDED.corrected_mistake_count,
			active_learning_days = EXCLUDED.active_learning_days,
			summary = EXCLUDED.summary
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query,
		snap.StudentID, snap.SnapshotDate, snap.TotalConceptsPractised,
		snap.IndependentSolves, snap.AssistedSolves, snap.MistakeCount,
		snap.CorrectedMistakeCount, snap.ActiveLearningDays, rawSummary,
	).Scan(&snap.ID, &snap.CreatedAt)
}

func (r *PostgresRepository) ListSnapshots(ctx context.Context, studentID string, days int) ([]LearningSnapshot, error) {
	if days <= 0 {
		days = 30
	}
	query := `
		SELECT id, student_id, snapshot_date::text, total_concepts_practised,
		       independent_solves, assisted_solves, mistake_count,
		       corrected_mistake_count, active_learning_days, summary, created_at
		FROM student_learning_snapshots
		WHERE student_id = $1
		ORDER BY snapshot_date DESC
		LIMIT $2`

	rows, err := r.pool.Query(ctx, query, studentID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []LearningSnapshot
	for rows.Next() {
		var s LearningSnapshot
		var rawSum []byte
		if err := rows.Scan(
			&s.ID, &s.StudentID, &s.SnapshotDate, &s.TotalConceptsPractised,
			&s.IndependentSolves, &s.AssistedSolves, &s.MistakeCount,
			&s.CorrectedMistakeCount, &s.ActiveLearningDays, &rawSum, &s.CreatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(rawSum, &s.Summary)
		list = append(list, s)
	}
	return list, nil
}

func (r *PostgresRepository) GetSummary(ctx context.Context, studentID string) (*MemorySummaryResponse, error) {
	query := `
		SELECT
			COALESCE(COUNT(id), 0) as total_concepts,
			COALESCE(SUM(independent_solves), 0) as ind_solves,
			COALESCE(SUM(assisted_solves), 0) as ass_solves,
			COALESCE(SUM(mistake_count), 0) as total_mistakes,
			COALESCE(SUM(corrected_mistakes), 0) as corr_mistakes,
			COALESCE(COUNT(CASE WHEN learning_state = 'mastered' THEN 1 END), 0) as mastered,
			COALESCE(COUNT(CASE WHEN learning_state = 'struggling' THEN 1 END), 0) as struggling
		FROM student_concept_memory
		WHERE student_id = $1`

	resp := &MemorySummaryResponse{
		RecentConcepts:   make([]ConceptMemory, 0),
		ActiveWeaknesses: make([]MistakeMemory, 0),
		KeyStrengths:     make([]ConceptStrength, 0),
	}

	err := r.pool.QueryRow(ctx, query, studentID).Scan(
		&resp.TotalConceptsPractised,
		&resp.IndependentSolves,
		&resp.AssistedSolves,
		&resp.TotalMistakes,
		&resp.CorrectedMistakes,
		&resp.MasteredConcepts,
		&resp.StrugglingConcepts,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	// Active learning days count
	dayQuery := `SELECT COALESCE(COUNT(DISTINCT occurred_at::date), 1) FROM student_learning_events WHERE student_id = $1`
	_ = r.pool.QueryRow(ctx, dayQuery, studentID).Scan(&resp.ActiveLearningDays)

	// Fetch top 5 recent concepts
	recList, _, _ := r.ListConceptMemories(ctx, studentID, ConceptFilter{Limit: 5})
	if recList != nil {
		resp.RecentConcepts = recList
	}

	// Fetch active weaknesses
	weaks, _ := r.ListMistakeMemories(ctx, studentID, "active")
	if weaks != nil {
		if len(weaks) > 5 {
			weaks = weaks[:5]
		}
		resp.ActiveWeaknesses = weaks
	}

	// Fetch key strengths
	strengths, _ := r.GetStrengths(ctx, studentID, 5)
	for _, cm := range strengths {
		resp.KeyStrengths = append(resp.KeyStrengths, ConceptStrength{
			ConceptKey:        cm.ConceptKey,
			ConceptName:       cm.ConceptName,
			TopicName:         cm.TopicName,
			IndependentSolves: cm.IndependentSolves,
			EstimatedMastery:  cm.EstimatedMastery,
			ConfidenceLevel:   cm.ConfidenceLevel,
			Reason:            fmt.Sprintf("Solved %d times independently without assistance.", cm.IndependentSolves),
		})
	}

	return resp, nil
}

func (r *PostgresRepository) GetTimeline(ctx context.Context, studentID string, days int) ([]TimelinePoint, error) {
	if days <= 0 {
		days = 14
	}

	query := `
		SELECT
			d::date::text as day,
			COALESCE(COUNT(CASE WHEN e.event_type = 'submission_accepted' AND (e.event_data->>'is_independent') = 'true' THEN 1 END), 0) as ind_solves,
			COALESCE(COUNT(CASE WHEN e.event_type = 'submission_accepted' AND (e.event_data->>'is_independent') = 'false' THEN 1 END), 0) as ass_solves,
			COALESCE(COUNT(CASE WHEN e.event_type = 'mistake_detected' THEN 1 END), 0) as mistakes,
			COALESCE(COUNT(CASE WHEN e.event_type = 'mistake_corrected' THEN 1 END), 0) as corrected,
			COALESCE(COUNT(CASE WHEN e.event_type = 'concept_practised' THEN 1 END), 0) as concepts
		FROM generate_series(CURRENT_DATE - ($2 * interval '1 day'), CURRENT_DATE, interval '1 day') d
		LEFT JOIN student_learning_events e
			ON e.student_id = $1::uuid AND e.occurred_at::date = d::date
		GROUP BY d
		ORDER BY d ASC`

	rows, err := r.pool.Query(ctx, query, studentID, days-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []TimelinePoint
	for rows.Next() {
		var pt TimelinePoint
		if err := rows.Scan(&pt.Date, &pt.IndependentSolves, &pt.AssistedSolves, &pt.MistakesDetected, &pt.MistakesCorrected, &pt.ConceptsPractised); err != nil {
			return nil, err
		}
		points = append(points, pt)
	}
	return points, nil
}

func formatConceptName(key string) string {
	parts := strings.Split(key, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
