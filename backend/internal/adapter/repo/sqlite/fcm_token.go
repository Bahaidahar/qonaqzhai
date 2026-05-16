package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"qonaqzhai-backend/internal/usecase/notification"
)

// FCMTokenRepo persists device push tokens.
type FCMTokenRepo struct {
	db    *sql.DB
	idGen interface{ New() string }
}

// NewFCMTokenRepo constructs an FCM token repository.
func NewFCMTokenRepo(db *sql.DB, idGen interface{ New() string }) *FCMTokenRepo {
	return &FCMTokenRepo{db: db, idGen: idGen}
}

// Register adds (or replaces) a token for user. Tokens are globally unique by FCM design.
func (r *FCMTokenRepo) Register(ctx context.Context, userID, token, platform string) error {
	// idempotent insert: try insert, fall back to update on unique conflict
	id := r.idGen.New()
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO fcm_tokens (id, user_id, token, platform) VALUES (?, ?, ?, ?)`,
		id, userID, token, platform,
	)
	if err == nil {
		return nil
	}
	if !isUniqueErr(err) {
		return fmt.Errorf("insert fcm token: %w", err)
	}
	_, err = r.db.ExecContext(ctx,
		`UPDATE fcm_tokens SET user_id = ?, platform = ? WHERE token = ?`,
		userID, platform, token,
	)
	if err != nil {
		return fmt.Errorf("update fcm token: %w", err)
	}
	return nil
}

// Unregister removes a token (e.g., after a 404 from FCM).
func (r *FCMTokenRepo) Unregister(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fcm_tokens WHERE token = ?`, token)
	return err
}

// TokensForUser returns all active FCM tokens registered for user.
func (r *FCMTokenRepo) TokensForUser(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT token FROM fcm_tokens WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

var _ notification.FCMTokenRepo = (*FCMTokenRepo)(nil)
