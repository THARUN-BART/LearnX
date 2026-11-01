package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/auth"
	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/database"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa"
	"github.com/THARUN-BART/LearnX/backend/internal/dsa/model"
	"github.com/THARUN-BART/LearnX/backend/internal/importer"
	"github.com/THARUN-BART/LearnX/backend/internal/middleware"
	"github.com/THARUN-BART/LearnX/backend/internal/otp"
	"github.com/THARUN-BART/LearnX/backend/internal/profile"
	"github.com/THARUN-BART/LearnX/backend/internal/server"
	"github.com/THARUN-BART/LearnX/backend/internal/supabase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("starting LearnX backend",
		slog.String("port", cfg.ServerPort),
		slog.String("env", cfg.ServerEnv),
		slog.String("otp_base_url", cfg.OTPServiceBaseURL),
		slog.String("app_name", cfg.OTPAppName),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var authRepo auth.Repository
	var profileRepo profile.Repository
	var dsaRepos dsa.Repositories
	var pool *pgxpool.Pool

	if cfg.SupabaseDBURL != "" {
		connected, err := database.NewPool(ctx, cfg.SupabaseDBURL)
		if err != nil {
			slog.Warn("could not connect to PostgreSQL database pool, falling back to memory repository for local development", slog.String("error", err.Error()))
			authRepo = auth.NewMemoryRepository()
			profileRepo = profile.NewMemoryRepository()
			dsaRepos = loadMemoryDSA()
		} else {
			pool = connected
			defer pool.Close()
			if err := database.Migrate(ctx, pool); err != nil {
				slog.Error("database migration error", slog.String("error", err.Error()))
			}
			authRepo = auth.NewPostgresRepository(pool)
			profileRepo = profile.NewPostgresRepository(pool)
			dsaRepos = dsa.NewPostgresRepositories(pool)
		}
	} else {
		slog.Warn("SUPABASE_DB_URL not configured: using in-memory store for development/testing")
		authRepo = auth.NewMemoryRepository()
		profileRepo = profile.NewMemoryRepository()
		dsaRepos = loadMemoryDSA()
	}

	otpClient, err := otp.NewClient(cfg.OTPServiceBaseURL, cfg.OTPAppName, cfg.OTPTimeout)
	if err != nil {
		slog.Error("failed to initialize OTP client", slog.String("error", err.Error()))
		os.Exit(1)
	}
	otpService := otp.NewService(otpClient, cfg.OTPCooldown, cfg.MaxOTPAttemptsPerWindow)

	supabaseClient := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey, cfg.SupabaseServiceRoleKey)

	authService := auth.NewService(
		authRepo,
		otpService,
		supabaseClient,
		cfg.EncryptionKey,
		cfg.PendingSignupExpiry,
		cfg.PasswordResetTokenExpiry,
	)
	profileService := profile.NewService(profileRepo)

	authHandler := auth.NewHandler(authService)
	profileHandler := profile.NewHandler(profileService)

	authMiddleware := middleware.NewAuthMiddleware(supabaseClient, cfg.SupabaseJWTSecret)
	rateLimiter := middleware.NewRateLimiter(60, 1*time.Minute)

	dsaHandlers := dsa.NewHandlers(dsaRepos)
	router := server.BuildRouter(cfg, authHandler, profileHandler, authMiddleware, rateLimiter, server.WithDSA(dsaHandlers))

	httpServer := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		slog.Info(fmt.Sprintf("LearnX backend listening on port %s", cfg.ServerPort))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

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

func loadMemoryDSA() dsa.Repositories {
	mem := model.NewMemory()
	for _, path := range []string{"data/strivers-a2z-sheet.json", "data/strivers-a2z-sheet.sample.json"} {
		seed, err := importer.LoadSeedFile(path)
		if err != nil {
			continue
		}
		if _, err := importer.ImportMemory(mem, seed); err != nil {
			slog.Warn("failed to load in-memory A2Z seed", slog.String("file", path), slog.String("error", err.Error()))
			continue
		}
		slog.Info("loaded in-memory Striver A2Z catalogue", slog.String("file", path))
		break
	}
	return dsa.NewMemoryRepositories(mem)
}
