package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/realtime/internal/domain"
	"qonaqzhai-backend/services/realtime/internal/ports"
)

// ThreadRepo persists DM threads and their messages.
type ThreadRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewThreadRepo constructs a thread repository.
func NewThreadRepo(db *sql.DB, idGen ports.IDGen) *ThreadRepo {
	return &ThreadRepo{db: db, idGen: idGen}
}

const threadCols = `id, booking_id, customer_id, vendor_id, created_at, updated_at`

// EnsureForBooking returns the existing thread for the booking or creates one.
func (r *ThreadRepo) EnsureForBooking(ctx context.Context, bookingID, customerID, vendorID string) (*domain.Thread, error) {
	existing, err := r.findByBooking(ctx, bookingID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return nil, err
	}
	id := r.idGen.New()
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO threads (id, booking_id, customer_id, vendor_id) VALUES ($1, $2, $3, $4)`,
		id, bookingID, customerID, vendorID,
	); err != nil {
		return nil, fmt.Errorf("insert thread: %w", err)
	}
	return r.FindByID(ctx, id)
}

// FindByID returns a thread by primary key.
func (r *ThreadRepo) FindByID(ctx context.Context, id string) (*domain.Thread, error) {
	return r.queryThread(ctx, `WHERE id = $1`, id)
}

func (r *ThreadRepo) findByBooking(ctx context.Context, bookingID string) (*domain.Thread, error) {
	return r.queryThread(ctx, `WHERE booking_id = $1`, bookingID)
}

// ListForUser returns threads where the user is either customer or vendor.
func (r *ThreadRepo) ListForUser(ctx context.Context, userID string) ([]*domain.Thread, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+threadCols+` FROM threads WHERE customer_id = $1 OR vendor_id = $2 ORDER BY updated_at DESC`,
		userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list threads: %w", err)
	}
	defer rows.Close()
	out := []*domain.Thread{}
	for rows.Next() {
		t, err := scanThread(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

// AddMessage inserts a new message and bumps the parent thread's updated_at.
func (r *ThreadRepo) AddMessage(ctx context.Context, m *domain.Message) (*domain.Message, error) {
	if m.ID == "" {
		m.ID = r.idGen.New()
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO thread_messages (id, thread_id, sender_id, text) VALUES ($1, $2, $3, $4)`,
		m.ID, m.ThreadID, m.SenderID, m.Text,
	); err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE threads SET updated_at = now() WHERE id = $1`, m.ThreadID); err != nil {
		return nil, fmt.Errorf("touch thread: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	row := r.db.QueryRowContext(ctx,
		`SELECT id, thread_id, sender_id, text, created_at FROM thread_messages WHERE id = $1`,
		m.ID,
	)
	var out domain.Message
	if err := row.Scan(&out.ID, &out.ThreadID, &out.SenderID, &out.Text, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMessages returns messages in a thread in chronological order.
func (r *ThreadRepo) ListMessages(ctx context.Context, threadID string) ([]*domain.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, thread_id, sender_id, text, created_at FROM thread_messages WHERE thread_id = $1 ORDER BY created_at ASC`,
		threadID,
	)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()
	out := []*domain.Message{}
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.SenderID, &m.Text, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *ThreadRepo) queryThread(ctx context.Context, where string, args ...any) (*domain.Thread, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+threadCols+` FROM threads `+where, args...)
	t, err := scanThread(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return t, err
}

func scanThread(s scanner) (*domain.Thread, error) {
	var t domain.Thread
	if err := s.Scan(&t.ID, &t.BookingID, &t.CustomerID, &t.VendorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

var _ ports.ThreadRepo = (*ThreadRepo)(nil)
