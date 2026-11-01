package submissions

import (
	"errors"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/execution"
)

var (
	ErrInvalidProblemID = errors.New("invalid problem_id: must not be empty")
	ErrEmptySourceCode  = errors.New("source_code cannot be empty")
	ErrCodeTooLarge     = errors.New("source_code exceeds maximum allowed size of 64KB")
	ErrInvalidLanguage  = errors.New("unsupported programming language")
)

func ValidateSubmissionRequest(req CreateSubmissionRequest) (execution.LanguageConfig, error) {
	if strings.TrimSpace(req.ProblemID) == "" {
		return execution.LanguageConfig{}, ErrInvalidProblemID
	}
	if strings.TrimSpace(req.SourceCode) == "" {
		return execution.LanguageConfig{}, ErrEmptySourceCode
	}
	if len(req.SourceCode) > execution.MaxSourceCodeBytes {
		return execution.LanguageConfig{}, ErrCodeTooLarge
	}

	langCfg, err := execution.GetRegistry().Resolve(req.Language)
	if err != nil {
		return execution.LanguageConfig{}, ErrInvalidLanguage
	}

	return langCfg, nil
}
