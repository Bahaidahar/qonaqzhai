package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// RefreshTokenRepo persists hashed refresh tokens.
type RefreshTokenRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewRefreshTokenRepo constructs a refresh token repository.
func NewRefreshTokenRepo(db *sql.DB, idGen usecase.IDGen) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db, idGen: idGen}
}

// Create inserts a new refresh token row.
func (r *RefreshTokenRepo) Create(ctx context.Context, t *domain.RefreshToken) error {
	if t.ID == "" {
		t.ID = r.idGen.New()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at) VALUES ($1, $2, $3, $4)`,
		t.ID, t.UserID, t.TokenHash, t.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// FindActiveByHash returns the active refresh token matching hash at the given clock time.
func (r *RefreshTokenRepo) FindActiveByHash(ctx context.Context, hash string, now time.Time) (*domain.RefreshToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		 FROM refresh_tokens
		 WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2`,
		hash, now,
	)
	var t domain.RefreshToken
	var revoked sql.NullTime
	if err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &revoked, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	if revoked.Valid {
		t.RevokedAt = &revoked.Time
	}
	return &t, nil
}

// Revoke marks a single refresh token as revoked.
func (r *RefreshTokenRepo) Revoke(ctx context.Context, id string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2`, at, id)
	return err
}

// RevokeAllForUser revokes every active refresh token belonging to user.
func (r *RefreshTokenRepo) RevokeAllForUser(ctx context.Context, userID string, at time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`,
		at, userID,
	)
	return err
}

var _ usecase.RefreshTokenRepo = (*RefreshTokenRepo)(nil)
