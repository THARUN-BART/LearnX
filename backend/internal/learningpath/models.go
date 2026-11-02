package learningpath

import "time"

// Node progression states
const (
	StateNotStarted           = "not_started"
	StateReadyToLearn         = "ready_to_learn"
	StateInProgress           = "in_progress"
	StatePracticing           = "practicing"
	StateNeedsRevision        = "needs_revision"
	StateMastered             = "mastered"
	StateLockedByPrerequisite = "locked_by_prerequisite"
)

// Action types
const (
	ActionLearnConcept        = "learn_concept"
	ActionReviseConcept       = "revise_concept"
	ActionSolveEasy           = "solve_easy"
	ActionPracticeMedium      = "practice_medium"
	ActionAttemptHard         = "attempt_hard"
	ActionReviewMistake       = "review_mistake"
	ActionCompletePrereq      = "complete_prerequisite"
)

// AdaptivePath models the high-level roadmap container.
type AdaptivePath struct {
	ID                     string     `json:"id"`
	StudentID              string     `json:"student_id"`
	CurrentTopicID         *string    `json:"current_topic_id,omitempty"`
	CurrentConceptKey      string     `json:"current_concept_key"`
	RecommendedDifficulty  string     `json:"recommended_difficulty"`
	OverallReadinessScore  float64    `json:"overall_readiness_score"`
	LearningStage          string     `json:"learning_stage"`
	Nodes                  []PathNode `json:"nodes"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// PathNode models a concept position on the adaptive curriculum graph.
type PathNode struct {
	ID                     string    `json:"id"`
	StudentID              string    `json:"student_id"`
	TopicID                *string   `json:"topic_id,omitempty"`
	TopicName              string    `json:"topic_name"`
	ConceptKey             string    `json:"concept_key"`
	ConceptName            string    `json:"concept_name"`
	NodeOrder              int       `json:"node_order"`
	State                  string    `json:"state"`
	PrerequisiteKeys       []string  `json:"prerequisite_concept_keys"`
	ReadinessScore         float64   `json:"readiness_score"`
	RecommendedDifficulty  string    `json:"recommended_difficulty"`
	RecommendedProblemID   *string   `json:"recommended_problem_id,omitempty"`
	RecommendedProblemName *string   `json:"recommended_problem_name,omitempty"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// LearningRecommendation answers "What should this student learn next, and why?"
type LearningRecommendation struct {
	ID                    string                 `json:"id"`
	StudentID             string                 `json:"student_id"`
	ActionType            string                 `json:"action_type"`
	ConceptKey            string                 `json:"concept_key"`
	ConceptName           string                 `json:"concept_name"`
	TopicID               *string                `json:"topic_id,omitempty"`
	TopicName             string                 `json:"topic_name,omitempty"`
	ProblemID             *string                `json:"problem_id,omitempty"`
	ProblemTitle          string                 `json:"problem_title,omitempty"`
	ProblemSlug           string                 `json:"problem_slug,omitempty"`
	RecommendedDifficulty string                 `json:"recommended_difficulty"`
	Headline              string                 `json:"headline"`
	Explanation           string                 `json:"explanation"`
	RelevantWeakness      string                 `json:"relevant_weakness,omitempty"`
	Evidence              map[string]interface{} `json:"evidence"`
	IsActive              bool                   `json:"is_active"`
	CompletedAt           *time.Time             `json:"completed_at,omitempty"`
	CreatedAt             time.Time              `json:"created_at"`
}

// PathProgressResponse returns statistics for dashboard visualization.
type PathProgressResponse struct {
	TotalConcepts    int       `json:"total_concepts"`
	MasteredCount    int       `json:"mastered_count"`
	InProgressCount  int       `json:"in_progress_count"`
	NeedsRevisionCount int     `json:"needs_revision_count"`
	ReadyToLearnCount int      `json:"ready_to_learn_count"`
	LockedCount      int       `json:"locked_count"`
	OverallProgress  float64   `json:"overall_progress_percentage"`
	CurrentStage     string    `json:"current_stage"`
	LastEvaluatedAt  time.Time `json:"last_evaluated_at"`
}
