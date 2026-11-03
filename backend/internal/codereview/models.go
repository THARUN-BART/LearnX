package codereview

import "time"

type IssueSeverity string

const (
	SeverityCritical IssueSeverity = "critical"
	SeverityWarning  IssueSeverity = "warning"
	SeverityInfo     IssueSeverity = "info"
)

type CodeIssue struct {
	ID             string        `json:"id"`
	Severity       IssueSeverity `json:"severity"` // critical, warning, info
	LineNumber     int           `json:"line_number"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
	Recommendation string        `json:"recommendation"`
	CodeSnippet    string        `json:"code_snippet,omitempty"`
	Category       string        `json:"category"` // syntax, performance, readability, safety, edge_case
	SourceType     string        `json:"source_type"` // static_analysis, execution_observation, ai_synthesis
}

type Scores struct {
	Overall         int `json:"overall"`
	Readability     int `json:"readability"`
	Efficiency      int `json:"efficiency"`
	Maintainability int `json:"maintainability"`
}

type ComplexityComparison struct {
	OriginalTime   string `json:"original_time"`
	OptimizedTime  string `json:"optimized_time"`
	OriginalSpace  string `json:"original_space"`
	OptimizedSpace string `json:"optimized_space"`
	Explanation    string `json:"explanation"`
}

type AlternativeApproach struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Space       string `json:"space"`
	TradeOff    string `json:"trade_off"`
}

type CodeReviewRequest struct {
	ProblemID    string `json:"problem_id,omitempty"`
	SubmissionID string `json:"submission_id,omitempty"`
	Language     string `json:"language"`
	SourceCode   string `json:"source_code"`
	Mode         string `json:"mode,omitempty"` // standard, deep
}

type CodeReviewResult struct {
	ID                          string                 `json:"id"`
	StudentID                   string                 `json:"student_id"`
	ProblemID                   string                 `json:"problem_id,omitempty"`
	ProblemTitle                string                 `json:"problem_title,omitempty"`
	SubmissionID                string                 `json:"submission_id,omitempty"`
	Language                    string                 `json:"language"`
	SourceCode                  string                 `json:"source_code"`
	Scores                      Scores                 `json:"scores"`
	CorrectnessObservations     []string               `json:"correctness_observations"`
	Issues                      []CodeIssue            `json:"issues"`
	OptimizedCode               string                 `json:"optimized_code"`
	ExplanationOfImprovements   string                 `json:"explanation_of_improvements"`
	Complexity                  ComplexityComparison   `json:"complexity"`
	AlternativeApproaches       []AlternativeApproach  `json:"alternative_approaches"`
	RecommendedVisualizerAlgoID string                 `json:"recommended_visualizer_algo_id,omitempty"`
	RecommendedVisualizerReason string                 `json:"recommended_visualizer_reason,omitempty"`
	CreatedAt                   time.Time              `json:"created_at"`
}

type ReviewHistoryItem struct {
	ID           string    `json:"id"`
	ProblemID    string    `json:"problem_id,omitempty"`
	ProblemTitle string    `json:"problem_title,omitempty"`
	Language     string    `json:"language"`
	OverallScore int       `json:"overall_score"`
	CreatedAt    time.Time `json:"created_at"`
}

type CodeComparison struct {
	OriginalCode  string `json:"original_code"`
	OptimizedCode string `json:"optimized_code"`
	Summary       string `json:"summary"`
	Complexity    ComplexityComparison `json:"complexity"`
}

