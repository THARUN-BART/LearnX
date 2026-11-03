package visualizer

import (
	"encoding/json"
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

func (h *Handler) ListAlgorithms(w http.ResponseWriter, r *http.Request) {
	algos := h.service.ListAlgorithms()
	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"algorithms": algos,
		"count":      len(algos),
	})
}

func (h *Handler) GetAlgorithm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	algo, err := h.service.GetAlgorithm(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"algorithm": algo,
	})
}

func (h *Handler) Execute(w http.ResponseWriter, r *http.Request) {
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AlgorithmID == "" {
		writeError(w, http.StatusBadRequest, "algorithm_id is required")
		return
	}

	resp, err := h.service.ExecuteAlgorithm(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"data":     resp,
		"metadata": resp, // backwards compatibility
	})
}

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	session, err := h.service.CreateSession(r.Context(), studentID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"session": session,
	})
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sessionID := chi.URLParam(r, "id")

	session, err := h.service.GetSession(r.Context(), sessionID, studentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "session not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"session": session,
	})
}

func (h *Handler) ListSessions(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	sessions, err := h.service.ListSessions(r.Context(), studentID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":  true,
		"sessions": sessions,
		"count":    len(sessions),
	})
}

func (h *Handler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	sessionID := chi.URLParam(r, "id")

	var req UpdateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.UpdateSession(r.Context(), sessionID, studentID, req); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "session updated",
	})
}

func (h *Handler) ExplainStep(w http.ResponseWriter, r *http.Request) {
	var req ExplainStepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.ExplainStep(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    resp,
	})
}

func (h *Handler) GetPracticeQuestions(w http.ResponseWriter, r *http.Request) {
	algoID := r.URL.Query().Get("algorithm_id")
	if algoID == "" {
		algoID = chi.URLParam(r, "id")
	}
	if algoID == "" {
		writeError(w, http.StatusBadRequest, "algorithm_id parameter is required")
		return
	}

	questions, err := h.service.GetPracticeQuestions(r.Context(), algoID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"questions": questions,
		"count":     len(questions),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{
		"success": false,
		"error":   msg,
	})
}

