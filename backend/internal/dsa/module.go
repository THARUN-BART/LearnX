package dsa

import (
	"github.com/THARUN-BART/LearnX/backend/internal/analytics"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/roadmap"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Roadmap  roadmap.Repository
	Topics   topics.Repository
	Problems problems.Repository
	Progress progress.Repository
}

func NewMemoryRepositories(mem *model.Memory) Repositories {
	return Repositories{
		Roadmap:  roadmap.NewMemoryRepository(mem),
		Topics:   topics.NewMemoryRepository(mem),
		Problems: problems.NewMemoryRepository(mem),
		Progress: progress.NewMemoryRepository(mem),
	}
}

func NewPostgresRepositories(pool *pgxpool.Pool) Repositories {
	return Repositories{
		Roadmap:  roadmap.NewPostgresRepository(pool),
		Topics:   topics.NewPostgresRepository(pool),
		Problems: problems.NewPostgresRepository(pool),
		Progress: progress.NewPostgresRepository(pool),
	}
}

func NewHandlers(repos Repositories) *server.DSAHandlers {
	roadmapSvc := roadmap.NewService(repos.Roadmap, repos.Topics)
	topicsSvc := topics.NewService(repos.Topics)
	problemsSvc := problems.NewService(repos.Problems, repos.Topics)
	progressSvc := progress.NewService(repos.Progress, repos.Problems)
	analyticsSvc := analytics.NewService(repos.Roadmap, repos.Topics, repos.Problems, repos.Progress)

	return &server.DSAHandlers{
		Roadmap:   roadmap.NewHandler(roadmapSvc),
		Topics:    topics.NewHandler(topicsSvc),
		Problems:  problems.NewHandler(problemsSvc),
		Progress:  progress.NewHandler(progressSvc),
		Analytics: analytics.NewHandler(analyticsSvc),
	}
}
