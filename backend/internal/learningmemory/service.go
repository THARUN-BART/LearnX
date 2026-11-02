package learningmemory

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
)

// Service coordinates all business logic for student learning memory.
type Service interface {
	// Ingest incoming signals
	RecordSubmissionEvent(ctx context.Context, studentID string, problem *model.Problem, subID string, accepted bool, failedReason string) error
	RecordHintEvent(ctx context.Context, studentID, problemID string, hintLevel int, hintContent string) error
	RecordMistakeEvent(ctx context.Context, studentID, problemID, submissionID, category, subtype, conceptKey string, confidence float64) error
	RecordWalkthroughEvent(ctx context.Context, studentID, problemID string) error

	// Query Memory
	GetSummary(ctx context.Context, studentID string) (*MemorySummaryResponse, error)
	ListConcepts(ctx context.Context, studentID string, filter ConceptFilter) ([]ConceptMemory, int, error)
	GetConceptDetail(ctx context.Context, studentID, conceptKey string) (*ConceptDetailResponse, error)
	GetStrengths(ctx context.Context, studentID string, limit int) ([]ConceptStrength, error)
	GetWeaknesses(ctx context.Context, studentID string, limit int) ([]ConceptWeakness, error)
	GetHistory(ctx context.Context, studentID string, filter EventFilter) ([]LearningEvent, int, error)
	GetTimeline(ctx context.Context, studentID string, days int) ([]TimelinePoint, error)
	GetSnapshots(ctx context.Context, studentID string, days int) ([]LearningSnapshot, error)
	GetStudentFacingSummary(ctx context.Context, studentID string) (*StudentFacingSummary, error)

	// Context generation for AI Tutor
	BuildAIContext(ctx context.Context, studentID, topicID, problemID string) (string, error)

	// Maintenance
	GenerateDailySnapshot(ctx context.Context, studentID string) error
}

type serviceImpl struct {
	repo         Repository
	problemsRepo problems.Repository
	logger       *slog.Logger
}

func NewService(repo Repository, problemsRepo problems.Repository, logger *slog.Logger) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		repo:         repo,
		problemsRepo: problemsRepo,
		logger:       logger,
	}
}

func (s *serviceImpl) RecordSubmissionEvent(
	ctx context.Context,
	studentID string,
	prob *model.Problem,
	subID string,
	accepted bool,
	failedReason string,
) error {
	now := time.Now()
	var topicID, subtopicID *string
	conceptKey := "general_problem_solving"

	if prob != nil {
		topicID = &prob.TopicID
		subtopicID = &prob.SubtopicID
		if len(prob.Tags) > 0 {
			conceptKey = sanitizeConceptKey(prob.Tags[0])
		} else {
			conceptKey = sanitizeConceptKey(prob.Slug)
		}
	}

	// 1. Determine if solve was independent vs assisted
	isIndependent := true
	if accepted && prob != nil {
		hasHint, err := s.repo.HasRecentHint(ctx, studentID, prob.ID, now)
		if err == nil && hasHint {
			isIndependent = false
		}
	}

	eventType := EventSubmissionFailed
	if accepted {
		eventType = EventSubmissionAccepted
	}

	eventData := map[string]interface{}{
		"is_accepted":    accepted,
		"is_independent": isIndependent,
		"failed_reason":  failedReason,
	}

	probIDStr := ""
	if prob != nil {
		probIDStr = prob.ID
	}

	// Record immutable learning event
	_ = s.repo.RecordEvent(ctx, &LearningEvent{
		StudentID:    studentID,
		EventType:    eventType,
		ProblemID:    &probIDStr,
		SubmissionID: &subID,
		TopicID:      topicID,
		SubtopicID:   subtopicID,
		ConceptKey:   conceptKey,
		EventData:    eventData,
		OccurredAt:   now,
	})

	// 2. Update Concept Memory
	cm, err := s.repo.GetConceptMemory(ctx, studentID, conceptKey)
	if err != nil || cm == nil {
		cm = &ConceptMemory{
			StudentID:        studentID,
			TopicID:          topicID,
			SubtopicID:       subtopicID,
			ConceptKey:       conceptKey,
			LearningState:    StateDiscovered,
			FirstPracticedAt: now,
			LastPracticedAt:  now,
		}
	}

	cm.TotalAttempts++
	cm.LastPracticedAt = now
	if accepted {
		cm.LastSuccessfulAt = &now
		if isIndependent {
			cm.IndependentSolves++
		} else {
			cm.AssistedSolves++
		}
	} else {
		cm.FailedAttempts++
	}

	// Evaluate provisional learning state & estimated mastery
	cm.LearningState, cm.EstimatedMastery, cm.ConfidenceLevel = evaluateConceptState(cm)
	if err := s.repo.UpsertConceptMemory(ctx, cm); err != nil {
		s.logger.Warn("failed to update concept memory", slog.String("error", err.Error()))
	}

	// 3. If accepted, check for mistake corrections
	if accepted {
		s.markMistakesCorrectedForConcept(ctx, studentID, conceptKey)
	}

	return nil
}

