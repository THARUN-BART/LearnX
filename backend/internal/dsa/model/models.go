package model

import "time"

const (
	DifficultyEasy        = "Easy"
	DifficultyMedium      = "Medium"
	DifficultyHard        = "Hard"
	DifficultyUnspecified = "Unspecified"

	StatusNotStarted    = "Not Started"
	StatusAttempted     = "Attempted"
	StatusSolved        = "Solved"
	StatusNeedsRevision = "Needs Revision"

	SourceA2ZSlug = "strivers-a2z-dsa-sheet"
)

type Roadmap struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	SourceName   string    `json:"source_name"`
	SourceURL    string    `json:"source_url"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TopicCount   int       `json:"topic_count,omitempty"`
	ProblemCount int       `json:"problem_count,omitempty"`
}

type Topic struct {
	ID            string    `json:"id"`
	RoadmapID     string    `json:"roadmap_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	TopicOrder    int       `json:"topic_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	SubtopicCount int       `json:"subtopic_count,omitempty"`
	ProblemCount  int       `json:"problem_count,omitempty"`
}

type Subtopic struct {
	ID            string    `json:"id"`
	TopicID       string    `json:"topic_id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	SubtopicOrder int       `json:"subtopic_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ProblemCount  int       `json:"problem_count,omitempty"`
}

type Problem struct {
	ID                 string            `json:"id"`
	TopicID            string            `json:"topic_id"`
	SubtopicID         string            `json:"subtopic_id"`
	Title              string            `json:"title"`
	Slug               string            `json:"slug"`
	Summary            string            `json:"summary"`
	Description        string            `json:"description,omitempty"`
	Difficulty         string            `json:"difficulty"`
	ExternalPlatform   string            `json:"external_platform"`
	ExternalProblemURL string            `json:"external_problem_url"`
	SourceSheet        string            `json:"source_sheet"`
	SourceProblemID    string            `json:"source_problem_id"`
	ProblemOrder       int               `json:"problem_order"`
	Tags               []string          `json:"tags"`
	StarterCode        map[string]string `json:"starter_code,omitempty"`
	IsActive           bool              `json:"is_active"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	TopicName          string            `json:"topic_name,omitempty"`
	SubtopicName       string            `json:"subtopic_name,omitempty"`
	ProgressStatus     string            `json:"progress_status,omitempty"`
}

type Progress struct {
	ID                 string     `json:"id"`
	StudentID          string     `json:"student_id"`
	ProblemID          string     `json:"problem_id"`
	Status             string     `json:"status"`
	TotalAttempts      int        `json:"total_attempts"`
	SuccessfulAttempts int        `json:"successful_attempts"`
	FirstAttemptAt     *time.Time `json:"first_attempt_at,omitempty"`
	LastAttemptAt      *time.Time `json:"last_attempt_at,omitempty"`
	SolvedAt           *time.Time `json:"solved_at,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

type ProblemFilter struct {
	Query      string
	Difficulty string
	TopicID    string
	SubtopicID string
	RoadmapID  string
	Status     string
	StudentID  string
	Sort       string
	Page       int
	Limit      int
}

type ProgressCounts struct {
	Total      int
	Solved     int
	Attempted  int
	Pending    int
	Revision   int
	Completion float64
}

type TopicProgress struct {
	TopicID    string  `json:"topic_id"`
	TopicName  string  `json:"topic_name"`
	TopicOrder int     `json:"topic_order"`
	Total      int     `json:"total_problems"`
	Solved     int     `json:"solved_problems"`
	Attempted  int     `json:"attempted_problems"`
	Pending    int     `json:"pending_problems"`
	Revision   int     `json:"needs_revision"`
	Completion float64 `json:"completion_percentage"`
}

type DifficultyStats struct {
	Difficulty         string  `json:"difficulty"`
	Total              int     `json:"total"`
	Solved             int     `json:"solved"`
	Attempted          int     `json:"attempted"`
	SuccessfulAttempts int     `json:"successful_attempts"`
	TotalAttempts      int     `json:"total_attempts"`
	Accuracy           float64 `json:"accuracy"`
}

type ActivityPoint struct {
	Date     string `json:"date"`
	Attempts int    `json:"attempts"`
	Solved   int    `json:"solved"`
}

func CompletionPercent(solved, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(solved) / float64(total) * 100
}

func Accuracy(successful, total int) float64 {
	if total <= 0 {
		return 0
	}
	return float64(successful) / float64(total) * 100
}

func ValidDifficulty(value string) bool {
	switch value {
	case "", DifficultyEasy, DifficultyMedium, DifficultyHard, DifficultyUnspecified:
		return true
	default:
		return false
	}
}

func ValidStatus(value string) bool {
	switch value {
	case "", StatusNotStarted, StatusAttempted, StatusSolved, StatusNeedsRevision:
		return true
	default:
		return false
	}
}
