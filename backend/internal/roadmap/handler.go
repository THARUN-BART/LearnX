package roadmap

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "failed to list roadmaps")
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"roadmaps": items,
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roadmap_id")
	item, err := h.service.Get(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"roadmap": item,
	})
}

func (h *Handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roadmap_id")
	items, err := h.service.ListTopics(r.Context(), id)
	if err != nil {
		writeErr(w, err)
		return
	}
	httputil.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"topics":  items,
	})
}

func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		httputil.RespondError(w, http.StatusNotFound, "roadmap not found")
	case errors.Is(err, model.ErrInvalidInput):
		httputil.RespondError(w, http.StatusBadRequest, err.Error())
	default:
		httputil.RespondError(w, http.StatusInternalServerError, "internal server error")
	}
}
