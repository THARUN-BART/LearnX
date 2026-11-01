package topics

import (
	"context"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
)

type Service interface {
	Get(ctx context.Context, id string) (*model.Topic, error)
	ListSubtopics(ctx context.Context, topicID string) ([]model.Subtopic, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get(ctx context.Context, id string) (*model.Topic, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, model.ErrInvalidInput
	}
	return s.repo.Get(ctx, id)
}

func (s *service) ListSubtopics(ctx context.Context, topicID string) ([]model.Subtopic, error) {
	topicID = strings.TrimSpace(topicID)
	if topicID == "" {
		return nil, model.ErrInvalidInput
	}
	topic, err := s.repo.Get(ctx, topicID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListSubtopics(ctx, topic.ID)
}
