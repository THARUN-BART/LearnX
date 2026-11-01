package assessment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
)

type Engine interface {
	Evaluate(ctx context.Context, lang execution.LanguageConfig, sourceCode string, tests []testcases.TestCase) (*EvaluationResult, error)
}

type DefaultEngine struct {
	client     execution.Client
	comparator *Comparator
	logger     *slog.Logger
}

func NewEngine(client execution.Client, logger *slog.Logger) *DefaultEngine {
	if logger == nil {
		logger = slog.Default()
	}
	return &DefaultEngine{
		client:     client,
		comparator: NewComparator(),
		logger:     logger,
	}
}

func (e *DefaultEngine) Evaluate(ctx context.Context, lang execution.LanguageConfig, sourceCode string, tests []testcases.TestCase) (*EvaluationResult, error) {
	result := &EvaluationResult{
		TotalTestCases:      len(tests),
		TestCaseEvaluations: make([]TestCaseEvaluation, 0, len(tests)),
		CompletedAt:         time.Now(),
	}

	if len(tests) == 0 {
		// If no test cases are registered for the problem, mark as Accepted default or Internal Error
		result.OverallStatus = execution.StatusAccepted
		return result, nil
	}

	overallStatus := execution.StatusAccepted
	firstErrorMessage := ""

	for _, tc := range tests {
		req := execution.ExecutionRequest{
			Language:       lang.Slug,
			LanguageID:     lang.ID,
			SourceCode:     sourceCode,
			Stdin:          tc.InputData,
			ExpectedOutput: tc.ExpectedOutput,
			CPUTimeLimit:   time.Duration(lang.DefaultTimeout) * time.Second,
			MemoryLimitKB:  lang.DefaultMemory,
		}

		execResp, err := e.client.Execute(ctx, req)
		if err != nil {
			e.logger.Error("execution provider error", slog.String("error", err.Error()))
			return nil, fmt.Errorf("assessment execution failed: %w", err)
		}

		tcEval := TestCaseEvaluation{
			TestCaseID:      tc.ID,
			IsHidden:        tc.IsHidden,
			ExecutionTimeMs: execResp.TimeMs,
			MemoryUsedKB:    execResp.MemoryKB,
			ActualOutput:    execResp.Stdout,
		}

		if execResp.TimeMs > result.MaxExecutionTimeMs {
			result.MaxExecutionTimeMs = execResp.TimeMs
		}
		if execResp.MemoryKB > result.MaxMemoryUsedKB {
			result.MaxMemoryUsedKB = execResp.MemoryKB
		}

		// Evaluate individual test case result
		switch execResp.Status {
		case execution.StatusCompilationError:
			tcEval.Status = execution.StatusCompilationError
			tcEval.ErrorMessage = execResp.CompileOutput
			if firstErrorMessage == "" {
				firstErrorMessage = execResp.CompileOutput
			}
			result.TestCaseEvaluations = append(result.TestCaseEvaluations, tcEval)
			result.FailedTestCases = len(tests)
			result.PassedTestCases = 0
			result.OverallStatus = execution.StatusCompilationError
			result.ErrorMessage = firstErrorMessage
			// Compilation error fails all test cases immediately
			return result, nil

		case execution.StatusRuntimeError:
			tcEval.Status = execution.StatusRuntimeError
			errMsg := execResp.Stderr
			if errMsg == "" {
				errMsg = execResp.Message
			}
			tcEval.ErrorMessage = errMsg
			if firstErrorMessage == "" {
				firstErrorMessage = errMsg
			}
			result.FailedTestCases++
			if overallStatus == execution.StatusAccepted || overallStatus == execution.StatusWrongAnswer {
				overallStatus = execution.StatusRuntimeError
			}

		case execution.StatusTimeLimitExceeded:
			tcEval.Status = execution.StatusTimeLimitExceeded
			tcEval.ErrorMessage = "Time Limit Exceeded"
			result.FailedTestCases++
			if overallStatus == execution.StatusAccepted || overallStatus == execution.StatusWrongAnswer {
				overallStatus = execution.StatusTimeLimitExceeded
			}

		case execution.StatusMemoryLimitExceeded:
			tcEval.Status = execution.StatusMemoryLimitExceeded
			tcEval.ErrorMessage = "Memory Limit Exceeded"
			result.FailedTestCases++
			if overallStatus == execution.StatusAccepted || overallStatus == execution.StatusWrongAnswer {
				overallStatus = execution.StatusMemoryLimitExceeded
			}

		case execution.StatusInternalError:
			tcEval.Status = execution.StatusInternalError
			tcEval.ErrorMessage = execResp.Message
			result.FailedTestCases++
			overallStatus = execution.StatusInternalError

		default:
			// Program completed: verify output correctness with whitespace normalization
			if e.comparator.Compare(execResp.Stdout, tc.ExpectedOutput) {
				tcEval.Status = execution.StatusAccepted
				result.PassedTestCases++
			} else {
				tcEval.Status = execution.StatusWrongAnswer
				tcEval.ErrorMessage = "Output did not match expected result"
				result.FailedTestCases++
				if overallStatus == execution.StatusAccepted {
					overallStatus = execution.StatusWrongAnswer
				}
			}
		}

		result.TestCaseEvaluations = append(result.TestCaseEvaluations, tcEval)
	}

	result.OverallStatus = overallStatus
	result.ErrorMessage = firstErrorMessage
	return result, nil
}
