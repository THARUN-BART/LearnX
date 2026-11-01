package progress

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/httputil"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	items, err := h.service.List(r.Context(), studentID)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"progress": items,
	})
}

func (h *Handler) GetProblem(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	problemID := chi.URLParam(r, "problem_id")
	item, err := h.service.Get(r.Context(), studentID, problemID)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"progress": item,
	})
}

func (h *Handler) Attempt(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	problemID := chi.URLParam(r, "problem_id")
	item, err := h.service.MarkAttempted(r.Context(), studentID, problemID)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"message":  "problem marked as attempted",
		"progress": item,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	studentID, ok := requireStudent(w, r)
	if !ok {
		return
	}
	problemID := chi.URLParam(r, "problem_id")
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Status == model.StatusSolved && !req.Reset {
		item, err := h.service.MarkSolved(r.Context(), studentID, problemID)
		if err != nil {
			writeErr(w, err)
			return
		}
		httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"success":  true,
			"message":  "problem marked as solved",
			"progress": item,
		})
		return
	}
	item, err := h.service.Update(r.Context(), studentID, problemID, req)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"message":  "progress updated",
		"progress": item,
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
		httputil.RespondError(w, http.StatusNotFound, "problem not found")
	case errors.Is(err, model.ErrInvalidInput):
		httputil.RespondError(w, http.StatusBadRequest, "invalid progress payload")
	default:
		httputil.RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
