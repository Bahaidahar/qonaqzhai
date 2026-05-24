package photo_test

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"sync"
	"testing"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
	"qonaqzhai-backend/services/core/internal/usecase/photo"
)

type memPhotos struct {
	mu sync.Mutex
	rows map[string]*domain.Photo
}

func newPhotos() *memPhotos { return &memPhotos{rows: map[string]*domain.Photo{}} }

func (m *memPhotos) Insert(_ context.Context, p *domain.Photo) (*domain.Photo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p.ID = "p1"
	cp := *p
	m.rows[p.ID] = &cp
	return &cp, nil
}
func (m *memPhotos) Find(_ context.Context, id string) (*domain.Photo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.rows[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *p
	return &cp, nil
}
func (m *memPhotos) Delete(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rows[id]; !ok {
		return errs.ErrNotFound
	}
	delete(m.rows, id)
	return nil
}
func (m *memPhotos) ListByVendor(context.Context, string) ([]*domain.Photo, error) { return nil, nil }

type memVendors struct{ owner, id string }

func (m *memVendors) Upsert(context.Context, string, domain.VendorInput) (*domain.Vendor, error) { return nil, nil }
func (m *memVendors) FindByID(_ context.Context, id string) (*domain.Vendor, error) {
	if id != m.id {
		return nil, errs.ErrNotFound
	}
	return &domain.Vendor{ID: m.id, UserID: m.owner}, nil
}
func (m *memVendors) FindByUserID(_ context.Context, userID string) (*domain.Vendor, error) {
	if userID != m.owner {
		return nil, errs.ErrNotFound
	}
	return &domain.Vendor{ID: m.id, UserID: m.owner}, nil
}
func (m *memVendors) FindByIDs(context.Context, []string) ([]*domain.Vendor, error) { return nil, nil }
func (m *memVendors) Search(context.Context, ports.VendorQuery) ([]*domain.Vendor, int, error) { return nil, 0, nil }
func (m *memVendors) UpdateStatus(context.Context, string, domain.VendorStatus) error { return nil }
func (m *memVendors) UpdateRating(context.Context, string, float64, int) error { return nil }
func (m *memVendors) CountByStatus(context.Context) (map[domain.VendorStatus]int, error) { return nil, nil }

func jpegBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func pngBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func svgBytes() []byte {
	return []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"></svg>`)
}

func TestUpload_AllowedMIMEs(t *testing.T) {
	svc := photo.New(photo.Deps{Photos: newPhotos(), Vendors: &memVendors{owner: "u", id: "v1"}})
	for _, payload := range [][]byte{jpegBytes(t), pngBytes(t)} {
		p, err := svc.Upload(context.Background(), "u", payload)
		if err != nil {
			t.Fatalf("upload: %v", err)
		}
		if p.MIME != "image/jpeg" && p.MIME != "image/png" {
			t.Fatalf("unexpected mime %q", p.MIME)
		}
	}
}

func TestUpload_SVGRejected(t *testing.T) {
	svc := photo.New(photo.Deps{Photos: newPhotos(), Vendors: &memVendors{owner: "u", id: "v1"}})
	_, err := svc.Upload(context.Background(), "u", svgBytes())
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("svg must be rejected as invalid, got %v", err)
	}
}

func TestUpload_TextRejected(t *testing.T) {
	svc := photo.New(photo.Deps{Photos: newPhotos(), Vendors: &memVendors{owner: "u", id: "v1"}})
	_, err := svc.Upload(context.Background(), "u", []byte("<script>alert(1)</script>"))
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("html payload must be rejected, got %v", err)
	}
}

func TestUpload_Empty(t *testing.T) {
	svc := photo.New(photo.Deps{Photos: newPhotos(), Vendors: &memVendors{owner: "u", id: "v1"}})
	_, err := svc.Upload(context.Background(), "u", nil)
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("empty body must be rejected, got %v", err)
	}
}

func TestUpload_TooLarge(t *testing.T) {
	svc := photo.New(photo.Deps{Photos: newPhotos(), Vendors: &memVendors{owner: "u", id: "v1"}})
	big := make([]byte, domain.MaxPhotoSize+1)
	for i := range big[:8] {
		big[i] = jpegBytes(t)[i]
	}
	_, err := svc.Upload(context.Background(), "u", big)
	if !errors.Is(err, errs.ErrTooLarge) {
		t.Fatalf("expected too large, got %v", err)
	}
}

func TestDelete_AuthorisesOwner(t *testing.T) {
	photos := newPhotos()
	svc := photo.New(photo.Deps{Photos: photos, Vendors: &memVendors{owner: "owner", id: "v1"}})
	p, err := svc.Upload(context.Background(), "owner", jpegBytes(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Delete(context.Background(), "stranger", p.ID); !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("stranger must be forbidden, got %v", err)
	}
	if err := svc.Delete(context.Background(), "owner", p.ID); err != nil {
		t.Fatalf("owner delete failed: %v", err)
	}
}
