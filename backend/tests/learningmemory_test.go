package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/learningmemory"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/go-chi/chi/v5"
)

func setupTestLearningMemoryService() (learningmemory.Service, *learningmemory.Handler) {
	memRepo := learningmemory.NewMemoryRepository()
	dsaMem := model.NewMemory()
	probRepo := problems.NewMemoryRepository(dsaMem)
	svc := learningmemory.NewService(memRepo, probRepo, nil)
	handler := learningmemory.NewHandler(svc)
	return svc, handler
}

func TestLearningMemory_IndependentVsAssistedSolves(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupTestLearningMemoryService()
	studentID := "stu-test-memory-1"

	prob := &model.Problem{
		ID:         "prob-101",
		Title:      "Two Sum",
		Slug:       "two-sum",
		TopicID:    "topic-arrays",
		SubtopicID: "subtopic-1",
		Tags:       []string{"hash_table", "arrays"},
	}

	// 1. Solve independently (no prior hints requested)
	err := svc.RecordSubmissionEvent(ctx, studentID, prob, "sub-1", true, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum, err := svc.GetSummary(ctx, studentID)
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}
	if sum.IndependentSolves != 1 {
		t.Errorf("expected 1 independent solve, got %d", sum.IndependentSolves)
	}
	if sum.AssistedSolves != 0 {
		t.Errorf("expected 0 assisted solves, got %d", sum.AssistedSolves)
	}

	// 2. Request a hint for prob-2
	prob2 := &model.Problem{
		ID:         "prob-102",
		Title:      "Binary Search",
		Slug:       "binary-search",
		TopicID:    "topic-bs",
		SubtopicID: "subtopic-2",
		Tags:       []string{"binary_search"},
	}

	err = svc.RecordHintEvent(ctx, studentID, prob2.ID, 1, "Consider using low and high pointers.")
	if err != nil {
		t.Fatalf("failed to record hint event: %v", err)
	}

	// Solve prob-2 after receiving hint -> Should be classified as assisted solve!
	err = svc.RecordSubmissionEvent(ctx, studentID, prob2, "sub-2", true, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum2, err := svc.GetSummary(ctx, studentID)
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}
	if sum2.IndependentSolves != 1 {
		t.Errorf("expected 1 independent solve, got %d", sum2.IndependentSolves)
	}
	if sum2.AssistedSolves != 1 {
		t.Errorf("expected 1 assisted solve, got %d", sum2.AssistedSolves)
	}
}

func TestLearningMemory_MistakeTrackingAndCorrection(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupTestLearningMemoryService()
	studentID := "stu-test-memory-2"

	conceptKey := "binary_search"

	// 1. Record 2 recurring mistakes
	_ = svc.RecordMistakeEvent(ctx, studentID, "prob-bs-1", "sub-1", "Boundary_Error", "off_by_one", conceptKey, 0.95)
	_ = svc.RecordMistakeEvent(ctx, studentID, "prob-bs-1", "sub-2", "Boundary_Error", "off_by_one", conceptKey, 0.95)

	weaks, err := svc.GetWeaknesses(ctx, studentID, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(weaks) != 1 {
		t.Fatalf("expected 1 active weakness, got %d", len(weaks))
	}
	if weaks[0].OccurrenceCount != 2 {
		t.Errorf("expected occurrence count 2, got %d", weaks[0].OccurrenceCount)
	}
	if weaks[0].Status != "active" {
		t.Errorf("expected status 'active', got %s", weaks[0].Status)
	}

	// 2. Submit accepted solve on that concept -> evidence of correction
	prob := &model.Problem{
		ID:         "prob-bs-1",
		Title:      "Binary Search",
		Slug:       "binary-search",
		Tags:       []string{"binary_search"},
		TopicID:    "topic-bs",
		SubtopicID: "sub-1",
	}
	_ = svc.RecordSubmissionEvent(ctx, studentID, prob, "sub-accepted-1", true, "")

	detail, err := svc.GetConceptDetail(ctx, studentID, conceptKey)
	if err != nil {
		t.Fatalf("failed to get concept detail: %v", err)
	}
	if detail.Concept.CorrectedMistakes < 1 {
		t.Errorf("expected at least 1 corrected mistake, got %d", detail.Concept.CorrectedMistakes)
	}
}

func TestLearningMemory_HTTPRoutesAndAuthorization(t *testing.T) {
	_, handler := setupTestLearningMemoryService()

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Mock authenticated student claims context
			claims := &middleware.UserClaims{
				UserID: "stu-authenticated-1",
				Email:  "student@learnx.com",
			}
			ctx := context.WithValue(r.Context(), middleware.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/api/v1/student/learning-memory", handler.GetSummary)
	r.Get("/api/v1/student/learning-memory/concepts", handler.ListConcepts)
	r.Get("/api/v1/student/learning-memory/concepts/{concept_key}", handler.GetConceptDetail)
	r.Get("/api/v1/student/learning-memory/timeline", handler.GetTimeline)

	// Test GET /api/v1/student/learning-memory
	req := httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-memory", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var sum learningmemory.MemorySummaryResponse
	if err := json.Unmarshal(w.Body.Bytes(), &sum); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	// Test GET /api/v1/student/learning-memory/timeline
	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-memory/timeline?days=7", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected status 200 for timeline, got %d", w2.Code)
	}
}

func TestLearningMemory_AITutorContextGrounding(t *testing.T) {
	ctx := context.Background()
	svc, _ := setupTestLearningMemoryService()
	studentID := "stu-grounding-1"

	// Record an independent solve for Recursion
	probRec := &model.Problem{
		ID:   "prob-rec-1",
		Slug: "recursion-subsets",
		Tags: []string{"recursion"},
	}
	_ = svc.RecordSubmissionEvent(ctx, studentID, probRec, "sub-1", true, "")

	// Record a recurring mistake for Dynamic Programming
	_ = svc.RecordMistakeEvent(ctx, studentID, "prob-dp-1", "sub-2", "State_Transition", "base_case_omitted", "dynamic_programming", 0.9)

	aiCtx, err := svc.BuildAIContext(ctx, studentID, "topic-dp", "prob-dp-1")
	if err != nil {
		t.Fatalf("unexpected error building AI context: %v", err)
	}

	if len(aiCtx) == 0 {
		t.Fatalf("expected non-empty AI context string")
	}

	// Validate evidence grounding
	if !containsStr(aiCtx, "STUDENT LEARNING MEMORY CONTEXT") {
		t.Errorf("missing header in AI context: %s", aiCtx)
	}
	if !containsStr(aiCtx, "dynamic_programming") && !containsStr(aiCtx, "Dynamic Programming") {
		t.Errorf("expected dynamic programming weakness mentioned in AI context: %s", aiCtx)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || time.Now().Unix() > 0 && len(sub) > 0 && containsSub(s, sub))
}

func containsSub(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
