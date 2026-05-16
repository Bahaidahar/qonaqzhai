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

// PasswordResetRepo persists hashed password reset tokens.
type PasswordResetRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewPasswordResetRepo constructs a password reset repository.
func NewPasswordResetRepo(db *sql.DB, idGen usecase.IDGen) *PasswordResetRepo {
	return &PasswordResetRepo{db: db, idGen: idGen}
}

// Create inserts a new reset token row.
func (r *PasswordResetRepo) Create(ctx context.Context, t *domain.PasswordResetToken) error {
	if t.ID == "" {
		t.ID = r.idGen.New()
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at) VALUES (?, ?, ?, ?)`,
		t.ID, t.UserID, t.TokenHash, t.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert password reset: %w", err)
	}
	return nil
}

// FindByHash returns the token row matching hash. Validity is checked by caller.
func (r *PasswordResetRepo) FindByHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, used_at, created_at
		 FROM password_reset_tokens WHERE token_hash = ?`, hash,
	)
	var t domain.PasswordResetToken
	var used sql.NullTime
	if err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &used, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	if used.Valid {
		t.UsedAt = &used.Time
	}
	return &t, nil
}

// MarkUsed records the reset token as consumed.
func (r *PasswordResetRepo) MarkUsed(ctx context.Context, id string, at time.Time) error {
	_, err := r.db.ExecContext(ctx, `UPDATE password_reset_tokens SET used_at = ? WHERE id = ?`, at, id)
	return err
}

var _ usecase.PasswordResetRepo = (*PasswordResetRepo)(nil)
