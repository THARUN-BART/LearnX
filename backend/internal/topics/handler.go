package topics

import (
	"errors"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/httputil"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "topic_id")
	if id == "" {
		httputil.RespondError(w, http.StatusBadRequest, "invalid topic_id")
		return
	}
	item, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err, "topic not found")
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"topic":   item,
	})
}

func (h *Handler) ListSubtopics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "topic_id")
	if id == "" {
		httputil.RespondError(w, http.StatusBadRequest, "invalid topic_id")
		return
	}
	items, err := h.service.ListSubtopics(r.Context(), id)
	if err != nil {
		writeErr(w, err, "topic not found")
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"subtopics": items,
	})
}

func writeErr(w http.ResponseWriter, err error, notFound string) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		httputil.RespondError(w, http.StatusNotFound, notFound)
	case errors.Is(err, model.ErrInvalidInput):
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
	default:
		httputil.RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
