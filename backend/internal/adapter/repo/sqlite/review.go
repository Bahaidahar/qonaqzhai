package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// ReviewRepo persists vendor reviews.
type ReviewRepo struct {
	db    *sql.DB
	idGen usecase.IDGen
}

// NewReviewRepo constructs a review repository.
func NewReviewRepo(db *sql.DB, idGen usecase.IDGen) *ReviewRepo {
	return &ReviewRepo{db: db, idGen: idGen}
}

// Create inserts a review. Unique constraint on booking_id enforces 1 review per booking.
func (r *ReviewRepo) Create(ctx context.Context, in *domain.Review) (*domain.Review, error) {
	if in.ID == "" {
		in.ID = r.idGen.New()
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO reviews (id, booking_id, customer_id, vendor_id, rating, text) VALUES ($1, $2, $3, $4, $5, $6)`,
		in.ID, in.BookingID, in.CustomerID, in.VendorID, in.Rating, in.Text,
	); err != nil {
		if isUniqueErr(err) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("insert review: %w", err)
	}
	return r.FindByID(ctx, in.ID)
}

// FindByID returns a review by id.
func (r *ReviewRepo) FindByID(ctx context.Context, id string) (*domain.Review, error) {
	return r.queryReview(ctx, `WHERE id = $1`, id)
}

// FindByBooking returns the (single) review attached to a booking.
func (r *ReviewRepo) FindByBooking(ctx context.Context, bookingID string) (*domain.Review, error) {
	return r.queryReview(ctx, `WHERE booking_id = $1`, bookingID)
}

// ListByVendor returns reviews for a vendor, newest first.
func (r *ReviewRepo) ListByVendor(ctx context.Context, vendorID string) ([]*domain.Review, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, booking_id, customer_id, vendor_id, rating, text, created_at FROM reviews WHERE vendor_id = $1 ORDER BY created_at DESC`,
		vendorID,
	)
	if err != nil {
		return nil, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()
	out := []*domain.Review{}
	for rows.Next() {
		rv, err := scanReview(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, rv)
	}
	return out, rows.Err()
}

// Delete removes a review by id.
func (r *ReviewRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reviews WHERE id = $1`, id)
	return err
}

// AggregateForVendor returns (average rating, count) for vendor.
func (r *ReviewRepo) AggregateForVendor(ctx context.Context, vendorID string) (float64, int, error) {
	var avg sql.NullFloat64
	var count int
	if err := r.db.QueryRowContext(ctx,
		`SELECT AVG(rating), COUNT(*) FROM reviews WHERE vendor_id = $1`, vendorID,
	).Scan(&avg, &count); err != nil {
		return 0, 0, err
	}
	if !avg.Valid {
		return 0, 0, nil
	}
	return avg.Float64, count, nil
}

func (r *ReviewRepo) queryReview(ctx context.Context, where string, args ...any) (*domain.Review, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, booking_id, customer_id, vendor_id, rating, text, created_at FROM reviews `+where,
		args...,
	)
	rv, err := scanReview(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return rv, err
}

func scanReview(s scanner) (*domain.Review, error) {
	var rv domain.Review
	if err := s.Scan(
		&rv.ID, &rv.BookingID, &rv.CustomerID, &rv.VendorID,
		&rv.Rating, &rv.Text, &rv.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &rv, nil
}

var _ usecase.ReviewRepo = (*ReviewRepo)(nil)
