package model

import (
	"crypto/rand"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type Memory struct {
	mu        sync.RWMutex
	Roadmaps  map[string]*Roadmap
	Topics    map[string]*Topic
	Subtopics map[string]*Subtopic
	Problems  map[string]*Problem
	Progress  map[string]*Progress
}

func (m *Memory) Lock()   { m.mu.Lock() }
func (m *Memory) Unlock() { m.mu.Unlock() }

func NewMemory() *Memory {
	return &Memory{
		Roadmaps:  make(map[string]*Roadmap),
		Topics:    make(map[string]*Topic),
		Subtopics: make(map[string]*Subtopic),
		Problems:  make(map[string]*Problem),
		Progress:  make(map[string]*Progress),
	}
}

func NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func progressKey(studentID, problemID string) string {
	return studentID + "::" + problemID
}

func cloneStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func (m *Memory) ListRoadmaps() []Roadmap {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Roadmap, 0, len(m.Roadmaps))
	for _, r := range m.Roadmaps {
		if !r.IsActive {
			continue
		}
		item := *r
		item.TopicCount = m.topicCountLocked(r.ID)
		item.ProblemCount = m.problemCountByRoadmapLocked(r.ID)
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (m *Memory) GetRoadmap(id string) (*Roadmap, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.Roadmaps[id]
	if !ok {
		for _, item := range m.Roadmaps {
			if item.Slug == id {
				r = item
				ok = true
				break
			}
		}
	}
	if !ok {
		return nil, false
	}
	cp := *r
	cp.TopicCount = m.topicCountLocked(r.ID)
	cp.ProblemCount = m.problemCountByRoadmapLocked(r.ID)
	return &cp, true
}

func (m *Memory) ListTopics(roadmapID string) []Topic {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Topic, 0)
	for _, t := range m.Topics {
		if t.RoadmapID != roadmapID {
			continue
		}
		item := *t
		item.SubtopicCount = m.subtopicCountLocked(t.ID)
		item.ProblemCount = m.problemCountByTopicLocked(t.ID)
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TopicOrder < out[j].TopicOrder })
	return out
}

func (m *Memory) GetTopic(id string) (*Topic, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.Topics[id]
	if !ok {
		return nil, false
	}
	cp := *t
	cp.SubtopicCount = m.subtopicCountLocked(t.ID)
	cp.ProblemCount = m.problemCountByTopicLocked(t.ID)
	return &cp, true
}

func (m *Memory) ListSubtopics(topicID string) []Subtopic {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Subtopic, 0)
	for _, s := range m.Subtopics {
		if s.TopicID != topicID {
			continue
		}
		item := *s
		item.ProblemCount = m.problemCountBySubtopicLocked(s.ID)
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SubtopicOrder < out[j].SubtopicOrder })
	return out
}

func (m *Memory) GetSubtopic(id string) (*Subtopic, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.Subtopics[id]
	if !ok {
		return nil, false
	}
	cp := *s
	cp.ProblemCount = m.problemCountBySubtopicLocked(s.ID)
	return &cp, true
}

func (m *Memory) GetProblem(id string) (*Problem, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.Problems[id]
	if !ok {
		for _, item := range m.Problems {
			if item.Slug == id {
				p = item
				ok = true
				break
			}
		}
	}
	if !ok {
		return nil, false
	}
	cp := *p
	cp.Tags = cloneStrings(p.Tags)
	if t, ok := m.Topics[p.TopicID]; ok {
		cp.TopicName = t.Name
	}
	if s, ok := m.Subtopics[p.SubtopicID]; ok {
		cp.SubtopicName = s.Name
	}
	return &cp, true
}

