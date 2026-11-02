package learningmemory

import (
	"context"
	"time"
)

// ConceptFilter specifies search options for concepts.
type ConceptFilter struct {
	TopicID       string
	LearningState string
	Search        string
	Limit         int
	Offset        int
}

// EventFilter specifies search options for learning events.
type EventFilter struct {
	EventType  string
	ConceptKey string
	ProblemID  string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

// Repository defines data access operations for student learning memory.
type Repository interface {
	// Concept Memory
	GetConceptMemory(ctx context.Context, studentID, conceptKey string) (*ConceptMemory, error)
	UpsertConceptMemory(ctx context.Context, cm *ConceptMemory) error
	ListConceptMemories(ctx context.Context, studentID string, filter ConceptFilter) ([]ConceptMemory, int, error)
	GetStrengths(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error)
	GetWeaknesses(ctx context.Context, studentID string, limit int) ([]ConceptMemory, error)

	// Learning Events
	RecordEvent(ctx context.Context, event *LearningEvent) error
	ListEvents(ctx context.Context, studentID string, filter EventFilter) ([]LearningEvent, int, error)
	HasRecentHint(ctx context.Context, studentID, problemID string, before time.Time) (bool, error)

	// Mistake Memory
	GetMistakeMemory(ctx context.Context, studentID, category, subtype, conceptKey string) (*MistakeMemory, error)
	UpsertMistakeMemory(ctx context.Context, mm *MistakeMemory) error
	ListMistakeMemories(ctx context.Context, studentID string, status string) ([]MistakeMemory, error)

	// Snapshots
	UpsertSnapshot(ctx context.Context, snap *LearningSnapshot) error
	ListSnapshots(ctx context.Context, studentID string, days int) ([]LearningSnapshot, error)

	// Aggregates
	GetSummary(ctx context.Context, studentID string) (*MemorySummaryResponse, error)
	GetTimeline(ctx context.Context, studentID string, days int) ([]TimelinePoint, error)
}
