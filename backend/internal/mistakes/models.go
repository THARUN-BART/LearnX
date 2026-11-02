package mistakes

import "time"

// Mistake Categories
const (
	CategoryCompilation   = "compilation"
	CategoryRuntime       = "runtime"
	CategoryLogical       = "logical"
	CategoryEdgeCase      = "edge_case"
	CategoryDataStructure = "data_structure"
	CategoryAlgorithmic   = "algorithmic"
	CategoryComplexity    = "complexity"
	CategoryUnknown       = "unknown"
)

// Confidence Levels
const (
	ConfidenceHigh                 = "High"
	ConfidenceMedium               = "Medium"
	ConfidenceLow                  = "Low"
	ConfidenceInsufficientEvidence = "Insufficient Evidence"
)

// Detection Methods
const (
	MethodDeterministic = "deterministic_execution"
	MethodStatic        = "static_analysis"
	MethodHistorical    = "historical_pattern"
	MethodRootCause     = "root_cause_inference"
)

// Diagnosis represents a detailed evaluation of why a submission failed.
type Diagnosis struct {
	ID                  string                 `json:"id"`
	StudentID           string                 `json:"student_id"`
	SubmissionID        string                 `json:"submission_id"`
	ProblemID           string                 `json:"problem_id"`
	Category            string                 `json:"category"`
	Subtype             string                 `json:"subtype"`
	Description         string                 `json:"description"`
	RelatedConcept      string                 `json:"related_concept"`
	DetectionMethod     string                 `json:"detection_method"`
	ConfidenceLevel     string                 `json:"confidence_level"`
	ConfidenceScore     float64                `json:"confidence_score"`
	Evidence            map[string]interface{} `json:"evidence"`
	SuggestedCorrection string                 `json:"suggested_correction"`
	CreatedAt           time.Time              `json:"created_at"`
}

// StudentMistakePattern aggregates recurring mistakes for a student.
type StudentMistakePattern struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	Category          string    `json:"category"`
	Subtype           string    `json:"subtype"`
	RelatedConcept    string    `json:"related_concept"`
	OccurrenceCount   int       `json:"occurrence_count"`
	FirstOccurrenceAt time.Time `json:"first_occurrence_at"`
	LastOccurrenceAt  time.Time `json:"last_occurrence_at"`
	CorrectionStatus  string    `json:"correction_status"` // "active", "corrected"
	UpdatedAt         time.Time `json:"updated_at"`
}

// MistakeSummary gives an overview of a student's mistake trends.
type MistakeSummary struct {
	StudentID              string                  `json:"student_id"`
	TotalMistakes          int                     `json:"total_mistakes"`
	CategoryCounts         map[string]int          `json:"category_counts"`
	RecurringPatterns      []StudentMistakePattern `json:"recurring_patterns"`
	RecentDiagnoses        []Diagnosis             `json:"recent_diagnoses"`
	NeedsAttentionConcepts []string                `json:"needs_attention_concepts"`
}

// MistakeFeedback captures student feedback on diagnosis accuracy.
type MistakeFeedback struct {
	ID              string    `json:"id"`
	DiagnosisID     string    `json:"diagnosis_id"`
	StudentID       string    `json:"student_id"`
	FeedbackType    string    `json:"feedback_type"` // helpful, unhelpful, inaccurate, too_vague
	StudentFeedback string    `json:"student_feedback,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}
