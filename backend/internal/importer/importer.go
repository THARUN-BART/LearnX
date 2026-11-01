package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedFile struct {
	Attribution map[string]any `json:"attribution"`
	Roadmap     SeedRoadmap    `json:"roadmap"`
	Topics      []SeedTopic    `json:"topics"`
}

type SeedRoadmap struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SourceName  string `json:"source_name"`
	SourceURL   string `json:"source_url"`
}

type SeedTopic struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	TopicOrder  int            `json:"topic_order"`
	Subtopics   []SeedSubtopic `json:"subtopics"`
}

type SeedSubtopic struct {
	Name          string        `json:"name"`
	Slug          string        `json:"slug"`
	Description   string        `json:"description"`
	SubtopicOrder int           `json:"subtopic_order"`
	Problems      []SeedProblem `json:"problems"`
}

type SeedProblem struct {
	Title              string   `json:"title"`
	Slug               string   `json:"slug"`
	Summary            string   `json:"summary"`
	Difficulty         string   `json:"difficulty"`
	ExternalPlatform   string   `json:"external_platform"`
	ExternalProblemURL string   `json:"external_problem_url"`
	SourceSheet        string   `json:"source_sheet"`
	SourceProblemID    string   `json:"source_problem_id"`
	ProblemOrder       int      `json:"problem_order"`
	Tags               []string `json:"tags"`
}

type Result struct {
	RoadmapID         string `json:"roadmap_id"`
	TopicsUpserted    int    `json:"topics_upserted"`
	SubtopicsUpserted int    `json:"subtopics_upserted"`
	ProblemsUpserted  int    `json:"problems_upserted"`
	ProblemsSkipped   int    `json:"problems_skipped"`
	DuplicatesFound   int    `json:"duplicates_found"`
}

func LoadSeedFile(path string) (*SeedFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read seed file: %w", err)
	}
	var seed SeedFile
	if err := json.Unmarshal(raw, &seed); err != nil {
		return nil, fmt.Errorf("parse seed file: %w", err)
	}
	if err := validateSeed(&seed); err != nil {
		return nil, err
	}
	return &seed, nil
}

func validateSeed(seed *SeedFile) error {
	if strings.TrimSpace(seed.Roadmap.Name) == "" || strings.TrimSpace(seed.Roadmap.Slug) == "" {
		return fmt.Errorf("%w: roadmap name and slug are required", model.ErrInvalidInput)
	}
	seenSource := map[string]struct{}{}
	seenSlug := map[string]struct{}{}
	for _, topic := range seed.Topics {
		if topic.Name == "" || topic.TopicOrder <= 0 {
			return fmt.Errorf("%w: topic name and order are required", model.ErrInvalidInput)
		}
		for _, sub := range topic.Subtopics {
			if sub.Name == "" || sub.SubtopicOrder <= 0 {
				return fmt.Errorf("%w: subtopic name and order are required", model.ErrInvalidInput)
			}
			for _, p := range sub.Problems {
				if strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.ExternalProblemURL) == "" {
					return fmt.Errorf("%w: problem title and canonical URL are required", model.ErrInvalidInput)
				}
				if p.SourceProblemID == "" {
					return fmt.Errorf("%w: source_problem_id is required to keep imports repeatable", model.ErrInvalidInput)
				}
				if _, ok := seenSource[p.SourceSheet+"|"+p.SourceProblemID]; ok {
					return fmt.Errorf("%w: duplicate source_problem_id %s", model.ErrConflict, p.SourceProblemID)
				}
				seenSource[p.SourceSheet+"|"+p.SourceProblemID] = struct{}{}
				if p.Slug != "" {
					if _, ok := seenSlug[p.Slug]; ok {
						return fmt.Errorf("%w: duplicate problem slug %s", model.ErrConflict, p.Slug)
					}
					seenSlug[p.Slug] = struct{}{}
				}
				if p.Difficulty == "" {
					p.Difficulty = model.DifficultyUnspecified
				}
				if !model.ValidDifficulty(p.Difficulty) || p.Difficulty == "" {
					return fmt.Errorf("%w: invalid difficulty %q", model.ErrInvalidInput, p.Difficulty)
				}
			}
		}
	}
	return nil
}

