package visualizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/ai"
)

type Service interface {
	ListAlgorithms() []AlgorithmMetadata
	GetAlgorithm(id string) (*AlgorithmMetadata, error)
	ExecuteAlgorithm(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error)
	CreateSession(ctx context.Context, studentID string, req CreateSessionRequest) (*VisualizerSession, error)
	GetSession(ctx context.Context, sessionID, studentID string) (*VisualizerSession, error)
	ListSessions(ctx context.Context, studentID string, limit, offset int) ([]VisualizerSession, error)
	UpdateSession(ctx context.Context, sessionID, studentID string, req UpdateSessionRequest) error
	ExplainStep(ctx context.Context, req ExplainStepRequest) (*ExplainStepResponse, error)
	GetPracticeQuestions(ctx context.Context, algoID string) ([]PracticeQuestion, error)
}

type serviceImpl struct {
	repo   Repository
	llm    ai.LLMProvider
	logger *slog.Logger
}

func NewService(repo Repository, llm ai.LLMProvider, logger *slog.Logger) Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		repo:   repo,
		llm:    llm,
		logger: logger,
	}
}

func (s *serviceImpl) ListAlgorithms() []AlgorithmMetadata {
	return AlgorithmsRegistry
}

func (s *serviceImpl) GetAlgorithm(id string) (*AlgorithmMetadata, error) {
	meta := GetAlgorithmByID(id)
	if meta == nil {
		return nil, fmt.Errorf("algorithm '%s' not found", id)
	}
	return meta, nil
}

func (s *serviceImpl) ExecuteAlgorithm(ctx context.Context, req ExecuteRequest) (*ExecuteResponse, error) {
	meta := GetAlgorithmByID(req.AlgorithmID)
	if meta == nil {
		return nil, fmt.Errorf("algorithm '%s' not found", req.AlgorithmID)
	}

	steps, err := GenerateExecutionSteps(req.AlgorithmID, req.CustomInput)
	if err != nil {
		return nil, err
	}

	return &ExecuteResponse{
		AlgorithmID: meta.ID,
		Name:        meta.Name,
		Category:    meta.Category,
		Complexity:  meta.Complexity,
		TotalSteps:  len(steps),
		Steps:       steps,
	}, nil
}

