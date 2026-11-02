package mistakes

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

// GetSubmissionDiagnosis handles GET /api/v1/submissions/{submission_id}/diagnosis
func (h *Handler) GetSubmissionDiagnosis(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	submissionID := chi.URLParam(r, "submission_id")
	if submissionID == "" {
		respondError(w, http.StatusBadRequest, "submission_id is required")
		return
	}

	diag, err := h.service.GetDiagnosisBySubmission(r.Context(), studentID, submissionID)
	if err != nil {
		if errors.Is(err, ErrDiagnosisNotFound) {
			respondError(w, http.StatusNotFound, "no diagnosis found for this submission (it may still be running, accepted, or unanalyzed)")
			return
		}
		if errors.Is(err, ErrUnauthorizedAccess) {
			respondError(w, http.StatusForbidden, "you do not have permission to view this diagnosis")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    diag,
	})
}

// GetStudentMistakes handles GET /api/v1/student/mistakes
func (h *Handler) GetStudentMistakes(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 20
	}

	diagnoses, err := h.service.GetStudentMistakes(r.Context(), studentID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    diagnoses,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetStudentMistakeSummary handles GET /api/v1/student/mistakes/summary
func (h *Handler) GetStudentMistakeSummary(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	summary, err := h.service.GetStudentMistakeSummary(r.Context(), studentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    summary,
	})
}

// GetStudentPatterns handles GET /api/v1/student/mistakes/patterns
func (h *Handler) GetStudentPatterns(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	patterns, err := h.service.GetStudentPatterns(r.Context(), studentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    patterns,
	})
}

type FeedbackRequest struct {
	FeedbackType    string `json:"feedback_type"`
	StudentFeedback string `json:"student_feedback"`
}

// SubmitFeedback handles POST /api/v1/student/mistakes/{diagnosis_id}/feedback
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	diagnosisID := chi.URLParam(r, "diagnosis_id")
	if diagnosisID == "" {
		respondError(w, http.StatusBadRequest, "diagnosis_id is required")
		return
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.FeedbackType == "" {
		respondError(w, http.StatusBadRequest, "feedback_type is required")
		return
	}

	if err := h.service.RecordFeedback(r.Context(), studentID, diagnosisID, req.FeedbackType, req.StudentFeedback); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "feedback recorded successfully",
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
