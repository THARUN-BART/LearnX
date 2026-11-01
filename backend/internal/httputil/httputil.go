package httputil

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

func IsUUID(value string) bool {
	return uuidPattern.MatchString(strings.TrimSpace(value))
}

func ParsePageLimit(r *http.Request) (page, limit int) {
	page = DefaultPage
	limit = DefaultLimit

	if raw := strings.TrimSpace(r.URL.Query().Get("page")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return page, limit
}

func Offset(page, limit int) int {
	return (page - 1) * limit
}

func TotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	if total == 0 {
		return 0
	}
	return (total + limit - 1) / limit
}

func NewPagination(page, limit, total int) Pagination {
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: TotalPages(total, limit),
	}
}
