// Package ports defines the interfaces the auth service usecase layer depends
// on. Adapter packages provide the concrete implementations.
package ports

import (
	"context"
	"time"

	"qonaqzhai-backend/services/auth/internal/domain"
)

// UserRepo persists users.
type UserRepo interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, opts ListUsersOpts) ([]*domain.User, error)
	FindByIDs(ctx context.Context, ids []string) ([]*domain.User, error)
	UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error
	UpdatePasswordHash(ctx context.Context, id, hash string) error
}

// ListUsersOpts are pagination + filter parameters for List.
type ListUsersOpts struct {
	Limit  int
	Offset int
	Role   string
}

// RefreshTokenRepo persists hashed refresh tokens.
type RefreshTokenRepo interface {
	Create(ctx context.Context, t *domain.RefreshToken) error
	FindActiveByHash(ctx context.Context, hash string, now time.Time) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id string, at time.Time) error
	RevokeAllForUser(ctx context.Context, userID string, at time.Time) error
}

// PasswordResetRepo persists hashed password reset tokens.
type PasswordResetRepo interface {
	Create(ctx context.Context, t *domain.PasswordResetToken) error
	FindByHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string, at time.Time) error
}

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(hash, plain string) error
}

// Clock abstracts the current time so tests can pin it.
type Clock interface{ Now() time.Time }

// IDGen generates new opaque entity IDs.
type IDGen interface{ New() string }

// Mailer delivers transactional email. nil-tolerant — services skip delivery
// when no mailer is configured.
type Mailer interface {
	Send(ctx context.Context, to, subject, htmlBody string) error
}
