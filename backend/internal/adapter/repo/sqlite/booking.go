package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// BookingRepo persists bookings.
type BookingRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewBookingRepo constructs a booking repository.
func NewBookingRepo(db *sql.DB, idGen usecase.IDGen) *BookingRepo {
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
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.CustomerID, b.VendorID, b.ServiceID, b.EventDate, b.GuestCount, b.Note,
		string(b.Status), b.Amount, b.PaymentID,
	); err != nil {
		return nil, fmt.Errorf("insert booking: %w", err)
	}
	return r.Find(ctx, b.ID)
}

// Find returns a booking by id.
func (r *BookingRepo) Find(ctx context.Context, id string) (*domain.Booking, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+bookingCols+` FROM bookings WHERE id = ?`,
		id,
	)
	b, err := scanBooking(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return b, err
}

// ListForCustomer returns bookings made by customer.
func (r *BookingRepo) ListForCustomer(ctx context.Context, customerID string) ([]*domain.Booking, error) {
	return r.list(ctx, `WHERE customer_id = ?`, customerID)
}

// ListForVendor returns bookings against vendor (by vendor row id, not user id).
func (r *BookingRepo) ListForVendor(ctx context.Context, vendorID string) ([]*domain.Booking, error) {
	return r.list(ctx, `WHERE vendor_id = ?`, vendorID)
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
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET status = ? WHERE id = ?`, string(status), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SetPayment associates a payment id with a booking (called from payment webhook).
func (r *BookingRepo) SetPayment(ctx context.Context, id, paymentID string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET payment_id = ? WHERE id = ?`, paymentID, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SetService associates a service id with a booking.
func (r *BookingRepo) SetService(ctx context.Context, id, serviceID string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET service_id = ? WHERE id = ?`, serviceID, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
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

// Ensure compile-time conformance with the BookingRepo port.
var _ usecase.BookingRepo = (*BookingRepo)(nil)
