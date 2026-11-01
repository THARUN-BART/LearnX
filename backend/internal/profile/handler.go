package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetProfile handles GET /api/v1/profile.
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized: user context not found",
		})
		return
	}

	p, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "profile not found",
			})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   "failed to fetch profile",
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"profile": p,
	})
}

// UpdateProfile handles PUT /api/v1/profile.
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	updated, err := h.service.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "profile not found",
			})
			return
		}
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "profile updated successfully",
		"profile": updated,
	})
}

// CompleteOnboarding handles POST /api/v1/profile/onboarding.
func (h *Handler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		respondJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"error":   "unauthorized: user context not found",
		})
		return
	}

	var req OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	p, err := h.service.CompleteOnboarding(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]interface{}{
				"success": false,
				"error":   "profile not found",
			})
			return
		}
		respondJSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "student onboarding completed successfully",
		"profile": p,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
