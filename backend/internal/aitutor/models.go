package aitutor

import "time"

// Progressive Hint Levels
const (
	HintLevelConcept     = 1
	HintLevelStructural  = 2
	HintLevelEdgeCase    = 3
	HintLevelBlueprint   = 4
	HintLevelCodeFix     = 5
)

// Request / Response structures for AI Tutor features

type HintRequest struct {
	ProblemID    string `json:"problem_id"`
	SubmissionID string `json:"submission_id,omitempty"`
	Level        int    `json:"level"` // 1 to 5
}

type HintResponse struct {
	ProblemID       string `json:"problem_id"`
	Level           int    `json:"level"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	GuidingQuestion string `json:"guiding_question,omitempty"`
	NextLevelAvailable bool `json:"next_level_available"`
	CreatedAt       time.Time `json:"created_at"`
}

type ExplainMistakeRequest struct {
	SubmissionID string `json:"submission_id"`
}

type ExplainMistakeResponse struct {
	SubmissionID       string   `json:"submission_id"`
	Category           string   `json:"category"`
	Subtype            string   `json:"subtype"`
	WhatWentWrong      string   `json:"what_went_wrong"`
	WhyItHappened      string   `json:"why_it_happened"`
	HowToFix           string   `json:"how_to_fix"`
	RelatedDSAConcepts []string `json:"related_dsa_concepts"`
}

type CodeWalkthroughRequest struct {
	ProblemID  string `json:"problem_id"`
	SourceCode string `json:"source_code"`
	Language   string `json:"language"`
}

type WalkthroughSection struct {
	Title       string `json:"title"`
	LineRange   string `json:"line_range,omitempty"`
	Observation string `json:"observation"`
	Evaluation  string `json:"evaluation"` // positive, warning, critical
}

type CodeWalkthroughResponse struct {
	OverallSummary string               `json:"overall_summary"`
	Sections       []WalkthroughSection `json:"sections"`
	KeyStrengths   []string             `json:"key_strengths"`
	AreasToImprove []string             `json:"areas_to_improve"`
}

type DryRunRequest struct {
	ProblemID    string `json:"problem_id"`
	SourceCode   string `json:"source_code"`
	Language     string `json:"language"`
	CustomInput  string `json:"custom_input,omitempty"`
}

type DryRunStep struct {
	StepNumber  int                    `json:"step_number"`
	LineNumber  int                    `json:"line_number"`
	Action      string                 `json:"action"`
	Variables   map[string]interface{} `json:"variables"`
	Explanation string                 `json:"explanation"`
}

type DryRunResponse struct {
	TestCaseInput string       `json:"test_case_input"`
	Steps         []DryRunStep `json:"steps"`
	FinalOutput   string       `json:"final_output"`
}

type ComplexityRequest struct {
	ProblemID  string `json:"problem_id"`
	SourceCode string `json:"source_code"`
	Language   string `json:"language"`
}

type ComplexityResponse struct {
	TimeComplexity     string `json:"time_complexity"`
	TimeExplanation    string `json:"time_explanation"`
	SpaceComplexity    string `json:"space_complexity"`
	SpaceExplanation   string `json:"space_explanation"`
	IsOptimal          bool   `json:"is_optimal"`
	OptimisationNotes  string `json:"optimisation_notes,omitempty"`
}

type AskRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	ProblemID      string `json:"problem_id,omitempty"`
	Question       string `json:"question"`
}

type AskResponse struct {
	ConversationID string    `json:"conversation_id"`
	Answer         string    `json:"answer"`
	CreatedAt      time.Time `json:"created_at"`
}

// Conversation and Message entities

type AIConversation struct {
	ID               string    `json:"id"`
	StudentID        string    `json:"student_id"`
	ProblemID        *string   `json:"problem_id,omitempty"`
	ConversationType string    `json:"conversation_type"`
	Title            string    `json:"title"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type AIMessage struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Role           string                 `json:"role"`
	Message        string                 `json:"message"`
	ResponseType   string                 `json:"response_type"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	TokenUsage     map[string]interface{} `json:"token_usage,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
}

type AIHintHistory struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	ProblemID    string    `json:"problem_id"`
	SubmissionID *string   `json:"submission_id,omitempty"`
	HintLevel    int       `json:"hint_level"`
	HintContent  string    `json:"hint_content"`
	CreatedAt    time.Time `json:"created_at"`
}
