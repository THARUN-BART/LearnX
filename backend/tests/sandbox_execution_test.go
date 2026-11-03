package tests

import (
	"context"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/assessment"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/importer"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
	"github.com/THARUN-BART/LearnX/backend/internal/workers"
)

func TestSandboxExecutionWithCuratedProblems(t *testing.T) {
	mem := model.NewMemory()
	seed, err := importer.LoadSeedFile("../data/strivers-a2z-sheet.json")
	if err != nil {
		t.Fatalf("failed to load seed file: %v", err)
	}
	if _, err := importer.ImportMemory(mem, seed); err != nil {
		t.Fatalf("failed to import memory: %v", err)
	}

	tcRepo := testcases.NewMemoryRepository()
	if err := importer.SeedMemoryTestCases(context.Background(), mem, seed, tcRepo); err != nil {
		t.Fatalf("failed to seed test cases: %v", err)
	}

	// Verify problem count: 19 modules * 3 problems = 57 problems
	if len(mem.Problems) != 57 {
		t.Fatalf("expected 57 curated problems, got %d", len(mem.Problems))
	}
	t.Logf("Successfully loaded %d curated problems into memory", len(mem.Problems))

	// Find "reverse-a-number" and "pattern-1"
	var reverseNumProb, patternProb *model.Problem
	for _, p := range mem.Problems {
		if p.Slug == "reverse-a-number" {
			reverseNumProb = p
		}
		if p.Slug == "pattern-1" {
			patternProb = p
		}
	}
	if reverseNumProb == nil {
		t.Fatal("reverse-a-number problem not found")
	}
	if patternProb == nil {
		t.Fatal("pattern-1 problem not found")
	}

	// Verify test cases were loaded
	tests, err := tcRepo.ListByProblem(context.Background(), reverseNumProb.ID, true)
	if err != nil || len(tests) == 0 {
		t.Fatalf("expected test cases for reverse-a-number, got %d, err: %v", len(tests), err)
	}
	t.Logf("Found %d test cases for reverse-a-number", len(tests))

	// Setup local runner sandbox & worker pool
	runner := execution.NewLocalRunner()
	engine := assessment.NewEngine(runner, nil)
	queue := workers.NewChannelJobQueue(100)
	subRepo := submissions.NewMemoryRepository()
	progRepo := progress.NewMemoryRepository(mem)

	pool := workers.NewWorkerPool(queue, subRepo, tcRepo, progRepo, engine, 2, nil)
	pool.Start()
	defer pool.Stop()

	// 1. Submit Correct Python solution for reverse-a-number
	correctPyCode := `import sys

def main():
    val = sys.stdin.read().strip()
    if not val:
        return
    n = int(val)
    if n < 0:
        print(-int(str(abs(n))[::-1]))
    else:
        print(int(str(n)[::-1]))

if __name__ == "__main__":
    main()
`
	sub1 := &submissions.Submission{
		StudentID:  "student-1",
		ProblemID:  reverseNumProb.ID,
		Language:   "python",
		SourceCode: correctPyCode,
		Status:     execution.StatusQueued,
	}
	if err := subRepo.Create(context.Background(), sub1); err != nil {
		t.Fatalf("create sub1: %v", err)
	}
	_ = queue.Enqueue(sub1.ID)

	// Poll sub1 for verdict
	var res1 *submissions.Submission
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		res1, _ = subRepo.GetByID(context.Background(), sub1.ID)
		if res1 != nil && res1.Status != execution.StatusQueued && res1.Status != execution.StatusProcessing {
			break
		}
	}
	if res1 == nil || res1.Status != execution.StatusAccepted {
		t.Fatalf("expected Accepted for correct Python solution, got status=%v, err=%v", res1.Status, res1.ErrorMessage)
	}
	if res1.PassedTestCases != res1.TotalTestCases || res1.TotalTestCases == 0 {
		t.Fatalf("expected all test cases to pass, got %d/%d", res1.PassedTestCases, res1.TotalTestCases)
	}
	t.Logf("Sandbox verified: Correct Python solution passed (%d/%d test cases)", res1.PassedTestCases, res1.TotalTestCases)

	// 2. Submit Correct C++ solution for pattern-1
	correctCppCode := `#include <iostream>
using namespace std;

int main() {
    int n;
    if (cin >> n) {
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                cout << "*";
            }
            cout << "\n";
        }
    }
    return 0;
}
`
	sub2 := &submissions.Submission{
		StudentID:  "student-1",
		ProblemID:  patternProb.ID,
		Language:   "cpp",
		SourceCode: correctCppCode,
		Status:     execution.StatusQueued,
	}
	if err := subRepo.Create(context.Background(), sub2); err != nil {
		t.Fatalf("create sub2: %v", err)
	}
	_ = queue.Enqueue(sub2.ID)

	var res2 *submissions.Submission
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		res2, _ = subRepo.GetByID(context.Background(), sub2.ID)
		if res2 != nil && res2.Status != execution.StatusQueued && res2.Status != execution.StatusProcessing {
			break
		}
	}
	if res2 == nil || res2.Status != execution.StatusAccepted {
		t.Fatalf("expected Accepted for correct C++ solution, got status=%v, err=%v", res2.Status, res2.ErrorMessage)
	}
	t.Logf("Sandbox verified: Correct C++ solution compiled & passed (%d/%d test cases)", res2.PassedTestCases, res2.TotalTestCases)

	// 3. Submit Wrong Python solution
	wrongPyCode := `print("completely wrong output")`
	sub3 := &submissions.Submission{
		StudentID:  "student-1",
		ProblemID:  reverseNumProb.ID,
		Language:   "python",
		SourceCode: wrongPyCode,
		Status:     execution.StatusQueued,
	}
	if err := subRepo.Create(context.Background(), sub3); err != nil {
		t.Fatalf("create sub3: %v", err)
	}
	_ = queue.Enqueue(sub3.ID)

	var res3 *submissions.Submission
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		res3, _ = subRepo.GetByID(context.Background(), sub3.ID)
		if res3 != nil && res3.Status != execution.StatusQueued && res3.Status != execution.StatusProcessing {
			break
		}
	}
	if res3 == nil || res3.Status != execution.StatusWrongAnswer {
		t.Fatalf("expected Wrong Answer, got status=%v", res3.Status)
	}
	t.Logf("Sandbox verified: Wrong solution correctly caught as Wrong Answer")

	// 4. Submit Compilation Error (syntax error)
	syntaxErrPyCode := `def broken(`
	sub4 := &submissions.Submission{
		StudentID:  "student-1",
		ProblemID:  reverseNumProb.ID,
		Language:   "python",
		SourceCode: syntaxErrPyCode,
		Status:     execution.StatusQueued,
	}
	if err := subRepo.Create(context.Background(), sub4); err != nil {
		t.Fatalf("create sub4: %v", err)
	}
	_ = queue.Enqueue(sub4.ID)

	var res4 *submissions.Submission
	for i := 0; i < 30; i++ {
		time.Sleep(200 * time.Millisecond)
		res4, _ = subRepo.GetByID(context.Background(), sub4.ID)
		if res4 != nil && res4.Status != execution.StatusQueued && res4.Status != execution.StatusProcessing {
			break
		}
	}
	if res4 == nil || (res4.Status != execution.StatusCompilationError && res4.Status != execution.StatusRuntimeError) {
		t.Fatalf("expected Compilation/Runtime Error for syntax error, got status=%v", res4.Status)
	}
	t.Logf("Sandbox verified: Syntax error correctly caught as %v", res4.Status)
}
