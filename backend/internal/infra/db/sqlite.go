// Package db opens the SQLite connection and applies migrations.
package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

// Migrations holds the embedded SQL migration files.
//
//go:embed all:migrations
var Migrations embed.FS

// Open opens the SQLite database at path, configures pragmas, and runs migrations.
// Returns a ready-to-use *sql.DB.
func Open(path string) (*sql.DB, error) {
	dsn := path + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	db.SetMaxOpenConns(1) // SQLite serializes writes; one conn avoids busy errors
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	if err := Migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

// Migrate applies all embedded *_up.sql migrations in order.
// Migration filenames are expected to look like `NNNN_name.up.sql`.
func Migrate(db *sql.DB) error {
	files, err := readMigrationFiles()
	if err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	for _, f := range files {
		var seen int
		if err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", f.version).Scan(&seen); err != nil {
			return fmt.Errorf("check migration %s: %w", f.version, err)
		}
		if seen > 0 {
			continue
		}
		if _, err := db.Exec(f.sql); err != nil {
			return fmt.Errorf("apply %s: %w", f.name, err)
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", f.version); err != nil {
			return fmt.Errorf("record %s: %w", f.version, err)
		}
	}
	return nil
}

type migrationFile struct {
	name    string
	version string
	sql     string
}

func readMigrationFiles() ([]migrationFile, error) {
	dir, err := fs.Sub(Migrations, "migrations")
	if err != nil {
		// Fallback: maybe embed didn't include subdir prefix; read root.
		dir = Migrations
	}
	entries, err := fs.ReadDir(dir, ".")
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
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
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		version := strings.TrimSuffix(name, ".up.sql")
		out = append(out, migrationFile{name: name, version: version, sql: string(body)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out, nil
}
