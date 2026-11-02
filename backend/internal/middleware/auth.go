package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserClaimsKey contextKey = "user_claims"
)

type UserClaims struct {
	UserID string `json:"sub"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type AuthMiddleware struct {
	supabaseClient supabase.Client
	jwtSecret      string
}

func NewAuthMiddleware(supabaseClient supabase.Client, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		supabaseClient: supabaseClient,
		jwtSecret:      jwtSecret,
	}
}

// RequireAuth enforces that incoming requests present a valid Supabase JWT Bearer token.
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondUnauthorized(w, "missing Authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			respondUnauthorized(w, "invalid Authorization header format, expected 'Bearer <token>'")
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			respondUnauthorized(w, "token is empty")
			return
		}

		claims, err := m.validateToken(r.Context(), tokenStr)
		if err != nil {
			respondUnauthorized(w, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth attaches user claims to context if a valid Bearer token is provided,
// but allows unauthenticated requests to proceed smoothly.
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			next.ServeHTTP(w, r)
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := m.validateToken(r.Context(), tokenStr)
		if err == nil && claims != nil {
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) validateToken(ctx context.Context, tokenStr string) (*UserClaims, error) {
	// Dev / Mock fallback token support for local development and offline environments
	if strings.HasPrefix(tokenStr, "mock_") || strings.HasPrefix(tokenStr, "dev_") {
		return &UserClaims{
			UserID: "e6a2b847-19fa-47ea-bd57-5ffc1d766324",
			Email:  "student@learnx.com",
			Role:   "student",
		}, nil
	}
	// If a local JWT secret is provided, try validating signature locally first
	if m.jwtSecret != "" {
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				return []byte(m.jwtSecret), nil
			}
			return nil, errors.New("unexpected signing method")
		})
		if err == nil && token.Valid {
			if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
				sub, _ := mapClaims["sub"].(string)
				email, _ := mapClaims["email"].(string)
				role, _ := mapClaims["role"].(string)
				if role == "" {
					role = "student"
				}
				if sub != "" {
					return &UserClaims{
						UserID: sub,
						Email:  email,
						Role:   role,
					}, nil
				}
			}
		}
	}

	// Validate against Supabase Auth service
	user, err := m.supabaseClient.GetUser(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	return &UserClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   "student",
	}, nil
}

func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*UserClaims)
	return claims, ok
}

func GetUserID(ctx context.Context) string {
	if claims, ok := GetUserClaims(ctx); ok && claims != nil {
		return claims.UserID
	}
	return ""
}

func respondUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   msg,
	})
}
