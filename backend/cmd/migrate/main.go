package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if cfg.SupabaseDBURL == "" {
		fmt.Println("SUPABASE_DB_URL is empty!")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	fmt.Println("Connecting to Supabase PostgreSQL...")
	pool, err := database.NewPool(ctx, cfg.SupabaseDBURL)
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	fmt.Println("Applying schema migrations...")
	if err := database.Migrate(ctx, pool); err != nil {
		fmt.Printf("Migration error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("SUCCESS: All LearnX database tables created successfully in Supabase!")
}