func ImportMemory(mem *model.Memory, seed *SeedFile) (*Result, error) {
	if mem == nil || seed == nil {
		return nil, model.ErrInvalidInput
	}
	now := time.Now().UTC()
	res := &Result{}

	mem.Lock()
	defer mem.Unlock()

	var roadmap *model.Roadmap
	for _, existing := range mem.Roadmaps {
		if existing.Slug == seed.Roadmap.Slug {
			roadmap = existing
			break
		}
	}
	if roadmap == nil {
		roadmap = &model.Roadmap{ID: model.NewID(), CreatedAt: now}
		mem.Roadmaps[roadmap.ID] = roadmap
	}
	roadmap.Name = seed.Roadmap.Name
	roadmap.Slug = seed.Roadmap.Slug
	roadmap.Description = seed.Roadmap.Description
	roadmap.SourceName = seed.Roadmap.SourceName
	roadmap.SourceURL = seed.Roadmap.SourceURL
	roadmap.IsActive = true
	roadmap.UpdatedAt = now
	res.RoadmapID = roadmap.ID

	topicByKey := map[string]*model.Topic{}
	for _, t := range mem.Topics {
		if t.RoadmapID == roadmap.ID {
			topicByKey[t.Slug] = t
		}
	}
	subByKey := map[string]*model.Subtopic{}
	for _, s := range mem.Subtopics {
		subByKey[s.TopicID+"|"+s.Slug] = s
	}
	problemBySource := map[string]*model.Problem{}
	for _, p := range mem.Problems {
		problemBySource[p.SourceSheet+"|"+p.SourceProblemID] = p
	}

	for _, tseed := range seed.Topics {
		topic := topicByKey[tseed.Slug]
		if topic == nil {
			topic = &model.Topic{ID: model.NewID(), RoadmapID: roadmap.ID, CreatedAt: now}
			mem.Topics[topic.ID] = topic
			topicByKey[tseed.Slug] = topic
		}
		topic.Name = tseed.Name
		topic.Slug = tseed.Slug
		topic.Description = tseed.Description
		topic.TopicOrder = tseed.TopicOrder
		topic.UpdatedAt = now
		res.TopicsUpserted++

		for _, sseed := range tseed.Subtopics {
			key := topic.ID + "|" + sseed.Slug
			sub := subByKey[key]
			if sub == nil {
				sub = &model.Subtopic{ID: model.NewID(), TopicID: topic.ID, CreatedAt: now}
				mem.Subtopics[sub.ID] = sub
				subByKey[key] = sub
			}
			sub.Name = sseed.Name
			sub.Slug = sseed.Slug
			sub.Description = sseed.Description
			sub.SubtopicOrder = sseed.SubtopicOrder
			sub.UpdatedAt = now
			res.SubtopicsUpserted++

			for _, pseed := range sseed.Problems {
				srcKey := pseed.SourceSheet + "|" + pseed.SourceProblemID
				prob := problemBySource[srcKey]
				if prob == nil {
					prob = &model.Problem{ID: model.NewID(), CreatedAt: now}
					mem.Problems[prob.ID] = prob
					problemBySource[srcKey] = prob
				} else {
					res.DuplicatesFound++
				}
				prob.TopicID = topic.ID
				prob.SubtopicID = sub.ID
				prob.Title = pseed.Title
				prob.Slug = pseed.Slug
				prob.Summary = pseed.Summary
				prob.Difficulty = normalizeDifficulty(pseed.Difficulty)
				prob.ExternalPlatform = pseed.ExternalPlatform
				prob.ExternalProblemURL = pseed.ExternalProblemURL
				prob.SourceSheet = pseed.SourceSheet
				if prob.SourceSheet == "" {
					prob.SourceSheet = seed.Roadmap.Slug
				}
				prob.SourceProblemID = pseed.SourceProblemID
				prob.ProblemOrder = pseed.ProblemOrder
				prob.Tags = append([]string{}, pseed.Tags...)
				prob.IsActive = true
				prob.UpdatedAt = now
				res.ProblemsUpserted++
			}
		}
	}

	return res, nil
}

