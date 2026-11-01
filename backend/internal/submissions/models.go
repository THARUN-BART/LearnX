package submissions

import "time"

type Submission struct {
	ID                   string     `json:"id"`
	StudentID            string     `json:"student_id"`
	ProblemID            string     `json:"problem_id"`
	Language             string     `json:"language"`
	SourceCode           string     `json:"source_code,omitempty"`
	Status               string     `json:"status"`
	ExecutionTimeMs      int        `json:"execution_time_ms"`
	MemoryUsedKB         int        `json:"memory_used_kb"`
	TotalTestCases       int        `json:"total_test_cases"`
	PassedTestCases      int        `json:"passed_test_cases"`
	FailedTestCases      int        `json:"failed_test_cases"`
	ErrorMessage         string     `json:"error_message,omitempty"`
	ExecutionProvider    string     `json:"execution_provider"`
	ProviderSubmissionID string     `json:"provider_submission_id,omitempty"`
	SubmittedAt          time.Time  `json:"submitted_at"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
}

type SubmissionTestResult struct {
	ID              string    `json:"id"`
	SubmissionID    string    `json:"submission_id"`
	TestCaseID      string    `json:"test_case_id"`
	Status          string    `json:"status"`
	ExecutionTimeMs int       `json:"execution_time_ms"`
	MemoryUsedKB    int       `json:"memory_used_kb"`
	ActualOutput    string    `json:"actual_output,omitempty"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	IsHidden        bool      `json:"is_hidden,omitempty"`
}

type CreateSubmissionRequest struct {
	ProblemID  string `json:"problem_id"`
	Language   string `json:"language"`
	SourceCode string `json:"source_code"`
}

type TestCasesSummary struct {
	Total  int `json:"total"`
	Passed int `json:"passed"`
	Failed int `json:"failed"`
}

type SubmissionSummaryResponse struct {
	SubmissionID    string           `json:"submission_id"`
	ProblemID       string           `json:"problem_id"`
	Status          string           `json:"status"`
	Language        string           `json:"language"`
	TestCases       TestCasesSummary `json:"test_cases"`
	ExecutionTimeMs int              `json:"execution_time_ms"`
	MemoryUsedKB    int              `json:"memory_used_kb"`
	ErrorMessage    string           `json:"error_message,omitempty"`
	SubmittedAt     time.Time        `json:"submitted_at"`
	CompletedAt     *time.Time       `json:"completed_at,omitempty"`
}

type PublicTestResult struct {
	TestCaseOrder   int    `json:"test_case_order"`
	Status          string `json:"status"`
	ExecutionTimeMs int    `json:"execution_time_ms"`
	MemoryUsedKB    int    `json:"memory_used_kb"`
	IsHidden        bool   `json:"is_hidden"`
	ActualOutput    string `json:"actual_output,omitempty"` // Omitted if hidden!
	ErrorMessage    string `json:"error_message,omitempty"`
}

type SubmissionDetailedResponse struct {
	SubmissionID    string             `json:"submission_id"`
	ProblemID       string             `json:"problem_id"`
	Status          string             `json:"status"`
	Language        string             `json:"language"`
	SourceCode      string             `json:"source_code,omitempty"`
	TestCases       TestCasesSummary   `json:"test_cases"`
	ExecutionTimeMs int                `json:"execution_time_ms"`
	MemoryUsedKB    int                `json:"memory_used_kb"`
	ErrorMessage    string             `json:"error_message,omitempty"`
	SubmittedAt     time.Time          `json:"submitted_at"`
	CompletedAt     *time.Time         `json:"completed_at,omitempty"`
	Results         []PublicTestResult `json:"results"`
}

type PaginatedSubmissions struct {
	Submissions []SubmissionSummaryResponse `json:"submissions"`
	Total       int                         `json:"total"`
	Page        int                         `json:"page"`
	Limit       int                         `json:"limit"`
	TotalPages  int                         `json:"total_pages"`
}
