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
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
	"github.com/go-chi/chi/v5"
)

func setupTestLearningPath() (learningpath.Service, *learningpath.Handler, learningmemory.Repository, *chi.Mux) {
	dsaMem := model.NewMemory()
	probRepo := problems.NewMemoryRepository(dsaMem)
	topicRepo := topics.NewMemoryRepository(dsaMem)
	memRepo := learningmemory.NewMemoryRepository()
	pathRepo := learningpath.NewMemoryRepository()

	pathSvc := learningpath.NewService(pathRepo, memRepo, probRepo, topicRepo, nil)
	pathHandler := learningpath.NewHandler(pathSvc)

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

	r.Route("/api/v1/student/learning-path", func(sub chi.Router) {
		sub.Get("/", pathHandler.GetLearningPath)
		sub.Post("/generate", pathHandler.RecalculatePath)
		sub.Get("/next-action", pathHandler.GetNextAction)
		sub.Get("/progress", pathHandler.GetPathProgress)
		sub.Get("/history", pathHandler.GetHistory)
		sub.Post("/recalculate", pathHandler.RecalculatePath)
	})

	return pathSvc, pathHandler, memRepo, r
}

func TestAdaptiveLearningPath_NewStudent(t *testing.T) {
	ctx := context.Background()
	pathSvc, _, _, _ := setupTestLearningPath()
	studentID := "stu-new-student-1"

	// New student has no history, should receive initial path
	path, err := pathSvc.GetLearningPath(ctx, studentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path == nil {
		t.Fatal("expected non-nil path")
	}

	if path.StudentID != studentID {
		t.Errorf("expected student ID %s, got %s", studentID, path.StudentID)
	}

	if len(path.Nodes) == 0 {
		t.Fatal("expected canonical curriculum nodes in path")
	}

	// First node (basics) must be ready_to_learn
	firstNode := path.Nodes[0]
	if firstNode.ConceptKey != "basics" {
		t.Errorf("expected first node to be basics, got %s", firstNode.ConceptKey)
	}
	if firstNode.State != learningpath.StateReadyToLearn {
		t.Errorf("expected basics to be ready_to_learn, got %s", firstNode.State)
	}

	// Dependent node (arrays or binary_search) must be locked_by_prerequisite
	foundLocked := false
	for _, n := range path.Nodes {
		if n.ConceptKey == "arrays" || n.ConceptKey == "binary_search" {
			if n.State == learningpath.StateLockedByPrerequisite {
				foundLocked = true
			}
		}
	}
	if !foundLocked {
		t.Errorf("expected dependent nodes to be locked_by_prerequisite")
	}

	// Next Best Action must recommend basics
	nextAction, err := pathSvc.GetNextAction(ctx, studentID)
	if err != nil {
		t.Fatalf("failed to get next action: %v", err)
	}
	if nextAction.ConceptKey != "basics" {
		t.Errorf("expected next action concept 'basics', got '%s'", nextAction.ConceptKey)
	}
	if nextAction.ActionType != learningpath.ActionLearnConcept {
		t.Errorf("expected action learn_concept, got %s", nextAction.ActionType)
	}
	if nextAction.RecommendedDifficulty != model.DifficultyEasy {
		t.Errorf("expected easy difficulty for new student, got %s", nextAction.RecommendedDifficulty)
	}
	if nextAction.Explanation == "" {
		t.Errorf("expected non-empty explanation")
	}
}

func TestAdaptiveLearningPath_PrerequisiteProgression(t *testing.T) {
	ctx := context.Background()
	pathSvc, _, memRepo, _ := setupTestLearningPath()
	studentID := "stu-progressing-student"

	now := time.Now().UTC()
	// 1. Simulate student mastering basics and sorting in learning memory
	err := memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:          studentID,
		ConceptKey:         "basics",
		ConceptName:        "Language & Logic Fundamentals",
		EstimatedMastery:   0.85,
		ConfidenceLevel:    "high",
		TotalAttempts:      10,
		IndependentSolves:  9,
		LearningState:      "mastered",
		FirstPracticedAt:   now,
		LastPracticedAt:    now,
	})
	if err != nil {
		t.Fatalf("failed to seed basics: %v", err)
	}

	err = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:          studentID,
		ConceptKey:         "sorting",
		ConceptName:        "Sorting Techniques",
		EstimatedMastery:   0.80,
		ConfidenceLevel:    "high",
		TotalAttempts:      8,
		IndependentSolves:  7,
		LearningState:      "mastered",
		FirstPracticedAt:   now,
		LastPracticedAt:    now,
	})
	if err != nil {
		t.Fatalf("failed to seed sorting: %v", err)
	}

	// Recalculate path
	path, err := pathSvc.RecalculatePath(ctx, studentID)
	if err != nil {
		t.Fatalf("recalculate failed: %v", err)
	}

	// Check basics and sorting are marked mastered
	var arraysNode *learningpath.PathNode
	for i := range path.Nodes {
		if path.Nodes[i].ConceptKey == "basics" {
			if path.Nodes[i].State != learningpath.StateMastered {
				t.Errorf("expected basics to be mastered, got %s", path.Nodes[i].State)
			}
		}
		if path.Nodes[i].ConceptKey == "sorting" {
			if path.Nodes[i].State != learningpath.StateMastered {
				t.Errorf("expected sorting to be mastered, got %s", path.Nodes[i].State)
			}
		}
		if path.Nodes[i].ConceptKey == "arrays" {
			arraysNode = &path.Nodes[i]
		}
	}

	if arraysNode == nil {
		t.Fatal("arrays node not found")
	}
	// With sorting mastered, arrays should be ready_to_learn
	if arraysNode.State != learningpath.StateReadyToLearn {
		t.Errorf("expected arrays to be ready_to_learn now, got %s", arraysNode.State)
	}

	// Next action should recommend arrays
	nextAction, err := pathSvc.GetNextAction(ctx, studentID)
	if err != nil {
		t.Fatalf("failed to get next action: %v", err)
	}
	if nextAction.ConceptKey != "arrays" {
		t.Errorf("expected next action arrays, got %s", nextAction.ConceptKey)
	}
}

