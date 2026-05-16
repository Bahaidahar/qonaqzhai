// Package vendor implements vendor-profile and photo management use cases.
package vendor

import (
	"context"
	"fmt"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Deps bundles vendor service collaborators.
type Deps struct {
	Vendors usecase.VendorRepo
	Photos  usecase.PhotoRepo
	Clock   usecase.Clock
}

// Service exposes vendor profile and photo operations.
type Service struct{ d Deps }

// New constructs a vendor Service.
func New(d Deps) *Service { return &Service{d: d} }

// Upsert creates or updates the calling user's vendor profile.
// Newly created profiles start in pending status.
func (s *Service) Upsert(ctx context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return s.d.Vendors.Upsert(ctx, userID, in)
}

// ByID returns the vendor by id.
func (s *Service) ByID(ctx context.Context, id string) (*domain.Vendor, error) {
	return s.d.Vendors.FindByID(ctx, id)
}

// MyVendor returns the vendor owned by userID.
func (s *Service) MyVendor(ctx context.Context, userID string) (*domain.Vendor, error) {
	return s.d.Vendors.FindByUserID(ctx, userID)
}

// PublicSearch executes a catalog search restricted to approved vendors.
func (s *Service) PublicSearch(ctx context.Context, q usecase.VendorQuery) ([]*domain.Vendor, int, error) {
	q.Status = domain.VendorApproved
	return s.d.Vendors.Search(ctx, q)
}

// AdminSearch returns vendors of any status — for moderation.
func (s *Service) AdminSearch(ctx context.Context, q usecase.VendorQuery) ([]*domain.Vendor, int, error) {
	return s.d.Vendors.Search(ctx, q)
}

// UpdateStatus is an admin-only operation that moves a vendor between
// pending → approved → rejected (or back to pending). Invalid statuses are rejected.
func (s *Service) UpdateStatus(ctx context.Context, id string, status domain.VendorStatus) (*domain.Vendor, error) {
	if !status.Valid() {
		return nil, fmt.Errorf("status: %w", domain.ErrInvalidInput)
	}
	if err := s.d.Vendors.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.d.Vendors.FindByID(ctx, id)
}

// UploadPhoto attaches a photo to the calling user's vendor profile.
func (s *Service) UploadPhoto(ctx context.Context, userID, mime string, data []byte) (*domain.Photo, error) {
	if !domain.ValidPhotoMIME(mime) {
		return nil, fmt.Errorf("mime: %w", domain.ErrInvalidInput)
	}
	if int64(len(data)) > domain.MaxPhotoSize {
		return nil, domain.ErrTooLarge
	}
	v, err := s.d.Vendors.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.d.Photos.Create(ctx, v.ID, mime, data)
}

// FindPhoto returns photo metadata + bytes (for streaming).
func (s *Service) FindPhoto(ctx context.Context, id string) (*domain.Photo, error) {
	return s.d.Photos.Find(ctx, id)
}

// DeletePhoto removes a photo if it belongs to the calling user's vendor.
func (s *Service) DeletePhoto(ctx context.Context, userID, photoID string) error {
	v, err := s.d.Vendors.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	p, err := s.d.Photos.Find(ctx, photoID)
	if err != nil {
		return err
	}
	if p.VendorID != v.ID {
		return domain.ErrForbidden
	}
	return s.d.Photos.Delete(ctx, photoID)
}
