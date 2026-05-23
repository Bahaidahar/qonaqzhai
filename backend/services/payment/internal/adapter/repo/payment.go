package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/payment/internal/domain"
	"qonaqzhai-backend/services/payment/internal/ports"
)

// PaymentRepo persists payment attempts.
type PaymentRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewPaymentRepo constructs a payment repository.
func NewPaymentRepo(db *sql.DB, idGen ports.IDGen) *PaymentRepo {
	return &PaymentRepo{db: db, idGen: idGen}
}

const paymentCols = `id, booking_id, user_id, card_id, amount, currency, status, provider_ref, created_at`

// Create inserts a payment attempt.
func (r *PaymentRepo) Create(ctx context.Context, p *domain.Payment) (*domain.Payment, error) {
	if p.ID == "" {
		p.ID = r.idGen.New()
	}
	if p.Status == "" {
		p.Status = domain.PaymentPending
	}
	if p.Currency == "" {
		p.Currency = "KZT"
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (id, booking_id, user_id, card_id, amount, currency, status, provider_ref)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		p.ID, p.BookingID, p.UserID, p.CardID, p.Amount, p.Currency, string(p.Status), p.ProviderRef,
	); err != nil {
		if isUniqueErr(err) {
			return nil, errs.ErrAlreadyExists
		}
		return nil, fmt.Errorf("insert payment: %w", err)
	}
	return r.Find(ctx, p.ID)
}

// Find returns a payment by id.
func (r *PaymentRepo) Find(ctx context.Context, id string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+paymentCols+` FROM payments WHERE id = $1`, id)
	p, err := scanPayment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return p, err
}

// FindByBooking returns the payment attached to a booking, if any.
func (r *PaymentRepo) FindByBooking(ctx context.Context, bookingID string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+paymentCols+` FROM payments WHERE booking_id = $1`, bookingID)
	p, err := scanPayment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return p, err
}

// ListByUser returns the user's payments newest-first.
func (r *PaymentRepo) ListByUser(ctx context.Context, userID string) ([]*domain.Payment, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+paymentCols+` FROM payments WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list payments: %w", err)
	}
	defer rows.Close()
	out := []*domain.Payment{}
	for rows.Next() {
		p, err := scanPayment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// UpdateStatus flips a payment to a new status.
func (r *PaymentRepo) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus) error {
	res, err := r.db.ExecContext(ctx, `UPDATE payments SET status = $1 WHERE id = $2`, string(status), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errs.ErrNotFound
	}
	return nil
}

func scanPayment(s scanner) (*domain.Payment, error) {
	var p domain.Payment
	var status string
	if err := s.Scan(&p.ID, &p.BookingID, &p.UserID, &p.CardID, &p.Amount, &p.Currency, &status, &p.ProviderRef, &p.CreatedAt); err != nil {
		return nil, err
	}
	p.Status = domain.PaymentStatus(status)
	return &p, nil
}

func isUniqueErr(err error) bool {
	if err == nil {
		return false
	}
	return contains(err.Error(), "23505")
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

var _ ports.PaymentRepo = (*PaymentRepo)(nil)
