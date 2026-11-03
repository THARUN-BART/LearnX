package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/learningmemory"
	"github.com/THARUN-BART/LearnX/backend/internal/learningpath"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/recommendations"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
	"github.com/go-chi/chi/v5"
)

func setupTestRecommendations() (recommendations.Service, *recommendations.Handler, *chi.Mux, problems.Repository, progress.Repository, mistakes.Repository) {
	dsaMem := model.NewMemory()
	probRepo := problems.NewMemoryRepository(dsaMem)
	topicRepo := topics.NewMemoryRepository(dsaMem)
	progressRepo := progress.NewMemoryRepository(dsaMem)
	mistakesRepo := mistakes.NewMemoryRepository()
	memoryRepo := learningmemory.NewMemoryRepository()
	memorySvc := learningmemory.NewService(memoryRepo, probRepo, nil)
	pathRepo := learningpath.NewMemoryRepository()
	pathSvc := learningpath.NewService(pathRepo, memoryRepo, probRepo, topicRepo, nil)
	recRepo := recommendations.NewMemoryRepository()

	recSvc := recommendations.NewService(
		recRepo,
		probRepo,
		progressRepo,
		memorySvc,
		mistakesRepo,
		pathSvc,
		nil,
	)
	recHandler := recommendations.NewHandler(recSvc)

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			studentID := req.Header.Get("X-Student-ID")
			if studentID != "" {
				claims := &middleware.UserClaims{
					UserID: studentID,
					Email:  studentID + "@learnx.test",
				}
				ctx := context.WithValue(req.Context(), middleware.UserClaimsKey, claims)
				req = req.WithContext(ctx)
			}
			next.ServeHTTP(w, req)
		})
	})

	r.Route("/api/v1/student/recommendations", func(sub chi.Router) {
		sub.Get("/", recHandler.GetRecommendations)
		sub.Get("/next", recHandler.GetNext)
		sub.Get("/weaknesses", recHandler.GetWeaknesses)
		sub.Get("/similar/{problem_id}", recHandler.GetSimilar)
		sub.Post("/refresh", recHandler.Refresh)
		sub.Post("/{id}/feedback", recHandler.SubmitFeedback)
		sub.Get("/history", recHandler.GetHistory)
	})

	return recSvc, recHandler, r, probRepo, progressRepo, mistakesRepo
}

func TestRecommendations_NewStudent(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _, _ := setupTestRecommendations()
	studentID := "stu-new-001"

	res, err := svc.GetCategorizedRecommendations(ctx, studentID)
	if err != nil {
		t.Fatalf("expected no error for new student, got: %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil recommendations response")
	}

	if res.TotalRecommendations == 0 {
		t.Log("Note: empty catalog in test memory repository, total is 0")
	} else {
		if res.TopPick == nil {
			t.Errorf("expected top pick when recommendations exist")
		}
	}
}

