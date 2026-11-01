package submissions

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
)

var (
	ErrProblemNotFound = errors.New("problem not found")
)

// JobQueue is an interface satisfied by any background worker queue.
type JobQueue interface {
	Enqueue(submissionID string) error
}

type Service interface {
	Submit(ctx context.Context, studentID string, req CreateSubmissionRequest) (*SubmissionSummaryResponse, error)
	GetStatus(ctx context.Context, studentID, submissionID string) (*SubmissionSummaryResponse, error)
	GetResult(ctx context.Context, studentID, submissionID string) (*SubmissionDetailedResponse, error)
	ListStudentSubmissions(ctx context.Context, studentID, problemID, status string, page, limit int) (*PaginatedSubmissions, error)
	GetProblemSubmissions(ctx context.Context, studentID, problemID string, page, limit int) (*PaginatedSubmissions, error)
	GetLatestProblemSubmission(ctx context.Context, studentID, problemID string) (*SubmissionSummaryResponse, error)
	Retry(ctx context.Context, studentID, submissionID string) (*SubmissionSummaryResponse, error)
}

type service struct {
	repo         Repository
	problemsRepo problems.Repository
	queue        JobQueue
}

func NewService(repo Repository, problemsRepo problems.Repository, queue JobQueue) Service {
	return &service{
		repo:         repo,
		problemsRepo: problemsRepo,
		queue:        queue,
	}
}

func (s *service) Submit(ctx context.Context, studentID string, req CreateSubmissionRequest) (*SubmissionSummaryResponse, error) {
	if strings.TrimSpace(studentID) == "" {
		return nil, ErrUnauthorized
	}

	langCfg, err := ValidateSubmissionRequest(req)
	if err != nil {
		return nil, err
	}

	// Verify problem exists
	prob, err := s.problemsRepo.Get(ctx, req.ProblemID)
	if err != nil {
		return nil, ErrProblemNotFound
	}

	sub := &Submission{
		StudentID:         studentID,
		ProblemID:         prob.ID,
		Language:          langCfg.Slug,
		SourceCode:        req.SourceCode,
		Status:            execution.StatusQueued,
		ExecutionProvider: "judge0",
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("create submission: %w", err)
	}

	// Queue asynchronous background assessment job
	if err := s.queue.Enqueue(sub.ID); err != nil {
		// Log warning, job can still be processed by recovery worker
		_ = err
	}

	return s.toSummaryResponse(sub), nil
}

func (s *service) GetStatus(ctx context.Context, studentID, submissionID string) (*SubmissionSummaryResponse, error) {
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.StudentID != studentID {
		return nil, ErrUnauthorized
	}
	return s.toSummaryResponse(sub), nil
}

func (s *service) GetResult(ctx context.Context, studentID, submissionID string) (*SubmissionDetailedResponse, error) {
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.StudentID != studentID {
		return nil, ErrUnauthorized
	}

	rawResults, err := s.repo.GetTestResultsBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, fmt.Errorf("get test results: %w", err)
	}

	publicResults := make([]PublicTestResult, 0, len(rawResults))
	for idx, res := range rawResults {
		item := PublicTestResult{
			TestCaseOrder:   idx + 1,
			Status:          res.Status,
			ExecutionTimeMs: res.ExecutionTimeMs,
			MemoryUsedKB:    res.MemoryUsedKB,
			IsHidden:        res.IsHidden,
		}

		if res.IsHidden {
			// Security requirement: Never leak hidden test case outputs or error traces to students!
			item.ActualOutput = ""
			if res.Status != execution.StatusAccepted {
				item.ErrorMessage = "Hidden test case failed"
			}
		} else {
			item.ActualOutput = res.ActualOutput
			item.ErrorMessage = res.ErrorMessage
		}
		publicResults = append(publicResults, item)
	}

	return &SubmissionDetailedResponse{
		SubmissionID: sub.ID,
		ProblemID:    sub.ProblemID,
		Status:       sub.Status,
		Language:     sub.Language,
		SourceCode:   sub.SourceCode,
		TestCases: TestCasesSummary{
			Total:  sub.TotalTestCases,
			Passed: sub.PassedTestCases,
			Failed: sub.FailedTestCases,
		},
		ExecutionTimeMs: sub.ExecutionTimeMs,
		MemoryUsedKB:    sub.MemoryUsedKB,
		ErrorMessage:    sub.ErrorMessage,
		SubmittedAt:     sub.SubmittedAt,
		CompletedAt:     sub.CompletedAt,
		Results:         publicResults,
	}, nil
}

func (s *service) ListStudentSubmissions(ctx context.Context, studentID, problemID, status string, page, limit int) (*PaginatedSubmissions, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	subs, total, err := s.repo.ListByStudent(ctx, studentID, problemID, status, limit, offset)
	if err != nil {
		return nil, err
	}

	summaries := make([]SubmissionSummaryResponse, 0, len(subs))
	for _, item := range subs {
		summaries = append(summaries, *s.toSummaryResponse(&item))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &PaginatedSubmissions{
		Submissions: summaries,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
	}, nil
}

func (s *service) GetProblemSubmissions(ctx context.Context, studentID, problemID string, page, limit int) (*PaginatedSubmissions, error) {
	return s.ListStudentSubmissions(ctx, studentID, problemID, "", page, limit)
}

func (s *service) GetLatestProblemSubmission(ctx context.Context, studentID, problemID string) (*SubmissionSummaryResponse, error) {
	sub, err := s.repo.GetLatestByStudentAndProblem(ctx, studentID, problemID)
	if err != nil {
		return nil, err
	}
	return s.toSummaryResponse(sub), nil
}

func (s *service) Retry(ctx context.Context, studentID, submissionID string) (*SubmissionSummaryResponse, error) {
	sub, err := s.repo.GetByID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if sub.StudentID != studentID {
		return nil, ErrUnauthorized
	}

	sub.Status = execution.StatusQueued
	sub.ErrorMessage = ""
	sub.CompletedAt = nil
	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, fmt.Errorf("reset submission for retry: %w", err)
	}

	_ = s.queue.Enqueue(sub.ID)
	return s.toSummaryResponse(sub), nil
}

func (s *service) toSummaryResponse(sub *Submission) *SubmissionSummaryResponse {
	return &SubmissionSummaryResponse{
		SubmissionID: sub.ID,
		ProblemID:    sub.ProblemID,
		Status:       sub.Status,
		Language:     sub.Language,
		TestCases: TestCasesSummary{
			Total:  sub.TotalTestCases,
			Passed: sub.PassedTestCases,
			Failed: sub.FailedTestCases,
		},
		ExecutionTimeMs: sub.ExecutionTimeMs,
		MemoryUsedKB:    sub.MemoryUsedKB,
		ErrorMessage:    sub.ErrorMessage,
		SubmittedAt:     sub.SubmittedAt,
		CompletedAt:     sub.CompletedAt,
	}
}
