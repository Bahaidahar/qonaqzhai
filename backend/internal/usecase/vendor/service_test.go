package vendor_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/inmem"
	"qonaqzhai-backend/internal/usecase/vendor"
)

func newSvc(t *testing.T) (*vendor.Service, *inmem.VendorRepo, *inmem.PhotoRepo) {
	t.Helper()
	vrepo := inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)
	prepo := inmem.NewPhotoRepo(inmem.NewSeqIDGen("p-").New, vrepo)
	return vendor.New(vendor.Deps{
		Vendors: vrepo,
		Photos:  prepo,
		Clock:   &inmem.FixedClock{T: time.Now()},
	}), vrepo, prepo
}

func TestUpsertCreatesPendingVendor(t *testing.T) {
	t.Parallel()
	svc, _, _ := newSvc(t)
	ctx := context.Background()
	v, err := svc.Upsert(ctx, "user-1", domain.VendorInput{Name: "Rixos", Category: "Venue", City: "Almaty", PriceFrom: 1_500_000})
	if err != nil {
		t.Fatal(err)
	}
	if v.Status != domain.VendorPending {
		t.Errorf("status=%s want pending", v.Status)
	}
	if v.UserID != "user-1" {
		t.Errorf("userId=%s", v.UserID)
	}
}

func TestUpsertNormalizesAndValidates(t *testing.T) {
	t.Parallel()
	svc, _, _ := newSvc(t)
	_, err := svc.Upsert(context.Background(), "u", domain.VendorInput{Name: "", Category: "Venue", City: "A"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUpsertIsIdempotentPerUser(t *testing.T) {
	t.Parallel()
	svc, repo, _ := newSvc(t)
	ctx := context.Background()
	v1, _ := svc.Upsert(ctx, "u", domain.VendorInput{Name: "A", Category: "Venue", City: "Almaty"})
	v2, _ := svc.Upsert(ctx, "u", domain.VendorInput{Name: "B", Category: "Venue", City: "Almaty"})
	if v1.ID != v2.ID {
		t.Error("upsert created a second row for same user")
	}
	all, _, _ := repo.Search(ctx, usecase.VendorQuery{})
	if len(all) != 1 {
		t.Errorf("got %d rows, want 1", len(all))
	}
	if v2.Name != "B" {
		t.Errorf("name not updated: %s", v2.Name)
	}
}

func TestPublicCatalogHidesNonApproved(t *testing.T) {
	t.Parallel()
	svc, repo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := svc.Upsert(ctx, "u1", domain.VendorInput{Name: "Pending", Category: "Venue", City: "Almaty"})
	approved, _ := svc.Upsert(ctx, "u2", domain.VendorInput{Name: "Approved", Category: "Venue", City: "Almaty"})
	_ = repo.UpdateStatus(ctx, approved.ID, domain.VendorApproved)
	_ = v // not approved

	items, total, err := svc.PublicSearch(ctx, usecase.VendorQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("got %d/%d, want 1/1", len(items), total)
	}
	if items[0].Name != "Approved" {
		t.Errorf("wrong vendor surfaced: %s", items[0].Name)
	}
}

func TestUploadPhotoEnforcesSizeAndMIME(t *testing.T) {
	t.Parallel()
	svc, _, _ := newSvc(t)
	ctx := context.Background()
	_, _ = svc.Upsert(ctx, "u", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})

	if _, err := svc.UploadPhoto(ctx, "u", "text/plain", []byte("hi")); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("non-image mime accepted: %v", err)
	}
	big := make([]byte, domain.MaxPhotoSize+1)
	if _, err := svc.UploadPhoto(ctx, "u", "image/png", big); !errors.Is(err, domain.ErrTooLarge) {
		t.Errorf("oversize accepted: %v", err)
	}
	p, err := svc.UploadPhoto(ctx, "u", "image/png", []byte{0x89, 0x50, 0x4e, 0x47})
	if err != nil {
		t.Fatal(err)
	}
	if p.Size != 4 {
		t.Errorf("size=%d", p.Size)
	}
}

func TestUploadPhotoRequiresVendorProfile(t *testing.T) {
	t.Parallel()
	svc, _, _ := newSvc(t)
	if _, err := svc.UploadPhoto(context.Background(), "nobody", "image/png", []byte{0x89}); !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("err=%v want ErrNotFound", err)
	}
}