func TestRecommendations_WithSeededProblems(t *testing.T) {
	ctx := context.Background()
	svc, _, _, probRepo, progressRepo, mistakesRepo := setupTestRecommendations()
	studentID := "stu-seed-002"

	// Seed problems in memory
	p1 := model.Problem{
		ID:         "prob-1",
		Title:      "Two Sum",
		Slug:       "two-sum",
		Difficulty: model.DifficultyEasy,
		TopicName:  "Arrays",
		Summary:    "Find two numbers that add up to target.",
		Tags:       []string{"Arrays", "Hash Map"},
		IsActive:   true,
	}
	p2 := model.Problem{
		ID:         "prob-2",
		Title:      "3Sum",
		Slug:       "3sum",
		Difficulty: model.DifficultyMedium,
		TopicName:  "Arrays",
		Summary:    "Find all unique triplets.",
		Tags:       []string{"Arrays", "Two Pointers"},
		IsActive:   true,
	}
	p3 := model.Problem{
		ID:         "prob-3",
		Title:      "Trapping Rain Water",
		Slug:       "trapping-rain-water",
		Difficulty: model.DifficultyHard,
		TopicName:  "Arrays",
		Summary:    "Compute trapped water.",
		Tags:       []string{"Two Pointers", "Stack"},
		IsActive:   true,
	}

	if mem, ok := probRepo.(*problems.MemoryRepository); ok {
		mem.DB.SetProblem(p1)
		mem.DB.SetProblem(p2)
		mem.DB.SetProblem(p3)
	}

	// Record a mistake in two pointers
	_ = mistakesRepo.SaveDiagnosis(ctx, &mistakes.Diagnosis{
		ID:              "diag-1",
		StudentID:       studentID,
		ProblemID:       "prob-2",
		Category:        "edge_case",
		Subtype:         "pointer_boundary_overflow",
		RelatedConcept:  "Arrays",
		ConfidenceLevel: "High",
		ConfidenceScore: 0.95,
		CreatedAt:       time.Now(),
	})

	_ = mistakesRepo.SavePattern(ctx, &mistakes.StudentMistakePattern{
		ID:                "pat-1",
		StudentID:         studentID,
		Category:          "edge_case",
		Subtype:           "pointer_boundary_overflow",
		RelatedConcept:    "Arrays",
		OccurrenceCount:   2,
		CorrectionStatus:  "active",
		FirstOccurrenceAt: time.Now(),
		LastOccurrenceAt:  time.Now(),
		UpdatedAt:         time.Now(),
	})

	// Mark p1 attempted
	_, _ = progressRepo.MarkAttempted(ctx, studentID, "prob-1")

	// Generate recommendations
	res, err := svc.GetCategorizedRecommendations(ctx, studentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.TotalRecommendations == 0 {
		t.Fatal("expected recommendations to be generated")
	}

	if res.TopPick == nil {
		t.Fatal("expected a TopPick problem")
	}

	// Verify categories are populated
	if len(res.Categories) == 0 {
		t.Fatal("expected categories map to have entries")
	}

	// Test GetNextRecommendation
	next, err := svc.GetNextRecommendation(ctx, studentID)
	if err != nil {
		t.Fatalf("expected next recommendation, got: %v", err)
	}
	if next == nil || next.ProblemID == "" {
		t.Errorf("expected valid problem in next recommendation")
	}

	// Test Feedback submission
	err = svc.SubmitFeedback(ctx, studentID, next.ID, recommendations.FeedbackRequest{
		Rating: recommendations.FeedbackHelpful,
		Notes:  "Very relevant to my current array practice!",
	})
	if err != nil {
		t.Errorf("expected feedback submission to succeed, got: %v", err)
	}

	// Test History
	history, err := svc.GetHistory(ctx, studentID, 10, 0)
	if err != nil {
		t.Fatalf("expected history, got: %v", err)
	}
	if history.Total == 0 {
		t.Errorf("expected at least 1 recommendation in history")
	}

	// Test Completion Hook
	svc.OnSubmissionEvaluated(ctx, studentID, next.ProblemID, "Accepted")

	// Test Refresh
	refreshed, err := svc.RefreshRecommendations(ctx, studentID)
	if err != nil {
		t.Fatalf("expected refresh to succeed, got: %v", err)
	}
	if refreshed == nil {
		t.Errorf("expected non-nil refreshed recommendations")
	}
}

func TestRecommendations_HTTPHandlers(t *testing.T) {
	_, _, router, probRepo, _, _ := setupTestRecommendations()
	studentID := "stu-http-test-001"

	p := model.Problem{
		ID:         "prob-h1",
		Title:      "Binary Search",
		Slug:       "binary-search",
		Difficulty: model.DifficultyEasy,
		TopicName:  "Binary Search",
		Summary:    "Search in sorted array.",
		Tags:       []string{"Binary Search"},
		IsActive:   true,
	}
	if mem, ok := probRepo.(*problems.MemoryRepository); ok {
		mem.DB.SetProblem(p)
	}

	// 1. GET /api/v1/student/recommendations
	req := httptest.NewRequest("GET", "/api/v1/student/recommendations", nil)
	req.Header.Set("X-Student-ID", studentID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d", w.Code)
	}

	var cat recommendations.CategorizedRecommendations
	if err := json.NewDecoder(w.Body).Decode(&cat); err != nil {
		t.Fatalf("failed decoding json: %v", err)
	}

	// 2. GET /api/v1/student/recommendations/next
	reqNext := httptest.NewRequest("GET", "/api/v1/student/recommendations/next", nil)
	reqNext.Header.Set("X-Student-ID", studentID)
	wNext := httptest.NewRecorder()
	router.ServeHTTP(wNext, reqNext)

	if wNext.Code != http.StatusOK {
		t.Fatalf("expected status 200 for next, got: %d", wNext.Code)
	}

	var nextRec recommendations.ProblemRecommendation
	_ = json.NewDecoder(wNext.Body).Decode(&nextRec)

	// 3. POST /api/v1/student/recommendations/{id}/feedback
	feedbackBody, _ := json.Marshal(recommendations.FeedbackRequest{
		Rating: "helpful",
		Notes:  "Great suggestion",
	})
	reqFeedback := httptest.NewRequest("POST", "/api/v1/student/recommendations/"+nextRec.ID+"/feedback", bytes.NewReader(feedbackBody))
	reqFeedback.Header.Set("X-Student-ID", studentID)
	wFeedback := httptest.NewRecorder()
	router.ServeHTTP(wFeedback, reqFeedback)

	if wFeedback.Code != http.StatusOK {
		t.Fatalf("expected status 200 for feedback, got: %d", wFeedback.Code)
	}

	// 4. GET /api/v1/student/recommendations/history
	reqHist := httptest.NewRequest("GET", "/api/v1/student/recommendations/history?limit=10", nil)
	reqHist.Header.Set("X-Student-ID", studentID)
	wHist := httptest.NewRecorder()
	router.ServeHTTP(wHist, reqHist)

	if wHist.Code != http.StatusOK {
		t.Fatalf("expected status 200 for history, got: %d", wHist.Code)
	}

	// 5. POST /api/v1/student/recommendations/refresh
	reqRefresh := httptest.NewRequest("POST", "/api/v1/student/recommendations/refresh", nil)
	reqRefresh.Header.Set("X-Student-ID", studentID)
	wRefresh := httptest.NewRecorder()
	router.ServeHTTP(wRefresh, reqRefresh)

	if wRefresh.Code != http.StatusOK {
		t.Fatalf("expected status 200 for refresh, got: %d", wRefresh.Code)
	}
}

