package pedagogy

// AnalyzeRequest represents the payload expected by the Pedagogy service /analyze endpoint.
type AnalyzeRequest struct {
	StudentCode string  `json:"student_code"`
	TestCode    string  `json:"test_code"`
	Skill       string  `json:"skill,omitempty"`
	HintLevel   string  `json:"hint_level,omitempty"`
	Timeout     float64 `json:"timeout,omitempty"`
}

// RunnerResult captures the execution output from Person_C.
type RunnerResult struct {
	Passed    bool   `json:"passed"`
	ErrorType string `json:"error_type,omitempty"`
	Output    string `json:"output,omitempty"`
	Traceback string `json:"traceback,omitempty"`
}

// Diagnosis contains the diagnosed misconception.
type Diagnosis struct {
	Skill           string  `json:"skill"`
	MisconceptionID string  `json:"misconception_id"`
	Misconception   string  `json:"misconception"`
	Confidence      float64 `json:"confidence"`
	Evidence        string  `json:"evidence"`
	DiagnosisSource string  `json:"diagnosis_source,omitempty"`
}

// Hint contains the generated hint.
type Hint struct {
	HintLevel  string `json:"hint_level"`
	Hint       string `json:"hint"`
	HintSource string `json:"hint_source,omitempty"`
	Error      string `json:"error,omitempty"`
}

// AnalyzeResponse is the complete response from /analyze.
type AnalyzeResponse struct {
	RunnerResult RunnerResult `json:"runner_result"`
	Diagnosis    Diagnosis    `json:"diagnosis"`
	Hint         Hint         `json:"hint"`
}

// DiagnoseResponse is the response from /diagnose.
type DiagnoseResponse struct {
	RunnerResult RunnerResult `json:"runner_result"`
	Diagnosis    Diagnosis    `json:"diagnosis"`
}

// HealthResponse is the response from /health.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
