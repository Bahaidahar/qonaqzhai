package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// BookingRepo persists bookings.
type BookingRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewBookingRepo constructs a booking repository.
func NewBookingRepo(db *sql.DB, idGen ports.IDGen) *BookingRepo {
	return &BookingRepo{db: db, idGen: idGen}
}

const bookingCols = `id, customer_id, vendor_id, service_id, event_date, guest_count, note, status, amount, payment_id, created_at`

// Create inserts a booking.
func (r *BookingRepo) Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error) {
	if b.ID == "" {
		b.ID = r.idGen.New()
	}
	if b.Status == "" {
		b.Status = domain.BookingPending
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO bookings (id, customer_id, vendor_id, service_id, event_date, guest_count, note, status, amount, payment_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		b.ID, b.CustomerID, b.VendorID, b.ServiceID, b.EventDate, b.GuestCount, b.Note,
		string(b.Status), b.Amount, b.PaymentID,
	); err != nil {
		return nil, fmt.Errorf("insert booking: %w", err)
	}
	return r.Find(ctx, b.ID)
}

// Find returns a booking by id.
func (r *BookingRepo) Find(ctx context.Context, id string) (*domain.Booking, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+bookingCols+` FROM bookings WHERE id = $1`, id)
	b, err := scanBooking(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return b, err
}

// ListForCustomer returns bookings made by customer.
func (r *BookingRepo) ListForCustomer(ctx context.Context, customerID string) ([]*domain.Booking, error) {
	return r.list(ctx, `WHERE customer_id = $1`, customerID)
}

// ListForVendor returns bookings against a vendor.
func (r *BookingRepo) ListForVendor(ctx context.Context, vendorID string) ([]*domain.Booking, error) {
	return r.list(ctx, `WHERE vendor_id = $1`, vendorID)
}

// ListAll returns every booking ordered by recency.
func (r *BookingRepo) ListAll(ctx context.Context) ([]*domain.Booking, error) {
	return r.list(ctx, ``)
}

func (r *BookingRepo) list(ctx context.Context, where string, args ...any) ([]*domain.Booking, error) {
	q := `SELECT ` + bookingCols + ` FROM bookings ` + where + ` ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("list bookings: %w", err)
	}
	defer rows.Close()
	out := []*domain.Booking{}
	for rows.Next() {
		b, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// UpdateStatus changes a booking's status.
func (r *BookingRepo) UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET status = $1 WHERE id = $2`, string(status), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// SetPayment records a payment id against a booking.
func (r *BookingRepo) SetPayment(ctx context.Context, id, paymentID string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET payment_id = $1 WHERE id = $2`, paymentID, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

// Stats returns aggregate counts in a single round-trip.
func (r *BookingRepo) Stats(ctx context.Context) (ports.BookingStats, error) {
	var out ports.BookingStats
	row := r.db.QueryRowContext(ctx, `
		SELECT
		  COUNT(*) AS total,
		  COUNT(*) FILTER (WHERE status = 'pending')  AS pending,
		  COUNT(*) FILTER (WHERE status = 'accepted') AS accepted,
		  COUNT(*) FILTER (WHERE status = 'paid')     AS paid,
		  COALESCE(SUM(amount) FILTER (WHERE status = 'paid'), 0) AS gmv
		FROM bookings`)
	if err := row.Scan(&out.Total, &out.Pending, &out.Accepted, &out.Paid, &out.GMV); err != nil {
		return out, fmt.Errorf("stats: %w", err)
	}
	return out, nil
}

func scanBooking(s scanner) (*domain.Booking, error) {
	var b domain.Booking
	var status string
	if err := s.Scan(
		&b.ID, &b.CustomerID, &b.VendorID, &b.ServiceID, &b.EventDate,
		&b.GuestCount, &b.Note, &status, &b.Amount, &b.PaymentID, &b.CreatedAt,
	); err != nil {
		return nil, err
	}
	b.Status = domain.BookingStatus(status)
	return &b, nil
}

var _ ports.BookingRepo = (*BookingRepo)(nil)
