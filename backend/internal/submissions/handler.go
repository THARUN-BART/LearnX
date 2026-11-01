package submissions

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Submit handles POST /api/v1/submissions
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid student session required")
		return
	}

	var req CreateSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.Submit(r.Context(), studentID, req)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// GetStatus handles GET /api/v1/submissions/{submission_id}
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	res, err := h.service.GetStatus(r.Context(), studentID, submissionID)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// GetResult handles GET /api/v1/submissions/{submission_id}/result
func (h *Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	res, err := h.service.GetResult(r.Context(), studentID, submissionID)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// ListStudentSubmissions handles GET /api/v1/student/submissions
func (h *Handler) ListStudentSubmissions(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	problemID := r.URL.Query().Get("problem_id")
	status := r.URL.Query().Get("status")

	res, err := h.service.ListStudentSubmissions(r.Context(), studentID, problemID, status, page, limit)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// ListProblemSubmissions handles GET /api/v1/problems/{problem_id}/submissions
func (h *Handler) ListProblemSubmissions(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID := chi.URLParam(r, "problem_id")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	res, err := h.service.GetProblemSubmissions(r.Context(), studentID, problemID, page, limit)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// GetLatestProblemSubmission handles GET /api/v1/problems/{problem_id}/submissions/latest
func (h *Handler) GetLatestProblemSubmission(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	problemID := chi.URLParam(r, "problem_id")
	res, err := h.service.GetLatestProblemSubmission(r.Context(), studentID, problemID)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

// Retry handles POST /api/v1/submissions/{submission_id}/retry
func (h *Handler) Retry(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	res, err := h.service.Retry(r.Context(), studentID, submissionID)
	if err != nil {
		mapError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		respondError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, ErrSubmissionNotFound), errors.Is(err, ErrProblemNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidProblemID), errors.Is(err, ErrEmptySourceCode),
		errors.Is(err, ErrCodeTooLarge), errors.Is(err, ErrInvalidLanguage):
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, err.Error())
	}
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
