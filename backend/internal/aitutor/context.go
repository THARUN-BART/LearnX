package aitutor

import (
	"context"
	"fmt"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
)

type LearningContext struct {
	Problem         *model.Problem
	SampleTestCases []testcases.TestCase
	Submission      *submissions.Submission
	Language        string
	Code            string
}

type ContextBuilder struct {
	problemsRepo    problems.Repository
	testCasesRepo   testcases.Repository
	submissionsRepo submissions.Repository
}

func NewContextBuilder(
	problemsRepo problems.Repository,
	testCasesRepo testcases.Repository,
	submissionsRepo submissions.Repository,
) *ContextBuilder {
	return &ContextBuilder{
		problemsRepo:    problemsRepo,
		testCasesRepo:   testCasesRepo,
		submissionsRepo: submissionsRepo,
	}
}

func (cb *ContextBuilder) Build(ctx context.Context, studentID, problemID, submissionID string) (*LearningContext, error) {
	lc := &LearningContext{}

	if problemID != "" && cb.problemsRepo != nil {
		p, err := cb.problemsRepo.Get(ctx, problemID)
		if err == nil && p != nil {
			lc.Problem = p
		}
	}

	// Always load ONLY non-hidden (sample) test cases!
	if problemID != "" && cb.testCasesRepo != nil {
		tcs, err := cb.testCasesRepo.ListByProblem(ctx, problemID, false)
		if err == nil {
			lc.SampleTestCases = tcs
		}
	}

	// Load submission if provided, or fetch latest submission for student and problem
	if submissionID != "" && cb.submissionsRepo != nil {
		sub, err := cb.submissionsRepo.GetByID(ctx, submissionID)
		if err == nil && sub != nil {
			lc.Submission = sub
			lc.Language = sub.Language
			lc.Code = sub.SourceCode
		}
	} else if studentID != "" && problemID != "" && cb.submissionsRepo != nil {
		sub, err := cb.submissionsRepo.GetLatestByStudentAndProblem(ctx, studentID, problemID)
		if err == nil && sub != nil {
			lc.Submission = sub
			lc.Language = sub.Language
			lc.Code = sub.SourceCode
		}
	}

	return lc, nil
}

func FormatSamples(samples []testcases.TestCase) string {
	if len(samples) == 0 {
		return "Input: standard sample input\nOutput: standard sample output"
	}
	var sb strings.Builder
	for i, tc := range samples {
		sb.WriteString(fmt.Sprintf("Example %d:\nInput:\n%s\nExpected Output:\n%s\n\n", i+1, tc.InputData, tc.ExpectedOutput))
	}
	return sb.String()
}
