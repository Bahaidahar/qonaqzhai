package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// NotificationRepo persists in-app notifications.
type NotificationRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewNotificationRepo constructs a notification repository.
func NewNotificationRepo(db *sql.DB, idGen usecase.IDGen) *NotificationRepo {
	return &NotificationRepo{db: db, idGen: idGen}
}

// Create inserts a notification row.
func (r *NotificationRepo) Create(ctx context.Context, n *domain.Notification) (*domain.Notification, error) {
	if n.ID == "" {
		n.ID = r.idGen.New()
	}
	if n.Status == "" {
		n.Status = "queued"
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO notifications (id, user_id, type, channel, title, body, status) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		n.ID, n.UserID, string(n.Type), string(n.Channel), n.Title, n.Body, n.Status,
	); err != nil {
		return nil, fmt.Errorf("insert notification: %w", err)
	}
	return r.findByID(ctx, n.ID)
}

// ListForUser returns recent notifications for user.
func (r *NotificationRepo) ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, type, channel, title, body, status, created_at FROM notifications
		 WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	out := []*domain.Notification{}
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

// MarkSent records successful delivery.
func (r *NotificationRepo) MarkSent(ctx context.Context, id string) error {
	return r.setStatus(ctx, id, "sent")
}

// MarkFailed records failed delivery.
func (r *NotificationRepo) MarkFailed(ctx context.Context, id string) error {
	return r.setStatus(ctx, id, "failed")
}

func (r *NotificationRepo) setStatus(ctx context.Context, id, status string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE notifications SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *NotificationRepo) findByID(ctx context.Context, id string) (*domain.Notification, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, type, channel, title, body, status, created_at FROM notifications WHERE id = $1`,
		id,
	)
	n, err := scanNotification(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return n, err
}

func scanNotification(s scanner) (*domain.Notification, error) {
	var n domain.Notification
	var typ, ch string
	if err := s.Scan(&n.ID, &n.UserID, &typ, &ch, &n.Title, &n.Body, &n.Status, &n.CreatedAt); err != nil {
		return nil, err
	}
	n.Type = domain.NotificationType(typ)
	n.Channel = domain.NotificationChannel(ch)
	return &n, nil
}

var _ usecase.NotificationRepo = (*NotificationRepo)(nil)