func (s *serviceImpl) RecordHintEvent(ctx context.Context, studentID, problemID string, hintLevel int, hintContent string) error {
	now := time.Now()
	eventData := map[string]interface{}{
		"hint_level":   hintLevel,
		"content_len":  len(hintContent),
	}

	var topicID, subtopicID *string
	conceptKey := "general_problem_solving"

	if s.problemsRepo != nil && problemID != "" {
		if p, err := s.problemsRepo.Get(ctx, problemID); err == nil && p != nil {
			topicID = &p.TopicID
			subtopicID = &p.SubtopicID
			if len(p.Tags) > 0 {
				conceptKey = sanitizeConceptKey(p.Tags[0])
			} else {
				conceptKey = sanitizeConceptKey(p.Slug)
			}
		}
	}

	_ = s.repo.RecordEvent(ctx, &LearningEvent{
		StudentID:  studentID,
		EventType:  EventHintRequested,
		ProblemID:  &problemID,
		TopicID:    topicID,
		SubtopicID: subtopicID,
		ConceptKey: conceptKey,
		EventData:  eventData,
		OccurredAt: now,
	})

	// Bump hint requests on concept memory
	cm, err := s.repo.GetConceptMemory(ctx, studentID, conceptKey)
	if err == nil && cm != nil {
		cm.HintRequests++
		cm.LastPracticedAt = now
		cm.LearningState, cm.EstimatedMastery, cm.ConfidenceLevel = evaluateConceptState(cm)
		_ = s.repo.UpsertConceptMemory(ctx, cm)
	}

	return nil
}

func (s *serviceImpl) RecordMistakeEvent(
	ctx context.Context,
	studentID, problemID, submissionID, category, subtype, conceptKey string,
	confidence float64,
) error {
	now := time.Now()
	if conceptKey == "" {
		conceptKey = "general_coding"
	}
	conceptKey = sanitizeConceptKey(conceptKey)

	_ = s.repo.RecordEvent(ctx, &LearningEvent{
		StudentID:    studentID,
		EventType:    EventMistakeDetected,
		ProblemID:    &problemID,
		SubmissionID: &submissionID,
		ConceptKey:   conceptKey,
		EventData: map[string]interface{}{
			"category":   category,
			"subtype":    subtype,
			"confidence": confidence,
		},
		OccurredAt: now,
	})

	// Upsert Mistake Memory
	mm, err := s.repo.GetMistakeMemory(ctx, studentID, category, subtype, conceptKey)
	if err != nil || mm == nil {
		mm = &MistakeMemory{
			StudentID:              studentID,
			Category:               category,
			Subtype:                subtype,
			ConceptKey:             conceptKey,
			OccurrenceCount:        1,
			FirstOccurrenceAt:      now,
			LastOccurrenceAt:       now,
			Status:                 "active",
			SupportingDiagnosisIDs: make([]string, 0),
		}
	} else {
		mm.OccurrenceCount++
		mm.LastOccurrenceAt = now
		mm.Status = "active" // Re-activate if it recurs
	}
	if submissionID != "" {
		mm.SupportingDiagnosisIDs = append(mm.SupportingDiagnosisIDs, submissionID)
	}
	_ = s.repo.UpsertMistakeMemory(ctx, mm)

	// Bump mistake count in Concept Memory
	cm, err := s.repo.GetConceptMemory(ctx, studentID, conceptKey)
	if err == nil && cm != nil {
		cm.MistakeCount++
		cm.LastPracticedAt = now
		cm.LearningState, cm.EstimatedMastery, cm.ConfidenceLevel = evaluateConceptState(cm)
		_ = s.repo.UpsertConceptMemory(ctx, cm)
	}

	return nil
}

