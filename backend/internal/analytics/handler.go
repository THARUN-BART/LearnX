package analytics

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/httputil"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	item, err := h.service.Overview(r.Context(), studentID, r.URL.Query().Get("roadmap_id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"overview": item,
	})
}

func (h *Handler) Topics(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	items, err := h.service.Topics(r.Context(), studentID, r.URL.Query().Get("roadmap_id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"topics":  items,
	})
}

func (h *Handler) Difficulty(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	items, err := h.service.Difficulty(r.Context(), studentID, r.URL.Query().Get("roadmap_id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"difficulty": items,
	})
}

func (h *Handler) Activity(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	days := 30
	if raw := r.URL.Query().Get("days"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			days = parsed
		}
	}
	item, err := h.service.Activity(r.Context(), studentID, days)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"activity": item,
	})
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	item, err := h.service.Summary(r.Context(), studentID, r.URL.Query().Get("roadmap_id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"summary": item,
	})
}

func requireStudent(w http.ResponseWriter, r *http.Request) (string, bool) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		httputil.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return "", false
	}
	return studentID, true
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		httputil.RespondError(w, http.StatusNotFound, "roadmap not found")
	default:
		httputil.RespondError(w, http.StatusInternalServerError, "failed to compute analytics")
	}
}
