package execution

import "time"

const (
	StatusQueued               = "Queued"
	StatusProcessing           = "Processing"
	StatusAccepted             = "Accepted"
	StatusWrongAnswer          = "Wrong Answer"
	StatusCompilationError     = "Compilation Error"
	StatusRuntimeError         = "Runtime Error"
	StatusTimeLimitExceeded    = "Time Limit Exceeded"
	StatusMemoryLimitExceeded  = "Memory Limit Exceeded"
	StatusInternalError        = "Internal Error"
)

// ExecutionRequest represents the code payload sent to the sandbox.
type ExecutionRequest struct {
	Language       string        `json:"language"`
	LanguageID     int           `json:"language_id"`
	SourceCode     string        `json:"source_code"`
	Stdin          string        `json:"stdin"`
	ExpectedOutput string        `json:"expected_output"`
	CPUTimeLimit   time.Duration `json:"cpu_time_limit"`
	MemoryLimitKB  int           `json:"memory_limit_kb"`
}

// ExecutionResponse represents the raw execution metrics from the sandbox.
type ExecutionResponse struct {
	Token          string `json:"token"`
	StatusID       int    `json:"status_id"`
	Status         string `json:"status"` // LearnX normalized status
	Stdout         string `json:"stdout"`
	Stderr         string `json:"stderr"`
	CompileOutput  string `json:"compile_output"`
	Message        string `json:"message"`
	TimeMs         int    `json:"time_ms"`
	MemoryKB       int    `json:"memory_kb"`
	ProviderStatus string `json:"provider_status"`
}
