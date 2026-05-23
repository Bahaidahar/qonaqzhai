package domain

import "time"

// RefreshToken is a long-lived authentication token stored hashed.
// The raw value is only returned to the client once at issuance.
type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

// Active reports whether the token may still be used.
func (t *RefreshToken) Active(now time.Time) bool {
	return t.RevokedAt == nil && now.Before(t.ExpiresAt)
}

// PasswordResetToken is a one-time, short-TTL token for password recovery.
type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// Valid reports whether the reset token may still be consumed.
func (t *PasswordResetToken) Valid(now time.Time) bool {
	return t.UsedAt == nil && now.Before(t.ExpiresAt)
}
