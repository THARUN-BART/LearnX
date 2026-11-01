package assessment

import "time"

type TestCaseEvaluation struct {
	TestCaseID      string `json:"test_case_id"`
	IsHidden        bool   `json:"is_hidden"`
	Status          string `json:"status"` // Accepted, Wrong Answer, Time Limit Exceeded, etc.
	ExecutionTimeMs int    `json:"execution_time_ms"`
	MemoryUsedKB    int    `json:"memory_used_kb"`
	ActualOutput    string `json:"actual_output,omitempty"` // populated only if safe / sample
	ErrorMessage    string `json:"error_message,omitempty"`
}

type EvaluationResult struct {
	OverallStatus       string               `json:"overall_status"`
	TotalTestCases      int                  `json:"total_test_cases"`
	PassedTestCases     int                  `json:"passed_test_cases"`
	FailedTestCases     int                  `json:"failed_test_cases"`
	MaxExecutionTimeMs  int                  `json:"max_execution_time_ms"`
	MaxMemoryUsedKB     int                  `json:"max_memory_used_kb"`
	ErrorMessage        string               `json:"error_message"`
	TestCaseEvaluations []TestCaseEvaluation `json:"test_case_evaluations"`
	CompletedAt         time.Time            `json:"completed_at"`
}
