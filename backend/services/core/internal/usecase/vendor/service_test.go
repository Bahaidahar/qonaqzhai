package vendor_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
	"qonaqzhai-backend/services/core/internal/usecase/vendor"
)

type memVendors struct {
	mu   sync.Mutex
	rows map[string]*domain.Vendor
}

func newMemVendors() *memVendors { return &memVendors{rows: map[string]*domain.Vendor{}} }

func (m *memVendors) Upsert(_ context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.rows {
		if v.UserID == userID {
			v.Name, v.Category, v.City, v.Description, v.PriceFrom = in.Name, in.Category, in.City, in.Description, in.PriceFrom
			cp := *v
			return &cp, nil
		}
	}
	id := "v" + strconv.Itoa(len(m.rows)+1)
	v := &domain.Vendor{
		ID: id, UserID: userID,
		Name: in.Name, Category: in.Category, City: in.City,
		Description: in.Description, PriceFrom: in.PriceFrom,
		Status: domain.VendorPending,
	}
	m.rows[id] = v
	cp := *v
	return &cp, nil
}
func (m *memVendors) FindByID(_ context.Context, id string) (*domain.Vendor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *v
	return &cp, nil
}
func (m *memVendors) FindByUserID(_ context.Context, userID string) (*domain.Vendor, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, v := range m.rows {
		if v.UserID == userID {
			cp := *v
			return &cp, nil
		}
	}
	return nil, errs.ErrNotFound
}
func (m *memVendors) FindByIDs(context.Context, []string) ([]*domain.Vendor, error) { return nil, nil }
func (m *memVendors) Search(_ context.Context, q ports.VendorQuery) ([]*domain.Vendor, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Vendor{}
	for _, v := range m.rows {
		if q.Status != "" && v.Status != q.Status {
			continue
		}
		cp := *v
		out = append(out, &cp)
	}
	return out, len(out), nil
}
func (m *memVendors) UpdateStatus(_ context.Context, id string, s domain.VendorStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	v.Status = s
	return nil
}
func (m *memVendors) UpdateRating(_ context.Context, id string, avg float64, count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	v.RatingAvg, v.RatingCount = avg, count
	return nil
}
func (m *memVendors) CountByStatus(context.Context) (map[domain.VendorStatus]int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := map[domain.VendorStatus]int{}
	for _, v := range m.rows {
		out[v.Status]++
	}
	return out, nil
}

type fakeReviews struct{ avg float64; count int; err error }

func (f *fakeReviews) Create(context.Context, *domain.Review) (*domain.Review, error) { return nil, nil }
func (f *fakeReviews) ListForVendor(context.Context, string, ports.Page) ([]*domain.Review, error) {
	return nil, nil
}
func (f *fakeReviews) FindByBooking(context.Context, string) (*domain.Review, error) { return nil, nil }
func (f *fakeReviews) AggregateForVendor(context.Context, string) (float64, int, error) {
	return f.avg, f.count, f.err
}

func newSvc() (*vendor.Service, *memVendors, *fakeReviews) {
	vs := newMemVendors()
	rv := &fakeReviews{}
	svc := vendor.New(vendor.Deps{Vendors: vs, Reviews: rv})
	return svc, vs, rv
}

// --- tests -------------------------------------------------------------------

func TestSubmit_DefaultsPending(t *testing.T) {
	svc, _, _ := newSvc()
	v, err := svc.Submit(context.Background(), "user-1", domain.VendorInput{
		Name: "Rixos", Category: "venue", City: "Almaty", PriceFrom: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if v.Status != domain.VendorPending {
		t.Fatalf("expected pending, got %s", v.Status)
	}
}

func TestSubmit_RequiredFields(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.Submit(context.Background(), "u", domain.VendorInput{Name: " "})
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestSubmit_NegativePrice(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.Submit(context.Background(), "u", domain.VendorInput{
		Name: "x", Category: "c", City: "Almaty", PriceFrom: -1,
	})
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestFindPublic_OnlyApproved(t *testing.T) {
	svc, vendors, _ := newSvc()
	v, _ := vendors.Upsert(context.Background(), "u", domain.VendorInput{
		Name: "x", Category: "c", City: "Almaty",
	})
	_, err := svc.FindPublic(context.Background(), v.ID)
	if !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("pending vendor must be hidden, got %v", err)
	}
	_ = vendors.UpdateStatus(context.Background(), v.ID, domain.VendorApproved)
	got, err := svc.FindPublic(context.Background(), v.ID)
	if err != nil || got.ID != v.ID {
		t.Fatalf("expected approved vendor, got %v err=%v", got, err)
	}
}

func TestSearch_DefaultsApproved(t *testing.T) {
	svc, vendors, _ := newSvc()
	a, _ := vendors.Upsert(context.Background(), "u1", domain.VendorInput{Name: "a", Category: "c", City: "Almaty"})
	_, _ = vendors.Upsert(context.Background(), "u2", domain.VendorInput{Name: "b", Category: "c", City: "Almaty"})
	_ = vendors.UpdateStatus(context.Background(), a.ID, domain.VendorApproved)
	got, total, err := svc.Search(context.Background(), ports.VendorQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(got) != 1 || got[0].ID != a.ID {
		t.Fatalf("expected only approved vendor, got total=%d items=%d", total, len(got))
	}
}

func TestSetStatus_InvalidRejected(t *testing.T) {
	svc, _, _ := newSvc()
	_, err := svc.SetStatus(context.Background(), "x", domain.VendorStatus("garbage"))
	if !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("want invalid input, got %v", err)
	}
}

func TestRecomputeRating(t *testing.T) {
	svc, vendors, reviews := newSvc()
	v, _ := vendors.Upsert(context.Background(), "u", domain.VendorInput{Name: "x", Category: "c", City: "Almaty"})
	reviews.avg, reviews.count = 4.5, 8
	if err := svc.RecomputeRating(context.Background(), v.ID); err != nil {
		t.Fatal(err)
	}
	got, _ := vendors.FindByID(context.Background(), v.ID)
	if got.RatingAvg != 4.5 || got.RatingCount != 8 {
		t.Fatalf("rating not persisted: %+v", got)
	}
}

func TestStats(t *testing.T) {
	svc, vendors, _ := newSvc()
	for i, st := range []domain.VendorStatus{domain.VendorPending, domain.VendorApproved, domain.VendorApproved} {
		v, _ := vendors.Upsert(context.Background(), "u"+strconv.Itoa(i),
			domain.VendorInput{Name: "x", Category: "c", City: "Almaty"})
		_ = vendors.UpdateStatus(context.Background(), v.ID, st)
	}
	stats, err := svc.Stats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if stats[domain.VendorApproved] != 2 || stats[domain.VendorPending] != 1 {
		t.Fatalf("stats wrong: %+v", stats)
	}
}
