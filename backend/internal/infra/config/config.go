// Package config loads environment-driven configuration.
package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config is the resolved runtime configuration.
type Config struct {
	Addr         string
	DatabaseURL  string
	CORSOrigin   string
	JWTSecret    string
	AccessTTL    time.Duration
	RefreshTTL   time.Duration
	ResetTTL     time.Duration
	GeminiAPIKey string
	GeminiModel  string

	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	FCMProjectID         string
	FCMServiceAccountKey string

	PayBoxMerchantID string
	PayBoxSecretKey  string
	PayBoxSandbox    bool

	BaseURL string

	RateLimitDisabled bool
	BcryptCost        int // 0 → bcrypt.DefaultCost; set lower (4) in tests for speed
}

// Load reads configuration from environment variables, providing sensible defaults.
// A fresh ephemeral JWT secret is generated if none is supplied.
func Load() (Config, error) {
	c := Config{
		Addr:         envOr("ADDR", ":8080"),
		DatabaseURL:  envOr("DATABASE_URL", "postgres://qonaqzhai:qonaqzhai@localhost:5433/qonaqzhai?sslmode=disable"),
		CORSOrigin:   envOr("CORS_ORIGIN", "http://localhost:3000"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		AccessTTL:    durOr("JWT_ACCESS_TTL", 15*time.Minute),
		RefreshTTL:   durOr("JWT_REFRESH_TTL", 30*24*time.Hour),
		ResetTTL:     durOr("PASSWORD_RESET_TTL", time.Hour),
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		GeminiModel:  envOr("GEMINI_MODEL", "gemini-2.0-flash"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     envOr("SMTP_PORT", "587"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		// Gmail app passwords are displayed with spaces; normalise just in case.
		SMTPPassword: strings.ReplaceAll(firstNonEmpty(os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_PASS")), " ", ""),
		SMTPFrom:     firstNonEmpty(os.Getenv("SMTP_FROM"), os.Getenv("FROM_EMAIL"), os.Getenv("SMTP_USER")),

		FCMProjectID:         os.Getenv("FCM_PROJECT_ID"),
		FCMServiceAccountKey: os.Getenv("FCM_SERVICE_ACCOUNT_KEY"),

		PayBoxMerchantID: os.Getenv("PAYBOX_MERCHANT_ID"),
		PayBoxSecretKey:  os.Getenv("PAYBOX_SECRET_KEY"),
		PayBoxSandbox:    envOr("PAYBOX_SANDBOX", "true") == "true",

		BaseURL: envOr("BASE_URL", "http://localhost:8080"),

		RateLimitDisabled: envOr("RATE_LIMIT_DISABLED", "false") == "true",
	}
	if c.JWTSecret == "" {
		secret, err := randomHex(32)
		if err != nil {
			return c, fmt.Errorf("generate jwt secret: %w", err)
		}
		c.JWTSecret = secret
	}
	return c, nil
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func durOr(key string, def time.Duration) time.Duration {
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

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