func (s *serviceImpl) CreateSession(ctx context.Context, studentID string, req CreateSessionRequest) (*VisualizerSession, error) {
	meta := GetAlgorithmByID(req.AlgorithmID)
	if meta == nil {
		return nil, fmt.Errorf("invalid algorithm_id '%s'", req.AlgorithmID)
	}

	session := &VisualizerSession{
		StudentID:      studentID,
		AlgorithmID:    req.AlgorithmID,
		Category:       meta.Category,
		InputConfig:    req.InputConfig,
		TotalSteps:     req.TotalSteps,
		StepsCompleted: 0,
		Completed:      false,
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *serviceImpl) GetSession(ctx context.Context, sessionID, studentID string) (*VisualizerSession, error) {
	return s.repo.GetSession(ctx, sessionID, studentID)
}

func (s *serviceImpl) ListSessions(ctx context.Context, studentID string, limit, offset int) ([]VisualizerSession, error) {
	return s.repo.ListSessions(ctx, studentID, limit, offset)
}

func (s *serviceImpl) UpdateSession(ctx context.Context, sessionID, studentID string, req UpdateSessionRequest) error {
	return s.repo.UpdateSession(ctx, sessionID, studentID, req.StepsCompleted, req.TimeSpentSeconds, req.Completed)
}

func (s *serviceImpl) ExplainStep(ctx context.Context, req ExplainStepRequest) (*ExplainStepResponse, error) {
	meta := GetAlgorithmByID(req.AlgorithmID)
	name := req.AlgorithmID
	if meta != nil {
		name = meta.Name
	}

	// If LLM provider is available, formulate Socratic algorithmic query
	if s.llm != nil {
		stateJSON, _ := json.Marshal(req.CurrentState)
		prompt := fmt.Sprintf(`You are an expert DSA Professor and algorithm visualizer tutor.
Algorithm: %s
Current Step: %d of %d
Operation: %s
Step Description: %s
State: %s
User Question: %s

Respond with a strictly formatted JSON object:
{
  "explanation": "Clear, intuitive breakdown of why this specific operation happens right now",
  "intuition": "Algorithmic invariant or reason behind this decision",
  "next_expected_move": "What pointer or data structure change is expected to follow",
  "key_takeaway": "One-sentence memorable rule of thumb for coding interviews"
}`, name, req.StepNumber, req.TotalSteps, req.Operation, req.Explanation, string(stateJSON), req.UserQuestion)

		llmResp, err := s.llm.Generate(ctx, ai.ChatCompletionRequest{
			Messages: []ai.ChatMessage{
				{Role: "system", Content: "You are LearnX's Algorithm Visualizer AI Tutor. Return pure JSON only."},
				{Role: "user", Content: prompt},
			},
			Temperature:    0.3,
			ResponseFormat: "json_object",
		})
		if err == nil && llmResp != nil && llmResp.Content != "" {
			var parsed ExplainStepResponse
			cleanContent := strings.TrimSpace(llmResp.Content)
			cleanContent = strings.TrimPrefix(cleanContent, "```json")
			cleanContent = strings.TrimPrefix(cleanContent, "```")
			cleanContent = strings.TrimSuffix(cleanContent, "```")
			if err := json.Unmarshal([]byte(cleanContent), &parsed); err == nil && parsed.Explanation != "" {
				return &parsed, nil
			}
		}
	}

	// Deterministic fallback explanation
	return &ExplainStepResponse{
		Explanation:      fmt.Sprintf("At Step %d of %s, the algorithm performs operation '%s'. %s", req.StepNumber, name, req.Operation, req.Explanation),
		Intuition:        "Algorithm invariants guarantee that each step monotonically reduces the remaining problem size or restores the required data structure property.",
		NextExpectedMove: "Continuing to the next sequential instruction or comparison in the algorithm's loop.",
		KeyTakeaway:      "Watch how pointers narrow down the active search window or how elements bubble into their final sorted positions.",
	}, nil
}

func (s *serviceImpl) GetPracticeQuestions(_ context.Context, algoID string) ([]PracticeQuestion, error) {
	meta := GetAlgorithmByID(algoID)
	if meta == nil {
		return nil, fmt.Errorf("algorithm '%s' not found", algoID)
	}

	// Return high-quality, targeted conceptual & complexity MCQs for this algorithm
	switch meta.Category {
	case CategorySorting:
		return []PracticeQuestion{
			{
				ID:                 fmt.Sprintf("%s-q1", algoID),
				AlgorithmID:        algoID,
				Question:           fmt.Sprintf("What is the worst-case time complexity of %s?", meta.Name),
				Options:            []string{"O(N)", "O(N log N)", "O(N²)", "O(2ⁿ)"},
				CorrectOptionIndex: func() int {
					if meta.ID == "merge-sort" || meta.ID == "heap-sort" {
						return 1
					}
					return 2
				}(),
				Explanation: fmt.Sprintf("%s has a worst-case time complexity of %s.", meta.Name, meta.Complexity.Worst),
				Difficulty:  "Easy",
			},
			{
				ID:                 fmt.Sprintf("%s-q2", algoID),
				AlgorithmID:        algoID,
				Question:           fmt.Sprintf("Is %s an in-place sorting algorithm?", meta.Name),
				Options:            []string{"Yes, it uses O(1) or O(log N) auxiliary space", "No, it requires O(N) extra space", "It depends on the compiler", "Only when sorted backwards"},
				CorrectOptionIndex: func() int {
					if meta.ID == "merge-sort" {
						return 1
					}
					return 0
				}(),
				Explanation: fmt.Sprintf("%s space complexity is %s.", meta.Name, meta.Complexity.Space),
				Difficulty:  "Medium",
			},
			{
				ID:                 fmt.Sprintf("%s-q3", algoID),
				AlgorithmID:        algoID,
				Question:           "Which property ensures that elements with identical keys maintain their relative order after sorting?",
				Options:            []string{"In-place", "Stability", "Adaptability", "Monotonicity"},
				CorrectOptionIndex: 1,
				Explanation:        "A stable sort preserves the initial relative order of equal keys.",
				Difficulty:         "Medium",
			},
		}, nil

	case CategorySearching:
		return []PracticeQuestion{
			{
				ID:                 fmt.Sprintf("%s-q1", algoID),
				AlgorithmID:        algoID,
				Question:           "What is the crucial prerequisite before executing Binary Search?",
				Options:            []string{"The array must have distinct elements only", "The array must be sorted monotonically", "The array must have an even length", "The target must be present in the middle"},
				CorrectOptionIndex: 1,
				Explanation:        "Binary Search relies on monotonic ordering to eliminate half of the remaining interval at each step.",
				Difficulty:         "Easy",
			},
			{
				ID:                 fmt.Sprintf("%s-q2", algoID),
				AlgorithmID:        algoID,
				Question:           "Why is `mid = low + (high - low) / 2` preferred over `mid = (low + high) / 2` in typed languages like C++ and Java?",
				Options:            []string{"It runs in O(1) instead of O(log N)", "It avoids 32-bit signed integer overflow when low + high exceeds INT_MAX", "It automatically rounds up", "It works for negative numbers only"},
				CorrectOptionIndex: 1,
				Explanation:        "If low and high are both close to 2^31 - 1, their sum overflows into a negative integer, causing out-of-bounds exceptions.",
				Difficulty:         "Hard",
			},
		}, nil

	case CategoryDP:
		return []PracticeQuestion{
			{
				ID:                 fmt.Sprintf("%s-q1", algoID),
				AlgorithmID:        algoID,
				Question:           "What two core characteristics make a problem suitable for Dynamic Programming?",
				Options:            []string{"Optimal Substructure and Overlapping Subproblems", "Divide and Conquer with Independent Subproblems", "Greedy Choice Property and Monotonicity", "Complete Search and Backtracking"},
				CorrectOptionIndex: 0,
				Explanation:        "Dynamic Programming requires optimal substructure (optimal solution composed of optimal sub-solutions) and overlapping subproblems (avoiding recalculating same subproblems).",
				Difficulty:         "Easy",
			},
			{
				ID:                 fmt.Sprintf("%s-q2", algoID),
				AlgorithmID:        algoID,
				Question:           "What is the key difference between Memoization (Top-Down) and Tabulation (Bottom-Up)?",
				Options:            []string{"Memoization uses iteration; Tabulation uses recursion", "Memoization uses recursive calls with cache lookup; Tabulation fills an iterative table starting from base cases", "Tabulation has worse time complexity", "Memoization uses zero memory"},
				CorrectOptionIndex: 1,
				Explanation:        "Top-down memoization adds a cache to recursion, whereas bottom-up tabulation builds the solution iteratively without recursion stack overhead.",
				Difficulty:         "Medium",
			},
		}, nil

	default:
		return []PracticeQuestion{
			{
				ID:                 fmt.Sprintf("%s-q1", algoID),
				AlgorithmID:        algoID,
				Question:           fmt.Sprintf("What is the primary problem solved by %s?", meta.Name),
				Options:            []string{meta.ProblemSolved, "Generates random permutations", "Compresses audio files", "Replaces hash tables in all cases"},
				CorrectOptionIndex: 0,
				Explanation:        meta.Description,
				Difficulty:         "Easy",
			},
			{
				ID:                 fmt.Sprintf("%s-q2", algoID),
				AlgorithmID:        algoID,
				Question:           fmt.Sprintf("What is the average time complexity of %s?", meta.Name),
				Options:            []string{meta.Complexity.Average, "O(N!)", "O(1)", "O(N³)"},
				CorrectOptionIndex: 0,
				Explanation:        fmt.Sprintf("The average complexity is %s and auxiliary space is %s.", meta.Complexity.Average, meta.Complexity.Space),
				Difficulty:         "Medium",
			},
		}, nil
	}
}

