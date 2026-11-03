package recommendations

import (
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
)

// Recommendation Categories matching pedagogical objectives
const (
	CategoryRecommendedForYou      = "recommended_for_you"
	CategoryImproveWeaknesses       = "improve_weaknesses"
	CategoryStrengthenFundamentals = "strengthen_fundamentals"
	CategoryChallengeYourself      = "challenge_yourself"
	CategorySimilarProblems        = "similar_problems"
	CategoryRevisionPractice       = "revision_practice"
)

// Lifecycle Status
const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusSkipped   = "skipped"
	StatusExpired   = "expired"
)

// Feedback Ratings
const (
	FeedbackHelpful      = "helpful"
	FeedbackNotRelevant  = "not_relevant"
	FeedbackTooEasy      = "too_easy"
	FeedbackTooDifficult = "too_difficult"
)

// ProblemRecommendation models a personalised suggestion with full pedagogical rationale.
type ProblemRecommendation struct {
	ID                    string     `json:"id"`
	StudentID             string     `json:"student_id"`
	ProblemID             string     `json:"problem_id"`
	Category              string     `json:"category"`
	RecommendedDifficulty string     `json:"recommended_difficulty"`
	PrimaryConcept        string     `json:"primary_concept"`
	RelatedConcepts       []string   `json:"related_concepts"`
	MatchScore            float64    `json:"match_score"`
	Headline              string     `json:"headline"`
	Explanation           string     `json:"explanation"`
	LearningObjective     string     `json:"learning_objective"`
	Status                string     `json:"status"`
	IsTopPick             bool       `json:"is_top_pick"`
	FeedbackRating        *string    `json:"feedback_rating,omitempty"`
	FeedbackNotes         *string    `json:"feedback_notes,omitempty"`
	CompletedAt           *time.Time `json:"completed_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`

	// Joined/Hydrated problem metadata for client UI
	ProblemTitle       string   `json:"problem_title,omitempty"`
	ProblemSlug        string   `json:"problem_slug,omitempty"`
	ProblemDifficulty  string   `json:"problem_difficulty,omitempty"`
	ProblemSummary     string   `json:"problem_summary,omitempty"`
	ExternalProblemURL string   `json:"external_problem_url,omitempty"`
	TopicName          string   `json:"topic_name,omitempty"`
	SubtopicName       string   `json:"subtopic_name,omitempty"`
	Tags               []string `json:"tags,omitempty"`
	AttemptStatus      string   `json:"attempt_status,omitempty"` // Not Started, Attempted, Solved, Needs Revision
}

// CategorizedRecommendations represents grouped recommendations for the dashboard.
type CategorizedRecommendations struct {
	TopPick               *ProblemRecommendation             `json:"top_pick,omitempty"`
	Categories            map[string][]ProblemRecommendation `json:"categories"`
	TotalRecommendations  int                                `json:"total_recommendations"`
	RecommendedDifficulty string                             `json:"recommended_difficulty"`
	PrimaryFocusConcept   string                             `json:"primary_focus_concept"`
	GeneratedAt           time.Time                          `json:"generated_at"`
}

// FeedbackRequest payload when a student rates a recommendation.
type FeedbackRequest struct {
	Rating string `json:"rating"` // helpful, not_relevant, too_easy, too_difficult
	Notes  string `json:"notes,omitempty"`
}

// HistoryResponse returned by GET /api/v1/student/recommendations/history
type HistoryResponse struct {
	History []ProblemRecommendation `json:"history"`
	Total   int                     `json:"total"`
	Limit   int                     `json:"limit"`
	Offset  int                     `json:"offset"`
}

// ProblemScoringCandidate is used internally during deterministic ranking.
type ProblemScoringCandidate struct {
	Problem               model.Problem
	PrimaryConcept        string
	RelatedConcepts       []string
	Category              string
	Score                 float64
	Headline              string
	Explanation           string
	LearningObjective     string
	RecommendedDifficulty string
	IsTopPick             bool
}

