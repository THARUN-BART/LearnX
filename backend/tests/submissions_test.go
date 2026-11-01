package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/assessment"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
	"github.com/THARUN-BART/LearnX/backend/internal/workers"
	"github.com/golang-jwt/jwt/v5"
)

func createTestToken(t *testing.T, secret, userID, email string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  "student",
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign test token: %v", err)
	}
	return tokenStr
}

type submissionTestContext struct {
	router          http.Handler
	submissionsRepo submissions.Repository
	testCasesRepo   testcases.Repository
	progressRepo    progress.Repository
	problemsRepo    problems.Repository
	mockClient      *execution.MockExecutionClient
	workerPool      *workers.WorkerPool
	jwtSecret       string
}

func setupSubmissionTest(t *testing.T) *submissionTestContext {
	t.Helper()
	catalog := loadSampleCatalog(t)
	dsaRepos := dsa.NewMemoryRepositories(catalog)

	submissionsRepo := submissions.NewMemoryRepository()
	testCasesRepo := testcases.NewMemoryRepository()
	mockClient := execution.NewMockExecutionClient()

	engine := assessment.NewEngine(mockClient, nil)
	jobQueue := workers.NewChannelJobQueue(500)

	workerPool := workers.NewWorkerPool(
		jobQueue,
		submissionsRepo,
		testCasesRepo,
		dsaRepos.Progress,
		engine,
		2,
		nil,
	)
	workerPool.Start()

	testCasesService := testcases.NewService(testCasesRepo)
	testCasesHandler := testcases.NewHandler(testCasesService)

	submissionsService := submissions.NewService(submissionsRepo, dsaRepos.Problems, jobQueue)
	submissionsHandler := submissions.NewHandler(submissionsService)

	jwtSecret := "test-jwt-secret-key-32-bytes-long!"
	authMiddleware := middleware.NewAuthMiddleware(nil, jwtSecret)

	cfg := &config.Config{ClientURL: "*"}

	router := server.BuildRouter(
		cfg,
		nil,
		nil,
		authMiddleware,
		nil,
		server.WithDSA(dsa.NewHandlers(dsaRepos)),
		server.WithSubmissions(&server.SubmissionHandlers{
			Submissions: submissionsHandler,
			TestCases:   testCasesHandler,
		}),
	)

	return &submissionTestContext{
		router:          router,
		submissionsRepo: submissionsRepo,
		testCasesRepo:   testCasesRepo,
		progressRepo:    dsaRepos.Progress,
		problemsRepo:    dsaRepos.Problems,
		mockClient:      mockClient,
		workerPool:      workerPool,
		jwtSecret:       jwtSecret,
	}
}

