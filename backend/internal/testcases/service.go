package testcases

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidTestCase = errors.New("invalid test case: problem ID and expected output are required")

type Service interface {
	GetSampleTestCases(ctx context.Context, problemID string) ([]PublicTestCase, error)
	CreateTestCase(ctx context.Context, problemID string, req CreateTestCaseRequest) (*TestCase, error)
	ListAllTestCasesForEngine(ctx context.Context, problemID string) ([]TestCase, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetSampleTestCases(ctx context.Context, problemID string) ([]PublicTestCase, error) {
	problemID = strings.TrimSpace(problemID)
	if problemID == "" {
		return nil, ErrInvalidTestCase
	}
	all, err := s.repo.ListByProblem(ctx, problemID, false)
	if err != nil {
		return nil, err
	}
	publicList := make([]PublicTestCase, 0, len(all))
	for _, tc := range all {
		publicList = append(publicList, PublicTestCase{
			ID:             tc.ID,
			ProblemID:      tc.ProblemID,
			InputData:      tc.InputData,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       false,
			TestOrder:      tc.TestOrder,
		})
	}
	return publicList, nil
}

func (s *service) CreateTestCase(ctx context.Context, problemID string, req CreateTestCaseRequest) (*TestCase, error) {
	problemID = strings.TrimSpace(problemID)
	if problemID == "" {
		return nil, ErrInvalidTestCase
	}
	tc := &TestCase{
		ProblemID:      problemID,
		InputData:      req.InputData,
		ExpectedOutput: req.ExpectedOutput,
		IsHidden:       req.IsHidden,
		TestOrder:      req.TestOrder,
	}
	if tc.TestOrder <= 0 {
		tc.TestOrder = 1
	}
	if err := s.repo.Create(ctx, tc); err != nil {
		return nil, err
	}
	return tc, nil
}

func (s *service) ListAllTestCasesForEngine(ctx context.Context, problemID string) ([]TestCase, error) {
	return s.repo.ListByProblem(ctx, problemID, true)
}
