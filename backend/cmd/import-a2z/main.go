package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/THARUN-BART/LearnX/backend/internal/config"
	"github.com/THARUN-BART/LearnX/backend/internal/database"
	"github.com/THARUN-BART/LearnX/backend/internal/importer"
)

func main() {
	filePath := flag.String("file", "data/strivers-a2z-sheet.json", "Path to the Striver A2Z seed JSON file")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if cfg.SupabaseDBURL == "" {
		fmt.Println("SUPABASE_DB_URL is required for catalogue import")
		os.Exit(1)
	}

	seed, err := importer.LoadSeedFile(*filePath)
	if err != nil {
		fmt.Printf("failed to load seed file: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.SupabaseDBURL)
	if err != nil {
		fmt.Printf("database connection error: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		fmt.Printf("migration error: %v\n", err)
		os.Exit(1)
	}

	result, err := importer.ImportPostgres(ctx, pool, seed)
	if err != nil {
		fmt.Printf("import failed: %v\n", err)
		os.Exit(1)
	}

	if err := importer.SeedPostgresTestCases(ctx, pool, seed); err != nil {
		fmt.Printf("seed test cases warning: %v\n", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	fmt.Println("Striver A2Z catalogue import completed (student progress was not deleted).")
	_ = enc.Encode(result)
}
