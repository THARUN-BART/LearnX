package codereview

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

func (h *Handler) ReviewCode(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CodeReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	res, err := h.service.ReviewCode(r.Context(), studentID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"review":  res,
		"data":    res,
	})
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	reviewID := chi.URLParam(r, "id")

	res, err := h.service.GetReview(r.Context(), reviewID, studentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "review not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"review":  res,
		"data":    res,
	})
}

func (h *Handler) ListReviews(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reviews, err := h.service.ListStudentReviews(r.Context(), studentID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"reviews": reviews,
		"count":   len(reviews),
	})
}

func (h *Handler) GetComparison(w http.ResponseWriter, r *http.Request) {
	studentID := middleware.GetUserID(r.Context())
	if studentID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	reviewID := chi.URLParam(r, "id")

	comp, err := h.service.GetComparison(r.Context(), reviewID, studentID)
	if err != nil {
		writeError(w, http.StatusNotFound, "review comparison not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success":    true,
		"comparison": comp,
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

