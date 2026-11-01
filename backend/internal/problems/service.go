package problems

import (
	"context"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
)

type Service interface {
	Get(ctx context.Context, idOrSlug string) (*model.Problem, error)
	List(ctx context.Context, filter model.ProblemFilter) ([]model.Problem, int, error)
	ListBySubtopic(ctx context.Context, subtopicID string, filter model.ProblemFilter) ([]model.Problem, int, error)
	Related(ctx context.Context, problemID string) ([]model.Problem, error)
}

type service struct {
	repo       Repository
	topicsRepo topics.Repository
}

func NewService(repo Repository, topicsRepo topics.Repository) Service {
	return &service{repo: repo, topicsRepo: topicsRepo}
}

func (s *service) Get(ctx context.Context, idOrSlug string) (*model.Problem, error) {
	idOrSlug = strings.TrimSpace(idOrSlug)
	if idOrSlug == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.Get(ctx, idOrSlug)
}

func (s *service) List(ctx context.Context, filter model.ProblemFilter) ([]model.Problem, int, error) {
	if err := normalizeFilter(&filter); err != nil {
		return nil, 0, err
	}
	return s.repo.List(ctx, filter)
}

func (s *service) ListBySubtopic(ctx context.Context, subtopicID string, filter model.ProblemFilter) ([]model.Problem, int, error) {
	subtopicID = strings.TrimSpace(subtopicID)
	if subtopicID == "" {
		return nil, 0, model.ErrInvalidInput
	}
	if _, err := s.topicsRepo.GetSubtopic(ctx, subtopicID); err != nil {
		return nil, 0, err
	}
	filter.SubtopicID = subtopicID
	if err := normalizeFilter(&filter); err != nil {
		return nil, 0, err
	}
	return s.repo.List(ctx, filter)
}

func (s *service) Related(ctx context.Context, problemID string) ([]model.Problem, error) {
	if _, err := s.repo.Get(ctx, problemID); err != nil {
		return nil, err
	}
	return s.repo.Related(ctx, problemID, 10)
}

func normalizeFilter(filter *model.ProblemFilter) error {
	filter.Query = strings.TrimSpace(filter.Query)
	filter.Difficulty = strings.TrimSpace(filter.Difficulty)
	filter.Status = strings.TrimSpace(filter.Status)
	filter.Sort = strings.TrimSpace(filter.Sort)
	if filter.Sort == "" {
		filter.Sort = "problem_order"
	}
	if filter.Sort != "problem_order" && filter.Sort != "title" {
		return model.ErrInvalidInput
	}
	if !model.ValidDifficulty(filter.Difficulty) {
		return model.ErrInvalidInput
	}
	if !model.ValidStatus(filter.Status) {
		return model.ErrInvalidInput
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return nil
}
