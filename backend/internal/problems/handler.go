package problems

import (
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
	filter := filterFromRequest(r)
	items, total, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"problems":   items,
		"pagination": httputil.NewPagination(filter.Page, filter.Limit, total),
	})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	filter := filterFromRequest(r)
	if stringsEmpty(filter.Query) {
		httputil.RespondError(w, http.StatusBadRequest, "query parameter q is required")
		return
	}
	items, total, err := h.service.List(r.Context(), filter)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"problems":   items,
		"pagination": httputil.NewPagination(filter.Page, filter.Limit, total),
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "problem_id")
	item, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	related, err := h.service.Related(r.Context(), item.ID)
	if err != nil {
		related = []model.Problem{}
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"problem": item,
		"related": related,
	})
}

func (h *Handler) ListBySubtopic(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "subtopic_id")
	if !httputil.IsUUID(id) {
		httputil.RespondError(w, http.StatusBadRequest, "invalid subtopic_id")
		return
	}
	filter := filterFromRequest(r)
	items, total, err := h.service.ListBySubtopic(r.Context(), id, filter)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":    true,
		"problems":   items,
		"pagination": httputil.NewPagination(filter.Page, filter.Limit, total),
	})
}

func filterFromRequest(r *http.Request) model.ProblemFilter {
	page, limit := httputil.ParsePageLimit(r)
	q := r.URL.Query()
	return model.ProblemFilter{
		Query:      q.Get("q"),
		Difficulty: q.Get("difficulty"),
		TopicID:    q.Get("topic_id"),
		SubtopicID: q.Get("subtopic_id"),
		RoadmapID:  q.Get("roadmap_id"),
		Status:     q.Get("status"),
		StudentID:  middleware.GetUserID(r.Context()),
		Sort:       q.Get("sort"),
		Page:       page,
		Limit:      limit,
	}
}

func stringsEmpty(v string) bool {
	return len(v) == 0
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		httputil.RespondError(w, http.StatusNotFound, "problem not found")
	case errors.Is(err, model.ErrInvalidInput):
		httputil.RespondError(w, http.StatusBadRequest, "invalid query parameters")
	default:
		httputil.RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
