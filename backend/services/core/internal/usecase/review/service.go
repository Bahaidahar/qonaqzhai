// Package review implements review submission.
package review

import (
	"context"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles review collaborators.
type Deps struct {
	Reviews  ports.ReviewRepo
	Bookings ports.BookingRepo
	Vendors  ports.VendorRepo
}

// Service exposes review operations.
type Service struct{ d Deps }

// New constructs a review Service.
func New(d Deps) *Service { return &Service{d: d} }

// Submit creates a review tied to a completed/paid booking owned by callerID.
func (s *Service) Submit(ctx context.Context, callerID, bookingID string, rating int, text string) (*domain.Review, error) {
	if !domain.ValidRating(rating) {
		return nil, fmt.Errorf("rating: %w", errs.ErrInvalidInput)
	}
	b, err := s.d.Bookings.Find(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID != callerID {
		return nil, errs.ErrForbidden
	}
	if b.Status != domain.BookingCompleted && b.Status != domain.BookingPaid {
		return nil, fmt.Errorf("cannot review booking in %s: %w", b.Status, errs.ErrConflict)
	}
	rv := &domain.Review{
		BookingID: b.ID, CustomerID: callerID, VendorID: b.VendorID,
		Rating: rating, Text: text,
	}
	created, err := s.d.Reviews.Create(ctx, rv)
	if err != nil {
		return nil, err
	}
	avg, count, err := s.d.Reviews.AggregateForVendor(ctx, b.VendorID)
	if err == nil {
		_ = s.d.Vendors.UpdateRating(ctx, b.VendorID, avg, count)
	}
	return created, nil
}

// ListForVendor returns the public review feed for a vendor.
func (s *Service) ListForVendor(ctx context.Context, vendorID string) ([]*domain.Review, error) {
	return s.d.Reviews.ListForVendor(ctx, vendorID)
}
