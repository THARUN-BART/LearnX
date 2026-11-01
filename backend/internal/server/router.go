package server

import (
	"encoding/json"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/analytics"
	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/roadmap"
	"github.com/THARUN-BART/LearnX/backend/internal/topics"
	"github.com/go-chi/chi/v5"
)

type DSAHandlers struct {
	Roadmap   *roadmap.Handler
	Topics    *topics.Handler
	Problems  *problems.Handler
	Progress  *progress.Handler
	Analytics *analytics.Handler
}

type RouterOption func(*routerOptions)

type routerOptions struct {
	dsa *DSAHandlers
}

func WithDSA(handlers *DSAHandlers) RouterOption {
	return func(opts *routerOptions) {
		opts.dsa = handlers
	}
}

func BuildRouter(
	cfg *config.Config,
	authHandler *auth.Handler,
	profileHandler *profile.Handler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
	opts ...RouterOption,
) http.Handler {
	options := &routerOptions{}
	for _, opt := range opts {
		opt(options)
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestLogger)
	origin := "*"
	if cfg != nil && cfg.ClientURL != "" {
		origin = cfg.ClientURL
	}
	r.Use(middleware.CORS(origin))
	if rateLimiter != nil {
		r.Use(rateLimiter.Limit)
	}

	health := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"status":  "ok",
			"service": "learnx-backend",
		})
	}

	r.Get("/health", health)
	r.Get("/api/v1/health", health)

	r.Route("/api/v1", func(r chi.Router) {
		if authHandler != nil {
			r.Post("/auth/signup", authHandler.Signup)
			r.Post("/auth/verify-signup-otp", authHandler.VerifySignupOTP)
			r.Post("/auth/login", authHandler.Login)
			r.Post("/auth/forgot-password", authHandler.ForgotPassword)
			r.Post("/auth/verify-reset-otp", authHandler.VerifyResetOTP)
			r.Post("/auth/reset-password", authHandler.ResetPassword)
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/auth/google", authHandler.GoogleOAuth)
			r.Get("/auth/callback", authHandler.OAuthCallback)
		}

		r.Group(func(r chi.Router) {
			if authMiddleware != nil {
				r.Use(authMiddleware.RequireAuth)
			}

			if profileHandler != nil {
				r.Get("/profile", profileHandler.GetProfile)
				r.Put("/profile", profileHandler.UpdateProfile)
				r.Post("/profile/onboarding", profileHandler.CompleteOnboarding)
			}

			if options.dsa != nil {
				registerDSARoutes(r, options.dsa)
			}
		})
	})

	return r
}

func registerDSARoutes(r chi.Router, dsa *DSAHandlers) {
	if dsa.Roadmap != nil {
		r.Get("/roadmaps", dsa.Roadmap.List)
		r.Get("/roadmaps/{roadmap_id}", dsa.Roadmap.Get)
		r.Get("/roadmaps/{roadmap_id}/topics", dsa.Roadmap.ListTopics)
	}
	if dsa.Topics != nil {
		r.Get("/topics/{topic_id}", dsa.Topics.Get)
		r.Get("/topics/{topic_id}/subtopics", dsa.Topics.ListSubtopics)
	}
	if dsa.Problems != nil {
		r.Get("/problems/search", dsa.Problems.Search)
		r.Get("/problems", dsa.Problems.List)
		r.Get("/problems/{problem_id}", dsa.Problems.Get)
		r.Get("/subtopics/{subtopic_id}/problems", dsa.Problems.ListBySubtopic)
	}
	if dsa.Progress != nil {
		r.Get("/student/progress", dsa.Progress.List)
		r.Get("/student/progress/problems/{problem_id}", dsa.Progress.GetProblem)
		r.Post("/student/progress/problems/{problem_id}/attempt", dsa.Progress.Attempt)
		r.Patch("/student/progress/problems/{problem_id}", dsa.Progress.Update)
	}
	if dsa.Analytics != nil {
		r.Get("/student/progress/summary", dsa.Analytics.Summary)
		r.Get("/student/progress/topics", dsa.Analytics.Topics)
		r.Get("/student/analytics/overview", dsa.Analytics.Overview)
		r.Get("/student/analytics/topics", dsa.Analytics.Topics)
		r.Get("/student/analytics/difficulty", dsa.Analytics.Difficulty)
		r.Get("/student/analytics/activity", dsa.Analytics.Activity)
	}
}
