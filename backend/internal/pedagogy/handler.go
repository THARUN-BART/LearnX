package pedagogy

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	client Client
}

func NewHandler(client Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	if h.client == nil {
		respondError(w, http.StatusServiceUnavailable, "pedagogy client not configured")
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.client.Analyze(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadGateway, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

func (h *Handler) Diagnose(w http.ResponseWriter, r *http.Request) {
	if h.client == nil {
		respondError(w, http.StatusServiceUnavailable, "pedagogy client not configured")
		return
	}

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.client.Diagnose(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusBadGateway, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
	})
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if h.client == nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"status":  "unconfigured",
		})
		return
	}

	res, err := h.client.Health(r.Context())
	if err != nil {
		respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"success": false,
			"status":  "unhealthy",
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    res,
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
