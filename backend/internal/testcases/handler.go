package testcases

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetSampleTestCases handles GET /api/v1/problems/{problem_id}/testcases/sample
func (h *Handler) GetSampleTestCases(w http.ResponseWriter, r *http.Request) {
	problemID := chi.URLParam(r, "problem_id")
	samples, err := h.service.GetSampleTestCases(r.Context(), problemID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    samples,
	})
}

// CreateTestCase handles POST /api/v1/admin/problems/{problem_id}/testcases
func (h *Handler) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	problemID := chi.URLParam(r, "problem_id")
	var req CreateTestCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tc, err := h.service.CreateTestCase(r.Context(), problemID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    tc,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
