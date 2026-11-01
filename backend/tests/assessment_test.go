package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/assessment"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
)

func TestAssessment_Comparator(t *testing.T) {
	comp := assessment.NewComparator()

	tests := []struct {
		name     string
		actual   string
		expected string
		matches  bool
	}{
		{
			name:     "exact match",
			actual:   "42",
			expected: "42",
			matches:  true,
		},
		{
			name:     "crlf vs lf normalization",
			actual:   "hello\r\nworld\r\n",
			expected: "hello\nworld",
			matches:  true,
		},
		{
			name:     "trailing line spaces and extra newlines",
			actual:   "1 2 3   \n4 5 6\n\n\n",
			expected: "1 2 3\n4 5 6",
			matches:  true,
		},
		{
			name:     "different values",
			actual:   "10",
			expected: "20",
			matches:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := comp.Compare(tt.actual, tt.expected); got != tt.matches {
				t.Errorf("Compare(%q, %q) = %v, expected %v", tt.actual, tt.expected, got, tt.matches)
			}
		})
	}
}

func TestAssessment_EngineEvaluationClassifications(t *testing.T) {
	langCfg := execution.LanguageConfig{
		ID:             54,
		Slug:           "cpp",
		DefaultTimeout: 2,
		DefaultMemory:  128000,
	}

	testCasesList := []testcases.TestCase{
		{ID: "tc-1", InputData: "2 3", ExpectedOutput: "5", IsHidden: false, TestOrder: 1},
		{ID: "tc-2", InputData: "10 20", ExpectedOutput: "30", IsHidden: true, TestOrder: 2},
	}

	ctx := context.Background()

	// 1. Accepted scenario
	t.Run("All Passed Accepted", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "#include <iostream>\nint main(){}", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusAccepted {
			t.Errorf("expected Accepted, got %s", res.OverallStatus)
		}
		if res.PassedTestCases != 2 || res.FailedTestCases != 0 {
			t.Errorf("expected 2 passed, 0 failed, got %d/%d", res.PassedTestCases, res.FailedTestCases)
		}
	})

	// 2. Compilation Error scenario
	t.Run("Compilation Error", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "COMPILE_ERROR_CODE", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusCompilationError {
			t.Errorf("expected Compilation Error, got %s", res.OverallStatus)
		}
		if res.ErrorMessage == "" {
			t.Error("expected error message to contain compiler output")
		}
	})

	// 3. Runtime Error scenario
	t.Run("Runtime Error", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "RUNTIME_ERROR_TRIGGER", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusRuntimeError {
			t.Errorf("expected Runtime Error, got %s", res.OverallStatus)
		}
	})

	// 4. Time Limit Exceeded scenario
	t.Run("Time Limit Exceeded", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "TLE_LOOP", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusTimeLimitExceeded {
			t.Errorf("expected Time Limit Exceeded, got %s", res.OverallStatus)
		}
	})

	// 5. Memory Limit Exceeded scenario
	t.Run("Memory Limit Exceeded", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "MLE_CODE", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusMemoryLimitExceeded {
			t.Errorf("expected Memory Limit Exceeded, got %s", res.OverallStatus)
		}
	})

	// 6. Wrong Answer scenario
	t.Run("Wrong Answer", func(t *testing.T) {
		mockClient := execution.NewMockExecutionClient()
		engine := assessment.NewEngine(mockClient, nil)

		res, err := engine.Evaluate(ctx, langCfg, "WRONG_ANSWER_CODE", testCasesList)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res.OverallStatus != execution.StatusWrongAnswer {
			t.Errorf("expected Wrong Answer, got %s", res.OverallStatus)
		}
		if res.FailedTestCases != 2 {
			t.Errorf("expected 2 failed test cases, got %d", res.FailedTestCases)
		}
	})
}

func TestExecution_Judge0HTTPClient(t *testing.T) {
	// Mock Judge0 HTTP server using net/http/httptest
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-RapidAPI-Key") != "test-api-key" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/submissions/" {
			out := "42\n"
			timeStr := "0.015"
			mem := 1540
			resp := map[string]interface{}{
				"token":  "token-xyz",
				"stdout": out,
				"time":   timeStr,
				"memory": mem,
				"status": map[string]interface{}{
					"id":          3,
					"description": "Accepted",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer ts.Close()

	client := execution.NewJudge0Client(execution.Judge0Config{
		BaseURL: ts.URL,
		APIKey:  "test-api-key",
		APIHost: "judge0.p.rapidapi.com",
		Timeout: 5 * time.Second,
	})

	req := execution.ExecutionRequest{
		Language:       "cpp",
		LanguageID:     54,
		SourceCode:     "#include <iostream>",
		Stdin:          "input",
		ExpectedOutput: "42",
	}

	res, err := client.Execute(context.Background(), req)
	if err != nil {
		t.Fatalf("Judge0 execution failed: %v", err)
	}

	if res.Status != execution.StatusAccepted {
		t.Errorf("expected Accepted, got %s", res.Status)
	}
	if res.TimeMs != 15 {
		t.Errorf("expected 15ms, got %d", res.TimeMs)
	}
	if res.MemoryKB != 1540 {
		t.Errorf("expected 1540KB, got %d", res.MemoryKB)
	}
}
