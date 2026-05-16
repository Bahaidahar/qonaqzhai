// Package review implements vendor review and rating use cases.
package review

import (
	"context"
	"fmt"
	"strings"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Deps bundles review service collaborators.
type Deps struct {
	Reviews  usecase.ReviewRepo
	Bookings usecase.BookingRepo
	Vendors  usecase.VendorRepo
	Clock    usecase.Clock
}

// Service exposes review submission, listing and moderation.
type Service struct{ d Deps }

// New constructs a review Service.
func New(d Deps) *Service { return &Service{d: d} }

// SubmitInput captures user-supplied fields for a new review.
type SubmitInput struct {
	BookingID string
	Rating    int
	Text      string
}

// Submit creates a review for a completed booking owned by the calling customer.
// Vendor's average rating + count are recomputed on success.
func (s *Service) Submit(ctx context.Context, customerID string, in SubmitInput) (*domain.Review, error) {
	if !domain.ValidRating(in.Rating) {
		return nil, fmt.Errorf("rating %d: %w", in.Rating, domain.ErrInvalidInput)
	}
	b, err := s.d.Bookings.Find(ctx, in.BookingID)
	if err != nil {
		return nil, err
	}
	if b.CustomerID != customerID {
		return nil, domain.ErrForbidden
	}
	if b.Status != domain.BookingCompleted {
		return nil, domain.ErrForbidden
	}
	r := &domain.Review{
		BookingID:  b.ID,
		CustomerID: customerID,
		VendorID:   b.VendorID,
		Rating:     in.Rating,
		Text:       strings.TrimSpace(in.Text),
		CreatedAt:  s.d.Clock.Now(),
	}
	created, err := s.d.Reviews.Create(ctx, r)
	if err != nil {
		return nil, err
	}
	if err := s.recomputeRating(ctx, b.VendorID); err != nil {
		return nil, err
	}
	return created, nil
}

// ListByVendor returns reviews for a vendor.
func (s *Service) ListByVendor(ctx context.Context, vendorID string) ([]*domain.Review, error) {
	return s.d.Reviews.ListByVendor(ctx, vendorID)
}

// AdminDelete removes a review and recomputes the vendor's rating aggregate.
func (s *Service) AdminDelete(ctx context.Context, reviewID string) error {
	r, err := s.d.Reviews.FindByID(ctx, reviewID)
	if err != nil {
		return err
	}
	if err := s.d.Reviews.Delete(ctx, reviewID); err != nil {
		return err
	}
	return s.recomputeRating(ctx, r.VendorID)
}

func (s *Service) recomputeRating(ctx context.Context, vendorID string) error {
	avg, count, err := s.d.Reviews.AggregateForVendor(ctx, vendorID)
	if err != nil {
		return err
	}
	return s.d.Vendors.UpdateRating(ctx, vendorID, avg, count)
}
