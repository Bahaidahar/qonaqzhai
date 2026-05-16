package db_test

import (
	"path/filepath"
	"testing"

	"qonaqzhai-backend/internal/infra/db"
)

func TestOpenAppliesAllMigrations(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	conn, err := db.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer conn.Close()

	tables := []string{
		"users", "vendors", "photos", "bookings", "reviews",
		"refresh_tokens", "password_reset_tokens", "notifications",
		"fcm_tokens", "vendors_fts", "schema_migrations",
	}
	for _, tbl := range tables {
		var n int
		if err := conn.QueryRow("SELECT COUNT(*) FROM "+tbl).Scan(&n); err != nil {
			t.Errorf("table %q missing or unreadable: %v", tbl, err)
		}
	}
}

func TestMigrationsAreIdempotent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	c1, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = c1.Close()
	c2, err := db.Open(path) // should not re-apply or error
	if err != nil {
		t.Fatalf("re-open: %v", err)
	}
	defer c2.Close()
	var n int
	if err := c2.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n < 1 {
		t.Errorf("expected at least 1 migration row, got %d", n)
	}
	// Reopen again and assert count is stable (no duplicate inserts).
	c3, err := db.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer c3.Close()
	var m int
	if err := c3.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&m); err != nil {
		t.Fatal(err)
	}
	if m != n {
		t.Errorf("migration count drifted: %d -> %d", n, m)
	}
}