func ImportPostgres(ctx context.Context, pool *pgxpool.Pool, seed *SeedFile) (*Result, error) {
	if pool == nil || seed == nil {
		return nil, model.ErrInvalidInput
	}

	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin import transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	res := &Result{}
	var roadmapID string
	err = tx.QueryRow(ctx, `
		INSERT INTO roadmaps (name, slug, description, source_name, source_url, is_active, updated_at)
		VALUES ($1, $2, $3, $4, $5, TRUE, NOW())
		ON CONFLICT (slug) DO UPDATE SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			source_name = EXCLUDED.source_name,
			source_url = EXCLUDED.source_url,
			is_active = TRUE,
			updated_at = NOW()
		RETURNING id;
	`, seed.Roadmap.Name, seed.Roadmap.Slug, seed.Roadmap.Description, seed.Roadmap.SourceName, seed.Roadmap.SourceURL).Scan(&roadmapID)
	if err != nil {
		return nil, fmt.Errorf("upsert roadmap: %w", err)
	}
	res.RoadmapID = roadmapID

	existingMap := make(map[string]struct{})
	if rows, err := tx.Query(ctx, `SELECT source_sheet || '|' || source_problem_id FROM problems`); err == nil {
		for rows.Next() {
			var k string
			if err := rows.Scan(&k); err == nil {
				existingMap[k] = struct{}{}
			}
		}
		rows.Close()
	}

	batch := &pgx.Batch{}

	for _, tseed := range seed.Topics {
		var topicID string
		err = tx.QueryRow(ctx, `
			INSERT INTO topics (roadmap_id, name, slug, description, topic_order, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (roadmap_id, slug) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				topic_order = EXCLUDED.topic_order,
				updated_at = NOW()
			RETURNING id;
		`, roadmapID, tseed.Name, tseed.Slug, tseed.Description, tseed.TopicOrder).Scan(&topicID)
		if err != nil {
			return nil, fmt.Errorf("upsert topic %s: %w", tseed.Slug, err)
		}
		res.TopicsUpserted++

		for _, sseed := range tseed.Subtopics {
			var subID string
			err = tx.QueryRow(ctx, `
				INSERT INTO subtopics (topic_id, name, slug, description, subtopic_order, updated_at)
				VALUES ($1, $2, $3, $4, $5, NOW())
				ON CONFLICT (topic_id, slug) DO UPDATE SET
					name = EXCLUDED.name,
					description = EXCLUDED.description,
					subtopic_order = EXCLUDED.subtopic_order,
					updated_at = NOW()
				RETURNING id;
			`, topicID, sseed.Name, sseed.Slug, sseed.Description, sseed.SubtopicOrder).Scan(&subID)
			if err != nil {
				return nil, fmt.Errorf("upsert subtopic %s: %w", sseed.Slug, err)
			}
			res.SubtopicsUpserted++

			for _, pseed := range sseed.Problems {
				sourceSheet := pseed.SourceSheet
				if sourceSheet == "" {
					sourceSheet = seed.Roadmap.Slug
				}
				tag := pseed.Tags
				if tag == nil {
					tag = []string{}
				}

				srcKey := sourceSheet + "|" + pseed.SourceProblemID
				if _, existed := existingMap[srcKey]; existed {
					res.DuplicatesFound++
				} else {
					existingMap[srcKey] = struct{}{}
				}

				batch.Queue(`
					INSERT INTO problems (
						topic_id, subtopic_id, title, slug, summary, difficulty,
						external_platform, external_problem_url, source_sheet, source_problem_id,
						problem_order, tags, is_active, updated_at
					) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,TRUE,NOW())
					ON CONFLICT (source_sheet, source_problem_id) DO UPDATE SET
						topic_id = EXCLUDED.topic_id,
						subtopic_id = EXCLUDED.subtopic_id,
						title = EXCLUDED.title,
						slug = EXCLUDED.slug,
						summary = EXCLUDED.summary,
						difficulty = EXCLUDED.difficulty,
						external_platform = EXCLUDED.external_platform,
						external_problem_url = EXCLUDED.external_problem_url,
						problem_order = EXCLUDED.problem_order,
						tags = EXCLUDED.tags,
						is_active = TRUE,
						updated_at = NOW();
				`, topicID, subID, pseed.Title, pseed.Slug, pseed.Summary, normalizeDifficulty(pseed.Difficulty),
					pseed.ExternalPlatform, pseed.ExternalProblemURL, sourceSheet, pseed.SourceProblemID,
					pseed.ProblemOrder, tag)

				res.ProblemsUpserted++
			}
		}
	}

	br := tx.SendBatch(ctx, batch)
	for i := 0; i < res.ProblemsUpserted; i++ {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return nil, fmt.Errorf("batch execute problem: %w", err)
		}
	}
	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit import transaction: %w", err)
	}
	return res, nil
}

func normalizeDifficulty(value string) string {
	switch strings.TrimSpace(value) {
	case model.DifficultyEasy, model.DifficultyMedium, model.DifficultyHard:
		return strings.TrimSpace(value)
	default:
		return model.DifficultyUnspecified
	}
}
