package mistakes

import (
	"context"
	"errors"
	"log/slog"

	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
)

var (
	ErrUnauthorizedAccess = errors.New("unauthorized access to diagnosis")
)

type Service interface {
	DiagnoseSubmission(ctx context.Context, sub *submissions.Submission, testResults []submissions.SubmissionTestResult, problemTitle, problemTopic string) (*Diagnosis, error)
	GetDiagnosisBySubmission(ctx context.Context, studentID, submissionID string) (*Diagnosis, error)
	GetDiagnosisByID(ctx context.Context, studentID, diagnosisID string) (*Diagnosis, error)
	GetStudentMistakes(ctx context.Context, studentID string, limit, offset int) ([]Diagnosis, error)
	GetStudentPatterns(ctx context.Context, studentID string) ([]StudentMistakePattern, error)
	GetStudentMistakeSummary(ctx context.Context, studentID string) (*MistakeSummary, error)
	RecordFeedback(ctx context.Context, studentID, diagnosisID, feedbackType, comment string) error
	MarkPatternCorrected(ctx context.Context, studentID, category, subtype string) error
}

type serviceImpl struct {
	repo     Repository
	detector *Detector
	logger   *slog.Logger
}

func NewService(repo Repository, detector *Detector, logger *slog.Logger) Service {
	if detector == nil {
		detector = NewDetector()
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &serviceImpl{
		repo:     repo,
		detector: detector,
		logger:   logger,
	}
}

func (s *serviceImpl) DiagnoseSubmission(
	ctx context.Context,
	sub *submissions.Submission,
	testResults []submissions.SubmissionTestResult,
	problemTitle, problemTopic string,
) (*Diagnosis, error) {
	if sub == nil || sub.Status == execution.StatusAccepted {
		return nil, nil
	}

	// 1. Idempotency check: see if a diagnosis was already saved for this submission
	existing, err := s.repo.GetDiagnosisBySubmissionID(ctx, sub.ID)
	if err == nil && existing != nil {
		return existing, nil
	}

	// 2. Load active patterns for historical context
	activePatterns, err := s.repo.GetActivePatterns(ctx, sub.StudentID)
	if err != nil {
		s.logger.Warn("could not load student active patterns", slog.String("student_id", sub.StudentID), slog.String("error", err.Error()))
	}

	// 3. Run multi-layered detection pipeline
	input := DetectionInput{
		Submission:     sub,
		TestResults:    testResults,
		ActivePatterns: activePatterns,
		ProblemTitle:   problemTitle,
		ProblemTopic:   problemTopic,
	}
	diagnosis := s.detector.Analyze(input)
	if diagnosis == nil {
		return nil, nil
	}

	// 4. Persist diagnosis
	if err := s.repo.SaveDiagnosis(ctx, diagnosis); err != nil {
		s.logger.Error("failed to persist mistake diagnosis", slog.String("submission_id", sub.ID), slog.String("error", err.Error()))
		return nil, err
	}

	// 5. Update student recurring patterns
	pattern := &StudentMistakePattern{
		StudentID:      sub.StudentID,
		Category:       diagnosis.Category,
		Subtype:        diagnosis.Subtype,
		RelatedConcept: diagnosis.RelatedConcept,
	}
	if err := s.repo.UpsertPattern(ctx, pattern); err != nil {
		s.logger.Warn("failed to update student mistake pattern", slog.String("student_id", sub.StudentID), slog.String("error", err.Error()))
	}

	s.logger.Info("mistake diagnosis generated",
		slog.String("submission_id", sub.ID),
		slog.String("category", diagnosis.Category),
		slog.String("subtype", diagnosis.Subtype),
		slog.String("confidence", diagnosis.ConfidenceLevel),
	)

	return diagnosis, nil
}

func (s *serviceImpl) GetDiagnosisBySubmission(ctx context.Context, studentID, submissionID string) (*Diagnosis, error) {
	diag, err := s.repo.GetDiagnosisBySubmissionID(ctx, submissionID)
	if err != nil {
		return nil, err
	}
	if studentID != "" && diag.StudentID != studentID {
		return nil, ErrUnauthorizedAccess
	}
	return diag, nil
}

func (s *serviceImpl) GetDiagnosisByID(ctx context.Context, studentID, diagnosisID string) (*Diagnosis, error) {
	diag, err := s.repo.GetDiagnosisByID(ctx, diagnosisID)
	if err != nil {
		return nil, err
	}
	if studentID != "" && diag.StudentID != studentID {
		return nil, ErrUnauthorizedAccess
	}
	return diag, nil
}

func (s *serviceImpl) GetStudentMistakes(ctx context.Context, studentID string, limit, offset int) ([]Diagnosis, error) {
	return s.repo.ListDiagnosesByStudent(ctx, studentID, limit, offset)
}

func (s *serviceImpl) GetStudentPatterns(ctx context.Context, studentID string) ([]StudentMistakePattern, error) {
	return s.repo.GetActivePatterns(ctx, studentID)
}

func (s *serviceImpl) GetStudentMistakeSummary(ctx context.Context, studentID string) (*MistakeSummary, error) {
	diagnoses, err := s.repo.ListDiagnosesByStudent(ctx, studentID, 100, 0)
	if err != nil {
		return nil, err
	}

	patterns, err := s.repo.GetActivePatterns(ctx, studentID)
	if err != nil {
		return nil, err
	}

	catCounts := make(map[string]int)
	conceptCounts := make(map[string]int)
	for _, d := range diagnoses {
		catCounts[d.Category]++
		if d.RelatedConcept != "" {
			conceptCounts[d.RelatedConcept]++
		}
	}

	// Concepts with 2 or more mistakes need attention
	var needsAttention []string
	for concept, count := range conceptCounts {
		if count >= 2 {
			needsAttention = append(needsAttention, concept)
		}
	}

	recent := diagnoses
	if len(recent) > 5 {
		recent = recent[:5]
	}

	return &MistakeSummary{
		StudentID:              studentID,
		TotalMistakes:          len(diagnoses),
		CategoryCounts:         catCounts,
		RecurringPatterns:      patterns,
		RecentDiagnoses:        recent,
		NeedsAttentionConcepts: needsAttention,
	}, nil
}

func (s *serviceImpl) RecordFeedback(ctx context.Context, studentID, diagnosisID, feedbackType, comment string) error {
	fb := &MistakeFeedback{
		DiagnosisID:     diagnosisID,
		StudentID:       studentID,
		FeedbackType:    feedbackType,
		StudentFeedback: comment,
	}
	return s.repo.SaveFeedback(ctx, fb)
}

func (s *serviceImpl) MarkPatternCorrected(ctx context.Context, studentID, category, subtype string) error {
	return s.repo.MarkPatternCorrected(ctx, studentID, category, subtype, "")
}
