package domain

import (
	"testing"
	"time"
)

func TestRefreshTokenActive(t *testing.T) {
	t.Parallel()
	now := time.Now()
	revoked := now.Add(-time.Hour)

	cases := []struct {
		name string
		tok  RefreshToken
		now  time.Time
		want bool
	}{
		{"fresh", RefreshToken{ExpiresAt: now.Add(time.Hour)}, now, true},
		{"expired", RefreshToken{ExpiresAt: now.Add(-time.Minute)}, now, false},
		{"revoked", RefreshToken{ExpiresAt: now.Add(time.Hour), RevokedAt: &revoked}, now, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.tok.Active(c.now); got != c.want {
				t.Errorf("Active=%v want %v", got, c.want)
			}
		})
	}
}

func TestPasswordResetTokenValid(t *testing.T) {
	t.Parallel()
	now := time.Now()
	used := now.Add(-time.Minute)

	cases := []struct {
		name string
		tok  PasswordResetToken
		want bool
	}{
		{"fresh", PasswordResetToken{ExpiresAt: now.Add(time.Hour)}, true},
		{"expired", PasswordResetToken{ExpiresAt: now.Add(-time.Minute)}, false},
		{"used", PasswordResetToken{ExpiresAt: now.Add(time.Hour), UsedAt: &used}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.tok.Valid(now); got != c.want {
				t.Errorf("Valid=%v want %v", got, c.want)
			}
		})
	}
}