func (s *serviceImpl) RecordWalkthroughEvent(ctx context.Context, studentID, problemID string) error {
	return s.repo.RecordEvent(ctx, &LearningEvent{
		StudentID:  studentID,
		EventType:  EventWalkthroughViewed,
		ProblemID:  &problemID,
		EventData:  map[string]interface{}{"action": "walkthrough_viewed"},
		OccurredAt: time.Now(),
	})
}

func (s *serviceImpl) markMistakesCorrectedForConcept(ctx context.Context, studentID, conceptKey string) {
	mistakes, err := s.repo.ListMistakeMemories(ctx, studentID, "active")
	if err != nil || len(mistakes) == 0 {
		return
	}

	now := time.Now()
	for _, m := range mistakes {
		if m.ConceptKey == conceptKey {
			m.CorrectedCount++
			m.LastCorrectedAt = &now
			if m.CorrectedCount >= m.OccurrenceCount {
				m.Status = "resolved"
			} else {
				m.Status = "improving"
			}
			_ = s.repo.UpsertMistakeMemory(ctx, &m)

			// Record event
			_ = s.repo.RecordEvent(ctx, &LearningEvent{
				StudentID:  studentID,
				EventType:  EventMistakeCorrected,
				ConceptKey: conceptKey,
				EventData: map[string]interface{}{
					"category": m.Category,
					"subtype":  m.Subtype,
					"status":   m.Status,
				},
				OccurredAt: now,
			})
		}
	}

	// Update corrected mistakes in Concept Memory
	cm, err := s.repo.GetConceptMemory(ctx, studentID, conceptKey)
	if err == nil && cm != nil {
		cm.CorrectedMistakes++
		cm.LearningState, cm.EstimatedMastery, cm.ConfidenceLevel = evaluateConceptState(cm)
		_ = s.repo.UpsertConceptMemory(ctx, cm)
	}
}

func (s *serviceImpl) GetSummary(ctx context.Context, studentID string) (*MemorySummaryResponse, error) {
	return s.repo.GetSummary(ctx, studentID)
}

func (s *serviceImpl) ListConcepts(ctx context.Context, studentID string, filter ConceptFilter) ([]ConceptMemory, int, error) {
	return s.repo.ListConceptMemories(ctx, studentID, filter)
}

func (s *serviceImpl) GetConceptDetail(ctx context.Context, studentID, conceptKey string) (*ConceptDetailResponse, error) {
	cm, err := s.repo.GetConceptMemory(ctx, studentID, conceptKey)
	if err != nil {
		// Return clean empty state if not yet practiced
		cm = &ConceptMemory{
			StudentID:        studentID,
			ConceptKey:       conceptKey,
			ConceptName:      formatConceptName(conceptKey),
			LearningState:    StateDiscovered,
			FirstPracticedAt: time.Now(),
			LastPracticedAt:  time.Now(),
		}
	}

	allMistakes, _ := s.repo.ListMistakeMemories(ctx, studentID, "")
	var conceptMistakes []MistakeMemory
	for _, m := range allMistakes {
		if m.ConceptKey == conceptKey {
			conceptMistakes = append(conceptMistakes, m)
		}
	}

	events, _, _ := s.repo.ListEvents(ctx, studentID, EventFilter{ConceptKey: conceptKey, Limit: 20})

	nextSteps := []string{
		"Practice an independent verification problem on this concept.",
		"Review boundary condition patterns before submitting.",
		"Ask AI Tutor to test your understanding with a targeted quiz.",
	}
	if cm.IndependentSolves >= 2 && cm.MistakeCount == 0 {
		nextSteps = []string{
			"Concept is well understood! Advance to medium/hard variations.",
			"Try timing your solves to improve competitive execution speed.",
		}
	}

	return &ConceptDetailResponse{
		Concept:   *cm,
		Mistakes:  conceptMistakes,
		Timeline:  events,
		NextSteps: nextSteps,
	}, nil
}

