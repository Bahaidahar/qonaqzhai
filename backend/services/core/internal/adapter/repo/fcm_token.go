package repo

import (
	"context"
	"database/sql"
	"fmt"

	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// FCMTokenRepo persists device tokens.
type FCMTokenRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewFCMTokenRepo constructs an FCM token repository.
func NewFCMTokenRepo(db *sql.DB, idGen ports.IDGen) *FCMTokenRepo {
	return &FCMTokenRepo{db: db, idGen: idGen}
}

// Upsert inserts or refreshes the user binding for an existing token.
func (r *FCMTokenRepo) Upsert(ctx context.Context, t *domain.FCMToken) error {
	if t.ID == "" {
		t.ID = r.idGen.New()
	}
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO fcm_tokens (id, user_id, token, platform)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform`,
		t.ID, t.UserID, t.Token, t.Platform,
	)
	if err != nil {
		return fmt.Errorf("upsert fcm: %w", err)
	}
	return nil
}

// Delete removes a token.
func (r *FCMTokenRepo) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM fcm_tokens WHERE token = $1`, token)
	return err
}

// ListByUsers returns tokens for the given user ids.
func (r *FCMTokenRepo) ListByUsers(ctx context.Context, userIDs []string) ([]*domain.FCMToken, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, token, platform, created_at FROM fcm_tokens WHERE user_id = ANY($1)`,
		userIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("list fcm: %w", err)
	}
	defer rows.Close()
	out := []*domain.FCMToken{}
	for rows.Next() {
		var t domain.FCMToken
		if err := rows.Scan(&t.ID, &t.UserID, &t.Token, &t.Platform, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

var _ ports.FCMTokenRepo = (*FCMTokenRepo)(nil)