func TestAdaptiveLearningPath_WeaknessAndRevision(t *testing.T) {
	ctx := context.Background()
	pathSvc, _, memRepo, _ := setupTestLearningPath()
	studentID := "stu-struggling-student"

	now := time.Now().UTC()
	// Student has mastered basics & sorting & arrays, but struggles with binary_search
	_ = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:        studentID,
		ConceptKey:       "basics",
		EstimatedMastery: 0.9,
		ConfidenceLevel:  "high",
		LearningState:    "mastered",
		FirstPracticedAt: now,
		LastPracticedAt:  now,
	})
	_ = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:        studentID,
		ConceptKey:       "sorting",
		EstimatedMastery: 0.85,
		ConfidenceLevel:  "high",
		LearningState:    "mastered",
		FirstPracticedAt: now,
		LastPracticedAt:  now,
	})
	_ = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:        studentID,
		ConceptKey:       "arrays",
		EstimatedMastery: 0.80,
		ConfidenceLevel:  "high",
		LearningState:    "mastered",
		FirstPracticedAt: now,
		LastPracticedAt:  now,
	})
	// Low mastery + active recurring mistake on binary_search
	_ = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:         studentID,
		ConceptKey:        "binary_search",
		EstimatedMastery:  0.25,
		ConfidenceLevel:   "low",
		TotalAttempts:     7,
		IndependentSolves: 2,
		LearningState:     "struggling",
		FirstPracticedAt:  now,
		LastPracticedAt:   now,
	})
	_ = memRepo.UpsertMistakeMemory(ctx, &learningmemory.MistakeMemory{
		StudentID:         studentID,
		ConceptKey:        "binary_search",
		Category:          "boundary_condition",
		Subtype:           "off_by_one_index",
		OccurrenceCount:   4,
		Status:            "active",
		FirstOccurrenceAt: now,
		LastOccurrenceAt:  now,
	})

	path, err := pathSvc.RecalculatePath(ctx, studentID)
	if err != nil {
		t.Fatalf("recalculate failed: %v", err)
	}

	// Verify binary_search node is marked needs_revision
	var bsNode *learningpath.PathNode
	for i := range path.Nodes {
		if path.Nodes[i].ConceptKey == "binary_search" {
			bsNode = &path.Nodes[i]
			break
		}
	}
	if bsNode == nil {
		t.Fatal("binary_search node not found")
	}
	if bsNode.State != learningpath.StateNeedsRevision {
		t.Errorf("expected binary_search to be needs_revision, got %s", bsNode.State)
	}

	// Next Action must be Revise Concept or Review Mistake
	nextAction, err := pathSvc.GetNextAction(ctx, studentID)
	if err != nil {
		t.Fatalf("failed to get next action: %v", err)
	}
	if nextAction.ConceptKey != "binary_search" {
		t.Errorf("expected revision for binary_search, got %s", nextAction.ConceptKey)
	}
	if nextAction.ActionType != learningpath.ActionReviseConcept && nextAction.ActionType != learningpath.ActionReviewMistake {
		t.Errorf("expected revise action, got %s", nextAction.ActionType)
	}
	if nextAction.RelevantWeakness == "" {
		t.Errorf("expected relevant weakness to be populated")
	}
}