func (s *serviceImpl) GetStrengths(ctx context.Context, studentID string, limit int) ([]ConceptStrength, error) {
	list, err := s.repo.GetStrengths(ctx, studentID, limit)
	if err != nil {
		return nil, err
	}
	var res []ConceptStrength
	for _, cm := range list {
		res = append(res, ConceptStrength{
			ConceptKey:        cm.ConceptKey,
			ConceptName:       cm.ConceptName,
			TopicName:         cm.TopicName,
			IndependentSolves: cm.IndependentSolves,
			EstimatedMastery:  cm.EstimatedMastery,
			ConfidenceLevel:   cm.ConfidenceLevel,
			Reason:            fmt.Sprintf("Solved %d problems independently with high accuracy.", cm.IndependentSolves),
		})
	}
	return res, nil
}

func (s *serviceImpl) GetWeaknesses(ctx context.Context, studentID string, limit int) ([]ConceptWeakness, error) {
	mistakes, err := s.repo.ListMistakeMemories(ctx, studentID, "active")
	if err != nil {
		return nil, err
	}
	if limit > 0 && len(mistakes) > limit {
		mistakes = mistakes[:limit]
	}

	var res []ConceptWeakness
	for _, m := range mistakes {
		res = append(res, ConceptWeakness{
			ConceptKey:      m.ConceptKey,
			ConceptName:     formatConceptName(m.ConceptKey),
			MistakeCategory: m.Category,
			MistakeSubtype:  m.Subtype,
			OccurrenceCount: m.OccurrenceCount,
			CorrectedCount:  m.CorrectedCount,
			Status:          m.Status,
			Reason:          fmt.Sprintf("Encountered %s (%s) %d times. Focus on reinforcing boundary cases.", m.Category, m.Subtype, m.OccurrenceCount),
		})
	}
	return res, nil
}

func (s *serviceImpl) GetHistory(ctx context.Context, studentID string, filter EventFilter) ([]LearningEvent, int, error) {
	return s.repo.ListEvents(ctx, studentID, filter)
}

func (s *serviceImpl) GetTimeline(ctx context.Context, studentID string, days int) ([]TimelinePoint, error) {
	return s.repo.GetTimeline(ctx, studentID, days)
}

func (s *serviceImpl) GetSnapshots(ctx context.Context, studentID string, days int) ([]LearningSnapshot, error) {
	return s.repo.ListSnapshots(ctx, studentID, days)
}

func (s *serviceImpl) GetStudentFacingSummary(ctx context.Context, studentID string) (*StudentFacingSummary, error) {
	sum, err := s.repo.GetSummary(ctx, studentID)
	if err != nil {
		return nil, err
	}

	headline := "Your Learning Journey is Underway"
	narrative := "Keep practicing problems to build cognitive mastery."
	if sum.TotalConceptsPractised > 0 {
		headline = fmt.Sprintf("You have actively practiced %d DSA concepts!", sum.TotalConceptsPractised)
		narrative = fmt.Sprintf("You've recorded %d independent solves and corrected %d mistakes so far.", sum.IndependentSolves, sum.CorrectedMistakes)
	}

	topStrengths := make([]string, 0)
	for _, s := range sum.KeyStrengths {
		topStrengths = append(topStrengths, s.ConceptName)
	}

	areasForFocus := make([]string, 0)
	for _, w := range sum.ActiveWeaknesses {
		areasForFocus = append(areasForFocus, fmt.Sprintf("%s (%s)", formatConceptName(w.ConceptKey), w.Subtype))
	}

	return &StudentFacingSummary{
		Headline:         headline,
		Narrative:        narrative,
		TopStrengths:     topStrengths,
		AreasForFocus:    areasForFocus,
		NextActionPrompt: "Recommended next step: solve a problem targeting your most recent focus area.",
	}, nil
}

