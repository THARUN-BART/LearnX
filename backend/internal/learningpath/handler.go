package learningpath

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
)

// Handler exposes the REST APIs for Adaptive Learning Path under /api/v1/student/learning-path.
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetLearningPath(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	path, err := h.svc.GetLearningPath(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch adaptive learning path"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(path)
}

func (h *Handler) GetNextAction(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rec, err := h.svc.GetNextAction(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to generate recommendation"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rec)
}

func (h *Handler) GetPathProgress(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	prog, err := h.svc.GetPathProgress(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch path progress"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(prog)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	history, total, err := h.svc.GetHistory(r.Context(), studentID, limit, offset)
	if err != nil {
		http.Error(w, `{"error":"failed to fetch recommendation history"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"history": history,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *Handler) RecalculatePath(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if strings.TrimSpace(studentID) == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}

	path, err := h.svc.RecalculatePath(r.Context(), studentID)
	if err != nil {
		http.Error(w, `{"error":"failed to recalculate learning path"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(path)
}
