package repo

import (
	"context"
	"database/sql"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// NotificationRepo persists in-app notifications.
type NotificationRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewNotificationRepo constructs a notification repository.
func NewNotificationRepo(db *sql.DB, idGen ports.IDGen) *NotificationRepo {
	return &NotificationRepo{db: db, idGen: idGen}
}

// Enqueue records a queued notification row.
func (r *NotificationRepo) Enqueue(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	if n.ID == "" {
		n.ID = r.idGen.New()
	}
	if n.Status == "" {
		n.Status = "queued"
	}
	if n.Channel == "" {
		n.Channel = domain.ChannelPush
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO notifications (id, user_id, type, channel, title, body, status)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		n.ID, n.UserID, n.Type, string(n.Channel), n.Title, n.Body, n.Status,
	); err != nil {
		return nil, fmt.Errorf("insert notification: %w", err)
	}
	return n, nil
}

// ListForUser returns the latest limit notifications for userID.
func (r *NotificationRepo) ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, type, channel, title, body, status, created_at
		 FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	out := []*domain.Notification{}
	for rows.Next() {
		var n domain.Notification
		var ch string
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &ch, &n.Title, &n.Body, &n.Status, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.Channel = domain.NotificationChannel(ch)
		out = append(out, &n)
	}
	return out, rows.Err()
}

// MarkSent flips a queued notification to sent.
func (r *NotificationRepo) MarkSent(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE notifications SET status = 'sent' WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

var _ ports.NotificationRepo = (*NotificationRepo)(nil)
