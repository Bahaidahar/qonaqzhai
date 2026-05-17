package db_test

import (
	"testing"

	"qonaqzhai-backend/internal/infra/db"
	"qonaqzhai-backend/internal/infra/db/testpg"
)

func TestOpenAppliesAllMigrations(t *testing.T) {
	t.Parallel()
	dsn := testpg.Start(t)
	conn, err := db.Open(dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer conn.Close()

	tables := []string{
		"users", "vendors", "photos", "bookings", "reviews",
		"refresh_tokens", "password_reset_tokens", "notifications",
		"fcm_tokens", "schema_migrations",
		"services", "chats", "chat_messages",
		"booking_threads", "thread_messages", "payment_cards",
		"audit_log",
	}
	for _, tbl := range tables {
		var n int
		if err := conn.QueryRow("SELECT COUNT(*) FROM " + tbl).Scan(&n); err != nil {
			t.Errorf("table %q missing or unreadable: %v", tbl, err)
		}
	}
}

func TestMigrationsAreIdempotent(t *testing.T) {
	t.Parallel()
	dsn := testpg.Start(t)
	c1, err := db.Open(dsn)
	if err != nil {
		t.Fatal(err)
	}
	_ = c1.Close()
	c2, err := db.Open(dsn) // should not re-apply or error
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
	c3, err := db.Open(dsn)
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
