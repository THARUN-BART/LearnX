package learningpath

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
)

var ErrNotFound = errors.New("record not found")

func newID() string {
	tok, err := crypto.GenerateRandomToken(16)
	if err != nil {
		return fmt.Sprintf("id-%d", time.Now().UnixNano())
	}
	return tok
}

// MemoryRepository provides a concurrent in-memory mock repository for tests & development.
type MemoryRepository struct {
	mu              sync.RWMutex
	paths           map[string]*AdaptivePath             // key: studentID
	nodes           map[string]map[string]*PathNode      // key: studentID -> conceptKey
	recommendations map[string][]*LearningRecommendation // key: studentID
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		paths:           make(map[string]*AdaptivePath),
		nodes:           make(map[string]map[string]*PathNode),
		recommendations: make(map[string][]*LearningRecommendation),
	}
}

func (m *MemoryRepository) GetAdaptivePath(ctx context.Context, studentID string) (*AdaptivePath, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.paths[studentID]
	if !ok {
		return nil, ErrNotFound
	}
	cp := *p
	return &cp, nil
}

func (m *MemoryRepository) UpsertAdaptivePath(ctx context.Context, path *AdaptivePath) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if path.ID == "" {
		path.ID = newID()
	}
	if path.CreatedAt.IsZero() {
		path.CreatedAt = time.Now()
	}
	path.UpdatedAt = time.Now()

	cp := *path
	m.paths[path.StudentID] = &cp
	return nil
}

func (m *MemoryRepository) GetPathNodes(ctx context.Context, studentID string) ([]PathNode, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	studentNodes, ok := m.nodes[studentID]
	if !ok {
		return make([]PathNode, 0), nil
	}

	res := make([]PathNode, 0, len(studentNodes))
	for _, n := range studentNodes {
		res = append(res, *n)
	}
	return res, nil
}

func (m *MemoryRepository) UpsertPathNode(ctx context.Context, node *PathNode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node.ID == "" {
		node.ID = newID()
	}
	node.UpdatedAt = time.Now()

	if _, ok := m.nodes[node.StudentID]; !ok {
		m.nodes[node.StudentID] = make(map[string]*PathNode)
	}
	cp := *node
	m.nodes[node.StudentID][node.ConceptKey] = &cp
	return nil
}

func (m *MemoryRepository) BatchUpsertPathNodes(ctx context.Context, studentID string, nodes []PathNode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.nodes[studentID]; !ok {
		m.nodes[studentID] = make(map[string]*PathNode)
	}
	for i := range nodes {
		n := nodes[i]
		if n.ID == "" {
			n.ID = newID()
		}
		n.UpdatedAt = time.Now()
		cp := n
		m.nodes[studentID][n.ConceptKey] = &cp
	}
	return nil
}

func (m *MemoryRepository) GetActiveRecommendation(ctx context.Context, studentID string) (*LearningRecommendation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	recs := m.recommendations[studentID]
	for i := len(recs) - 1; i >= 0; i-- {
		if recs[i].IsActive {
			cp := *recs[i]
			return &cp, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MemoryRepository) CreateRecommendation(ctx context.Context, rec *LearningRecommendation) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if rec.ID == "" {
		rec.ID = newID()
	}
	if rec.CreatedAt.IsZero() {
		rec.CreatedAt = time.Now()
	}

	// Deactivate any currently active recommendations
	for _, r := range m.recommendations[rec.StudentID] {
		r.IsActive = false
	}

	cp := *rec
	m.recommendations[rec.StudentID] = append(m.recommendations[rec.StudentID], &cp)
	return nil
}

func (m *MemoryRepository) CompleteRecommendation(ctx context.Context, studentID, recID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for _, r := range m.recommendations[studentID] {
		if r.ID == recID {
			r.IsActive = false
			r.CompletedAt = &now
			return nil
		}
	}
	return nil
}

func (m *MemoryRepository) ListRecommendationHistory(ctx context.Context, studentID string, limit, offset int) ([]LearningRecommendation, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	recs := m.recommendations[studentID]
	total := len(recs)

	var reversed []LearningRecommendation
	for i := len(recs) - 1; i >= 0; i-- {
		reversed = append(reversed, *recs[i])
	}

	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if limit <= 0 || end > total {
		end = total
	}

	return reversed[start:end], total, nil
}

var _ Repository = (*MemoryRepository)(nil)
