package recommendations

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// Handler exposes REST endpoints for Personalised Problem Recommendations under /api/v1/student/recommendations.
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GET /api/v1/student/recommendations
func (h *Handler) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	res, err := h.svc.GetCategorizedRecommendations(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to generate problem recommendations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

// GET /api/v1/student/recommendations/next
func (h *Handler) GetNext(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rec, err := h.svc.GetNextRecommendation(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"no next recommendation available"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rec)
}

// GET /api/v1/student/recommendations/weaknesses
func (h *Handler) GetWeaknesses(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	weaknesses, err := h.svc.GetWeaknessRecommendations(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch weakness recommendations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"weakness_recommendations": weaknesses,
		"count":                    len(weaknesses),
	})
}

// GET /api/v1/student/recommendations/similar/{problem_id}
func (h *Handler) GetSimilar(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	problemID := chi.URLParam(r, "problem_id")
	if strings.TrimSpace(problemID) == "" {
		http.Error(w, `{"error":"problem_id is required"}`, http.StatusBadRequest)
		return
	}

	similar, err := h.svc.GetSimilarRecommendations(r.Context(), studentID, problemID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch similar problem recommendations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"similar_problems": similar,
		"count":            len(similar),
	})
}

// POST /api/v1/student/recommendations/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	fresh, err := h.svc.RefreshRecommendations(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to refresh problem recommendations"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(fresh)
}

// POST /api/v1/student/recommendations/{id}/feedback
func (h *Handler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	recID := chi.URLParam(r, "id")
	if strings.TrimSpace(recID) == "" {
		http.Error(w, `{"error":"recommendation id is required"}`, http.StatusBadRequest)
		return
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if err := h.svc.SubmitFeedback(r.Context(), studentID, recID, req); err != nil {
		http.Error(w, `{"error":"failed to save recommendation feedback"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Feedback recorded successfully",
	})
}

// GET /api/v1/student/recommendations/history
func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	history, err := h.svc.GetHistory(r.Context(), studentID, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch recommendation history"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(history)
}

