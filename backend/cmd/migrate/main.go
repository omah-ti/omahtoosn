package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/db"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/migrations"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	cfg := config.Load()
	pool, err := db.New(cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "db open failed:", err)
		os.Exit(1)
	}
	defer pool.Close()

	dir := filepath.Clean(envString("MIGRATIONS_PATH", "./migrations"))
	if err := migrations.Apply(ctx, pool, dir); err != nil {
		fmt.Fprintln(os.Stderr, "migration failed:", err)
		os.Exit(1)
	}
	fmt.Println("Migrations completed.")
}

func envString(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
