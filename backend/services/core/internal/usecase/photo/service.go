// Package photo implements vendor photo upload + retrieval.
package photo

import (
	"context"
	"fmt"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles photo collaborators.
type Deps struct {
	Photos  ports.PhotoRepo
	Vendors ports.VendorRepo
}

// Service exposes photo operations.
type Service struct{ d Deps }

// New constructs a photo Service.
func New(d Deps) *Service { return &Service{d: d} }

// Upload validates + stores a vendor photo. The caller must own vendorID. The
// MIME is sniffed from the byte stream, NOT trusted from the Content-Type
// header (the client controls headers). Only jpeg/png/webp/gif are accepted.
func (s *Service) Upload(ctx context.Context, vendorUserID string, data []byte) (*domain.Photo, error) {
	v, err := s.d.Vendors.FindByUserID(ctx, vendorUserID)
	if err != nil {
		return nil, err
	}
	size := int64(len(data))
	if size == 0 {
		return nil, fmt.Errorf("empty: %w", errs.ErrInvalidInput)
	}
	if size > domain.MaxPhotoSize {
		return nil, errs.ErrTooLarge
	}
	mime := domain.DetectPhotoMIME(data)
	if mime == "" {
		return nil, fmt.Errorf("not an allowed image type: %w", errs.ErrInvalidInput)
	}
	return s.d.Photos.Insert(ctx, &domain.Photo{
		VendorID: v.ID, MIME: mime, Size: size, Data: data,
	})
}

// Get returns a photo by id (used by the public serve endpoint).
func (s *Service) Get(ctx context.Context, id string) (*domain.Photo, error) {
	return s.d.Photos.Find(ctx, id)
}

// Delete removes a photo. The caller must own the vendor.
func (s *Service) Delete(ctx context.Context, vendorUserID, photoID string) error {
	p, err := s.d.Photos.Find(ctx, photoID)
	if err != nil {
		return err
	}
	v, err := s.d.Vendors.FindByID(ctx, p.VendorID)
	if err != nil {
		return err
	}
	if v.UserID != vendorUserID {
		return errs.ErrForbidden
	}
	return s.d.Photos.Delete(ctx, photoID)
}
