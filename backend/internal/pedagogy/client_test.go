package pedagogy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPedagogyClient_Analyze(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/analyze" {
			http.NotFound(w, r)
			return
		}

		var req AnalyzeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := AnalyzeResponse{
			RunnerResult: RunnerResult{
				Passed:    false,
				ErrorType: "TIMEOUT",
				Traceback: "Code execution exceeded the time limit.",
			},
			Diagnosis: Diagnosis{
				Skill:           req.Skill,
				MisconceptionID: "MIS-04",
				Misconception:   "Infinite Loop",
				Confidence:      0.85,
				Evidence:        "Code execution exceeded the time limit.",
				DiagnosisSource: "rule_based",
			},
			Hint: Hint{
				HintLevel:  req.HintLevel,
				Hint:       "Your loop isn't terminating. Is your pointer moving forward with each step?",
				HintSource: "person_a_catalogue",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, 2*time.Second, nil)
	ctx := context.Background()

	res, err := client.Analyze(ctx, AnalyzeRequest{
		StudentCode: "while True: pass",
		TestCode:    "assert True",
		Skill:       "SKL-02",
		HintLevel:   "L1",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if res.Diagnosis.MisconceptionID != "MIS-04" {
		t.Errorf("expected MIS-04, got %s", res.Diagnosis.MisconceptionID)
	}

	if res.Hint.HintSource != "person_a_catalogue" {
		t.Errorf("expected person_a_catalogue, got %s", res.Hint.HintSource)
	}
}

func TestPedagogyClient_Diagnose(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/diagnose" {
			http.NotFound(w, r)
			return
		}

		resp := DiagnoseResponse{
			RunnerResult: RunnerResult{
				Passed:    false,
				ErrorType: "RUNTIME_ERROR",
			},
			Diagnosis: Diagnosis{
				Skill:           "SKL-07",
				MisconceptionID: "MIS-02",
				Misconception:   "Dereferencing Null Pointer",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, 2*time.Second, nil)
	ctx := context.Background()

	res, err := client.Diagnose(ctx, AnalyzeRequest{
		StudentCode: "curr = None; print(curr.next)",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Diagnosis.MisconceptionID != "MIS-02" {
		t.Errorf("expected MIS-02, got %s", res.Diagnosis.MisconceptionID)
	}
}

func TestPedagogyClient_Health(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:  "healthy",
			Service: "LearnX Pedagogy Diagnostic API",
		})
	}))
	defer mockServer.Close()

	client := NewClient(mockServer.URL, 2*time.Second, nil)
	health, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("expected healthy, got error: %v", err)
	}
	if health.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", health.Status)
	}
}
