package server

import (
	"encoding/json"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/go-chi/chi/v5"
)

// BuildRouter configures all HTTP routes, middlewares, and handlers.
func BuildRouter(
	cfg *config.Config,
	authHandler *auth.Handler,
	profileHandler *profile.Handler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
) http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLogger)
	r.Use(middleware.CORS(cfg.ClientURL))

	// Health Checks
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "healthy",
			"app":     cfg.OTPAppName,
			"service": "LearnX Backend Auth Service",
		})
	}
	r.Get("/health", healthHandler)
	r.Get("/api/v1/health", healthHandler)

	// API v1 Routes
	r.Route("/api/v1", func(v1 chi.Router) {
		// Authentication Routes (rate-limited)
		v1.Route("/auth", func(authRouter chi.Router) {
			if rateLimiter != nil {
				authRouter.Use(rateLimiter.Limit)
			}

			authRouter.Post("/signup", authHandler.Signup)
			authRouter.Post("/verify-signup-otp", authHandler.VerifySignupOTP)
			authRouter.Post("/login", authHandler.Login)
			authRouter.Post("/forgot-password", authHandler.ForgotPassword)
			authRouter.Post("/verify-reset-otp", authHandler.VerifyResetOTP)
			authRouter.Post("/reset-password", authHandler.ResetPassword)
			authRouter.Get("/google", authHandler.GoogleOAuth)
			authRouter.Get("/callback", authHandler.OAuthCallback)

			// Protected auth endpoints
			authRouter.With(authMiddleware.RequireAuth).Post("/logout", authHandler.Logout)
		})

		// Student Profile & Onboarding Routes (Protected)
		v1.Route("/profile", func(profileRouter chi.Router) {
			profileRouter.Use(authMiddleware.RequireAuth)
			profileRouter.Get("/", profileHandler.GetProfile)
			profileRouter.Put("/", profileHandler.UpdateProfile)
			profileRouter.Post("/onboarding", profileHandler.CompleteOnboarding)
		})

		// Convenience alias for onboarding: POST /api/v1/onboarding
		v1.With(authMiddleware.RequireAuth).Post("/onboarding", profileHandler.CompleteOnboarding)
	})

	return r
}
