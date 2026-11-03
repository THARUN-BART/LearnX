package server

import (
	"encoding/json"
	"net/http"

	"github.com/THARUN-BART/LearnX/backend/internal/aitutor"
	"github.com/THARUN-BART/LearnX/backend/internal/analytics"
	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/learningmemory"
	"github.com/THARUN-BART/LearnX/backend/internal/learningpath"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/mistakes"
	"github.com/THARUN-BART/LearnX/backend/internal/pedagogy"
	"github.com/THARUN-BART/LearnX/backend/internal/problems"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/progress"
	"github.com/THARUN-BART/LearnX/backend/internal/recommendations"
	"github.com/THARUN-BART/LearnX/backend/internal/roadmap"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
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

type SubmissionHandlers struct {
	Submissions *submissions.Handler
	TestCases   *testcases.Handler
}

type IntelligenceHandlers struct {
	Mistakes        *mistakes.Handler
	AITutor         *aitutor.Handler
	Pedagogy        *pedagogy.Handler
	PersonC         *pedagogy.Handler // Deprecated alias for backwards compatibility
	LearningMemory  *learningmemory.Handler
	LearningPath    *learningpath.Handler
	Recommendations *recommendations.Handler
}

type RouterOption func(*routerOptions)

type routerOptions struct {
	dsa          *DSAHandlers
	submissions  *SubmissionHandlers
	intelligence *IntelligenceHandlers
}

func WithDSA(handlers *DSAHandlers) RouterOption {
	return func(opts *routerOptions) {
		opts.dsa = handlers
	}
}

func WithSubmissions(handlers *SubmissionHandlers) RouterOption {
	return func(opts *routerOptions) {
		opts.submissions = handlers
	}
}

func WithIntelligence(handlers *IntelligenceHandlers) RouterOption {
	return func(opts *routerOptions) {
		opts.intelligence = handlers
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

		// Public & Guest-accessible educational catalog routes (with OptionalAuth for token enrichment)
		r.Group(func(r chi.Router) {
			if authMiddleware != nil {
				r.Use(authMiddleware.OptionalAuth)
			}

			// Public roadmap browsing (singular alias for seamless guest access)
			if options.dsa != nil && options.dsa.Roadmap != nil {
				r.Get("/roadmap", options.dsa.Roadmap.List)
				r.Get("/roadmap/{roadmap_id}", options.dsa.Roadmap.Get)
				r.Get("/roadmap/{roadmap_id}/topics", options.dsa.Roadmap.ListTopics)
			}

			// Public topics and problems catalog
			if options.dsa != nil {
				if options.dsa.Topics != nil {
					r.Get("/topics/{topic_id}", options.dsa.Topics.Get)
					r.Get("/topics/{topic_id}/subtopics", options.dsa.Topics.ListSubtopics)
				}
				if options.dsa.Problems != nil {
					r.Get("/problems/search", options.dsa.Problems.Search)
					r.Get("/problems", options.dsa.Problems.List)
					r.Get("/problems/{problem_id}", options.dsa.Problems.Get)
					r.Get("/subtopics/{subtopic_id}/problems", options.dsa.Problems.ListBySubtopic)
				}
			}

			// Public sample testcases for problem preview
			if options.submissions != nil && options.submissions.TestCases != nil {
				r.Get("/problems/{problem_id}/testcases/sample", options.submissions.TestCases.GetSampleTestCases)
			}

			// Guest-capable AI Tutor endpoints (supports both authenticated students & guests)
			if options.intelligence != nil && options.intelligence.AITutor != nil {
				r.Post("/ai/walkthrough", options.intelligence.AITutor.WalkthroughCode)
				r.Post("/ai/dry-run", options.intelligence.AITutor.GenerateDryRun)
				r.Post("/ai/complexity", options.intelligence.AITutor.AnalyzeComplexity)
				r.Post("/ai/ask", options.intelligence.AITutor.AskTutor)
			}

			if options.intelligence != nil {
				pedagogyHandler := options.intelligence.Pedagogy
				if pedagogyHandler == nil {
					pedagogyHandler = options.intelligence.PersonC
				}
				if pedagogyHandler != nil {
					// Canonical pedagogy routes
					r.Post("/ai/pedagogy/analyze", pedagogyHandler.Analyze)
					r.Post("/ai/pedagogy/diagnose", pedagogyHandler.Diagnose)
					r.Get("/ai/pedagogy/health", pedagogyHandler.Health)

					// Backwards-compatible aliases
					r.Post("/ai/person-c/analyze", pedagogyHandler.Analyze)
					r.Post("/ai/person-c/diagnose", pedagogyHandler.Diagnose)
					r.Get("/ai/person-c/health", pedagogyHandler.Health)
				}
			}
		})

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

			if options.submissions != nil {
				registerSubmissionRoutes(r, options.submissions)
			}

			if options.intelligence != nil {
				registerIntelligenceRoutes(r, options.intelligence)
			}
		})
	})

	return r
}

