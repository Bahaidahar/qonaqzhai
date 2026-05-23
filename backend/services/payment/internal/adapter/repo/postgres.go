// Package repo holds the PostgreSQL persistence adapters for the payment service.
package repo

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed all:migrations
var Migrations embed.FS

// Open opens a Postgres connection and applies pending migrations.
func Open(dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("postgres DSN required")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open pgx: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	if err := Migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// Migrate applies all embedded *.up.sql migrations in lexical order.
func Migrate(db *sql.DB) error {
	files, err := readMigrationFiles()
	if err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	for _, f := range files {
		var seen int
		if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", f.version).Scan(&seen); err != nil {
			return fmt.Errorf("check %s: %w", f.version, err)
		}
		if seen > 0 {
			continue
		}
		if _, err := db.Exec(f.sql); err != nil {
			return fmt.Errorf("apply %s: %w", f.name, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", f.version); err != nil {
			return fmt.Errorf("record %s: %w", f.version, err)
		}
	}
	return nil
}

type migrationFile struct{ name, version, sql string }

func readMigrationFiles() ([]migrationFile, error) {
	dir, err := fs.Sub(Migrations, "migrations")
	if err != nil {
		dir = Migrations
	}
	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, err
	}
	var out []migrationFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		body, err := fs.ReadFile(dir, name)
		if err != nil {
			return nil, err
		}
		out = append(out, migrationFile{name: name, version: strings.TrimSuffix(name, ".up.sql"), sql: string(body)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out, nil
}

type scanner interface{ Scan(dest ...any) error }
