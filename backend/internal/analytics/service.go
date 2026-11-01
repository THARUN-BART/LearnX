package analytics

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/roadmap"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
)

type Service interface {
	Overview(ctx context.Context, studentID, roadmapID string) (*Overview, error)
	Topics(ctx context.Context, studentID, roadmapID string) ([]model.TopicProgress, error)
	Difficulty(ctx context.Context, studentID, roadmapID string) ([]model.DifficultyStats, error)
	Activity(ctx context.Context, studentID string, days int) (*ActivityResponse, error)
	Summary(ctx context.Context, studentID, roadmapID string) (*Summary, error)
}

type service struct {
	roadmaps roadmap.Repository
	topics   topics.Repository
	problems problems.Repository
	progress progress.Repository
}

func NewService(roadmaps roadmap.Repository, topicsRepo topics.Repository, problemsRepo problems.Repository, progressRepo progress.Repository) Service {
	return &service{
		roadmaps: roadmaps,
		topics:   topicsRepo,
		problems: problemsRepo,
		progress: progressRepo,
	}
}

func (s *service) Overview(ctx context.Context, studentID, roadmapID string) (*Overview, error) {
	summary, err := s.compute(ctx, studentID, roadmapID)
	if err != nil {
		return nil, err
	}
	return &summary.Overview, nil
}

func (s *service) Topics(ctx context.Context, studentID, roadmapID string) ([]model.TopicProgress, error) {
	summary, err := s.compute(ctx, studentID, roadmapID)
	if err != nil {
		return nil, err
	}
	return summary.Topics, nil
}

func (s *service) Difficulty(ctx context.Context, studentID, roadmapID string) ([]model.DifficultyStats, error) {
	summary, err := s.compute(ctx, studentID, roadmapID)
	if err != nil {
		return nil, err
	}
	return summary.Difficulty, nil
}

func (s *service) Summary(ctx context.Context, studentID, roadmapID string) (*Summary, error) {
	return s.compute(ctx, studentID, roadmapID)
}

func (s *service) Activity(ctx context.Context, studentID string, days int) (*ActivityResponse, error) {
	if days <= 0 {
		days = 30
	}
	if days > 90 {
		days = 90
	}
	records, err := s.progress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	start := now.AddDate(0, 0, -days+1)
	bucket := map[string]*model.ActivityPoint{}
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i).Format("2006-01-02")
		bucket[day] = &model.ActivityPoint{Date: day}
	}
	for _, rec := range records {
		if rec.LastAttemptAt != nil {
			day := rec.LastAttemptAt.UTC().Format("2006-01-02")
			if point, ok := bucket[day]; ok {
				point.Attempts += rec.TotalAttempts
			}
		}
		if rec.SolvedAt != nil {
			day := rec.SolvedAt.UTC().Format("2006-01-02")
			if point, ok := bucket[day]; ok {
				point.Solved++
			}
		}
	}
	out := make([]model.ActivityPoint, 0, days)
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i).Format("2006-01-02")
		out = append(out, *bucket[day])
	}
	return &ActivityResponse{Days: days, Activity: out}, nil
}

func (s *service) compute(ctx context.Context, studentID, roadmapID string) (*Summary, error) {
	filter := model.ProblemFilter{Limit: 100000, Page: 1, StudentID: studentID}
	overview := Overview{}
	if strings.TrimSpace(roadmapID) != "" {
		rm, err := s.roadmaps.Get(ctx, roadmapID)
		if err != nil {
			return nil, err
		}
		filter.RoadmapID = rm.ID
		overview.RoadmapID = rm.ID
		overview.RoadmapName = rm.Name
	} else {
		roadmaps, err := s.roadmaps.List(ctx)
		if err != nil {
			return nil, err
		}
		if len(roadmaps) == 1 {
			filter.RoadmapID = roadmaps[0].ID
			overview.RoadmapID = roadmaps[0].ID
			overview.RoadmapName = roadmaps[0].Name
		}
	}

	items, _, err := s.problems.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	progressRows, err := s.progress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}
	progressByProblem := map[string]model.Progress{}
	for _, row := range progressRows {
		progressByProblem[row.ProblemID] = row
	}

	topicAgg := map[string]*model.TopicProgress{}
	diffAgg := map[string]*model.DifficultyStats{
		model.DifficultyEasy:        {Difficulty: model.DifficultyEasy},
		model.DifficultyMedium:      {Difficulty: model.DifficultyMedium},
		model.DifficultyHard:        {Difficulty: model.DifficultyHard},
		model.DifficultyUnspecified: {Difficulty: model.DifficultyUnspecified},
	}

	for _, problem := range items {
		overview.Total++
		status := model.StatusNotStarted
		rec, hasProgress := progressByProblem[problem.ID]
		if hasProgress {
			status = rec.Status
		}
		switch status {
		case model.StatusSolved:
			overview.Solved++
		case model.StatusAttempted:
			overview.Attempted++
		case model.StatusNeedsRevision:
			overview.Revision++
			overview.Attempted++
		default:
			overview.Pending++
		}

		tp, ok := topicAgg[problem.TopicID]
		if !ok {
			tp = &model.TopicProgress{TopicID: problem.TopicID, TopicName: problem.TopicName}
			if topic, err := s.topics.Get(ctx, problem.TopicID); err == nil && topic != nil {
				tp.TopicName = topic.Name
				tp.TopicOrder = topic.TopicOrder
			}
			topicAgg[problem.TopicID] = tp
		}
		tp.Total++
		switch status {
		case model.StatusSolved:
			tp.Solved++
		case model.StatusAttempted:
			tp.Attempted++
		case model.StatusNeedsRevision:
			tp.Revision++
			tp.Attempted++
		default:
			tp.Pending++
		}

		ds := diffAgg[problem.Difficulty]
		if ds == nil {
			ds = &model.DifficultyStats{Difficulty: problem.Difficulty}
			diffAgg[problem.Difficulty] = ds
		}
		ds.Total++
		if hasProgress {
			ds.TotalAttempts += rec.TotalAttempts
			ds.SuccessfulAttempts += rec.SuccessfulAttempts
			if rec.Status == model.StatusSolved {
				ds.Solved++
			} else if rec.Status == model.StatusAttempted || rec.Status == model.StatusNeedsRevision {
				ds.Attempted++
			}
		}
	}

	overview.Completion = model.CompletionPercent(overview.Solved, overview.Total)

	topicsOut := make([]model.TopicProgress, 0, len(topicAgg))
	for _, tp := range topicAgg {
		tp.Completion = model.CompletionPercent(tp.Solved, tp.Total)
		topicsOut = append(topicsOut, *tp)
	}
	sort.Slice(topicsOut, func(i, j int) bool { return topicsOut[i].TopicOrder < topicsOut[j].TopicOrder })

	diffOut := make([]model.DifficultyStats, 0, 4)
	for _, key := range []string{model.DifficultyEasy, model.DifficultyMedium, model.DifficultyHard, model.DifficultyUnspecified} {
		ds := diffAgg[key]
		ds.Accuracy = model.Accuracy(ds.SuccessfulAttempts, ds.TotalAttempts)
		diffOut = append(diffOut, *ds)
	}

	return &Summary{
		Overview:   overview,
		Topics:     topicsOut,
		Difficulty: diffOut,
	}, nil
}
