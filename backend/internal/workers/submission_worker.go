package workers

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/assessment"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
)

type WorkerPool struct {
	queue           JobQueue
	submissionsRepo submissions.Repository
	testCasesRepo   testcases.Repository
	progressRepo    progress.Repository
	mistakesService mistakes.Service
	memoryService   interface {
		RecordSubmissionEvent(ctx context.Context, studentID string, problem *model.Problem, subID string, accepted bool, failedReason string) error
		RecordMistakeEvent(ctx context.Context, studentID, problemID, submissionID, category, subtype, conceptKey string, confidence float64) error
	}
	pathService interface {
		Recalculate(ctx context.Context, studentID string)
	}
	recService interface {
		OnSubmissionEvaluated(ctx context.Context, studentID, problemID, status string)
	}
	problemsRepo problems.Repository
	engine       assessment.Engine
	concurrency  int
	logger       *slog.Logger
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
}

func (p *WorkerPool) SetMistakesService(service mistakes.Service) {
	p.mistakesService = service
}

func (p *WorkerPool) SetMemoryService(service interface {
	RecordSubmissionEvent(ctx context.Context, studentID string, problem *model.Problem, subID string, accepted bool, failedReason string) error
	RecordMistakeEvent(ctx context.Context, studentID, problemID, submissionID, category, subtype, conceptKey string, confidence float64) error
}, probRepo problems.Repository) {
	p.memoryService = service
	p.problemsRepo = probRepo
}

func (p *WorkerPool) SetPathService(service interface {
	Recalculate(ctx context.Context, studentID string)
}) {
	p.pathService = service
}

func (p *WorkerPool) SetRecommendationsService(service interface {
	OnSubmissionEvaluated(ctx context.Context, studentID, problemID, status string)
}) {
	p.recService = service
}

func NewWorkerPool(
	queue JobQueue,
	submissionsRepo submissions.Repository,
	testCasesRepo testcases.Repository,
	progressRepo progress.Repository,
	engine assessment.Engine,
	concurrency int,
	logger *slog.Logger,
) *WorkerPool {
	if concurrency <= 0 {
		concurrency = 4
	}
	if logger == nil {
		logger = slog.Default()
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		queue:           queue,
		submissionsRepo: submissionsRepo,
		testCasesRepo:   testCasesRepo,
		progressRepo:    progressRepo,
		engine:          engine,
		concurrency:     concurrency,
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (p *WorkerPool) Start() {
	p.logger.Info("starting submission worker pool", slog.Int("concurrency", p.concurrency))
	for i := 0; i < p.concurrency; i++ {
		p.wg.Add(1)
		go p.worker(i + 1)
	}

	// Recover any pending submissions left in database from a previous crash/restart
	go p.recoverPendingJobs()
}

func (p *WorkerPool) Stop() {
	p.logger.Info("stopping submission worker pool...")
	p.cancel()
	p.wg.Wait()
	p.logger.Info("submission worker pool stopped cleanly")
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case <-p.ctx.Done():
			return
		case submissionID, ok := <-p.queue.Channel():
			if !ok {
				return
			}
			p.processSubmission(submissionID)
		}
	}
}

