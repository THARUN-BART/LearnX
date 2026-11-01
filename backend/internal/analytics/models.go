package analytics

import "github.com/THARUN-BART/LearnX/backend/internal/dsa/model"

type Overview struct {
	RoadmapID   string  `json:"roadmap_id,omitempty"`
	RoadmapName string  `json:"roadmap_name,omitempty"`
	Total       int     `json:"total_problems"`
	Solved      int     `json:"solved_problems"`
	Attempted   int     `json:"attempted_problems"`
	Pending     int     `json:"pending_problems"`
	Revision    int     `json:"needs_revision"`
	Completion  float64 `json:"completion_percentage"`
}

type Summary struct {
	Overview   Overview                `json:"overview"`
	Topics     []model.TopicProgress   `json:"topics"`
	Difficulty []model.DifficultyStats `json:"difficulty"`
}

type ActivityResponse struct {
	Days     int                   `json:"days"`
	Activity []model.ActivityPoint `json:"activity"`
}
