package aitutor

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

// RequestHint handles POST /api/v1/ai/hint
func (h *Handler) RequestHint(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req HintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ProblemID == "" {
		respondError(w, http.StatusBadRequest, "problem_id is required")
		return
	}

	hint, err := h.service.GetProgressiveHint(r.Context(), studentID, req)
	if err != nil {
		if errors.Is(err, ErrInvalidHintLevel) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, ErrProblemNotFound) {
			respondError(w, http.StatusNotFound, "problem not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    hint,
	})
}

// GetUnlockedHints handles GET /api/v1/problems/{problem_id}/hints
func (h *Handler) GetUnlockedHints(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	problemID := chi.URLParam(r, "problem_id")
	if problemID == "" {
		respondError(w, http.StatusBadRequest, "problem_id is required")
		return
	}

	hints, err := h.service.GetUnlockedHints(r.Context(), studentID, problemID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    hints,
	})
}

// ExplainMistake handles POST /api/v1/ai/explain-mistake
func (h *Handler) ExplainMistake(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req ExplainMistakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.SubmissionID == "" {
		respondError(w, http.StatusBadRequest, "submission_id is required")
		return
	}

	resp, err := h.service.ExplainMistake(r.Context(), studentID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

// WalkthroughCode handles POST /api/v1/ai/walkthrough
func (h *Handler) WalkthroughCode(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req CodeWalkthroughRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.WalkthroughCode(r.Context(), studentID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

// GenerateDryRun handles POST /api/v1/ai/dry-run
func (h *Handler) GenerateDryRun(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req DryRunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.GenerateDryRun(r.Context(), studentID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

// AnalyzeComplexity handles POST /api/v1/ai/complexity
func (h *Handler) AnalyzeComplexity(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req ComplexityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.AnalyzeComplexity(r.Context(), studentID, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

// AskTutor handles POST /api/v1/ai/ask
func (h *Handler) AskTutor(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	var req AskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Question == "" {
		respondError(w, http.StatusBadRequest, "question cannot be empty")
		return
	}

	resp, err := h.service.AskTutor(r.Context(), studentID, req)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "you do not have access to this conversation")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    resp,
	})
}

// ListConversations handles GET /api/v1/ai/conversations
func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
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

	convs, err := h.service.ListConversations(r.Context(), studentID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    convs,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetConversationMessages handles GET /api/v1/ai/conversations/{conversation_id}
func (h *Handler) GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		respondError(w, http.StatusUnauthorized, "unauthorized: valid session required")
		return
	}

	conversationID := chi.URLParam(r, "conversation_id")
	if conversationID == "" {
		respondError(w, http.StatusBadRequest, "conversation_id is required")
		return
	}

	msgs, err := h.service.GetConversationMessages(r.Context(), studentID, conversationID)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			respondError(w, http.StatusForbidden, "unauthorized")
			return
		}
		if errors.Is(err, ErrConversationNotFound) {
			respondError(w, http.StatusNotFound, "conversation not found")
			return
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    msgs,
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
