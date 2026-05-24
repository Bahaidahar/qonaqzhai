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

// ReviewRepo persists reviews.
type ReviewRepo struct {
	db    *sql.DB
	idGen ports.IDGen
}

// NewReviewRepo constructs a review repository.
func NewReviewRepo(db *sql.DB, idGen ports.IDGen) *ReviewRepo {
	return &ReviewRepo{db: db, idGen: idGen}
}

const reviewCols = `id, booking_id, customer_id, vendor_id, rating, text, created_at`

// Create inserts a review. Unique-violation on booking_id surfaces as ErrAlreadyExists.
func (r *ReviewRepo) Create(ctx context.Context, in *domain.Review) (*domain.Review, error) {
	if in.ID == "" {
		in.ID = r.idGen.New()
	}
	if _, err := r.db.ExecContext(ctx,
		`INSERT INTO reviews (id, booking_id, customer_id, vendor_id, rating, text)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		in.ID, in.BookingID, in.CustomerID, in.VendorID, in.Rating, in.Text,
	); err != nil {
		if isUniqueErr(err) {
			return nil, errs.ErrAlreadyExists
		}
		return nil, fmt.Errorf("insert review: %w", err)
	}
	return r.find(ctx, in.ID)
}

// ListForVendor returns paginated reviews for vendorID, newest first.
func (r *ReviewRepo) ListForVendor(ctx context.Context, vendorID string, p ports.Page) ([]*domain.Review, error) {
	p = p.Clamp()
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+reviewCols+` FROM reviews WHERE vendor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		vendorID, p.Limit, p.Offset,
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

// FindByBooking returns the review attached to a booking, if any.
func (r *ReviewRepo) FindByBooking(ctx context.Context, bookingID string) (*domain.Review, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+reviewCols+` FROM reviews WHERE booking_id = $1`, bookingID,
	)
	rv, err := scanReview(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return rv, err
}

// AggregateForVendor returns the average rating + count for vendorID.
func (r *ReviewRepo) AggregateForVendor(ctx context.Context, vendorID string) (float64, int, error) {
	var avg sql.NullFloat64
	var count int
	if err := r.db.QueryRowContext(ctx,
		`SELECT AVG(rating)::float, COUNT(*) FROM reviews WHERE vendor_id = $1`,
		vendorID,
	).Scan(&avg, &count); err != nil {
		return 0, 0, fmt.Errorf("agg reviews: %w", err)
	}
	if !avg.Valid {
		return 0, count, nil
	}
	return avg.Float64, count, nil
}

func (r *ReviewRepo) find(ctx context.Context, id string) (*domain.Review, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+reviewCols+` FROM reviews WHERE id = $1`, id)
	rv, err := scanReview(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	return rv, err
}

func scanReview(s scanner) (*domain.Review, error) {
	var rv domain.Review
	if err := s.Scan(&rv.ID, &rv.BookingID, &rv.CustomerID, &rv.VendorID, &rv.Rating, &rv.Text, &rv.CreatedAt); err != nil {
		return nil, err
	}
	return &rv, nil
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

var _ ports.ReviewRepo = (*ReviewRepo)(nil)
