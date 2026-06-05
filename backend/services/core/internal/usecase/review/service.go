// Package review implements review submission.
package review

import (
	"context"
	"fmt"
	"log/slog"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles review collaborators.
type Deps struct {
	Reviews  ports.ReviewRepo
	Bookings ports.BookingRepo
	Vendors  ports.VendorRepo
	Logger   *slog.Logger
}

// Service exposes review operations.
type Service struct{ d Deps }

// New constructs a review Service. A nil logger falls back to slog.Default().
func New(d Deps) *Service {
	if d.Logger == nil {
		d.Logger = slog.Default()
	}
	return &Service{d: d}
}

// Submit creates a review tied to a completed/paid booking owned by callerID.
// Rating aggregation failures are logged but do not roll back the review — the
// next successful Submit will recompute.
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
	if err != nil {
		s.d.Logger.Warn("rating aggregate failed", "vendor", b.VendorID, "err", err)
		return created, nil
	}
	if err := s.d.Vendors.UpdateRating(ctx, b.VendorID, avg, count); err != nil {
		s.d.Logger.Warn("rating update failed", "vendor", b.VendorID, "err", err)
	}
	return created, nil
}

// ListForVendor returns the paginated public review feed for a vendor.
func (s *Service) ListForVendor(ctx context.Context, vendorID string, p ports.Page) ([]*domain.Review, error) {
	return s.d.Reviews.ListForVendor(ctx, vendorID, p)
}

// Delete removes a review (admin moderation) and recomputes the vendor's
// rating aggregate. Aggregate failures are logged, not fatal.
func (s *Service) Delete(ctx context.Context, id string) error {
	rv, err := s.d.Reviews.Find(ctx, id)
	if err != nil {
		return err
	}
	if err := s.d.Reviews.Delete(ctx, id); err != nil {
		return err
	}
	avg, count, err := s.d.Reviews.AggregateForVendor(ctx, rv.VendorID)
	if err != nil {
		s.d.Logger.Warn("rating aggregate failed", "vendor", rv.VendorID, "err", err)
		return nil
	}
	if err := s.d.Vendors.UpdateRating(ctx, rv.VendorID, avg, count); err != nil {
		s.d.Logger.Warn("rating update failed", "vendor", rv.VendorID, "err", err)
	}
	return nil
}