func (p *WorkerPool) processSubmission(submissionID string) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error("panic in submission worker", slog.Any("panic", r), slog.String("submission_id", submissionID))
			now := time.Now()
			_ = p.submissionsRepo.Update(context.Background(), &submissions.Submission{
				ID:           submissionID,
				Status:       execution.StatusInternalError,
				ErrorMessage: fmt.Sprintf("internal worker crash: %v", r),
				CompletedAt:  &now,
			})
		}
	}()

	ctx, cancel := context.WithTimeout(p.ctx, 2*time.Minute)
	defer cancel()

	sub, err := p.submissionsRepo.GetByID(ctx, submissionID)
	if err != nil {
		p.logger.Error("failed to load submission for processing", slog.String("id", submissionID), slog.String("error", err.Error()))
		return
	}

	// Guard against duplicate processing
	if sub.Status != execution.StatusQueued && sub.Status != execution.StatusProcessing {
		p.logger.Info("submission already finalized, skipping", slog.String("id", submissionID), slog.String("status", sub.Status))
		return
	}

	// Mark status as Processing
	sub.Status = execution.StatusProcessing
	if err := p.submissionsRepo.Update(ctx, sub); err != nil {
		p.logger.Warn("could not update submission to Processing", slog.String("error", err.Error()))
	}

	// Resolve language config
	langCfg, err := execution.GetRegistry().Resolve(sub.Language)
	if err != nil {
		p.markFailed(ctx, sub, execution.StatusInternalError, err.Error())
		return
	}

	// Load test cases for problem (including hidden)
	tests, err := p.testCasesRepo.ListByProblem(ctx, sub.ProblemID, true)
	if err != nil {
		p.markFailed(ctx, sub, execution.StatusInternalError, fmt.Sprintf("failed to load test cases: %v", err))
		return
	}

	// Evaluate code using assessment engine
	evalResult, err := p.engine.Evaluate(ctx, langCfg, sub.SourceCode, tests)
	if err != nil {
		p.markFailed(ctx, sub, execution.StatusInternalError, fmt.Sprintf("assessment failure: %v", err))
		return
	}

	// Save individual test results
	testResults := make([]submissions.SubmissionTestResult, 0, len(evalResult.TestCaseEvaluations))
	for _, tcEval := range evalResult.TestCaseEvaluations {
		testResults = append(testResults, submissions.SubmissionTestResult{
			SubmissionID:    sub.ID,
			TestCaseID:      tcEval.TestCaseID,
			Status:          tcEval.Status,
			ExecutionTimeMs: tcEval.ExecutionTimeMs,
			MemoryUsedKB:    tcEval.MemoryUsedKB,
			ActualOutput:    tcEval.ActualOutput,
			ErrorMessage:    tcEval.ErrorMessage,
			IsHidden:        tcEval.IsHidden,
		})
	}
	if err := p.submissionsRepo.SaveTestResults(ctx, testResults); err != nil {
		p.logger.Error("failed saving individual test results", slog.String("error", err.Error()))
	}

	// Update submission with final evaluation status and metrics
	completedAt := evalResult.CompletedAt
	sub.Status = evalResult.OverallStatus
	sub.ExecutionTimeMs = evalResult.MaxExecutionTimeMs
	sub.MemoryUsedKB = evalResult.MaxMemoryUsedKB
	sub.TotalTestCases = evalResult.TotalTestCases
	sub.PassedTestCases = evalResult.PassedTestCases
	sub.FailedTestCases = evalResult.FailedTestCases
	sub.ErrorMessage = evalResult.ErrorMessage
	sub.CompletedAt = &completedAt

	if err := p.submissionsRepo.Update(ctx, sub); err != nil {
		p.logger.Error("failed updating submission final status", slog.String("error", err.Error()))
	}

	// Update student progress in the existing progress module
	p.updateStudentProgress(ctx, sub)

	// Trigger intelligent mistake detection if submission was not accepted
	if p.mistakesService != nil {
		if sub.Status != execution.StatusAccepted {
			diag, err := p.mistakesService.DiagnoseSubmission(ctx, sub, testResults, "", "")
			if err != nil {
				p.logger.Warn("mistake diagnosis failed", slog.String("submission_id", sub.ID), slog.String("error", err.Error()))
			} else if diag != nil {
				p.logger.Info("mistake diagnosed",
					slog.String("submission_id", sub.ID),
					slog.String("category", diag.Category),
					slog.String("subtype", diag.Subtype),
					slog.String("confidence", diag.ConfidenceLevel),
				)
				if p.memoryService != nil {
					_ = p.memoryService.RecordMistakeEvent(ctx, sub.StudentID, sub.ProblemID, sub.ID, diag.Category, diag.Subtype, diag.RelatedConcept, diag.ConfidenceScore)
				}
			}
		} else {
			_ = p.mistakesService.MarkPatternCorrected(ctx, sub.StudentID, "", "")
		}
	}

	// Trigger Student Learning Memory ingestion
	if p.memoryService != nil {
		var prob *model.Problem
		if p.problemsRepo != nil {
			prob, _ = p.problemsRepo.Get(ctx, sub.ProblemID)
		}
		isAccepted := sub.Status == execution.StatusAccepted
		_ = p.memoryService.RecordSubmissionEvent(ctx, sub.StudentID, prob, sub.ID, isAccepted, sub.ErrorMessage)
	}

	// Trigger Adaptive Learning Path re-evaluation
	if p.pathService != nil {
		p.pathService.Recalculate(ctx, sub.StudentID)
	}

	// Trigger Problem Recommendations feedback update
	if p.recService != nil {
		p.recService.OnSubmissionEvaluated(ctx, sub.StudentID, sub.ProblemID, sub.Status)
	}

	p.logger.Info("submission completed",
		slog.String("submission_id", sub.ID),
		slog.String("status", sub.Status),
		slog.Int("passed", sub.PassedTestCases),
		slog.Int("total", sub.TotalTestCases),
	)
}

func (p *WorkerPool) updateStudentProgress(ctx context.Context, sub *submissions.Submission) {
	if p.progressRepo == nil {
		return
	}

	if sub.Status == execution.StatusAccepted {
		_, err := p.progressRepo.MarkSolved(ctx, sub.StudentID, sub.ProblemID)
		if err != nil {
			p.logger.Error("failed to mark problem as solved in progress",
				slog.String("student_id", sub.StudentID),
				slog.String("problem_id", sub.ProblemID),
				slog.String("error", err.Error()),
			)
		}
	} else {
		_, err := p.progressRepo.MarkAttempted(ctx, sub.StudentID, sub.ProblemID)
		if err != nil {
			p.logger.Error("failed to mark problem as attempted in progress",
				slog.String("student_id", sub.StudentID),
				slog.String("problem_id", sub.ProblemID),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (p *WorkerPool) markFailed(ctx context.Context, sub *submissions.Submission, status, errMsg string) {
	now := time.Now()
	sub.Status = status
	sub.ErrorMessage = errMsg
	sub.CompletedAt = &now
	_ = p.submissionsRepo.Update(ctx, sub)
}

func (p *WorkerPool) recoverPendingJobs() {
	time.Sleep(1 * time.Second)
	pending, err := p.submissionsRepo.GetPendingSubmissions(context.Background(), 50)
	if err != nil {
		return
	}
	for _, sub := range pending {
		_ = p.queue.Enqueue(sub.ID)
	}
}