func (m *Memory) ListProblems(filter ProblemFilter) ([]Problem, int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	matched := make([]Problem, 0)
	q := strings.ToLower(strings.TrimSpace(filter.Query))
	for _, p := range m.Problems {
		if !p.IsActive {
			continue
		}
		if filter.TopicID != "" && p.TopicID != filter.TopicID {
			continue
		}
		if filter.SubtopicID != "" && p.SubtopicID != filter.SubtopicID {
			continue
		}
		if filter.RoadmapID != "" {
			topic, ok := m.Topics[p.TopicID]
			if !ok || topic.RoadmapID != filter.RoadmapID {
				continue
			}
		}
		if filter.Difficulty != "" && !strings.EqualFold(p.Difficulty, filter.Difficulty) {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(p.Title), q) && !strings.Contains(strings.ToLower(p.Slug), q) {
			continue
		}

		cp := *p
		cp.Tags = cloneStrings(p.Tags)
		if t, ok := m.Topics[p.TopicID]; ok {
			cp.TopicName = t.Name
		}
		if s, ok := m.Subtopics[p.SubtopicID]; ok {
			cp.SubtopicName = s.Name
		}
		if filter.StudentID != "" {
			if rec, ok := m.Progress[progressKey(filter.StudentID, p.ID)]; ok {
				cp.ProgressStatus = rec.Status
			} else {
				cp.ProgressStatus = StatusNotStarted
			}
			if filter.Status != "" && !strings.EqualFold(cp.ProgressStatus, filter.Status) {
				continue
			}
		} else if filter.Status != "" {
			continue
		}
		matched = append(matched, cp)
	}

	sort.Slice(matched, func(i, j int) bool {
		if filter.Sort == "title" {
			return strings.ToLower(matched[i].Title) < strings.ToLower(matched[j].Title)
		}
		if matched[i].ProblemOrder == matched[j].ProblemOrder {
			return matched[i].Title < matched[j].Title
		}
		return matched[i].ProblemOrder < matched[j].ProblemOrder
	})

	total := len(matched)
	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = total
		if limit == 0 {
			return matched, total
		}
	}
	start := (page - 1) * limit
	if start >= total {
		return []Problem{}, total
	}
	end := start + limit
	if end > total {
		end = total
	}
	return matched[start:end], total
}

