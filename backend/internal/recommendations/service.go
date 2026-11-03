package recommendations

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/learningmemory"
	"github.com/THARUN-BART/LearnX/backend/internal/learningpath"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
)

// Service coordinates problem recommendations for students.
type Service interface {
	GetCategorizedRecommendations(ctx context.Context, studentID string) (*CategorizedRecommendations, error)
	GetNextRecommendation(ctx context.Context, studentID string) (*ProblemRecommendation, error)
	GetWeaknessRecommendations(ctx context.Context, studentID string) ([]ProblemRecommendation, error)
	GetSimilarRecommendations(ctx context.Context, studentID, problemID string) ([]ProblemRecommendation, error)
	RefreshRecommendations(ctx context.Context, studentID string) (*CategorizedRecommendations, error)
	SubmitFeedback(ctx context.Context, studentID, recommendationID string, req FeedbackRequest) error
	GetHistory(ctx context.Context, studentID string, limit, offset int) (*HistoryResponse, error)
	OnSubmissionEvaluated(ctx context.Context, studentID, problemID, status string)
}

type serviceImpl struct {
	recRepo           Repository
	problemsRepo      problems.Repository
	progressRepo      progress.Repository
	learningMemorySvc learningmemory.Service
	mistakesRepo      mistakes.Repository
	learningPathSvc   learningpath.Service
	logger            *slog.Logger
}

func NewService(
	recRepo Repository,
	problemsRepo problems.Repository,
	progressRepo progress.Repository,
	learningMemorySvc learningmemory.Service,
	mistakesRepo mistakes.Repository,
	learningPathSvc learningpath.Service,
	logger *slog.Logger,
) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		recRepo:           recRepo,
		problemsRepo:      problemsRepo,
		progressRepo:      progressRepo,
		learningMemorySvc: learningMemorySvc,
		mistakesRepo:      mistakesRepo,
		learningPathSvc:   learningPathSvc,
		logger:            logger,
	}
}

func (s *serviceImpl) GetCategorizedRecommendations(ctx context.Context, studentID string) (*CategorizedRecommendations, error) {
	// 1. Check if fresh active recommendations already exist
	cached, err := s.recRepo.GetActiveRecommendations(ctx, studentID)
	if err == nil && len(cached) > 0 {
		return s.formatCategorized(cached), nil
	}

	// 2. Generate new deterministic recommendations
	return s.generateRecommendations(ctx, studentID)
}

func (s *serviceImpl) GetNextRecommendation(ctx context.Context, studentID string) (*ProblemRecommendation, error) {
	rec, err := s.recRepo.GetTopPick(ctx, studentID)
	if err == nil && rec != nil {
		return rec, nil
	}

	// If none cached, regenerate
	categorized, err := s.generateRecommendations(ctx, studentID)
	if err != nil {
		return nil, err
	}
	if categorized.TopPick != nil {
		return categorized.TopPick, nil
	}
	return nil, ErrNotFound
}

