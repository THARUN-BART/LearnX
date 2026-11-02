package learningmemory

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

// MemoryRepository is an in-memory thread-safe implementation of Repository for tests & local dev.
type MemoryRepository struct {
	mu        sync.RWMutex
	concepts  map[string]*ConceptMemory   // key: studentID + ":" + conceptKey
	events    []*LearningEvent            // append-only
	mistakes  map[string]*MistakeMemory   // key: studentID + ":" + category + ":" + subtype + ":" + conceptKey
	snapshots map[string]*LearningSnapshot // key: studentID + ":" + date
}

// NewMemoryRepository constructs a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		concepts:  make(map[string]*ConceptMemory),
		events:    make([]*LearningEvent, 0),
		mistakes:  make(map[string]*MistakeMemory),
		snapshots: make(map[string]*LearningSnapshot),
	}
}

func (m *MemoryRepository) GetConceptMemory(ctx context.Context, studentID, conceptKey string) (*ConceptMemory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	k := studentID + ":" + conceptKey
	if cm, ok := m.concepts[k]; ok {
		cp := *cm
		return &cp, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryRepository) UpsertConceptMemory(ctx context.Context, cm *ConceptMemory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if cm.ID == "" {
		cm.ID = newID()
	}
	if cm.CreatedAt.IsZero() {
		cm.CreatedAt = time.Now()
	}
	cm.UpdatedAt = time.Now()

	k := cm.StudentID + ":" + cm.ConceptKey
	cp := *cm
	m.concepts[k] = &cp
	return nil
}

func (m *MemoryRepository) ListConceptMemories(ctx context.Context, studentID string, filter ConceptFilter) ([]ConceptMemory, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var matched []ConceptMemory
	for _, cm := range m.concepts {
		if cm.StudentID != studentID {
			continue
		}
		if filter.TopicID != "" && (cm.TopicID == nil || *cm.TopicID != filter.TopicID) {
			continue
		}
		if filter.LearningState != "" && cm.LearningState != filter.LearningState {
			continue
		}
		if filter.Search != "" {
			s := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(cm.ConceptKey), s) && !strings.Contains(strings.ToLower(cm.ConceptName), s) {
				continue
			}
		}
		matched = append(matched, *cm)
	}

	total := len(matched)
	// Apply pagination
	start := filter.Offset
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if filter.Limit <= 0 || end > total {
		end = total
	}

	return matched[start:end], total, nil
}

func (m *MemoryRepository) GetStrengths(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []ConceptMemory
	for _, cm := range m.concepts {
		if cm.StudentID == studentID && (cm.IndependentSolves > 0 || cm.LearningState == StateMastered) {
			list = append(list, *cm)
		}
	}
	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	return list, nil
}

func (m *MemoryRepository) GetWeaknesses(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []ConceptMemory
	for _, cm := range m.concepts {
		if cm.StudentID == studentID && (cm.MistakeCount > cm.CorrectedMistakes || cm.LearningState == StateStruggling) {
			list = append(list, *cm)
		}
	}
	if limit > 0 && len(list) > limit {
		list = list[:limit]
	}
	return list, nil
}

