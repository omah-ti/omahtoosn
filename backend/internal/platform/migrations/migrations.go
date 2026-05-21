package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Apply(ctx context.Context, pool *pgxpool.Pool, dir string) error {
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
		isAlreadyPresent, err := alreadyPresent(ctx, pool, version)
		if err != nil {
			return err
		}
		if isAlreadyPresent {
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

func alreadyPresent(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
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
