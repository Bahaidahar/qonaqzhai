// Package vendor implements vendor profile lifecycle and search use cases.
package vendor

import (
	"context"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles vendor usecase collaborators.
type Deps struct {
	Vendors ports.VendorRepo
	Reviews ports.ReviewRepo
}

// Service exposes vendor operations.
type Service struct{ d Deps }

// New constructs a vendor Service.
func New(d Deps) *Service { return &Service{d: d} }

// Submit upserts the caller's vendor profile. New profiles start in pending
// state awaiting admin moderation.
func (s *Service) Submit(ctx context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error) {
	in.Normalize()
	if in.Name == "" || in.Category == "" || in.City == "" {
		return nil, fmt.Errorf("required fields: %w", errs.ErrInvalidInput)
	}
	if in.PriceFrom < 0 {
		return nil, fmt.Errorf("price: %w", errs.ErrInvalidInput)
	}
	return s.d.Vendors.Upsert(ctx, userID, in)
}

// FindPublic returns the vendor by id only when approved.
func (s *Service) FindPublic(ctx context.Context, id string) (*domain.Vendor, error) {
	v, err := s.d.Vendors.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !v.IsPublic() {
		return nil, errs.ErrNotFound
	}
	return v, nil
}

// FindByID returns any vendor (admin/owner views).
func (s *Service) FindByID(ctx context.Context, id string) (*domain.Vendor, error) {
	return s.d.Vendors.FindByID(ctx, id)
}

// FindByUserID returns the vendor owned by userID.
func (s *Service) FindByUserID(ctx context.Context, userID string) (*domain.Vendor, error) {
	return s.d.Vendors.FindByUserID(ctx, userID)
}

// Search returns the paginated public catalog with filters applied.
func (s *Service) Search(ctx context.Context, q ports.VendorQuery) ([]*domain.Vendor, int, error) {
	if q.Status == "" {
		q.Status = domain.VendorApproved
	}
	return s.d.Vendors.Search(ctx, q)
}

// AdminSearch is like Search but without the implicit approved filter.
func (s *Service) AdminSearch(ctx context.Context, q ports.VendorQuery) ([]*domain.Vendor, int, error) {
	return s.d.Vendors.Search(ctx, q)
}

// SetStatus is the admin moderation endpoint.
func (s *Service) SetStatus(ctx context.Context, id string, status domain.VendorStatus) (*domain.Vendor, error) {
	if !status.Valid() {
		return nil, errs.ErrInvalidInput
	}
	if err := s.d.Vendors.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.d.Vendors.FindByID(ctx, id)
}

// RecomputeRating refreshes the cached rating aggregate from reviews.
// Called by the review usecase after Create.
func (s *Service) RecomputeRating(ctx context.Context, vendorID string) error {
	avg, count, err := s.d.Reviews.AggregateForVendor(ctx, vendorID)
	if err != nil {
		return err
	}
	return s.d.Vendors.UpdateRating(ctx, vendorID, avg, count)
}

// ListByIDs is a batch accessor used by the gRPC server.
func (s *Service) ListByIDs(ctx context.Context, ids []string) ([]*domain.Vendor, error) {
	return s.d.Vendors.FindByIDs(ctx, ids)
}

// Stats returns vendor counts by status.
func (s *Service) Stats(ctx context.Context) (map[domain.VendorStatus]int, error) {
	return s.d.Vendors.CountByStatus(ctx)
}