func TestSubmissions_FullLifecycleAndProgress(t *testing.T) {
	ctx := setupSubmissionTest(t)
	defer ctx.workerPool.Stop()

	studentID := "student-uuid-111"
	token := createTestToken(t, ctx.jwtSecret, studentID, "student1@learnx.com")

	// Get a sample problem from the loaded catalogue
	allProblems, _, err := ctx.problemsRepo.List(context.Background(), model.ProblemFilter{})
	if err != nil || len(allProblems) == 0 {
		t.Fatalf("no problems in test catalogue: %v", err)
	}
	targetProblem := allProblems[0]

	// Register 1 sample test case and 1 hidden test case
	_ = ctx.testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ProblemID:      targetProblem.ID,
		InputData:      "2 2",
		ExpectedOutput: "4",
		IsHidden:       false,
		TestOrder:      1,
	})
	_ = ctx.testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ProblemID:      targetProblem.ID,
		InputData:      "100 200",
		ExpectedOutput: "300",
		IsHidden:       true,
		TestOrder:      2,
	})

	// 1. Submit valid solution (Accepted)
	body, _ := json.Marshal(map[string]interface{}{
		"problem_id":  targetProblem.ID,
		"language":    "cpp",
		"source_code": "#include <iostream>\nint main() { std::cout << 4; }",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx.router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202 Accepted on submit, got %d: %s", w.Code, w.Body.String())
	}

	var submitResp struct {
		Success bool `json:"success"`
		Data    struct {
			SubmissionID string `json:"submission_id"`
			Status       string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &submitResp)
	subID := submitResp.Data.SubmissionID
	if subID == "" {
		t.Fatal("expected submission_id to be returned")
	}

	// 2. Poll submission status until worker completes
	var finalStatus string
	for i := 0; i < 20; i++ {
		time.Sleep(50 * time.Millisecond)
		statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/"+subID, nil)
		statusReq.Header.Set("Authorization", "Bearer "+token)
		wStatus := httptest.NewRecorder()
		ctx.router.ServeHTTP(wStatus, statusReq)

		var statusResp struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		_ = json.Unmarshal(wStatus.Body.Bytes(), &statusResp)
		if statusResp.Data.Status != execution.StatusQueued && statusResp.Data.Status != execution.StatusProcessing {
			finalStatus = statusResp.Data.Status
			break
		}
	}

	if finalStatus != execution.StatusAccepted {
		t.Errorf("expected final status Accepted, got %s", finalStatus)
	}

	// 3. Verify student problem progress was updated to Solved!
	prog, err := ctx.progressRepo.Get(context.Background(), studentID, targetProblem.ID)
	if err != nil {
		t.Fatalf("get progress error: %v", err)
	}
	if prog.Status != model.StatusSolved {
		t.Errorf("expected progress status 'Solved', got %s", prog.Status)
	}
	if prog.TotalAttempts < 1 || prog.SuccessfulAttempts < 1 {
		t.Errorf("expected attempts to be incremented, got total: %d, successful: %d", prog.TotalAttempts, prog.SuccessfulAttempts)
	}

	// 4. Submit an incorrect solution (Wrong Answer)
	bodyWA, _ := json.Marshal(map[string]interface{}{
		"problem_id":  targetProblem.ID,
		"language":    "cpp",
		"source_code": "WRONG_ANSWER_CODE",
	})
	reqWA := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(bodyWA))
	reqWA.Header.Set("Authorization", "Bearer "+token)
	reqWA.Header.Set("Content-Type", "application/json")
	wWA := httptest.NewRecorder()
	ctx.router.ServeHTTP(wWA, reqWA)

	var submitWAResp struct {
		Data struct {
			SubmissionID string `json:"submission_id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wWA.Body.Bytes(), &submitWAResp)
	subWAID := submitWAResp.Data.SubmissionID

	// Wait for worker
	for i := 0; i < 20; i++ {
		time.Sleep(50 * time.Millisecond)
		sReq := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/"+subWAID, nil)
		sReq.Header.Set("Authorization", "Bearer "+token)
		wS := httptest.NewRecorder()
		ctx.router.ServeHTTP(wS, sReq)
		var sResp struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}
		_ = json.Unmarshal(wS.Body.Bytes(), &sResp)
		if sResp.Data.Status != execution.StatusQueued && sResp.Data.Status != execution.StatusProcessing {
			break
		}
	}

	// Rule: Subsequent wrong answers must NOT reset 'Solved' status back to 'Attempted'!
	progAfterWA, _ := ctx.progressRepo.Get(context.Background(), studentID, targetProblem.ID)
	if progAfterWA.Status != model.StatusSolved {
		t.Errorf("expected problem to remain 'Solved' after later wrong answer, got: %s", progAfterWA.Status)
	}
	if progAfterWA.TotalAttempts < 2 {
		t.Errorf("expected total_attempts to increment to 2, got %d", progAfterWA.TotalAttempts)
	}
}

func TestSubmissions_HiddenTestCaseSecurity(t *testing.T) {
	ctx := setupSubmissionTest(t)
	defer ctx.workerPool.Stop()

	studentID := "student-uuid-sec"
	token := createTestToken(t, ctx.jwtSecret, studentID, "sec@learnx.com")

	allProblems, _, _ := ctx.problemsRepo.List(context.Background(), model.ProblemFilter{})
	targetProblem := allProblems[0]

	_ = ctx.testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ProblemID:      targetProblem.ID,
		InputData:      "VISIBLE_INPUT",
		ExpectedOutput: "VISIBLE_OUTPUT",
		IsHidden:       false,
		TestOrder:      1,
	})
	_ = ctx.testCasesRepo.Create(context.Background(), &testcases.TestCase{
		ProblemID:      targetProblem.ID,
		InputData:      "SUPER_SECRET_HIDDEN_INPUT_999",
		ExpectedOutput: "SUPER_SECRET_HIDDEN_OUTPUT_999",
		IsHidden:       true,
		TestOrder:      2,
	})

	body, _ := json.Marshal(map[string]interface{}{
		"problem_id":  targetProblem.ID,
		"language":    "python",
		"source_code": "print('VISIBLE_OUTPUT')",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx.router.ServeHTTP(w, req)

	var submitResp struct {
		Data struct {
			SubmissionID string `json:"submission_id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &submitResp)
	subID := submitResp.Data.SubmissionID

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	// Fetch detailed result
	resReq := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/"+subID+"/result", nil)
	resReq.Header.Set("Authorization", "Bearer "+token)
	wRes := httptest.NewRecorder()
	ctx.router.ServeHTTP(wRes, resReq)

	if wRes.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on get result, got %d: %s", wRes.Code, wRes.Body.String())
	}

	bodyStr := wRes.Body.String()
	// CRITICAL SECURITY REQUIREMENT: Hidden test inputs or outputs MUST NEVER be exposed!
	if strings.Contains(bodyStr, "SUPER_SECRET_HIDDEN_INPUT_999") {
		t.Errorf("SECURITY BREACH: Hidden test case input leaked in public response!")
	}
	if strings.Contains(bodyStr, "SUPER_SECRET_HIDDEN_OUTPUT_999") {
		t.Errorf("SECURITY BREACH: Hidden test case expected output leaked in public response!")
	}

	// Verify sample testcases endpoint does not leak hidden testcases
	sampleReq := httptest.NewRequest(http.MethodGet, "/api/v1/problems/"+targetProblem.ID+"/testcases/sample", nil)
	sampleReq.Header.Set("Authorization", "Bearer "+token)
	wSample := httptest.NewRecorder()
	ctx.router.ServeHTTP(wSample, sampleReq)

	if strings.Contains(wSample.Body.String(), "SUPER_SECRET_HIDDEN") {
		t.Errorf("SECURITY BREACH: sample testcase endpoint exposed hidden test case!")
	}
}

func TestSubmissions_StudentOwnershipIsolation(t *testing.T) {
	ctx := setupSubmissionTest(t)
	defer ctx.workerPool.Stop()

	studentA := "student-uuid-A"
	tokenA := createTestToken(t, ctx.jwtSecret, studentA, "studentA@learnx.com")

	studentB := "student-uuid-B"
	tokenB := createTestToken(t, ctx.jwtSecret, studentB, "studentB@learnx.com")

	allProblems, _, _ := ctx.problemsRepo.List(context.Background(), model.ProblemFilter{})
	targetProblem := allProblems[0]

	// Student A submits
	body, _ := json.Marshal(map[string]interface{}{
		"problem_id":  targetProblem.ID,
		"language":    "go",
		"source_code": "package main\nfunc main() {}",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	ctx.router.ServeHTTP(w, req)

	var submitResp struct {
		Data struct {
			SubmissionID string `json:"submission_id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &submitResp)
	subID := submitResp.Data.SubmissionID

	// Student B attempts to access Student A's submission status
	statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/"+subID, nil)
	statusReq.Header.Set("Authorization", "Bearer "+tokenB)
	wStatus := httptest.NewRecorder()
	ctx.router.ServeHTTP(wStatus, statusReq)

	if wStatus.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when Student B accesses Student A's submission, got: %d", wStatus.Code)
	}

	// Student B attempts to access Student A's result
	resultReq := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/"+subID+"/result", nil)
	resultReq.Header.Set("Authorization", "Bearer "+tokenB)
	wResult := httptest.NewRecorder()
	ctx.router.ServeHTTP(wResult, resultReq)

	if wResult.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden when Student B accesses Student A's result, got: %d", wResult.Code)
	}
}

func TestSubmissions_ValidationErrors(t *testing.T) {
	ctx := setupSubmissionTest(t)
	defer ctx.workerPool.Stop()

	token := createTestToken(t, ctx.jwtSecret, "student-val", "val@learnx.com")

	allProblems, _, _ := ctx.problemsRepo.List(context.Background(), model.ProblemFilter{})
	validProblemID := allProblems[0].ID

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "missing problem_id",
			payload: map[string]interface{}{
				"problem_id":  "",
				"language":    "cpp",
				"source_code": "int main(){}",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty source code",
			payload: map[string]interface{}{
				"problem_id":  validProblemID,
				"language":    "cpp",
				"source_code": "   ",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unsupported language",
			payload: map[string]interface{}{
				"problem_id":  validProblemID,
				"language":    "brainfuck",
				"source_code": "++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "oversized source code (>64KB)",
			payload: map[string]interface{}{
				"problem_id":  validProblemID,
				"language":    "cpp",
				"source_code": strings.Repeat("A", 65*1024),
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			ctx.router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestSubmissions_HistoryAndRetry(t *testing.T) {
	ctx := setupSubmissionTest(t)
	defer ctx.workerPool.Stop()

	studentID := "student-history-1"
	token := createTestToken(t, ctx.jwtSecret, studentID, "history@learnx.com")

	allProblems, _, _ := ctx.problemsRepo.List(context.Background(), model.ProblemFilter{})
	targetProblem := allProblems[0]

	// Create 3 submissions
	for i := 1; i <= 3; i++ {
		body, _ := json.Marshal(map[string]interface{}{
			"problem_id":  targetProblem.ID,
			"language":    "go",
			"source_code": fmt.Sprintf("package main\n// submission %d", i),
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/submissions", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		ctx.router.ServeHTTP(w, req)
	}

	time.Sleep(150 * time.Millisecond)

	// 1. Get student submission history
	hReq := httptest.NewRequest(http.MethodGet, "/api/v1/student/submissions?limit=10&page=1", nil)
	hReq.Header.Set("Authorization", "Bearer "+token)
	wH := httptest.NewRecorder()
	ctx.router.ServeHTTP(wH, hReq)

	if wH.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", wH.Code, wH.Body.String())
	}
	var histResp struct {
		Data struct {
			Total       int `json:"total"`
			Submissions []struct {
				SubmissionID string `json:"submission_id"`
			} `json:"submissions"`
		} `json:"data"`
	}
	_ = json.Unmarshal(wH.Body.Bytes(), &histResp)
	if histResp.Data.Total != 3 {
		t.Errorf("expected 3 total submissions, got %d", histResp.Data.Total)
	}

	// 2. Get latest problem submission
	latestReq := httptest.NewRequest(http.MethodGet, "/api/v1/problems/"+targetProblem.ID+"/submissions/latest", nil)
	latestReq.Header.Set("Authorization", "Bearer "+token)
	wLatest := httptest.NewRecorder()
	ctx.router.ServeHTTP(wLatest, latestReq)

	if wLatest.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", wLatest.Code, wLatest.Body.String())
	}

	// 3. Retry a submission
	firstSubID := histResp.Data.Submissions[0].SubmissionID
	retryReq := httptest.NewRequest(http.MethodPost, "/api/v1/submissions/"+firstSubID+"/retry", nil)
	retryReq.Header.Set("Authorization", "Bearer "+token)
	wRetry := httptest.NewRecorder()
	ctx.router.ServeHTTP(wRetry, retryReq)

	if wRetry.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on retry, got %d: %s", wRetry.Code, wRetry.Body.String())
	}
}