func (s *serviceImpl) GetWeaknessRecommendations(ctx context.Context, studentID string) ([]ProblemRecommendation, error) {
	recs, err := s.recRepo.GetByCategory(ctx, studentID, CategoryImproveWeaknesses)
	if err == nil && len(recs) > 0 {
		return recs, nil
	}

	// Generate fresh if not available
	categorized, err := s.generateRecommendations(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return categorized.Categories[CategoryImproveWeaknesses], nil
}

func (s *serviceImpl) GetSimilarRecommendations(ctx context.Context, studentID, problemID string) ([]ProblemRecommendation, error) {
	prob, err := s.problemsRepo.Get(ctx, problemID)
	if err != nil || prob == nil {
		return nil, fmt.Errorf("reference problem not found")
	}

	// Fetch related problems from DSA database
	related, err := s.problemsRepo.Related(ctx, problemID, 8)
	if err != nil {
		return nil, err
	}

	// Filter and convert to recommendations
	var recs []ProblemRecommendation
	for _, p := range related {
		if p.ID == problemID {
			continue
		}
		concept := p.TopicName
		if concept == "" {
			concept = "Data Structures & Algorithms"
		}
		recs = append(recs, ProblemRecommendation{
			StudentID:             studentID,
			ProblemID:             p.ID,
			Category:              CategorySimilarProblems,
			RecommendedDifficulty: p.Difficulty,
			PrimaryConcept:        concept,
			RelatedConcepts:       p.Tags,
			MatchScore:            82.0,
			Headline:              fmt.Sprintf("Similar Pattern: %s", p.Title),
			Explanation:           fmt.Sprintf("Reinforces similar algorithmic concepts as %s with alternative variations.", prob.Title),
			LearningObjective:     fmt.Sprintf("Master the %s pattern across diverse problem constraints.", concept),
			Status:                StatusActive,
			ProblemTitle:          p.Title,
			ProblemSlug:           p.Slug,
			ProblemDifficulty:     p.Difficulty,
			ProblemSummary:        p.Summary,
			ExternalProblemURL:    p.ExternalProblemURL,
			TopicName:             p.TopicName,
			SubtopicName:          p.SubtopicName,
			Tags:                  p.Tags,
		})
	}
	return recs, nil
}

func (s *serviceImpl) RefreshRecommendations(ctx context.Context, studentID string) (*CategorizedRecommendations, error) {
	_ = s.recRepo.ClearActive(ctx, studentID)
	return s.generateRecommendations(ctx, studentID)
}

func (s *serviceImpl) SubmitFeedback(ctx context.Context, studentID, recommendationID string, req FeedbackRequest) error {
	rating := strings.ToLower(strings.TrimSpace(req.Rating))
	switch rating {
	case FeedbackHelpful, FeedbackNotRelevant, FeedbackTooEasy, FeedbackTooDifficult:
	default:
		return fmt.Errorf("invalid feedback rating: must be one of helpful, not_relevant, too_easy, too_difficult")
	}
	return s.recRepo.SaveFeedback(ctx, studentID, recommendationID, rating, req.Notes)
}

func (s *serviceImpl) GetHistory(ctx context.Context, studentID string, limit, offset int) (*HistoryResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	items, total, err := s.recRepo.GetHistory(ctx, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	return &HistoryResponse{
		History: items,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

func (s *serviceImpl) OnSubmissionEvaluated(ctx context.Context, studentID, problemID, status string) {
	if status == "Accepted" {
		_ = s.recRepo.MarkCompleted(ctx, studentID, problemID)
	}
}

// -------------------------------------------------------------------------
// Deterministic Scoring & Generation Engine
// -------------------------------------------------------------------------

func (s *serviceImpl) generateRecommendations(ctx context.Context, studentID string) (*CategorizedRecommendations, error) {
	s.logger.Info("generating personalised problem recommendations", slog.String("student_id", studentID))

	// 1. Gather all student signals
	progressList, _ := s.progressRepo.ListByStudent(ctx, studentID)
	progressMap := make(map[string]*model.Progress)
	for i := range progressList {
		progressMap[progressList[i].ProblemID] = &progressList[i]
	}

	var weaknesses []learningmemory.ConceptWeakness
	var conceptMemories []learningmemory.ConceptMemory
	if s.learningMemorySvc != nil {
		weaknesses, _ = s.learningMemorySvc.GetWeaknesses(ctx, studentID, 10)
		conceptMemories, _, _ = s.learningMemorySvc.ListConcepts(ctx, studentID, learningmemory.ConceptFilter{Limit: 20})
	}

	var mistakePatterns []mistakes.StudentMistakePattern
	if s.mistakesRepo != nil {
		mistakePatterns, _ = s.mistakesRepo.GetActivePatterns(ctx, studentID)
	}

	var nextPathAction *learningpath.LearningRecommendation
	if s.learningPathSvc != nil {
		nextPathAction, _ = s.learningPathSvc.GetNextAction(ctx, studentID)
	}

	// 2. Fetch problems catalog from database
	allProblems, _, err := s.problemsRepo.List(ctx, model.ProblemFilter{
		Limit: 150,
		Sort:  "problem_order",
	})
	if err != nil || len(allProblems) == 0 {
		return s.emptyCategorized(studentID), nil
	}

	// 3. Score problems deterministically across categories
	categories := map[string][]ProblemRecommendation{
		CategoryRecommendedForYou:      {},
		CategoryImproveWeaknesses:       {},
		CategoryStrengthenFundamentals: {},
		CategoryChallengeYourself:      {},
		CategoryRevisionPractice:       {},
	}

	var candidates []ProblemRecommendation
	var topCandidate *ProblemRecommendation

	for _, p := range allProblems {
		prog := progressMap[p.ID]
		isSolved := prog != nil && prog.Status == model.StatusSolved
		isAttempted := prog != nil && (prog.Status == model.StatusAttempted || prog.Status == model.StatusNeedsRevision)

		// Concept name
		concept := p.TopicName
		if concept == "" {
			concept = "Foundational DSA"
		}

		// A. Weakness Candidate Check
		isWeaknessMatch := false
		weaknessReason := ""
		for _, w := range weaknesses {
			if strings.EqualFold(w.ConceptKey, p.TopicName) || strings.EqualFold(w.ConceptName, p.TopicName) {
				isWeaknessMatch = true
				weaknessReason = fmt.Sprintf("Addresses recurring %s mistakes (%s)", w.MistakeSubtype, w.MistakeCategory)
				break
			}
		}
		if !isWeaknessMatch {
			for _, m := range mistakePatterns {
				if m.CorrectionStatus == "active" && (strings.EqualFold(m.RelatedConcept, p.TopicName) || strings.Contains(strings.ToLower(p.Summary), strings.ToLower(m.Category))) {
					isWeaknessMatch = true
					weaknessReason = fmt.Sprintf("Reinforces boundary handling and edge-case validation for %s", m.Category)
					break
				}
			}
		}

		if isWeaknessMatch && (!isSolved || isAttempted) {
			rec := ProblemRecommendation{
				StudentID:             studentID,
				ProblemID:             p.ID,
				Category:              CategoryImproveWeaknesses,
				RecommendedDifficulty: p.Difficulty,
				PrimaryConcept:        concept,
				RelatedConcepts:       p.Tags,
				MatchScore:            92.0,
				Headline:              fmt.Sprintf("Overcome Weakness: %s", p.Title),
				Explanation:           fmt.Sprintf("Your recent submissions showed recurring patterns in this area. %s", weaknessReason),
				LearningObjective:     fmt.Sprintf("Eliminate failure modes and reinforce correct invariant checking in %s.", concept),
				Status:                StatusActive,
				ProblemTitle:          p.Title,
				ProblemSlug:           p.Slug,
				ProblemDifficulty:     p.Difficulty,
				ProblemSummary:        p.Summary,
				ExternalProblemURL:    p.ExternalProblemURL,
				TopicName:             p.TopicName,
				SubtopicName:          p.SubtopicName,
				Tags:                  p.Tags,
				AttemptStatus:         progStatus(prog),
			}
			categories[CategoryImproveWeaknesses] = append(categories[CategoryImproveWeaknesses], rec)
			candidates = append(candidates, rec)
		}

		// B. Revision Practice (Attempted but not solved, or marked Needs Revision)
		if isAttempted || (prog != nil && prog.Status == model.StatusNeedsRevision) {
			rec := ProblemRecommendation{
				StudentID:             studentID,
				ProblemID:             p.ID,
				Category:              CategoryRevisionPractice,
				RecommendedDifficulty: p.Difficulty,
				PrimaryConcept:        concept,
				RelatedConcepts:       p.Tags,
				MatchScore:            88.0,
				Headline:              fmt.Sprintf("Complete Revision: %s", p.Title),
				Explanation:           "You previously attempted this problem. Re-attempting it independently solidifies conceptual understanding.",
				LearningObjective:     "Achieve an independent accepted solution to close previous comprehension gaps.",
				Status:                StatusActive,
				ProblemTitle:          p.Title,
				ProblemSlug:           p.Slug,
				ProblemDifficulty:     p.Difficulty,
				ProblemSummary:        p.Summary,
				ExternalProblemURL:    p.ExternalProblemURL,
				TopicName:             p.TopicName,
				SubtopicName:          p.SubtopicName,
				Tags:                  p.Tags,
				AttemptStatus:         progStatus(prog),
			}
			categories[CategoryRevisionPractice] = append(categories[CategoryRevisionPractice], rec)
			candidates = append(candidates, rec)
		}

		// C. Strengthen Fundamentals (Easy problems that are not yet solved)
		if p.Difficulty == model.DifficultyEasy && !isSolved {
			rec := ProblemRecommendation{
				StudentID:             studentID,
				ProblemID:             p.ID,
				Category:              CategoryStrengthenFundamentals,
				RecommendedDifficulty: model.DifficultyEasy,
				PrimaryConcept:        concept,
				RelatedConcepts:       p.Tags,
				MatchScore:            84.0,
				Headline:              fmt.Sprintf("Strengthen Foundations: %s", p.Title),
				Explanation:           "Mastering clean implementation and time complexity on fundamental problems prevents subtle bugs in advanced topics.",
				LearningObjective:     fmt.Sprintf("Establish rapid algorithmic intuition for core %s operations.", concept),
				Status:                StatusActive,
				ProblemTitle:          p.Title,
				ProblemSlug:           p.Slug,
				ProblemDifficulty:     p.Difficulty,
				ProblemSummary:        p.Summary,
				ExternalProblemURL:    p.ExternalProblemURL,
				TopicName:             p.TopicName,
				SubtopicName:          p.SubtopicName,
				Tags:                  p.Tags,
				AttemptStatus:         progStatus(prog),
			}
			categories[CategoryStrengthenFundamentals] = append(categories[CategoryStrengthenFundamentals], rec)
			candidates = append(candidates, rec)
		}

		// D. Challenge Yourself (Hard or top Medium problems for progressing topics)
		if (p.Difficulty == model.DifficultyHard || p.Difficulty == model.DifficultyMedium) && !isSolved {
			rec := ProblemRecommendation{
				StudentID:             studentID,
				ProblemID:             p.ID,
				Category:              CategoryChallengeYourself,
				RecommendedDifficulty: p.Difficulty,
				PrimaryConcept:        concept,
				RelatedConcepts:       p.Tags,
				MatchScore:            86.0,
				Headline:              fmt.Sprintf("Push Limits: %s", p.Title),
				Explanation:           "Challenge yourself with optimal asymptotic requirements, non-trivial constraints, and multi-concept combinations.",
				LearningObjective:     "Cultivate deep algorithmic problem-solving and rigorous edge-case analysis.",
				Status:                StatusActive,
				ProblemTitle:          p.Title,
				ProblemSlug:           p.Slug,
				ProblemDifficulty:     p.Difficulty,
				ProblemSummary:        p.Summary,
				ExternalProblemURL:    p.ExternalProblemURL,
				TopicName:             p.TopicName,
				SubtopicName:          p.SubtopicName,
				Tags:                  p.Tags,
				AttemptStatus:         progStatus(prog),
			}
			categories[CategoryChallengeYourself] = append(categories[CategoryChallengeYourself], rec)
			candidates = append(candidates, rec)
		}

		// E. Recommended For You (Current curriculum node or unpracticed topic problem)
		if !isSolved {
			score := 80.0
			headline := fmt.Sprintf("Next Step: %s", p.Title)
			explanation := fmt.Sprintf("Aligned with your current curriculum trajectory in %s.", concept)

			if nextPathAction != nil && (nextPathAction.ProblemID != nil && *nextPathAction.ProblemID == p.ID) {
				score = 98.0
				headline = fmt.Sprintf("Top Priority: %s", p.Title)
				explanation = nextPathAction.Explanation
			}

			rec := ProblemRecommendation{
				StudentID:             studentID,
				ProblemID:             p.ID,
				Category:              CategoryRecommendedForYou,
				RecommendedDifficulty: p.Difficulty,
				PrimaryConcept:        concept,
				RelatedConcepts:       p.Tags,
				MatchScore:            score,
				Headline:              headline,
				Explanation:           explanation,
				LearningObjective:     fmt.Sprintf("Develop strong algorithmic patterns for %s.", concept),
				Status:                StatusActive,
				ProblemTitle:          p.Title,
				ProblemSlug:           p.Slug,
				ProblemDifficulty:     p.Difficulty,
				ProblemSummary:        p.Summary,
				ExternalProblemURL:    p.ExternalProblemURL,
				TopicName:             p.TopicName,
				SubtopicName:          p.SubtopicName,
				Tags:                  p.Tags,
				AttemptStatus:         progStatus(prog),
			}
			categories[CategoryRecommendedForYou] = append(categories[CategoryRecommendedForYou], rec)
			candidates = append(candidates, rec)
		}
	}

	// Limit each category to top 5
	for cat := range categories {
		sort.Slice(categories[cat], func(i, j int) bool {
			return categories[cat][i].MatchScore > categories[cat][j].MatchScore
		})
		if len(categories[cat]) > 5 {
			categories[cat] = categories[cat][:5]
		}
	}

	// Identify Top Pick candidate (highest match score across all active)
	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].MatchScore > candidates[j].MatchScore
		})
		topCandidate = &candidates[0]
		topCandidate.IsTopPick = true
	}

	// Flatten and save all active recommendations to repository
	var toSave []ProblemRecommendation
	seenIDs := make(map[string]bool)
	if topCandidate != nil {
		toSave = append(toSave, *topCandidate)
		seenIDs[topCandidate.ProblemID] = true
	}
	for _, recList := range categories {
		for _, r := range recList {
			if !seenIDs[r.ProblemID] {
				toSave = append(toSave, r)
				seenIDs[r.ProblemID] = true
			}
		}
	}

	if err := s.recRepo.SaveRecommendations(ctx, studentID, toSave); err != nil {
		s.logger.Warn("failed to persist problem recommendations", slog.String("error", err.Error()))
	}

	primaryFocus := "Arrays & Two Pointers"
	if len(conceptMemories) > 0 {
		primaryFocus = conceptMemories[0].ConceptName
	} else if nextPathAction != nil {
		primaryFocus = nextPathAction.ConceptName
	}

	return &CategorizedRecommendations{
		TopPick:               topCandidate,
		Categories:            categories,
		TotalRecommendations:  len(toSave),
		RecommendedDifficulty: model.DifficultyMedium,
		PrimaryFocusConcept:   primaryFocus,
		GeneratedAt:           time.Now(),
	}, nil
}

