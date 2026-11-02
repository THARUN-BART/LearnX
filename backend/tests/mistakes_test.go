package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/crypto"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
)

func TestMistakeDetection_DeterministicLayers(t *testing.T) {
	detector := mistakes.NewDetector()

	t.Run("Compilation Error Syntax", func(t *testing.T) {
		sub := &submissions.Submission{
			ID:           crypto.GenerateUUID(),
			StudentID:    "student-1",
			ProblemID:    "prob-1",
			Language:     "python",
			Status:       execution.StatusCompilationError,
			ErrorMessage: "SyntaxError: invalid syntax on line 4",
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag == nil {
			t.Fatal("expected non-nil diagnosis")
		}
		if diag.Category != mistakes.CategoryCompilation {
			t.Errorf("expected compilation category, got %s", diag.Category)
		}
		if diag.Subtype != "syntax_error" {
			t.Errorf("expected syntax_error subtype, got %s", diag.Subtype)
		}
		if diag.ConfidenceLevel != mistakes.ConfidenceHigh {
			t.Errorf("expected High confidence, got %s", diag.ConfidenceLevel)
		}
	})

	t.Run("Runtime Error Out of Bounds", func(t *testing.T) {
		sub := &submissions.Submission{
			ID:           crypto.GenerateUUID(),
			StudentID:    "student-1",
			ProblemID:    "prob-1",
			Language:     "cpp",
			Status:       execution.StatusRuntimeError,
			ErrorMessage: "fatal error: index out of range [5] with length 5 (SIGSEGV)",
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag == nil {
			t.Fatal("expected non-nil diagnosis")
		}
		if diag.Category != mistakes.CategoryRuntime {
			t.Errorf("expected runtime category, got %s", diag.Category)
		}
		if diag.Subtype != "out_of_bounds" {
			t.Errorf("expected out_of_bounds subtype, got %s", diag.Subtype)
		}
	})

	t.Run("Time Limit Exceeded", func(t *testing.T) {
		sub := &submissions.Submission{
			ID:              crypto.GenerateUUID(),
			StudentID:       "student-1",
			ProblemID:       "prob-1",
			Language:        "java",
			Status:          execution.StatusTimeLimitExceeded,
			ExecutionTimeMs: 2050,
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag == nil {
			t.Fatal("expected non-nil diagnosis")
		}
		if diag.Category != mistakes.CategoryComplexity {
			t.Errorf("expected complexity category, got %s", diag.Category)
		}
		if diag.Subtype != "time_limit_exceeded" {
			t.Errorf("expected time_limit_exceeded subtype, got %s", diag.Subtype)
		}
	})
}

func TestMistakeDetection_StaticCodeAnalysis(t *testing.T) {
	detector := mistakes.NewDetector()

	t.Run("Off-by-one loop bound", func(t *testing.T) {
		code := `
		for (int i = 0; i <= arr.length; i++) {
			sum += arr[i];
		}`
		sub := &submissions.Submission{
			ID:           crypto.GenerateUUID(),
			StudentID:    "student-1",
			ProblemID:    "prob-1",
			Language:     "java",
			Status:       execution.StatusRuntimeError,
			ErrorMessage: "java.lang.ArrayIndexOutOfBoundsException: Index 5 out of bounds for length 5",
			SourceCode:   code,
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag.Subtype != "off_by_one" {
			t.Errorf("expected off_by_one detection, got %s", diag.Subtype)
		}
		if diag.DetectionMethod != mistakes.MethodStatic {
			t.Errorf("expected static_analysis detection method, got %s", diag.DetectionMethod)
		}
	})

	t.Run("Integer midpoint overflow", func(t *testing.T) {
		code := `
		int low = 0, high = n - 1;
		while (low <= high) {
			int mid = (low + high) / 2;
			if (nums[mid] == target) return mid;
		}`
		sub := &submissions.Submission{
			ID:         crypto.GenerateUUID(),
			StudentID:  "student-1",
			ProblemID:  "prob-1",
			Language:   "cpp",
			Status:     execution.StatusWrongAnswer,
			SourceCode: code,
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag.Subtype != "integer_overflow" {
			t.Errorf("expected integer_overflow subtype, got %s", diag.Subtype)
		}
		if diag.Category != mistakes.CategoryEdgeCase {
			t.Errorf("expected edge_case category, got %s", diag.Category)
		}
	})

	t.Run("Chained pointer null traversal", func(t *testing.T) {
		code := `
		ListNode* slow = head;
		ListNode* fast = head.next.next;
		`
		sub := &submissions.Submission{
			ID:           crypto.GenerateUUID(),
			StudentID:    "student-1",
			ProblemID:    "prob-1",
			Language:     "cpp",
			Status:       execution.StatusRuntimeError,
			ErrorMessage: "SIGSEGV null pointer dereference",
			SourceCode:   code,
		}
		diag := detector.Analyze(mistakes.DetectionInput{Submission: sub})
		if diag.Subtype != "null_pointer_traversal" {
			t.Errorf("expected null_pointer_traversal subtype, got %s", diag.Subtype)
		}
		if diag.Category != mistakes.CategoryDataStructure {
			t.Errorf("expected data_structure category, got %s", diag.Category)
		}
	})
}

func TestMistakeDetection_HistoricalPatternsAndRootCause(t *testing.T) {
	detector := mistakes.NewDetector()

	t.Run("Historical Pattern Recognition and Confidence Boost", func(t *testing.T) {
		activePatterns := []mistakes.StudentMistakePattern{
			{
				StudentID:       "student-1",
				Category:        mistakes.CategoryLogical,
				Subtype:         "off_by_one",
				OccurrenceCount: 3,
			},
		}
		code := `for i := 0; i <= len(arr); i++ { sum += arr[i] }`
		sub := &submissions.Submission{
			ID:           crypto.GenerateUUID(),
			StudentID:    "student-1",
			ProblemID:    "prob-1",
			Language:     "go",
			Status:       execution.StatusRuntimeError,
			ErrorMessage: "runtime error: index out of range",
			SourceCode:   code,
		}

		diag := detector.Analyze(mistakes.DetectionInput{
			Submission:     sub,
			ActivePatterns: activePatterns,
		})
		if !diag.Evidence["recurring_mistake_detected"].(bool) {
			t.Errorf("expected recurring_mistake_detected to be true")
		}
		if diag.ConfidenceScore < 0.90 {
			t.Errorf("expected confidence score >= 0.90 with historical pattern, got %f", diag.ConfidenceScore)
		}
	})

	t.Run("Root Cause Misconception Binary Search Infinite Loop", func(t *testing.T) {
		code := `
		func search(nums []int, target int) int {
			low, high := 0, len(nums)-1
			for low <= high {
				mid := low + (high-low)/2
				if nums[mid] == target {
					return mid
				} else if nums[mid] < target {
					low = mid // missing + 1 causes infinite loop
				} else {
					high = mid - 1
				}
			}
			return -1
		}`
		sub := &submissions.Submission{
			ID:              crypto.GenerateUUID(),
			StudentID:       "student-1",
			ProblemID:       "prob-bs",
			Language:        "go",
			Status:          execution.StatusTimeLimitExceeded,
			ExecutionTimeMs: 2000,
			SourceCode:      code,
		}
		diag := detector.Analyze(mistakes.DetectionInput{
			Submission:   sub,
			ProblemTitle: "Binary Search",
			ProblemTopic: "Binary Search",
		})
		if diag.Subtype != "binary_search_boundary" {
			t.Errorf("expected binary_search_boundary subtype, got %s", diag.Subtype)
		}
		if diag.DetectionMethod != mistakes.MethodRootCause {
			t.Errorf("expected root_cause_inference method, got %s", diag.DetectionMethod)
		}
	})
}

func TestMistakes_ServiceAndHTTPIntegration(t *testing.T) {
	memRepo := mistakes.NewMemoryRepository()
	detector := mistakes.NewDetector()
	svc := mistakes.NewService(memRepo, detector, nil)
	handler := mistakes.NewHandler(svc)

	studentID := "stu-100"
	probID := "prob-test"

	// 1. Diagnose a submission
	sub := &submissions.Submission{
		ID:           "sub-test-1",
		StudentID:    studentID,
		ProblemID:    probID,
		Language:     "python",
		Status:       execution.StatusCompilationError,
		ErrorMessage: "SyntaxError: unexpected EOF while parsing",
	}
	diag, err := svc.DiagnoseSubmission(context.Background(), sub, nil, "Two Sum", "Arrays")
	if err != nil {
		t.Fatalf("DiagnoseSubmission error: %v", err)
	}
	if diag.Subtype != "syntax_error" {
		t.Errorf("expected syntax_error, got %s", diag.Subtype)
	}

	// 2. Idempotency test: Diagnose again should return existing
	diag2, err := svc.DiagnoseSubmission(context.Background(), sub, nil, "Two Sum", "Arrays")
	if err != nil || diag2.ID != diag.ID {
		t.Errorf("idempotency failed: expected ID %s, got %s", diag.ID, diag2.ID)
	}

	// Setup HTTP server with Auth and Mistakes routes
	authRepo := auth.NewMemoryRepository()
	sbClient := newRouterMockSupabase()
	sbClient.users["stu-100@learnx.com"] = &supabase.SupabaseUser{
		ID:    studentID,
		Email: "stu-100@learnx.com",
	}
	sbClient.users["stu-200@learnx.com"] = &supabase.SupabaseUser{
		ID:    "stu-200",
		Email: "stu-200@learnx.com",
	}

	authSvc := auth.NewService(authRepo, otp.NewService(&routerMockOTPClient{}, time.Minute, 5), sbClient, []byte("32-byte-key-for-aes-256-encryption"), 10*time.Minute, 15*time.Minute)
	profileSvc := profile.NewService(profile.NewMemoryRepository())
	authMW := middleware.NewAuthMiddleware(sbClient, "")

	router := server.BuildRouter(
		nil,
		auth.NewHandler(authSvc),
		profile.NewHandler(profileSvc),
		authMW,
		nil,
		server.WithIntelligence(&server.IntelligenceHandlers{
			Mistakes: handler,
		}),
	)

	token := "mock-jwt-token-for-" + studentID
	tokenOther := "mock-jwt-token-for-stu-200"

	// 3. GET /api/v1/submissions/{submission_id}/diagnosis
	req := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/sub-test-1/diagnosis", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var res struct {
		Success bool               `json:"success"`
		Data    mistakes.Diagnosis `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if res.Data.Subtype != "syntax_error" {
		t.Errorf("expected syntax_error in HTTP response, got %s", res.Data.Subtype)
	}

	// 4. Student Isolation: Another student cannot access
	reqOther := httptest.NewRequest(http.MethodGet, "/api/v1/submissions/sub-test-1/diagnosis", nil)
	reqOther.Header.Set("Authorization", "Bearer "+tokenOther)
	recOther := httptest.NewRecorder()
	router.ServeHTTP(recOther, reqOther)
	if recOther.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for unauthorized student, got %d", recOther.Code)
	}

	// 5. GET /api/v1/student/mistakes/summary
	reqSum := httptest.NewRequest(http.MethodGet, "/api/v1/student/mistakes/summary", nil)
	reqSum.Header.Set("Authorization", "Bearer "+token)
	recSum := httptest.NewRecorder()
	router.ServeHTTP(recSum, reqSum)
	if recSum.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recSum.Code)
	}

	// 6. POST /api/v1/student/mistakes/{diagnosis_id}/feedback
	fbPayload := map[string]string{
		"feedback_type":    "helpful",
		"student_feedback": "Helped me spot the missing parenthesis immediately!",
	}
	fbBytes, _ := json.Marshal(fbPayload)
	reqFB := httptest.NewRequest(http.MethodPost, "/api/v1/student/mistakes/"+diag.ID+"/feedback", bytes.NewReader(fbBytes))
	reqFB.Header.Set("Authorization", "Bearer "+token)
	recFB := httptest.NewRecorder()
	router.ServeHTTP(recFB, reqFB)
	if recFB.Code != http.StatusOK {
		t.Fatalf("expected 200 for feedback submission, got %d", recFB.Code)
	}
}