func registerIntelligenceRoutes(r chi.Router, intel *IntelligenceHandlers) {
	if intel.Mistakes != nil {
		r.Get("/student/mistakes", intel.Mistakes.GetStudentMistakes)
		r.Get("/student/mistakes/summary", intel.Mistakes.GetStudentMistakeSummary)
		r.Get("/student/mistakes/patterns", intel.Mistakes.GetStudentPatterns)
		r.Get("/submissions/{submission_id}/diagnosis", intel.Mistakes.GetSubmissionDiagnosis)
		r.Post("/student/mistakes/{diagnosis_id}/feedback", intel.Mistakes.SubmitFeedback)
	}

	if intel.AITutor != nil {
		r.Post("/ai/hint", intel.AITutor.RequestHint)
		r.Get("/problems/{problem_id}/hints", intel.AITutor.GetUnlockedHints)
		r.Post("/ai/explain-mistake", intel.AITutor.ExplainMistake)
		r.Get("/ai/conversations", intel.AITutor.ListConversations)
		r.Get("/ai/conversations/{conversation_id}", intel.AITutor.GetConversationMessages)
	}

	if intel.LearningMemory != nil {
		r.Get("/student/learning-memory", intel.LearningMemory.GetSummary)
		r.Get("/student/learning-memory/summary", intel.LearningMemory.GetStudentFacingSummary)
		r.Get("/student/learning-memory/concepts", intel.LearningMemory.ListConcepts)
		r.Get("/student/learning-memory/concepts/{concept_key}", intel.LearningMemory.GetConceptDetail)
		r.Get("/student/learning-memory/strengths", intel.LearningMemory.GetStrengths)
		r.Get("/student/learning-memory/weaknesses", intel.LearningMemory.GetWeaknesses)
		r.Get("/student/learning-memory/history", intel.LearningMemory.GetHistory)
		r.Get("/student/learning-memory/timeline", intel.LearningMemory.GetTimeline)
		r.Get("/student/learning-memory/snapshots", intel.LearningMemory.GetSnapshots)
	}

	if intel.LearningPath != nil {
		r.Get("/student/learning-path", intel.LearningPath.GetLearningPath)
		r.Post("/student/learning-path/generate", intel.LearningPath.RecalculatePath)
		r.Get("/student/learning-path/next-action", intel.LearningPath.GetNextAction)
		r.Get("/student/learning-path/progress", intel.LearningPath.GetPathProgress)
		r.Get("/student/learning-path/history", intel.LearningPath.GetHistory)
		r.Post("/student/learning-path/recalculate", intel.LearningPath.RecalculatePath)
	}

	if intel.Recommendations != nil {
		r.Get("/student/recommendations", intel.Recommendations.GetRecommendations)
		r.Get("/student/recommendations/next", intel.Recommendations.GetNext)
		r.Get("/student/recommendations/weaknesses", intel.Recommendations.GetWeaknesses)
		r.Get("/student/recommendations/similar/{problem_id}", intel.Recommendations.GetSimilar)
		r.Post("/student/recommendations/refresh", intel.Recommendations.Refresh)
		r.Post("/student/recommendations/{id}/feedback", intel.Recommendations.SubmitFeedback)
		r.Get("/student/recommendations/history", intel.Recommendations.GetHistory)
	}
}

func registerSubmissionRoutes(r chi.Router, s *SubmissionHandlers) {
	if s.Submissions != nil {
		r.Post("/submissions", s.Submissions.Submit)
		r.Get("/submissions/{submission_id}", s.Submissions.GetStatus)
		r.Get("/submissions/{submission_id}/result", s.Submissions.GetResult)
		r.Post("/submissions/{submission_id}/retry", s.Submissions.Retry)

		r.Get("/student/submissions", s.Submissions.ListStudentSubmissions)
		r.Get("/problems/{problem_id}/submissions", s.Submissions.ListProblemSubmissions)
		r.Get("/problems/{problem_id}/submissions/latest", s.Submissions.GetLatestProblemSubmission)
	}
	if s.TestCases != nil {
		r.Get("/problems/{problem_id}/testcases/sample", s.TestCases.GetSampleTestCases)
		r.Post("/admin/problems/{problem_id}/testcases", s.TestCases.CreateTestCase)
	}
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