func (m *MemoryRepository) RecordEvent(ctx context.Context, event *LearningEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if event.ID == "" {
		event.ID = newID()
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	cp := *event
	m.events = append(m.events, &cp)
	return nil
}

func (m *MemoryRepository) ListEvents(ctx context.Context, studentID string, filter EventFilter) ([]LearningEvent, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var matched []LearningEvent
	for i := len(m.events) - 1; i >= 0; i-- {
		e := m.events[i]
		if e.StudentID != studentID {
			continue
		}
		if filter.EventType != "" && e.EventType != filter.EventType {
			continue
		}
		if filter.ConceptKey != "" && e.ConceptKey != filter.ConceptKey {
			continue
		}
		if filter.ProblemID != "" && (e.ProblemID == nil || *e.ProblemID != filter.ProblemID) {
			continue
		}
		if filter.StartDate != nil && e.OccurredAt.Before(*filter.StartDate) {
			continue
		}
		if filter.EndDate != nil && e.OccurredAt.After(*filter.EndDate) {
			continue
		}
		matched = append(matched, *e)
	}

	total := len(matched)
	start := filter.Offset
	if start > total {
		start = total
	}
	end := start + filter.Limit
	if filter.Limit <= 0 || end > total {
		end = total
	}

	return matched[start:end], total, nil
}

func (m *MemoryRepository) HasRecentHint(ctx context.Context, studentID, problemID string, before time.Time) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Look back 2 hours before the given time
	threshold := before.Add(-2 * time.Hour)
	for _, e := range m.events {
		if e.StudentID == studentID && e.EventType == EventHintRequested && e.ProblemID != nil && *e.ProblemID == problemID {
			if e.OccurredAt.After(threshold) && e.OccurredAt.Before(before.Add(1*time.Minute)) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (m *MemoryRepository) GetMistakeMemory(ctx context.Context, studentID, category, subtype, conceptKey string) (*MistakeMemory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	k := fmt.Sprintf("%s:%s:%s:%s", studentID, category, subtype, conceptKey)
	if mm, ok := m.mistakes[k]; ok {
		cp := *mm
		return &cp, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryRepository) UpsertMistakeMemory(ctx context.Context, mm *MistakeMemory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if mm.ID == "" {
		mm.ID = newID()
	}
	if mm.CreatedAt.IsZero() {
		mm.CreatedAt = time.Now()
	}
	mm.UpdatedAt = time.Now()

	k := fmt.Sprintf("%s:%s:%s:%s", mm.StudentID, mm.Category, mm.Subtype, mm.ConceptKey)
	cp := *mm
	m.mistakes[k] = &cp
	return nil
}

func (m *MemoryRepository) ListMistakeMemories(ctx context.Context, studentID string, status string) ([]MistakeMemory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []MistakeMemory
	for _, mm := range m.mistakes {
		if mm.StudentID == studentID {
			if status != "" && mm.Status != status {
				continue
			}
			list = append(list, *mm)
		}
	}
	return list, nil
}

func (m *MemoryRepository) UpsertSnapshot(ctx context.Context, snap *LearningSnapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if snap.ID == "" {
		snap.ID = newID()
	}
	if snap.CreatedAt.IsZero() {
		snap.CreatedAt = time.Now()
	}

	k := snap.StudentID + ":" + snap.SnapshotDate
	cp := *snap
	m.snapshots[k] = &cp
	return nil
}

func (m *MemoryRepository) ListSnapshots(ctx context.Context, studentID string, days int) ([]LearningSnapshot, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var list []LearningSnapshot
	for _, s := range m.snapshots {
		if s.StudentID == studentID {
			list = append(list, *s)
		}
	}
	return list, nil
}

func (m *MemoryRepository) GetSummary(ctx context.Context, studentID string) (*MemorySummaryResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	resp := &MemorySummaryResponse{
		RecentConcepts:   make([]ConceptMemory, 0),
		ActiveWeaknesses: make([]MistakeMemory, 0),
		KeyStrengths:     make([]ConceptStrength, 0),
	}

	for _, cm := range m.concepts {
		if cm.StudentID != studentID {
			continue
		}
		resp.TotalConceptsPractised++
		resp.IndependentSolves += cm.IndependentSolves
		resp.AssistedSolves += cm.AssistedSolves
		resp.TotalMistakes += cm.MistakeCount
		resp.CorrectedMistakes += cm.CorrectedMistakes

		if cm.LearningState == StateMastered {
			resp.MasteredConcepts++
		}
		if cm.LearningState == StateStruggling {
			resp.StrugglingConcepts++
		}

		if cm.IndependentSolves >= 2 && cm.FailedAttempts == 0 {
			resp.KeyStrengths = append(resp.KeyStrengths, ConceptStrength{
				ConceptKey:        cm.ConceptKey,
				ConceptName:       cm.ConceptName,
				IndependentSolves: cm.IndependentSolves,
				EstimatedMastery:  cm.EstimatedMastery,
				ConfidenceLevel:   cm.ConfidenceLevel,
				Reason:            fmt.Sprintf("Solved %d times independently without assistance.", cm.IndependentSolves),
			})
		}
		resp.RecentConcepts = append(resp.RecentConcepts, *cm)
	}

	for _, mm := range m.mistakes {
		if mm.StudentID == studentID && mm.Status == "active" {
			resp.ActiveWeaknesses = append(resp.ActiveWeaknesses, *mm)
		}
	}

	resp.ActiveLearningDays = 1
	if len(resp.RecentConcepts) > 5 {
		resp.RecentConcepts = resp.RecentConcepts[:5]
	}

	return resp, nil
}

func (m *MemoryRepository) GetTimeline(ctx context.Context, studentID string, days int) ([]TimelinePoint, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if days <= 0 {
		days = 14
	}

	// Bucket by date (YYYY-MM-DD)
	buckets := make(map[string]*TimelinePoint)
	now := time.Now()
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i).Format("2006-01-02")
		buckets[d] = &TimelinePoint{Date: d}
	}

	for _, e := range m.events {
		if e.StudentID != studentID {
			continue
		}
		d := e.OccurredAt.Format("2006-01-02")
		pt, exists := buckets[d]
		if !exists {
			continue
		}
		switch e.EventType {
		case EventSubmissionAccepted:
			if isInd, _ := e.EventData["is_independent"].(bool); isInd {
				pt.IndependentSolves++
			} else {
				pt.AssistedSolves++
			}
		case EventMistakeDetected:
			pt.MistakesDetected++
		case EventMistakeCorrected:
			pt.MistakesCorrected++
		case EventConceptPractised:
			pt.ConceptsPractised++
		}
	}

	res := make([]TimelinePoint, 0, len(buckets))
	for i := days - 1; i >= 0; i-- {
		d := now.AddDate(0, 0, -i).Format("2006-01-02")
		if pt, ok := buckets[d]; ok {
			res = append(res, *pt)
		}
	}
	return res, nil
}

// Compile check
var _ Repository = (*MemoryRepository)(nil)