func TestAdaptiveLearningPath_RESTEndpoints(t *testing.T) {
	_, _, _, router := setupTestLearningPath()
	studentID := "stu-api-user-99"

	// 1. GET /api/v1/student/learning-path
	req := httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-path", nil)
	req.Header.Set("X-Student-ID", studentID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET learning-path failed: status %d, body: %s", w.Code, w.Body.String())
	}
	var pathData learningpath.AdaptivePath
	if err := json.Unmarshal(w.Body.Bytes(), &pathData); err != nil {
		t.Fatalf("failed to parse get response: %v", err)
	}
	if len(pathData.Nodes) == 0 {
		t.Errorf("expected non-empty path nodes")
	}

	// 2. GET /api/v1/student/learning-path/next-action
	req = httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-path/next-action", nil)
	req.Header.Set("X-Student-ID", studentID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET next-action failed: status %d", w.Code)
	}
	var actionData learningpath.LearningRecommendation
	if err := json.Unmarshal(w.Body.Bytes(), &actionData); err != nil {
		t.Fatalf("failed to parse action response: %v", err)
	}
	if actionData.ConceptKey != "basics" {
		t.Errorf("expected basics, got %s", actionData.ConceptKey)
	}

	// 3. GET /api/v1/student/learning-path/progress
	req = httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-path/progress", nil)
	req.Header.Set("X-Student-ID", studentID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET progress failed: status %d", w.Code)
	}
	var progData learningpath.PathProgressResponse
	if err := json.Unmarshal(w.Body.Bytes(), &progData); err != nil {
		t.Fatalf("failed to parse progress response: %v", err)
	}
	if progData.TotalConcepts == 0 {
		t.Errorf("expected > 0 total concepts")
	}

	// 4. POST /api/v1/student/learning-path/recalculate
	req = httptest.NewRequest(http.MethodPost, "/api/v1/student/learning-path/recalculate", bytes.NewBuffer([]byte("{}")))
	req.Header.Set("X-Student-ID", studentID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST recalculate failed: status %d", w.Code)
	}

	// 5. GET /api/v1/student/learning-path/history
	req = httptest.NewRequest(http.MethodGet, "/api/v1/student/learning-path/history?limit=10&offset=0", nil)
	req.Header.Set("X-Student-ID", studentID)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("GET history failed: status %d", w.Code)
	}
	var histResp struct {
		History []learningpath.LearningRecommendation `json:"history"`
		Total   int                                   `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &histResp); err != nil {
		t.Fatalf("failed to parse history response: %v", err)
	}
	if histResp.Total == 0 {
		t.Errorf("expected at least 1 recommendation in history, got 0")
	}
}

func TestAdaptiveLearningPath_DataIsolation(t *testing.T) {
	ctx := context.Background()
	pathSvc, _, memRepo, _ := setupTestLearningPath()

	studentA := "student-alpha"
	studentB := "student-beta"

	now := time.Now().UTC()
	// Student A masters basics
	_ = memRepo.UpsertConceptMemory(ctx, &learningmemory.ConceptMemory{
		StudentID:        studentA,
		ConceptKey:       "basics",
		EstimatedMastery: 0.95,
		ConfidenceLevel:  "high",
		LearningState:    "mastered",
		FirstPracticedAt: now,
		LastPracticedAt:  now,
	})
	pathA, err := pathSvc.RecalculatePath(ctx, studentA)
	if err != nil {
		t.Fatalf("student A path error: %v", err)
	}

	// Student B has zero progress
	pathB, err := pathSvc.GetLearningPath(ctx, studentB)
	if err != nil {
		t.Fatalf("student B path error: %v", err)
	}

	if pathA.Nodes[0].State != learningpath.StateMastered {
		t.Errorf("expected student A basics to be mastered")
	}
	if pathB.Nodes[0].State != learningpath.StateReadyToLearn {
		t.Errorf("expected student B basics to be ready_to_learn, got %s", pathB.Nodes[0].State)
	}
}
