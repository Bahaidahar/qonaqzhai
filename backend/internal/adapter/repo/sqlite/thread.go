package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// ThreadRepo persists booking DM threads.
type ThreadRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewThreadRepo constructs a thread repository.
func NewThreadRepo(db *sql.DB, idGen usecase.IDGen) *ThreadRepo {
	return &ThreadRepo{db: db, idGen: idGen}
}

// CreateForBooking inserts a thread row for the given booking (idempotent — re-use existing).
func (r *ThreadRepo) CreateForBooking(ctx context.Context, bookingID, customerID, vendorID string) (*domain.BookingThread, error) {
	if existing, err := r.FindByBooking(ctx, bookingID); err == nil {
		return existing, nil
	} else if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	id := r.idGen.New()
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO booking_threads (id, booking_id, customer_id, vendor_id) VALUES (?, ?, ?, ?)`,
		id, bookingID, customerID, vendorID,
	); err != nil {
		return nil, fmt.Errorf("insert thread: %w", err)
	}
	return r.FindByID(ctx, id)
}

// FindByID looks up a thread by primary key.
func (r *ThreadRepo) FindByID(ctx context.Context, id string) (*domain.BookingThread, error) {
	return r.queryThread(ctx, `WHERE id = ?`, id)
}

// FindByBooking looks up the thread attached to a booking.
func (r *ThreadRepo) FindByBooking(ctx context.Context, bookingID string) (*domain.BookingThread, error) {
	return r.queryThread(ctx, `WHERE booking_id = ?`, bookingID)
}

// ListForUser returns threads where the user is either customer or vendor, newest first.
func (r *ThreadRepo) ListForUser(ctx context.Context, userID string) ([]*domain.BookingThread, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, booking_id, customer_id, vendor_id, created_at, updated_at
		 FROM booking_threads WHERE customer_id = ? OR vendor_id = ?
		 ORDER BY updated_at DESC`,
		userID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list threads: %w", err)
	}
	defer rows.Close()
	out := []*domain.BookingThread{}
	for rows.Next() {
		var t domain.BookingThread
		if err := rows.Scan(&t.ID, &t.BookingID, &t.CustomerID, &t.VendorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

// Touch bumps updated_at on the thread (called on each new message).
func (r *ThreadRepo) Touch(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE booking_threads SET updated_at = CURRENT_TIMESTAMP WHERE id = ?`, id,
	)
	return err
}

// AddMessage inserts a message and bumps the thread's updated_at.
func (r *ThreadRepo) AddMessage(ctx context.Context, m *domain.ThreadMessage) (*domain.ThreadMessage, error) {
	if m.ID == "" {
		m.ID = r.idGen.New()
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO thread_messages (id, thread_id, sender_id, text) VALUES (?, ?, ?, ?)`,
		m.ID, m.ThreadID, m.SenderID, m.Text,
	); err != nil {
		return nil, fmt.Errorf("insert thread message: %w", err)
	}
	_ = r.Touch(ctx, m.ThreadID)
	row := r.db.QueryRowContext(ctx,
		`SELECT id, thread_id, sender_id, text, created_at FROM thread_messages WHERE id = ?`, m.ID,
	)
	var out domain.ThreadMessage
	if err := row.Scan(&out.ID, &out.ThreadID, &out.SenderID, &out.Text, &out.CreatedAt); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListMessages returns ordered messages of a thread.
func (r *ThreadRepo) ListMessages(ctx context.Context, threadID string) ([]*domain.ThreadMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, thread_id, sender_id, text, created_at FROM thread_messages
		 WHERE thread_id = ? ORDER BY created_at ASC`, threadID,
	)
	if err != nil {
		return nil, fmt.Errorf("list thread messages: %w", err)
	}
	defer rows.Close()
	out := []*domain.ThreadMessage{}
	for rows.Next() {
		var m domain.ThreadMessage
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.SenderID, &m.Text, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *ThreadRepo) queryThread(ctx context.Context, where string, args ...any) (*domain.BookingThread, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, booking_id, customer_id, vendor_id, created_at, updated_at FROM booking_threads `+where,
		args...,
	)
	var t domain.BookingThread
	if err := row.Scan(&t.ID, &t.BookingID, &t.CustomerID, &t.VendorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

var _ usecase.ThreadRepo = (*ThreadRepo)(nil)