func (s *serviceImpl) BuildAIContext(ctx context.Context, studentID, topicID, problemID string) (string, error) {
	if studentID == "" {
		return "", nil
	}

	// Fetch top 3 active weaknesses & top 3 strengths
	weaknesses, _ := s.GetWeaknesses(ctx, studentID, 3)
	strengths, _ := s.GetStrengths(ctx, studentID, 3)

	var sb strings.Builder
	sb.WriteString("STUDENT LEARNING MEMORY CONTEXT:\n")

	if len(strengths) > 0 {
		sb.WriteString("- Demonstrated Strengths: ")
		names := make([]string, 0, len(strengths))
		for _, st := range strengths {
			names = append(names, fmt.Sprintf("%s (%d solves)", st.ConceptName, st.IndependentSolves))
		}
		sb.WriteString(strings.Join(names, ", ") + "\n")
	}

	if len(weaknesses) > 0 {
		sb.WriteString("- Known Recurring Weaknesses: ")
		names := make([]string, 0, len(weaknesses))
		for _, w := range weaknesses {
			names = append(names, fmt.Sprintf("%s [%s: %s]", w.ConceptName, w.MistakeCategory, w.MistakeSubtype))
		}
		sb.WriteString(strings.Join(names, ", ") + "\n")
	}

	sb.WriteString("Guidance: Adapt your teaching to address known weaknesses gently while acknowledging demonstrated strengths. Do not state these metadata points explicitly unless directly asked.")

	return sb.String(), nil
}

func (s *serviceImpl) GenerateDailySnapshot(ctx context.Context, studentID string) error {
	sum, err := s.repo.GetSummary(ctx, studentID)
	if err != nil {
		return err
	}

	snap := &LearningSnapshot{
		StudentID:              studentID,
		SnapshotDate:           time.Now().Format("2006-01-02"),
		TotalConceptsPractised: sum.TotalConceptsPractised,
		IndependentSolves:      sum.IndependentSolves,
		AssistedSolves:         sum.AssistedSolves,
		MistakeCount:           sum.TotalMistakes,
		CorrectedMistakeCount:  sum.CorrectedMistakes,
		ActiveLearningDays:     sum.ActiveLearningDays,
		Summary: map[string]interface{}{
			"mastered_count":   sum.MasteredConcepts,
			"struggling_count": sum.StrugglingConcepts,
		},
	}
	return s.repo.UpsertSnapshot(ctx, snap)
}

func evaluateConceptState(cm *ConceptMemory) (state string, mastery float64, conf string) {
	// Provisional estimate labeled clearly
	totalSolves := cm.IndependentSolves + cm.AssistedSolves
	if totalSolves == 0 && cm.FailedAttempts == 0 {
		return StateDiscovered, 0.0, "low"
	}

	score := float64(cm.IndependentSolves*100 + cm.AssistedSolves*50 - cm.FailedAttempts*20)
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	mastery = math.Round(score*10) / 10

	// Confidence level based on attempts
	if cm.TotalAttempts < 3 {
		conf = "low"
	} else if cm.TotalAttempts < 8 {
		conf = "moderate"
	} else {
		conf = "high"
	}

	// State calculation
	if cm.MistakeCount > cm.CorrectedMistakes && cm.FailedAttempts > cm.IndependentSolves {
		state = StateStruggling
	} else if cm.IndependentSolves >= 3 && cm.MistakeCount <= cm.CorrectedMistakes {
		state = StateMastered
	} else if cm.IndependentSolves > 0 {
		state = StateReinforcing
	} else {
		state = StatePracticing
	}
	return state, mastery, conf
}

func sanitizeConceptKey(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}
