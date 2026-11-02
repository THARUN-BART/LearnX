package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/aitutor"
	"github.com/THARUN-BART/LearnX/backend/internal/pedagogy"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/go-chi/chi/v5"
)

type mockPedagogyServer struct {
	server *httptest.Server
}

func newMockPedagogyServer() *mockPedagogyServer {
	m := &mockPedagogyServer{}
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(pedagogy.HealthResponse{
			Status:  "healthy",
			Service: "LearnX Pedagogy Diagnostic API",
		})
	})

	r.Post("/analyze", func(w http.ResponseWriter, r *http.Request) {
		var req pedagogy.AnalyzeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp := pedagogy.AnalyzeResponse{
			RunnerResult: pedagogy.RunnerResult{
				Passed:    false,
				ErrorType: "TIMEOUT",
				Traceback: "Code execution exceeded the time limit.",
			},
			Diagnosis: pedagogy.Diagnosis{
				Skill:           req.Skill,
				MisconceptionID: "MIS-04",
				Misconception:   "Infinite Loop",
				Confidence:      0.85,
				Evidence:        "Code execution exceeded the time limit.",
				DiagnosisSource: "rule_based",
			},
			Hint: pedagogy.Hint{
				HintLevel:  req.HintLevel,
				Hint:       "Your loop isn't terminating. Is your pointer moving forward with each step?",
				HintSource: "pedagogy_catalogue",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Post("/diagnose", func(w http.ResponseWriter, r *http.Request) {
		var req pedagogy.AnalyzeRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp := pedagogy.DiagnoseResponse{
			RunnerResult: pedagogy.RunnerResult{
				Passed:    false,
				ErrorType: "RUNTIME_ERROR",
			},
			Diagnosis: pedagogy.Diagnosis{
				Skill:           req.Skill,
				MisconceptionID: "MIS-02",
				Misconception:   "Dereferencing Null Pointer",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	m.server = httptest.NewServer(r)
	return m
}

func TestPedagogyRoutes_EndToEnd(t *testing.T) {
	mockPedagogy := newMockPedagogyServer()
	defer mockPedagogy.server.Close()

	client := pedagogy.NewClient(mockPedagogy.server.URL, 2*time.Second, nil)
	pedagogyHandler := pedagogy.NewHandler(client)

	intelHandlers := &server.IntelligenceHandlers{
		Pedagogy: pedagogyHandler,
	}

	router := server.BuildRouter(
		nil,
		nil,
		nil,
		nil,
		nil,
		server.WithIntelligence(intelHandlers),
	)

	// 1. Test GET /api/v1/ai/pedagogy/health
	req := httptest.NewRequest(http.MethodGet, "/api/v1/ai/pedagogy/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from Pedagogy health, got %d: %s", rec.Code, rec.Body.String())
	}

	var healthRes struct {
		Success bool                    `json:"success"`
		Data    pedagogy.HealthResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &healthRes); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}
	if healthRes.Data.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", healthRes.Data.Status)
	}

	// 2. Test POST /api/v1/ai/pedagogy/analyze
	analyzeReqBody := pedagogy.AnalyzeRequest{
		StudentCode: "while True:\n    pass",
		TestCode:    "print('Test')",
		Skill:       "SKL-02",
		HintLevel:   "L1",
	}
	jsonBytes, _ := json.Marshal(analyzeReqBody)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/ai/pedagogy/analyze", bytes.NewReader(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from Pedagogy analyze, got %d: %s", rec.Code, rec.Body.String())
	}

	var analyzeRes struct {
		Success bool                     `json:"success"`
		Data    pedagogy.AnalyzeResponse `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &analyzeRes); err != nil {
		t.Fatalf("failed to decode analyze response: %v", err)
	}

	if analyzeRes.Data.Diagnosis.MisconceptionID != "MIS-04" {
		t.Errorf("expected misconception MIS-04, got %s", analyzeRes.Data.Diagnosis.MisconceptionID)
	}

	// 3. Test Backwards-compatible alias /api/v1/ai/person-c/health
	reqAlias := httptest.NewRequest(http.MethodGet, "/api/v1/ai/person-c/health", nil)
	recAlias := httptest.NewRecorder()
	router.ServeHTTP(recAlias, reqAlias)

	if recAlias.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from backwards-compatible alias, got %d", recAlias.Code)
	}
}

func TestPedagogy_ProgressiveHintDelegation(t *testing.T) {
	mockPedagogy := newMockPedagogyServer()
	defer mockPedagogy.server.Close()

	client := pedagogy.NewClient(mockPedagogy.server.URL, 2*time.Second, nil)

	repo := aitutor.NewMemoryRepository()
	ctxBuilder := aitutor.NewContextBuilder(nil, nil, nil)
	aiSvc := aitutor.NewService(nil, repo, nil, ctxBuilder, nil)
	aiSvc.SetPedagogyClient(client)

	ctx := context.Background()
	studentID := "test-student-123"

	_, err := aiSvc.GetProgressiveHint(ctx, studentID, aitutor.HintRequest{
		ProblemID: "non-existent-problem",
		Level:     1,
	})
	if err != aitutor.ErrProblemNotFound {
		t.Errorf("expected ErrProblemNotFound, got %v", err)
	}
}
