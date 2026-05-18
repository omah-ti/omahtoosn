package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/migrations"
)

func main() {
	var runAPI bool
	var skipDocker bool
	var allowNonLocalDB bool
	var seedMode string

	flag.BoolVar(&runAPI, "run", false, "run API after setup")
	flag.BoolVar(&skipDocker, "skip-docker", false, "skip docker compose up")
	flag.BoolVar(&allowNonLocalDB, "allow-nonlocal-db", false, "allow setup against a non-local database")
	flag.StringVar(&seedMode, "seed", "demo", "seed mode: demo, omahtoosn, none")
	flag.Parse()

	root, err := backendRoot()
	must(err)
	must(os.Chdir(root))

	must(ensureEnvFile(root))

	if !skipDocker {
		must(runDockerCompose(root))
	}

	cfg := config.Load()
	if !allowNonLocalDB {
		must(ensureLocalDatabaseURL(cfg.DatabaseURL))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	pool := mustConnectDB(ctx, cfg.DatabaseURL)
	defer pool.Close()

	must(migrations.Apply(ctx, pool, filepath.Join(root, "migrations")))
	must(runSeed(ctx, pool, root, seedMode))

	fmt.Println("Backend setup selesai.")
	fmt.Println("API URL: http://localhost:" + cfg.AppPort)
	fmt.Println("Swagger: http://localhost:" + cfg.AppPort + "/swagger/index.html")

	if runAPI {
		fmt.Println("Menjalankan API. Tekan Ctrl+C untuk berhenti.")
		cmd := exec.Command("go", "run", "./cmd/api")
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		must(cmd.Run())
	}
}

func backendRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for _, candidate := range []string{cwd, filepath.Join(cwd, "backend")} {
		if exists(filepath.Join(candidate, "go.mod")) && exists(filepath.Join(candidate, "docker-compose.yml")) {
			return candidate, nil
		}
	}
	return "", errors.New("jalankan dari repo root atau folder backend")
}

func ensureEnvFile(root string) error {
	envPath := filepath.Join(root, ".env")
	if exists(envPath) {
		fmt.Println(".env sudah ada, tidak ditimpa.")
		return nil
	}

	content, err := os.ReadFile(filepath.Join(root, ".env.example"))
	if err != nil {
		return fmt.Errorf("read .env.example: %w", err)
	}
	if err := os.WriteFile(envPath, content, 0600); err != nil {
		return fmt.Errorf("write .env: %w", err)
	}
	fmt.Println(".env dibuat dari .env.example.")
	return nil
}

func runDockerCompose(root string) error {
	fmt.Println("Menjalankan PostgreSQL Docker...")

	if _, err := exec.LookPath("docker"); err == nil {
		cmd := exec.Command("docker", "compose", "up", "-d")
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	if _, err := exec.LookPath("docker-compose"); err != nil {
		return fmt.Errorf("docker compose tidak ditemukan. Install Docker Desktop atau jalankan ulang dengan --skip-docker jika memakai PostgreSQL lokal")
	}
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker-compose up failed: %w", err)
	}
	return nil
}

func ensureLocalDatabaseURL(databaseURL string) error {
	host := ""
	if config, err := pgxpool.ParseConfig(databaseURL); err == nil && config.ConnConfig != nil {
		host = strings.TrimSpace(config.ConnConfig.Host)
	}
	switch strings.ToLower(host) {
	case "", "localhost", "127.0.0.1", "::1":
		return nil
	default:
		return fmt.Errorf("DATABASE_URL host %q bukan localhost. Ubah .env ke DB lokal atau pakai --allow-nonlocal-db jika benar-benar sengaja", host)
	}
}

func mustConnectDB(ctx context.Context, databaseURL string) *pgxpool.Pool {
	deadline := time.Now().Add(60 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		pool, err := pgxpool.New(ctx, databaseURL)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
			if err == nil {
				fmt.Println("Database siap.")
				return pool
			}
			pool.Close()
		}
		lastErr = err
		time.Sleep(2 * time.Second)
	}
	must(fmt.Errorf("database belum siap: %w", lastErr))
	return nil
}

func runSeed(ctx context.Context, pool *pgxpool.Pool, root string, mode string) error {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "", "none", "skip":
		fmt.Println("Seed dilewati.")
		return nil
	case "demo":
		sql, err := os.ReadFile(filepath.Join(root, "seeds", "demo_seed.sql"))
		if err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("seed demo: %w", err)
		}
		fmt.Println("Seed demo selesai.")
		return nil
	case "omahtoosn":
		cmd := exec.Command("go", "run", "./cmd/seed_questions")
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "TRYOUT_STATUS=ongoing")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Println("Menjalankan seed OmahTOOSN...")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("seed omahtoosn: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("seed mode tidak dikenal: %s", mode)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func must(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "ERROR:", err)
	os.Exit(1)
}
