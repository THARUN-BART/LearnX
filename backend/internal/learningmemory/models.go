package learningmemory

import (
	"encoding/json"
	"time"
)

// LearningState describes the student's lifecycle with a concept.
const (
	StateDiscovered  = "discovered"
	StatePracticing  = "practicing"
	StateReinforcing = "reinforcing"
	StateMastered    = "mastered"
	StateStruggling  = "struggling"
)

// EventType enumerations for append-only audit trail.
const (
	EventProblemStarted      = "problem_started"
	EventSubmissionFailed    = "submission_failed"
	EventSubmissionAccepted  = "submission_accepted"
	EventHintRequested       = "hint_requested"
	EventMistakeDetected     = "mistake_detected"
	EventMistakeCorrected    = "mistake_corrected"
	EventWalkthroughViewed   = "walkthrough_requested"
	EventRevisionCompleted   = "revision_completed"
	EventConceptPractised    = "concept_practised"
	EventAssessmentCompleted = "assessment_completed"
)

// ConceptMemory models the student's persistent learning state for one DSA concept.
type ConceptMemory struct {
	ID                 string     `json:"id"`
	StudentID          string     `json:"student_id"`
	TopicID            *string    `json:"topic_id,omitempty"`
	SubtopicID         *string    `json:"subtopic_id,omitempty"`
	ConceptKey         string     `json:"concept_key"`
	ConceptName        string     `json:"concept_name"`
	TopicName          string     `json:"topic_name,omitempty"`
	LearningState      string     `json:"learning_state"`
	EstimatedMastery   float64    `json:"estimated_mastery"` // Explicitly labeled provisional estimate
	ConfidenceLevel    string     `json:"confidence_level"`   // low, moderate, high
	TotalAttempts      int        `json:"total_attempts"`
	IndependentSolves  int        `json:"independent_solves"`
	AssistedSolves     int        `json:"assisted_solves"`
	FailedAttempts     int        `json:"failed_attempts"`
	HintRequests       int        `json:"hint_requests"`
	MistakeCount       int        `json:"mistake_count"`
	CorrectedMistakes  int        `json:"corrected_mistakes"`
	FirstPracticedAt   time.Time  `json:"first_practiced_at"`
	LastPracticedAt    time.Time  `json:"last_practiced_at"`
	LastSuccessfulAt   *time.Time `json:"last_successful_at,omitempty"`
	NextReviewAt       *time.Time `json:"next_review_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// LearningEvent models an immutable learning signal.
type LearningEvent struct {
	ID           string                 `json:"id"`
	StudentID    string                 `json:"student_id"`
	EventType    string                 `json:"event_type"`
	ProblemID    *string                `json:"problem_id,omitempty"`
	SubmissionID *string                `json:"submission_id,omitempty"`
	DiagnosisID  *string                `json:"diagnosis_id,omitempty"`
	TopicID      *string                `json:"topic_id,omitempty"`
	SubtopicID   *string                `json:"subtopic_id,omitempty"`
	ConceptKey   string                 `json:"concept_key,omitempty"`
	EventData    map[string]interface{} `json:"event_data"`
	OccurredAt   time.Time              `json:"occurred_at"`
	CreatedAt    time.Time              `json:"created_at"`
}

// MistakeMemory tracks recurring patterns in the student's cognitive profile.
type MistakeMemory struct {
	ID                    string    `json:"id"`
	StudentID             string    `json:"student_id"`
	Category              string    `json:"category"`
	Subtype               string    `json:"subtype"`
	ConceptKey            string    `json:"concept_key"`
	OccurrenceCount       int       `json:"occurrence_count"`
	CorrectedCount        int       `json:"corrected_count"`
	FirstOccurrenceAt     time.Time `json:"first_occurrence_at"`
	LastOccurrenceAt      time.Time `json:"last_occurrence_at"`
	LastCorrectedAt       *time.Time `json:"last_corrected_at,omitempty"`
	Status                string    `json:"status"` // active, improving, resolved
	SupportingDiagnosisIDs []string  `json:"supporting_diagnosis_ids"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// LearningSnapshot stores aggregated daily milestones.
type LearningSnapshot struct {
	ID                    string                 `json:"id"`
	StudentID             string                 `json:"student_id"`
	SnapshotDate          string                 `json:"snapshot_date"`
	TotalConceptsPractised int                   `json:"total_concepts_practised"`
	IndependentSolves     int                    `json:"independent_solves"`
	AssistedSolves        int                    `json:"assisted_solves"`
	MistakeCount          int                    `json:"mistake_count"`
	CorrectedMistakeCount int                    `json:"corrected_mistake_count"`
	ActiveLearningDays    int                    `json:"active_learning_days"`
	Summary               map[string]interface{} `json:"summary"`
	CreatedAt             time.Time              `json:"created_at"`
}

// MemorySummaryResponse contains the high-level dashboard telemetry.
type MemorySummaryResponse struct {
	TotalConceptsPractised int               `json:"total_concepts_practised"`
	IndependentSolves     int               `json:"independent_solves"`
	AssistedSolves        int               `json:"assisted_solves"`
	TotalMistakes         int               `json:"total_mistakes"`
	CorrectedMistakes     int               `json:"corrected_mistakes"`
	ActiveLearningDays    int               `json:"active_learning_days"`
	MasteredConcepts      int               `json:"mastered_concepts"`
	StrugglingConcepts    int               `json:"struggling_concepts"`
	RecentConcepts        []ConceptMemory   `json:"recent_concepts"`
	ActiveWeaknesses      []MistakeMemory   `json:"active_weaknesses"`
	KeyStrengths          []ConceptStrength `json:"key_strengths"`
}

// ConceptStrength models a concept backed by independent solve evidence.
type ConceptStrength struct {
	ConceptKey        string  `json:"concept_key"`
	ConceptName       string  `json:"concept_name"`
	TopicName         string  `json:"topic_name"`
	IndependentSolves int     `json:"independent_solves"`
	EstimatedMastery  float64 `json:"estimated_mastery"`
	ConfidenceLevel   string  `json:"confidence_level"`
	Reason            string  `json:"reason"`
}

// ConceptWeakness models a concept with recurring mistake evidence.
type ConceptWeakness struct {
	ConceptKey      string  `json:"concept_key"`
	ConceptName     string  `json:"concept_name"`
	TopicName       string  `json:"topic_name"`
	MistakeCategory string  `json:"mistake_category"`
	MistakeSubtype  string  `json:"mistake_subtype"`
	OccurrenceCount int     `json:"occurrence_count"`
	CorrectedCount  int     `json:"corrected_count"`
	Status          string  `json:"status"`
	Reason          string  `json:"reason"`
}

// TimelinePoint for charts.
type TimelinePoint struct {
	Date              string `json:"date"`
	IndependentSolves int    `json:"independent_solves"`
	AssistedSolves    int    `json:"assisted_solves"`
	MistakesDetected  int    `json:"mistakes_detected"`
	MistakesCorrected int    `json:"mistakes_corrected"`
	ConceptsPractised int    `json:"concepts_practised"`
}

// StudentFacingSummary represents human-readable cognitive progress.
type StudentFacingSummary struct {
	Headline         string   `json:"headline"`
	Narrative        string   `json:"narrative"`
	TopStrengths     []string `json:"top_strengths"`
	AreasForFocus    []string `json:"areas_for_focus"`
	NextActionPrompt string   `json:"next_action_prompt"`
}

// ConceptDetailResponse includes concept telemetry, associated mistakes, and event audit.
type ConceptDetailResponse struct {
	Concept   ConceptMemory   `json:"concept"`
	Mistakes  []MistakeMemory `json:"mistakes"`
	Timeline  []LearningEvent `json:"timeline"`
	NextSteps []string        `json:"next_steps"`
}

// JSONBytes helper for marshaling map.
func ToJSONMap(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}