func (m *Memory) RelatedProblems(problemID string, limit int) []Problem {
	p, ok := m.GetProblem(problemID)
	if !ok {
		return nil
	}
	items, _ := m.ListProblems(ProblemFilter{SubtopicID: p.SubtopicID, Limit: 0})
	out := make([]Problem, 0)
	for _, item := range items {
		if item.ID == p.ID {
			continue
		}
		out = append(out, item)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func (m *Memory) GetProgress(studentID, problemID string) (*Progress, bool) {
	problem, ok := m.GetProblem(problemID)
	if !ok {
		return nil, false
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.Progress[progressKey(studentID, problem.ID)]
	if !ok {
		return nil, false
	}
	cp := *p
	return &cp, true
}

func (m *Memory) ListProgress(studentID string) []Progress {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Progress, 0)
	for _, p := range m.Progress {
		if p.StudentID != studentID {
			continue
		}
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out
}

func (m *Memory) MarkAttempted(studentID, problemID string) (*Progress, error) {
	problem, ok := m.GetProblem(problemID)
	if !ok {
		return nil, ErrNotFound
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	key := progressKey(studentID, problem.ID)
	problemID = problem.ID
	existing, ok := m.Progress[key]
	if !ok {
		rec := &Progress{
			ID:             NewID(),
			StudentID:      studentID,
			ProblemID:      problemID,
			Status:         StatusAttempted,
			TotalAttempts:  1,
			FirstAttemptAt: &now,
			LastAttemptAt:  &now,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		m.Progress[key] = rec
		cp := *rec
		return &cp, nil
	}
	existing.TotalAttempts++
	existing.LastAttemptAt = &now
	if existing.FirstAttemptAt == nil {
		existing.FirstAttemptAt = &now
	}
	if existing.Status == StatusNotStarted {
		existing.Status = StatusAttempted
	}
	existing.UpdatedAt = now
	cp := *existing
	return &cp, nil
}

func (m *Memory) MarkSolved(studentID, problemID string) (*Progress, error) {
	problem, ok := m.GetProblem(problemID)
	if !ok {
		return nil, ErrNotFound
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	key := progressKey(studentID, problem.ID)
	problemID = problem.ID
	existing, ok := m.Progress[key]
	if !ok {
		rec := &Progress{
			ID:                 NewID(),
			StudentID:          studentID,
			ProblemID:          problemID,
			Status:             StatusSolved,
			TotalAttempts:      1,
			SuccessfulAttempts: 1,
			FirstAttemptAt:     &now,
			LastAttemptAt:      &now,
			SolvedAt:           &now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		m.Progress[key] = rec
		cp := *rec
		return &cp, nil
	}
	existing.TotalAttempts++
	existing.SuccessfulAttempts++
	existing.Status = StatusSolved
	existing.LastAttemptAt = &now
	if existing.FirstAttemptAt == nil {
		existing.FirstAttemptAt = &now
	}
	if existing.SolvedAt == nil {
		existing.SolvedAt = &now
	}
	existing.UpdatedAt = now
	cp := *existing
	return &cp, nil
}

func (m *Memory) UpdateProgress(studentID, problemID, status string, reset bool) (*Progress, error) {
	problem, ok := m.GetProblem(problemID)
	if !ok {
		return nil, ErrNotFound
	}
	problemID = problem.ID
	if !ValidStatus(status) || status == "" {
		return nil, ErrInvalidInput
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	key := progressKey(studentID, problemID)
	existing, ok := m.Progress[key]
	if !ok {
		existing = &Progress{
			ID:        NewID(),
			StudentID: studentID,
			ProblemID: problemID,
			CreatedAt: now,
		}
		m.Progress[key] = existing
	}
	if reset || status == StatusNotStarted {
		existing.Status = StatusNotStarted
		existing.TotalAttempts = 0
		existing.SuccessfulAttempts = 0
		existing.FirstAttemptAt = nil
		existing.LastAttemptAt = nil
		existing.SolvedAt = nil
		existing.UpdatedAt = now
		cp := *existing
		return &cp, nil
	}
	existing.Status = status
	existing.UpdatedAt = now
	if status == StatusSolved && existing.SolvedAt == nil {
		existing.SolvedAt = &now
		if existing.SuccessfulAttempts == 0 {
			existing.SuccessfulAttempts = 1
		}
		if existing.TotalAttempts == 0 {
			existing.TotalAttempts = 1
		}
	}
	if existing.FirstAttemptAt == nil && status != StatusNotStarted {
		existing.FirstAttemptAt = &now
		existing.LastAttemptAt = &now
	}
	cp := *existing
	return &cp, nil
}

func (m *Memory) topicCountLocked(roadmapID string) int {
	n := 0
	for _, t := range m.Topics {
		if t.RoadmapID == roadmapID {
			n++
		}
	}
	return n
}

func (m *Memory) subtopicCountLocked(topicID string) int {
	n := 0
	for _, s := range m.Subtopics {
		if s.TopicID == topicID {
			n++
		}
	}
	return n
}

func (m *Memory) problemCountByRoadmapLocked(roadmapID string) int {
	n := 0
	for _, p := range m.Problems {
		if !p.IsActive {
			continue
		}
		t, ok := m.Topics[p.TopicID]
		if ok && t.RoadmapID == roadmapID {
			n++
		}
	}
	return n
}

func (m *Memory) problemCountByTopicLocked(topicID string) int {
	n := 0
	for _, p := range m.Problems {
		if p.IsActive && p.TopicID == topicID {
			n++
		}
	}
	return n
}

func (m *Memory) problemCountBySubtopicLocked(subtopicID string) int {
	n := 0
	for _, p := range m.Problems {
		if p.IsActive && p.SubtopicID == subtopicID {
			n++
		}
	}
	return n
}
