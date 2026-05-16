package domain

import (
	"strings"
	"time"
)

// UserStatus is the lifecycle state of a user account.
type UserStatus string

const (
	UserActive    UserStatus = "active"
	UserSuspended UserStatus = "suspended"
)

// User is the authenticated principal.
type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	Role         Role       `json:"role"`
	Status       UserStatus `json:"status"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"createdAt"`
}

// IsActive reports whether the user may authenticate.
func (u *User) IsActive() bool { return u.Status == UserActive }

// NormalizeEmail returns the canonical email form: lowercase, trimmed.
func NormalizeEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ValidEmail performs lightweight email shape validation
// sufficient for signup-time rejection. Full deliverability is out of scope.
func ValidEmail(s string) bool {
	if s == "" || !strings.Contains(s, "@") {
		return false
	}
	parts := strings.SplitN(s, "@", 2)
	return parts[0] != "" && strings.Contains(parts[1], ".")
}

// MinPasswordLength is the minimum acceptable password length at signup.
const MinPasswordLength = 8

// ValidPassword checks length only — strength scoring belongs elsewhere.
func ValidPassword(p string) bool { return len(p) >= MinPasswordLength }

// DefaultName derives a display name from the email local part when none provided.
func DefaultName(name, email string) string {
	name = strings.TrimSpace(name)
	if name != "" {
		return name
	}
	if !strings.Contains(email, "@") {
		return email
	}
	return strings.SplitN(email, "@", 2)[0]
}
