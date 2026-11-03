package learningpath

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/learningmemory"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
)

// CurriculumConcept defines the standard prerequisite hierarchy in DSA.
type CurriculumConcept struct {
	ConceptKey   string
	ConceptName  string
	TopicSlug    string
	TopicName    string
	Order        int
	Prerequisite []string
	Difficulty   string
}

// Canonical curriculum sequence matching Striver A2Z sheet pillars.
var canonicalCurriculum = []CurriculumConcept{
	{
		ConceptKey:   "basics",
		ConceptName:  "Language & Logic Fundamentals",
		TopicSlug:    "basics",
		TopicName:    "Learn the Basics",
		Order:        1,
		Prerequisite: []string{},
		Difficulty:   model.DifficultyEasy,
	},
	{
		ConceptKey:   "sorting",
		ConceptName:  "Sorting Techniques",
		TopicSlug:    "sorting",
		TopicName:    "Sorting Techniques",
		Order:        2,
		Prerequisite: []string{"basics"},
		Difficulty:   model.DifficultyEasy,
	},
	{
		ConceptKey:   "arrays",
		ConceptName:  "Arrays & Linear Iteration",
		TopicSlug:    "arrays",
		TopicName:    "Arrays",
		Order:        3,
		Prerequisite: []string{"sorting"},
		Difficulty:   model.DifficultyEasy,
	},
	{
		ConceptKey:   "binary_search",
		ConceptName:  "Binary Search & Bounds",
		TopicSlug:    "binary-search",
		TopicName:    "Binary Search",
		Order:        4,
		Prerequisite: []string{"arrays"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "strings",
		ConceptName:  "String Manipulation & Parsing",
		TopicSlug:    "strings",
		TopicName:    "Strings",
		Order:        5,
		Prerequisite: []string{"arrays"},
		Difficulty:   model.DifficultyEasy,
	},
	{
		ConceptKey:   "linked_list",
		ConceptName:  "Linked List Operations",
		TopicSlug:    "linked-list",
		TopicName:    "Linked List",
		Order:        6,
		Prerequisite: []string{"basics"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "recursion",
		ConceptName:  "Recursion & Backtracking",
		TopicSlug:    "recursion",
		TopicName:    "Recursion",
		Order:        7,
		Prerequisite: []string{"basics"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "bit_manipulation",
		ConceptName:  "Bit Manipulation",
		TopicSlug:    "bit-manipulation",
		TopicName:    "Bit Manipulation",
		Order:        8,
		Prerequisite: []string{"basics"},
		Difficulty:   model.DifficultyEasy,
	},
	{
		ConceptKey:   "stack_queue",
		ConceptName:  "Stacks & Queues",
		TopicSlug:    "stack-queue",
		TopicName:    "Stack and Queues",
		Order:        9,
		Prerequisite: []string{"arrays", "linked_list"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "sliding_window",
		ConceptName:  "Sliding Window & Two Pointer",
		TopicSlug:    "sliding-window",
		TopicName:    "Sliding Window & Two Pointer",
		Order:        10,
		Prerequisite: []string{"arrays"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "heaps",
		ConceptName:  "Heaps & Priority Queues",
		TopicSlug:    "heaps",
		TopicName:    "Heaps",
		Order:        11,
		Prerequisite: []string{"arrays"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "greedy",
		ConceptName:  "Greedy Algorithms",
		TopicSlug:    "greedy",
		TopicName:    "Greedy Algorithms",
		Order:        12,
		Prerequisite: []string{"sorting", "arrays"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "binary_trees",
		ConceptName:  "Binary Trees & Traversals",
		TopicSlug:    "binary-trees",
		TopicName:    "Binary Trees",
		Order:        13,
		Prerequisite: []string{"recursion", "stack_queue"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "bst",
		ConceptName:  "Binary Search Trees",
		TopicSlug:    "bst",
		TopicName:    "Binary Search Trees",
		Order:        14,
		Prerequisite: []string{"binary_trees", "binary_search"},
		Difficulty:   model.DifficultyMedium,
	},
	{
		ConceptKey:   "graphs",
		ConceptName:  "Graphs & Shortest Path",
		TopicSlug:    "graphs",
		TopicName:    "Graphs",
		Order:        15,
		Prerequisite: []string{"recursion", "stack_queue"},
		Difficulty:   model.DifficultyHard,
	},
	{
		ConceptKey:   "dp",
		ConceptName:  "Dynamic Programming",
		TopicSlug:    "dynamic-programming",
		TopicName:    "Dynamic Programming",
		Order:        16,
		Prerequisite: []string{"recursion", "arrays"},
		Difficulty:   model.DifficultyHard,
	},
	{
		ConceptKey:   "trie",
		ConceptName:  "Trie (Prefix Tree)",
		TopicSlug:    "trie",
		TopicName:    "Trie",
		Order:        17,
		Prerequisite: []string{"strings", "binary_trees"},
		Difficulty:   model.DifficultyHard,
	},
	{
		ConceptKey:   "strings_advanced",
		ConceptName:  "Advanced String Algorithms",
		TopicSlug:    "strings-advanced",
		TopicName:    "Advanced Strings",
		Order:        18,
		Prerequisite: []string{"trie", "strings"},
		Difficulty:   model.DifficultyHard,
	},
}

// Service provides high-level business methods for the adaptive learning path.
type Service interface {
	GetLearningPath(ctx context.Context, studentID string) (*AdaptivePath, error)
	GetNextAction(ctx context.Context, studentID string) (*LearningRecommendation, error)
	GetPathProgress(ctx context.Context, studentID string) (*PathProgressResponse, error)
	GetHistory(ctx context.Context, studentID string, limit, offset int) ([]LearningRecommendation, int, error)
	RecalculatePath(ctx context.Context, studentID string) (*AdaptivePath, error)
	Recalculate(ctx context.Context, studentID string)
}

type serviceImpl struct {
	repo         Repository
	memoryRepo   learningmemory.Repository
	problemsRepo problems.Repository
	topicsRepo   topics.Repository
	logger       *slog.Logger
}

func NewService(
	repo Repository,
	memoryRepo learningmemory.Repository,
	problemsRepo problems.Repository,
	topicsRepo topics.Repository,
	logger *slog.Logger,
) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		repo:         repo,
		memoryRepo:   memoryRepo,
		problemsRepo: problemsRepo,
		topicsRepo:   topicsRepo,
		logger:       logger,
	}
}

func (s *serviceImpl) GetLearningPath(ctx context.Context, studentID string) (*AdaptivePath, error) {
	path, err := s.repo.GetAdaptivePath(ctx, studentID)
	if err == nil && path != nil && len(path.Nodes) > 0 {
		return path, nil
	}
	// Initial path generation
	return s.RecalculatePath(ctx, studentID)
}

func (s *serviceImpl) GetNextAction(ctx context.Context, studentID string) (*LearningRecommendation, error) {
	rec, err := s.repo.GetActiveRecommendation(ctx, studentID)
	if err == nil && rec != nil {
		return rec, nil
	}
	// If no active recommendation, recalculate and generate one
	_, err = s.RecalculatePath(ctx, studentID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetActiveRecommendation(ctx, studentID)
}

func (s *serviceImpl) GetPathProgress(ctx context.Context, studentID string) (*PathProgressResponse, error) {
	path, err := s.GetLearningPath(ctx, studentID)
	if err != nil {
		return nil, err
	}

	resp := &PathProgressResponse{
		TotalConcepts:   len(path.Nodes),
		CurrentStage:    path.LearningStage,
		LastEvaluatedAt: path.UpdatedAt,
	}

	var weightedScore float64
	for _, n := range path.Nodes {
		switch n.State {
		case StateMastered:
			resp.MasteredCount++
			weightedScore += 100.0
		case StateInProgress:
			resp.InProgressCount++
			weightedScore += 75.0
		case StatePracticing:
			resp.InProgressCount++
			weightedScore += 50.0
		case StateNeedsRevision:
			resp.NeedsRevisionCount++
			weightedScore += 35.0
		case StateReadyToLearn:
			resp.ReadyToLearnCount++
			weightedScore += 10.0
		case StateLockedByPrerequisite, StateNotStarted:
			resp.LockedCount++
		}
	}

	if resp.TotalConcepts > 0 {
		pct := weightedScore / float64(resp.TotalConcepts)
		if pct > 100.0 {
			pct = 100.0
		}
		resp.OverallProgress = math.Round(pct*10) / 10
	}

	return resp, nil
}

func (s *serviceImpl) GetHistory(ctx context.Context, studentID string, limit, offset int) ([]LearningRecommendation, int, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.ListRecommendationHistory(ctx, studentID, limit, offset)
}

func (s *serviceImpl) Recalculate(ctx context.Context, studentID string) {
	_, err := s.RecalculatePath(ctx, studentID)
	if err != nil {
		s.logger.Warn("async recalculation failed for student", slog.String("student_id", studentID), slog.String("error", err.Error()))
	}
}

func (s *serviceImpl) RecalculatePath(ctx context.Context, studentID string) (*AdaptivePath, error) {
	s.logger.Info("recalculating adaptive learning path", slog.String("student_id", studentID))

	// 1. Fetch student's learning memory concept states & mistakes
	conceptMemories := make(map[string]*learningmemory.ConceptMemory)
	if s.memoryRepo != nil {
		cmList, _, err := s.memoryRepo.ListConceptMemories(ctx, studentID, learningmemory.ConceptFilter{Limit: 200})
		if err == nil {
			for i := range cmList {
				conceptMemories[cmList[i].ConceptKey] = &cmList[i]
			}
		}
	}

	mistakesList := make([]learningmemory.MistakeMemory, 0)
	if s.memoryRepo != nil {
		m, err := s.memoryRepo.ListMistakeMemories(ctx, studentID, "active")
		if err == nil {
			mistakesList = m
		}
	}

	// Pre-fetch problems pool once to prevent N+1 database queries
	problemsByDifficulty := make(map[string][]model.Problem)
	problemsByTopicSlug := make(map[string][]model.Problem)
	if s.problemsRepo != nil {
		allProbs, _, err := s.problemsRepo.List(ctx, model.ProblemFilter{Limit: 250})
		if err == nil {
			for _, p := range allProbs {
				problemsByDifficulty[p.Difficulty] = append(problemsByDifficulty[p.Difficulty], p)
				topicKey := strings.ToLower(strings.TrimSpace(p.TopicName))
				if topicKey != "" {
					problemsByTopicSlug[topicKey] = append(problemsByTopicSlug[topicKey], p)
				}
			}
		}
	}

	// 2. Evaluate each node along curriculum DAG
	nodes := make([]PathNode, 0, len(canonicalCurriculum))
	nodeStates := make(map[string]string)

	var nextRecCandidate *LearningRecommendation
	var currentActiveConceptKey string = "basics"
	var currentTopicID *string
	recommendedDifficulty := model.DifficultyEasy

	for _, c := range canonicalCurriculum {
		state := StateNotStarted
		readinessScore := 50.0
		diffToRecommend := c.Difficulty

		// Check if prerequisites are satisfied
		prereqsSatisfied := true
		for _, prereqKey := range c.Prerequisite {
			pState, exists := nodeStates[prereqKey]
			if !exists || (pState != StateMastered && pState != StateInProgress && pState != StatePracticing) {
				prereqsSatisfied = false
				break
			}
		}

		cm, hasMemory := conceptMemories[c.ConceptKey]

		if !prereqsSatisfied {
			state = StateLockedByPrerequisite
			readinessScore = 20.0
		} else if hasMemory {
			// Memory exists -> reflect actual state
			if cm.LearningState == learningmemory.StateMastered {
				state = StateMastered
				readinessScore = 100.0
				diffToRecommend = model.DifficultyHard
			} else if cm.LearningState == learningmemory.StateStruggling || (cm.MistakeCount > cm.CorrectedMistakes) {
				state = StateNeedsRevision
				readinessScore = 35.0
				diffToRecommend = model.DifficultyEasy
			} else if cm.IndependentSolves > 0 {
				state = StateInProgress
				readinessScore = 75.0
				diffToRecommend = model.DifficultyMedium
			} else {
				state = StatePracticing
				readinessScore = 55.0
				diffToRecommend = model.DifficultyEasy
			}
		} else {
			state = StateReadyToLearn
			readinessScore = 60.0
		}

		nodeStates[c.ConceptKey] = state

		// Check candidate problems for this concept from pre-fetched pool
		var recommendedProblemID *string
		var recommendedProblemTitle *string
		var recommendedProblemSlug *string
		var topicID *string

		var candidateProblem *model.Problem
		cTopicKey := strings.ToLower(strings.TrimSpace(c.TopicName))
		cSlugKey := strings.ToLower(strings.TrimSpace(c.TopicSlug))

		if probList, ok := problemsByTopicSlug[cTopicKey]; ok && len(probList) > 0 {
			for i := range probList {
				if probList[i].Difficulty == diffToRecommend {
					candidateProblem = &probList[i]
					break
				}
			}
			if candidateProblem == nil {
				candidateProblem = &probList[0]
			}
		} else if probList, ok := problemsByTopicSlug[cSlugKey]; ok && len(probList) > 0 {
			candidateProblem = &probList[0]
		}

		if candidateProblem == nil {
			if diffList, ok := problemsByDifficulty[diffToRecommend]; ok && len(diffList) > 0 {
				candidateProblem = &diffList[0]
			}
		}

		if candidateProblem != nil {
			recommendedProblemID = &candidateProblem.ID
			recommendedProblemTitle = &candidateProblem.Title
			recommendedProblemSlug = &candidateProblem.Slug
			if candidateProblem.TopicID != "" {
				topicID = &candidateProblem.TopicID
			}
		}

		node := PathNode{
			StudentID:              studentID,
			TopicID:                topicID,
			TopicName:              c.TopicName,
			ConceptKey:             c.ConceptKey,
			ConceptName:            c.ConceptName,
			NodeOrder:              c.Order,
			State:                  state,
			PrerequisiteKeys:       c.Prerequisite,
			ReadinessScore:         readinessScore,
			RecommendedDifficulty:  diffToRecommend,
			RecommendedProblemID:   recommendedProblemID,
			RecommendedProblemName: recommendedProblemTitle,
			UpdatedAt:              time.Now(),
		}
		nodes = append(nodes, node)

		// 3. Next Action Prioritization:
		// Priority 1: Needs Revision (Active mistake or struggling)
		if nextRecCandidate == nil && state == StateNeedsRevision {
			currentActiveConceptKey = c.ConceptKey
			currentTopicID = topicID
			recommendedDifficulty = diffToRecommend

			// Find specific mistake pattern if present
			mistakeSubtype := "edge cases"
			for _, m := range mistakesList {
				if m.ConceptKey == c.ConceptKey {
					mistakeSubtype = m.Subtype
					break
				}
			}

			probIDStr := ""
			if recommendedProblemID != nil {
				probIDStr = *recommendedProblemID
			}
			probTitleStr := ""
			if recommendedProblemTitle != nil {
				probTitleStr = *recommendedProblemTitle
			}
			probSlugStr := ""
			if recommendedProblemSlug != nil {
				probSlugStr = *recommendedProblemSlug
			}

			nextRecCandidate = &LearningRecommendation{
				StudentID:             studentID,
				ActionType:            ActionReviseConcept,
				ConceptKey:            c.ConceptKey,
				ConceptName:           c.ConceptName,
				TopicID:               topicID,
				TopicName:             c.TopicName,
				ProblemID:             &probIDStr,
				ProblemTitle:          probTitleStr,
				ProblemSlug:           probSlugStr,
				RecommendedDifficulty: diffToRecommend,
				Headline:              fmt.Sprintf("Reinforce %s", c.ConceptName),
				Explanation:           fmt.Sprintf("Your recent submissions for %s showed difficulty with %s. Practicing targeted %s problems will establish a stronger foundation.", c.ConceptName, mistakeSubtype, diffToRecommend),
				RelevantWeakness:      fmt.Sprintf("Struggling with %s (%s)", mistakeSubtype, c.ConceptName),
				Evidence: map[string]interface{}{
					"learning_state":    state,
					"mistakes":          mistakeSubtype,
					"relevant_weakness": fmt.Sprintf("Struggling with %s (%s)", mistakeSubtype, c.ConceptName),
				},
				IsActive:  true,
				CreatedAt: time.Now(),
			}
		}

		// Priority 2: Ready to learn / In Progress
		if nextRecCandidate == nil && (state == StateInProgress || state == StateReadyToLearn || state == StatePracticing) {
			currentActiveConceptKey = c.ConceptKey
			currentTopicID = topicID
			recommendedDifficulty = diffToRecommend

			action := ActionPracticeMedium
			if state == StateReadyToLearn {
				action = ActionLearnConcept
			} else if diffToRecommend == model.DifficultyEasy {
				action = ActionSolveEasy
			} else if diffToRecommend == model.DifficultyHard {
				action = ActionAttemptHard
			}

			probIDStr := ""
			if recommendedProblemID != nil {
				probIDStr = *recommendedProblemID
			}
			probTitleStr := ""
			if recommendedProblemTitle != nil {
				probTitleStr = *recommendedProblemTitle
			}
			probSlugStr := ""
			if recommendedProblemSlug != nil {
				probSlugStr = *recommendedProblemSlug
			}

			nextRecCandidate = &LearningRecommendation{
				StudentID:             studentID,
				ActionType:            action,
				ConceptKey:            c.ConceptKey,
				ConceptName:           c.ConceptName,
				TopicID:               topicID,
				TopicName:             c.TopicName,
				ProblemID:             &probIDStr,
				ProblemTitle:          probTitleStr,
				ProblemSlug:           probSlugStr,
				RecommendedDifficulty: diffToRecommend,
				Headline:              fmt.Sprintf("Master %s", c.ConceptName),
				Explanation:           fmt.Sprintf("All prerequisites for %s are verified. Attempting %s difficulty problems will advance your algorithmic readiness.", c.ConceptName, diffToRecommend),
				Evidence: map[string]interface{}{
					"learning_state": state,
					"difficulty":     diffToRecommend,
				},
				IsActive:  true,
				CreatedAt: time.Now(),
			}
		}
	}

	// Fallback if all mastered
	if nextRecCandidate == nil {
		nextRecCandidate = &LearningRecommendation{
			StudentID:             studentID,
			ActionType:            ActionAttemptHard,
			ConceptKey:            "dp",
			ConceptName:           "Dynamic Programming",
			RecommendedDifficulty: model.DifficultyHard,
			Headline:              "Master Advanced Concepts",
			Explanation:           "You have demonstrated proficiency across standard foundational topics! Challenge yourself with advanced Hard problems.",
			Evidence: map[string]interface{}{
				"stage": "Advanced Mastery",
			},
			IsActive:  true,
			CreatedAt: time.Now(),
		}
	}

	// 4. Determine high-level stage
	stage := "Foundations & Array Iteration"
	if nodeStates["binary_trees"] == StateMastered {
		stage = "Non-Linear Data Structures & Trees"
	} else if nodeStates["binary_search"] == StateMastered {
		stage = "Divide & Conquer, Search, & Linear Structures"
	}

	var totalReadiness float64
	for _, n := range nodes {
		totalReadiness += n.ReadinessScore
	}
	avgReadiness := 70.0
	if len(nodes) > 0 {
		avgReadiness = math.Round((totalReadiness/float64(len(nodes)))*10) / 10
	}

	adaptivePath := &AdaptivePath{
		StudentID:             studentID,
		CurrentTopicID:        currentTopicID,
		CurrentConceptKey:     currentActiveConceptKey,
		RecommendedDifficulty: recommendedDifficulty,
		OverallReadinessScore: avgReadiness,
		LearningStage:         stage,
		Nodes:                 nodes,
		UpdatedAt:             time.Now(),
	}

	// 5. Persist path, nodes and active recommendation
	_ = s.repo.UpsertAdaptivePath(ctx, adaptivePath)
	_ = s.repo.BatchUpsertPathNodes(ctx, studentID, nodes)
	_ = s.repo.CreateRecommendation(ctx, nextRecCandidate)

	return adaptivePath, nil
}
