package learningmemory

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

// Handler exposes the REST APIs for Student Learning Memory under /api/v1/student/learning-memory.
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	summary, err := h.svc.GetSummary(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch learning summary"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}

func (h *Handler) ListConcepts(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	filter := ConceptFilter{
		TopicID:       q.Get("topic_id"),
		LearningState: q.Get("learning_state"),
		Search:        q.Get("search"),
		Limit:         limit,
		Offset:        offset,
	}

	concepts, total, err := h.svc.ListConcepts(r.Context(), studentID, filter)
	if err != nil {
		http.Error(w, `{"error":"failed to list concepts"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"concepts": concepts,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *Handler) GetConceptDetail(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	conceptKey := chi.URLParam(r, "concept_key")
	if strings.TrimSpace(conceptKey) == "" {
		http.Error(w, `{"error":"concept_key is required"}`, http.StatusBadRequest)
		return
	}

	detail, err := h.svc.GetConceptDetail(r.Context(), studentID, conceptKey)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch concept details"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(detail)
}

func (h *Handler) GetStrengths(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 5
	}

	strengths, err := h.svc.GetStrengths(r.Context(), studentID, limit)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch strengths"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"strengths": strengths,
	})
}

func (h *Handler) GetWeaknesses(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 5
	}

	weaknesses, err := h.svc.GetWeaknesses(r.Context(), studentID, limit)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch weaknesses"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"weaknesses": weaknesses,
	})
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	filter := EventFilter{
		EventType:  q.Get("event_type"),
		ConceptKey: q.Get("concept_key"),
		ProblemID:  q.Get("problem_id"),
		Limit:      limit,
		Offset:     offset,
	}

	if startStr := q.Get("start_date"); startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			filter.StartDate = &t
		}
	}
	if endStr := q.Get("end_date"); endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			filter.EndDate = &t
		}
	}

	events, total, err := h.svc.GetHistory(r.Context(), studentID, filter)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch learning history"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *Handler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 || days > 90 {
		days = 14
	}

	timeline, err := h.svc.GetTimeline(r.Context(), studentID, days)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch timeline"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"timeline": timeline,
	})
}

func (h *Handler) GetSnapshots(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days <= 0 {
		days = 30
	}

	snapshots, err := h.svc.GetSnapshots(r.Context(), studentID, days)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch snapshots"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"snapshots": snapshots,
	})
}

func (h *Handler) GetStudentFacingSummary(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	summary, err := h.svc.GetStudentFacingSummary(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to generate student summary"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(summary)
}
