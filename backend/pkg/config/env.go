// Package config holds small environment-variable helpers used by service main
// packages. Service-specific config structs live in each service.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// EnvOr returns the value of key or def when unset.
func EnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// MustEnv returns the value of key or panics if unset. Use only for secrets
// that must not silently fall back to a default in production startup.
func MustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %s is empty", key))
	}
	return v
}

// DurationEnv parses an env var as a Go duration (e.g. "30m", "24h") or
// returns def when unset / malformed.
func DurationEnv(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}

// BoolEnv reports whether the env var is one of "1", "true", "yes".
func BoolEnv(key string, def bool) bool {
	switch os.Getenv(key) {
	case "1", "true", "yes":
		return true
	case "0", "false", "no":
		return false
	}
	return def
}

// FirstNonEmpty returns the first non-empty string from values.
func FirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// RandomHex returns a hex-encoded random byte string of length 2*n.
func RandomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("random: %w", err)
	}
	return hex.EncodeToString(b), nil
}