func (s *serviceImpl) formatCategorized(items []ProblemRecommendation) *CategorizedRecommendations {
	categories := map[string][]ProblemRecommendation{
		CategoryRecommendedForYou:      {},
		CategoryImproveWeaknesses:       {},
		CategoryStrengthenFundamentals: {},
		CategoryChallengeYourself:      {},
		CategoryRevisionPractice:       {},
	}

	var topPick *ProblemRecommendation
	for _, item := range items {
		if item.IsTopPick && topPick == nil {
			topPick = &item
		}
		if _, ok := categories[item.Category]; ok {
			categories[item.Category] = append(categories[item.Category], item)
		} else {
			categories[CategoryRecommendedForYou] = append(categories[CategoryRecommendedForYou], item)
		}
	}

	primary := "Data Structures & Algorithms"
	diff := model.DifficultyMedium
	if topPick != nil {
		primary = topPick.PrimaryConcept
		diff = topPick.RecommendedDifficulty
	}

	return &CategorizedRecommendations{
		TopPick:               topPick,
		Categories:            categories,
		TotalRecommendations:  len(items),
		RecommendedDifficulty: diff,
		PrimaryFocusConcept:   primary,
		GeneratedAt:           time.Now(),
	}
}

func (s *serviceImpl) emptyCategorized(studentID string) *CategorizedRecommendations {
	return &CategorizedRecommendations{
		Categories: map[string][]ProblemRecommendation{
			CategoryRecommendedForYou:      {},
			CategoryImproveWeaknesses:       {},
			CategoryStrengthenFundamentals: {},
			CategoryChallengeYourself:      {},
			CategoryRevisionPractice:       {},
		},
		TotalRecommendations:  0,
		RecommendedDifficulty: model.DifficultyEasy,
		PrimaryFocusConcept:   "Arrays & Basics",
		GeneratedAt:           time.Now(),
	}
}

func progStatus(p *model.Progress) string {
	if p == nil {
		return "Not Started"
	}
	return p.Status
}

