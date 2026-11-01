package progress

import (
	"context"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
)

type Service interface {
	Get(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	List(ctx context.Context, studentID string) ([]model.Progress, error)
	MarkAttempted(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	MarkSolved(ctx context.Context, studentID, problemID string) (*model.Progress, error)
	Update(ctx context.Context, studentID, problemID string, req UpdateRequest) (*model.Progress, error)
}

type service struct {
	repo         Repository
	problemsRepo problems.Repository
}

func NewService(repo Repository, problemsRepo problems.Repository) Service {
	return &service{repo: repo, problemsRepo: problemsRepo}
}

func (s *service) Get(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(problemID) == "" {
		return nil, model.ErrInvalidInput
	}
	if _, err := s.problemsRepo.Get(ctx, problemID); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx, studentID, problemID)
}

func (s *service) List(ctx context.Context, studentID string) ([]model.Progress, error) {
	if strings.TrimSpace(studentID) == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.ListByStudent(ctx, studentID)
}

func (s *service) MarkAttempted(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(problemID) == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.MarkAttempted(ctx, studentID, problemID)
}

func (s *service) MarkSolved(ctx context.Context, studentID, problemID string) (*model.Progress, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(problemID) == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.MarkSolved(ctx, studentID, problemID)
}

func (s *service) Update(ctx context.Context, studentID, problemID string, req UpdateRequest) (*model.Progress, error) {
	if strings.TrimSpace(studentID) == "" || strings.TrimSpace(problemID) == "" {
		return nil, model.ErrInvalidInput
	}
	req.Status = strings.TrimSpace(req.Status)
	if req.Reset {
		req.Status = model.StatusNotStarted
	}
	if !model.ValidStatus(req.Status) || req.Status == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.UpdateStatus(ctx, studentID, problemID, req.Status, req.Reset)
}
