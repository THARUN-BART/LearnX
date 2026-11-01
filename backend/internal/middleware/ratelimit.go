package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type clientRecord struct {
	tokens     int
	lastRefill time.Time
}

type RateLimiter struct {
	mu         sync.Mutex
	records    map[string]*clientRecord
	limit      int
	window     time.Duration
	refillRate float64
}

// NewRateLimiter creates a token-bucket rate limiter.
// limit: maximum burst allowable.
// window: period over which limit tokens are refreshed.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		records:    make(map[string]*clientRecord),
		limit:      limit,
		window:     window,
		refillRate: float64(limit) / window.Seconds(),
	}

	// Periodic cleanup of stale client records every 5 minutes
	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, rec := range rl.records {
			if now.Sub(rec.lastRefill) > 2*rl.window {
				delete(rl.records, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	rec, exists := rl.records[key]
	if !exists {
		rl.records[key] = &clientRecord{
			tokens:     rl.limit - 1,
			lastRefill: now,
		}
		return true
	}

	// Calculate tokens to refill
	elapsed := now.Sub(rec.lastRefill).Seconds()
	tokensToAdd := int(elapsed * rl.refillRate)
	if tokensToAdd > 0 {
		rec.tokens += tokensToAdd
		if rec.tokens > rl.limit {
			rec.tokens = rl.limit
		}
		rec.lastRefill = now
	}

	if rec.tokens > 0 {
		rec.tokens--
		return true
	}

	return false
}

// Limit returns an HTTP middleware enforcing rate limiting based on client IP.
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetClientIP(r)
		if !rl.Allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "60")
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "rate limit exceeded, please try again later",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GetClientIP extracts the real client IP from X-Forwarded-For or RemoteAddr.
func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
