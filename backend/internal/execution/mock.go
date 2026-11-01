package execution

import (
	"context"
	"strings"
	"sync"
)

type MockExecutionClient struct {
	mu            sync.Mutex
	ForcedStatus  string
	CustomExecute func(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error)
	Calls         []ExecutionRequest
}

func NewMockExecutionClient() *MockExecutionClient {
	return &MockExecutionClient{
		Calls: make([]ExecutionRequest, 0),
	}
}

func (m *MockExecutionClient) ProviderName() string {
	return "mock"
}

func (m *MockExecutionClient) Execute(ctx context.Context, req ExecutionRequest) (*ExecutionResponse, error) {
	m.mu.Lock()
	m.Calls = append(m.Calls, req)
	forced := m.ForcedStatus
	custom := m.CustomExecute
	m.mu.Unlock()

	if custom != nil {
		return custom(ctx, req)
	}

	code := req.SourceCode

	// Built-in mock rules based on source code content or forced status
	if forced != "" {
		switch forced {
		case StatusCompilationError:
			return &ExecutionResponse{
				Status:        StatusCompilationError,
				CompileOutput: "syntax error: unexpected token",
			}, nil
		case StatusRuntimeError:
			return &ExecutionResponse{
				Status: StatusRuntimeError,
				Stderr: "panic: runtime error: index out of range",
			}, nil
		case StatusTimeLimitExceeded:
			return &ExecutionResponse{
				Status: StatusTimeLimitExceeded,
				TimeMs: 5000,
			}, nil
		case StatusMemoryLimitExceeded:
			return &ExecutionResponse{
				Status:   StatusMemoryLimitExceeded,
				MemoryKB: 512000,
			}, nil
		case StatusInternalError:
			return nil, ErrExecutionServiceUnavailable
		}
	}

	// Code content inspection rules
	if strings.Contains(code, "COMPILE_ERROR") || strings.Contains(code, "syntax_error") {
		return &ExecutionResponse{
			Status:        StatusCompilationError,
			CompileOutput: "error: expected ';' at end of statement",
		}, nil
	}
	if strings.Contains(code, "RUNTIME_ERROR") || strings.Contains(code, "panic(") {
		return &ExecutionResponse{
			Status: StatusRuntimeError,
			Stderr: "Exception in thread 'main' java.lang.NullPointerException",
			TimeMs: 25,
		}, nil
	}
	if strings.Contains(code, "TLE") || strings.Contains(code, "infinite_loop") {
		return &ExecutionResponse{
			Status: StatusTimeLimitExceeded,
			TimeMs: 5000,
		}, nil
	}
	if strings.Contains(code, "MLE") {
		return &ExecutionResponse{
			Status:   StatusMemoryLimitExceeded,
			MemoryKB: 512000,
		}, nil
	}
	if strings.Contains(code, "WRONG_ANSWER") {
		return &ExecutionResponse{
			Status:   StatusAccepted, // Judge0 executes and returns stdout
			Stdout:   "incorrect output 999",
			TimeMs:   35,
			MemoryKB: 2048,
		}, nil
	}

	// Default behavior: mimic ideal execution returning the expected output
	expected := req.ExpectedOutput
	if expected == "" {
		expected = "1"
	}

	return &ExecutionResponse{
		Token:    "mock-token-123",
		StatusID: 3,
		Status:   StatusAccepted,
		Stdout:   expected,
		TimeMs:   20,
		MemoryKB: 1024,
	}, nil
}
