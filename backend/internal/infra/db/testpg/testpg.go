// Package testpg spawns a one-shot Postgres container for integration tests.
// Skips the test (not fail) if Docker is unreachable or the published port
// is not routable from the host (e.g. some colima networking modes).
package testpg

import (
	"context"
	"net"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Start returns a Postgres DSN backed by a fresh container.
// The container is killed on test cleanup.
func Start(t *testing.T) string {
	t.Helper()
	ctx := context.Background()
	c, err := tcpg.Run(ctx,
		"postgres:16-alpine",
		tcpg.WithDatabase("test"),
		tcpg.WithUsername("test"),
		tcpg.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		if isDockerUnavailable(err) {
			t.Skipf("docker not available: %v", err)
		}
		t.Fatalf("postgres container: %v", err)
	}
	t.Cleanup(func() { _ = c.Terminate(ctx) })
	dsn, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("dsn: %v", err)
	}
	// Some Docker hosts (e.g. colima with default networking) report the published
	// port on localhost but don't actually forward it. Probe once before handing
	// the DSN to the caller — skip rather than time out the migration step.
	if host, port := hostPortFromDSN(dsn); host != "" {
		conn, perr := net.DialTimeout("tcp", net.JoinHostPort(host, port), 2*time.Second)
		if perr != nil {
			t.Skipf("docker container port %s:%s not reachable from host (colima networking?): %v", host, port, perr)
		}
		_ = conn.Close()
	}
	return dsn
}

func hostPortFromDSN(dsn string) (host, port string) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", ""
	}
	return u.Hostname(), u.Port()
}

func isDockerUnavailable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, sub := range []string{
		"Cannot connect to the Docker daemon",
		"docker.sock",
		"rootless Docker not found",
		"no such file",
	} {
		if strings.Contains(msg, sub) {
			return true
		}
	}
	return false
}
