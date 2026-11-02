package aitutor

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrConversationNotFound = errors.New("ai conversation not found")
	ErrMessageNotFound      = errors.New("ai message not found")
)

type Repository interface {
	CreateConversation(ctx context.Context, c *AIConversation) error
	GetConversation(ctx context.Context, id string) (*AIConversation, error)
	ListConversationsByStudent(ctx context.Context, studentID string, limit, offset int) ([]AIConversation, error)
	UpdateConversation(ctx context.Context, c *AIConversation) error

	SaveMessage(ctx context.Context, m *AIMessage) error
	ListMessagesByConversation(ctx context.Context, conversationID string) ([]AIMessage, error)

	SaveHintHistory(ctx context.Context, h *AIHintHistory) error
	ListHintHistory(ctx context.Context, studentID, problemID string) ([]AIHintHistory, error)
	GetLatestHintLevel(ctx context.Context, studentID, problemID string) (int, error)
}

// PostgresRepository persists AI tutor conversations and hint progress in Postgres.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateConversation(ctx context.Context, c *AIConversation) error {
	if c.ID == "" {
		c.ID = crypto.GenerateUUID()
	}
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	query := `
		INSERT INTO ai_conversations (id, student_id, problem_id, conversation_type, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := r.pool.Exec(ctx, query, c.ID, c.StudentID, c.ProblemID, c.ConversationType, c.Title, c.CreatedAt, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) GetConversation(ctx context.Context, id string) (*AIConversation, error) {
	query := `
		SELECT id, student_id, problem_id, conversation_type, title, created_at, updated_at
		FROM ai_conversations
		WHERE id = $1;
	`
	var c AIConversation
	err := r.pool.QueryRow(ctx, query, id).Scan(&c.ID, &c.StudentID, &c.ProblemID, &c.ConversationType, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrConversationNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PostgresRepository) ListConversationsByStudent(ctx context.Context, studentID string, limit, offset int) ([]AIConversation, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
		SELECT id, student_id, problem_id, conversation_type, title, created_at, updated_at
		FROM ai_conversations
		WHERE student_id = $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3;
	`
	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []AIConversation
	for rows.Next() {
		var c AIConversation
		if err := rows.Scan(&c.ID, &c.StudentID, &c.ProblemID, &c.ConversationType, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}

func (r *PostgresRepository) UpdateConversation(ctx context.Context, c *AIConversation) error {
	c.UpdatedAt = time.Now()
	query := `
		UPDATE ai_conversations
		SET title = $2, updated_at = $3
		WHERE id = $1;
	`
	_, err := r.pool.Exec(ctx, query, c.ID, c.Title, c.UpdatedAt)
	return err
}

func (r *PostgresRepository) SaveMessage(ctx context.Context, m *AIMessage) error {
	if m.ID == "" {
		m.ID = crypto.GenerateUUID()
	}
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	metaJSON, _ := json.Marshal(m.Metadata)
	tokensJSON, _ := json.Marshal(m.TokenUsage)

	query := `
		INSERT INTO ai_messages (id, conversation_id, role, message, response_type, metadata, token_usage, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := r.pool.Exec(ctx, query, m.ID, m.ConversationID, m.Role, m.Message, m.ResponseType, metaJSON, tokensJSON, m.CreatedAt)
	return err
}

func (r *PostgresRepository) ListMessagesByConversation(ctx context.Context, conversationID string) ([]AIMessage, error) {
	query := `
		SELECT id, conversation_id, role, message, response_type, metadata, token_usage, created_at
		FROM ai_messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC;
	`
	rows, err := r.pool.Query(ctx, query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []AIMessage
	for rows.Next() {
		var m AIMessage
		var metaBytes, tokenBytes []byte
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Message, &m.ResponseType, &metaBytes, &tokenBytes, &m.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(metaBytes, &m.Metadata)
		_ = json.Unmarshal(tokenBytes, &m.TokenUsage)
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (r *PostgresRepository) SaveHintHistory(ctx context.Context, h *AIHintHistory) error {
	if h.ID == "" {
		h.ID = crypto.GenerateUUID()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO ai_hint_history (id, student_id, problem_id, submission_id, hint_level, hint_content, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (student_id, problem_id, hint_level) DO UPDATE SET
			hint_content = EXCLUDED.hint_content,
			created_at = EXCLUDED.created_at;
	`
	_, err := r.pool.Exec(ctx, query, h.ID, h.StudentID, h.ProblemID, h.SubmissionID, h.HintLevel, h.HintContent, h.CreatedAt)
	return err
}

func (r *PostgresRepository) ListHintHistory(ctx context.Context, studentID, problemID string) ([]AIHintHistory, error) {
	query := `
		SELECT id, student_id, problem_id, submission_id, hint_level, hint_content, created_at
		FROM ai_hint_history
		WHERE student_id = $1 AND problem_id = $2
		ORDER BY hint_level ASC;
	`
	rows, err := r.pool.Query(ctx, query, studentID, problemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hints []AIHintHistory
	for rows.Next() {
		var h AIHintHistory
		if err := rows.Scan(&h.ID, &h.StudentID, &h.ProblemID, &h.SubmissionID, &h.HintLevel, &h.HintContent, &h.CreatedAt); err != nil {
			return nil, err
		}
		hints = append(hints, h)
	}
	return hints, rows.Err()
}

func (r *PostgresRepository) GetLatestHintLevel(ctx context.Context, studentID, problemID string) (int, error) {
	query := `
		SELECT COALESCE(MAX(hint_level), 0)
		FROM ai_hint_history
		WHERE student_id = $1 AND problem_id = $2;
	`
	var maxLevel int
	err := r.pool.QueryRow(ctx, query, studentID, problemID).Scan(&maxLevel)
	return maxLevel, err
}

// MemoryRepository implements Repository in-memory for testing.
type MemoryRepository struct {
	mu            sync.RWMutex
	conversations map[string]*AIConversation
	messages      map[string][]*AIMessage
	hintHistory   map[string][]*AIHintHistory
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		conversations: make(map[string]*AIConversation),
		messages:      make(map[string][]*AIMessage),
		hintHistory:   make(map[string][]*AIHintHistory),
	}
}

func (m *MemoryRepository) CreateConversation(ctx context.Context, c *AIConversation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID == "" {
		c.ID = crypto.GenerateUUID()
	}
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now
	m.conversations[c.ID] = c
	return nil
}

func (m *MemoryRepository) GetConversation(ctx context.Context, id string) (*AIConversation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.conversations[id]
	if !ok {
		return nil, ErrConversationNotFound
	}
	return c, nil
}

func (m *MemoryRepository) ListConversationsByStudent(ctx context.Context, studentID string, limit, offset int) ([]AIConversation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var list []AIConversation
	for _, c := range m.conversations {
		if c.StudentID == studentID {
			list = append(list, *c)
		}
	}
	if offset > len(list) {
		return []AIConversation{}, nil
	}
	list = list[offset:]
	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	return list, nil
}

func (m *MemoryRepository) UpdateConversation(ctx context.Context, c *AIConversation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	c.UpdatedAt = time.Now()
	m.conversations[c.ID] = c
	return nil
}

func (m *MemoryRepository) SaveMessage(ctx context.Context, msg *AIMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if msg.ID == "" {
		msg.ID = crypto.GenerateUUID()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}
	m.messages[msg.ConversationID] = append(m.messages[msg.ConversationID], msg)
	return nil
}

func (m *MemoryRepository) ListMessagesByConversation(ctx context.Context, conversationID string) ([]AIMessage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	msgs, ok := m.messages[conversationID]
	if !ok {
		return []AIMessage{}, nil
	}
	res := make([]AIMessage, len(msgs))
	for i, mItem := range msgs {
		res[i] = *mItem
	}
	return res, nil
}

func (m *MemoryRepository) SaveHintHistory(ctx context.Context, h *AIHintHistory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if h.ID == "" {
		h.ID = crypto.GenerateUUID()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}
	key := h.StudentID + ":" + h.ProblemID
	list := m.hintHistory[key]
	// Replace if same hint level
	replaced := false
	for i, existing := range list {
		if existing.HintLevel == h.HintLevel {
			list[i] = h
			replaced = true
			break
		}
	}
	if !replaced {
		list = append(list, h)
	}
	m.hintHistory[key] = list
	return nil
}

func (m *MemoryRepository) ListHintHistory(ctx context.Context, studentID, problemID string) ([]AIHintHistory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := studentID + ":" + problemID
	list, ok := m.hintHistory[key]
	if !ok {
		return []AIHintHistory{}, nil
	}
	res := make([]AIHintHistory, len(list))
	for i, h := range list {
		res[i] = *h
	}
	return res, nil
}

func (m *MemoryRepository) GetLatestHintLevel(ctx context.Context, studentID, problemID string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := studentID + ":" + problemID
	list := m.hintHistory[key]
	maxLevel := 0
	for _, h := range list {
		if h.HintLevel > maxLevel {
			maxLevel = h.HintLevel
		}
	}
	return maxLevel, nil
}
