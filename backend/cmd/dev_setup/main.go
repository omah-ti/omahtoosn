package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/config"
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

	must(applyMigrations(ctx, pool, filepath.Join(root, "migrations")))
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

func applyMigrations(ctx context.Context, pool *pgxpool.Pool, dir string) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS dev_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.up.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(filepath.Base(file), ".up.sql")
		var exists bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM dev_migrations WHERE version = $1)`, version).Scan(&exists); err != nil {
			return err
		}
		if exists {
			continue
		}
		alreadyPresent, err := migrationAlreadyPresent(ctx, pool, version)
		if err != nil {
			return err
		}
		if alreadyPresent {
			if _, err := pool.Exec(ctx, `INSERT INTO dev_migrations (version) VALUES ($1) ON CONFLICT DO NOTHING`, version); err != nil {
				return err
			}
			fmt.Println("Migration marked applied:", filepath.Base(file))
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply migration %s: %w", filepath.Base(file), err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO dev_migrations (version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
		fmt.Println("Migration applied:", filepath.Base(file))
	}
	return nil
}

func migrationAlreadyPresent(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	switch {
	case strings.HasPrefix(version, "000001"):
		users, err := tableExists(ctx, pool, "users")
		if err != nil {
			return false, err
		}
		tryouts, err := tableExists(ctx, pool, "tryouts")
		if err != nil {
			return false, err
		}
		return users && tryouts, nil
	case strings.HasPrefix(version, "000002"):
		return tableExists(ctx, pool, "password_reset_tokens")
	default:
		return false, nil
	}
}

func tableExists(ctx context.Context, pool *pgxpool.Pool, name string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = 'public' AND table_name = $1
		)`, name).Scan(&exists)
	return exists, err
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
