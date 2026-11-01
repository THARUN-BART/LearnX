package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/assessment"
	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/database"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/execution"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/submissions"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
	"github.com/THARUN-BART/LearnX/backend/internal/testcases"
	"github.com/THARUN-BART/LearnX/backend/internal/workers"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("starting LearnX authentication service",
		slog.String("port", cfg.ServerPort),
		slog.String("env", cfg.ServerEnv),
		slog.String("otp_base_url", cfg.OTPServiceBaseURL),
		slog.String("app_name", cfg.OTPAppName),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Database Repository
	var authRepo auth.Repository
	var profileRepo profile.Repository
	var dsaRepos dsa.Repositories
	var submissionsRepo submissions.Repository
	var testCasesRepo testcases.Repository

	if cfg.SupabaseDBURL != "" {
		pool, err := database.NewPool(ctx, cfg.SupabaseDBURL)
		if err != nil {
			slog.Warn("could not connect to PostgreSQL database pool, falling back to memory repository for local development", slog.String("error", err.Error()))
			memRepo := auth.NewMemoryRepository()
			authRepo = memRepo
			memProfileRepo := profile.NewMemoryRepository()
			profileRepo = memProfileRepo
			dsaRepos = dsa.NewMemoryRepositories(model.NewMemory())
			submissionsRepo = submissions.NewMemoryRepository()
			testCasesRepo = testcases.NewMemoryRepository()
		} else {
			defer pool.Close()
			if err := database.Migrate(ctx, pool); err != nil {
				slog.Error("database migration error", slog.String("error", err.Error()))
			}
			authRepo = auth.NewPostgresRepository(pool)
			profileRepo = profile.NewPostgresRepository(pool)
			dsaRepos = dsa.NewPostgresRepositories(pool)
			submissionsRepo = submissions.NewPostgresRepository(pool)
			testCasesRepo = testcases.NewPostgresRepository(pool)
		}
	} else {
		slog.Warn("SUPABASE_DB_URL not configured: using in-memory store for development/testing")
		memRepo := auth.NewMemoryRepository()
		authRepo = memRepo
		memProfileRepo := profile.NewMemoryRepository()
		profileRepo = memProfileRepo
		dsaRepos = dsa.NewMemoryRepositories(model.NewMemory())
		submissionsRepo = submissions.NewMemoryRepository()
		testCasesRepo = testcases.NewMemoryRepository()
	}

	// Code Execution Client
	var execClient execution.Client
	switch strings.ToLower(cfg.ExecutionProvider) {
	case "wandbox":
		execClient = execution.NewWandboxClient(cfg.ExecutionTimeout)
	case "local":
		execClient = execution.NewLocalRunner()
	case "mock":
		execClient = execution.NewMockExecutionClient()
	default:
		execClient = execution.NewJudge0Client(execution.Judge0Config{
			BaseURL: cfg.Judge0APIURL,
			APIKey:  cfg.Judge0APIKey,
			APIHost: cfg.Judge0APIHost,
			Timeout: cfg.ExecutionTimeout,
			Logger:  logger,
		})
	}
	slog.Info("code execution provider configured", slog.String("provider", execClient.ProviderName()))

	// Assessment Engine & Submission Worker Pool
	assessmentEngine := assessment.NewEngine(execClient, logger)
	jobQueue := workers.NewChannelJobQueue(cfg.SubmissionQueueSize)
	workerPool := workers.NewWorkerPool(
		jobQueue,
		submissionsRepo,
		testCasesRepo,
		dsaRepos.Progress,
		assessmentEngine,
		cfg.SubmissionWorkerConcurrency,
		logger,
	)
	workerPool.Start()
	defer workerPool.Stop()

	// Submissions & TestCases Services & Handlers
	submissionsSvc := submissions.NewService(submissionsRepo, dsaRepos.Problems, jobQueue)
	submissionsHandler := submissions.NewHandler(submissionsSvc)
	testCasesSvc := testcases.NewService(testCasesRepo)
	testCasesHandler := testcases.NewHandler(testCasesSvc)

	subHandlers := &server.SubmissionHandlers{
		Submissions: submissionsHandler,
		TestCases:   testCasesHandler,
	}

	// DSA Handlers
	dsaHandlers := dsa.NewHandlers(dsaRepos)

	// Initialize OTP Client and Service
	otpClient, err := otp.NewClient(cfg.OTPServiceBaseURL, cfg.OTPAppName, cfg.OTPTimeout)
	if err != nil {
		slog.Error("failed to initialize OTP client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	otpService := otp.NewService(otpClient, cfg.OTPCooldown, cfg.MaxOTPAttemptsPerWindow)

	// Initialize Supabase Client
	supabaseClient := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseServiceRoleKey)

	// Initialize Services
	authService := auth.NewService(
		authRepo,
		otpService,
		supabaseClient,
		cfg.EncryptionKey,
		cfg.PendingSignupExpiry,
		cfg.PasswordResetTokenExpiry,
	)
	profileService := profile.NewService(profileRepo)

	// Initialize Handlers
	authHandler := auth.NewHandler(authService)
	profileHandler := profile.NewHandler(profileService)

	// Initialize Middlewares
	authMiddleware := middleware.NewAuthMiddleware(supabaseClient, cfg.SupabaseJWTSecret)
	rateLimiter := middleware.NewRateLimiter(60, 1*time.Minute) // 60 requests per minute burst limit

	// Build HTTP Router with DSA & Submissions enabled
	router := server.BuildRouter(
		cfg,
		authHandler,
		profileHandler,
		authMiddleware,
		rateLimiter,
		server.WithDSA(dsaHandlers),
		server.WithSubmissions(subHandlers),
	)

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context for graceful shutdown
	serverErr := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("LearnX Auth Server listening on port %s", cfg.ServerPort))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// Listen for OS termination signals
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		slog.Error("server startup error", slog.String("error", err.Error()))
	case sig := <-shutdownSignal:
		slog.Info("shutdown signal received, gracefully shutting down...", slog.String("signal", sig.String()))
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("server forced shutdown error", slog.String("error", err.Error()))
		}
		slog.Info("server exited cleanly")
	}
}
