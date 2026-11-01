package roadmap

import (
	"context"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
)

type Service interface {
	List(ctx context.Context) ([]model.Roadmap, error)
	Get(ctx context.Context, idOrSlug string) (*model.Roadmap, error)
	ListTopics(ctx context.Context, roadmapID string) ([]model.Topic, error)
}

type service struct {
	repo       Repository
	topicsRepo topics.Repository
}

func NewService(repo Repository, topicsRepo topics.Repository) Service {
	return &service{repo: repo, topicsRepo: topicsRepo}
}

func (s *service) List(ctx context.Context) ([]model.Roadmap, error) {
	return s.repo.List(ctx)
}

func (s *service) Get(ctx context.Context, idOrSlug string) (*model.Roadmap, error) {
	idOrSlug = strings.TrimSpace(idOrSlug)
	if idOrSlug == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.Get(ctx, idOrSlug)
}

func (s *service) ListTopics(ctx context.Context, roadmapID string) ([]model.Topic, error) {
	roadmapID = strings.TrimSpace(roadmapID)
	if roadmapID == "" {
		return nil, model.ErrInvalidInput
	}
	rm, err := s.repo.Get(ctx, roadmapID)
	if err != nil {
		return nil, err
	}
	return s.topicsRepo.ListByRoadmap(ctx, rm.ID)
}
