package review_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"qonaqzhai-backend/pkg/errs"
	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
	"qonaqzhai-backend/services/core/internal/usecase/review"
)

type memReviews struct {
	mu   sync.Mutex
	rows []*domain.Review
}

func (m *memReviews) Create(_ context.Context, r *domain.Review) (*domain.Review, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.rows {
		if e.BookingID == r.BookingID {
			return nil, errs.ErrAlreadyExists
		}
	}
	cp := *r
	cp.ID = "r1"
	m.rows = append(m.rows, &cp)
	return &cp, nil
}
func (m *memReviews) ListForVendor(_ context.Context, id string, _ ports.Page) ([]*domain.Review, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.Review{}
	for _, r := range m.rows {
		if r.VendorID == id {
			cp := *r
			out = append(out, &cp)
		}
	}
	return out, nil
}
func (m *memReviews) FindByBooking(context.Context, string) (*domain.Review, error) { return nil, errs.ErrNotFound }
func (m *memReviews) AggregateForVendor(_ context.Context, id string) (float64, int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var sum, n float64
	for _, r := range m.rows {
		if r.VendorID == id {
			sum += float64(r.Rating)
			n++
		}
	}
	if n == 0 {
		return 0, 0, nil
	}
	return sum / n, int(n), nil
}

type memBookings struct{ b *domain.Booking }

func (m *memBookings) Create(context.Context, *domain.Booking) (*domain.Booking, error) { return nil, nil }
func (m *memBookings) Find(context.Context, string) (*domain.Booking, error) {
	if m.b == nil {
		return nil, errs.ErrNotFound
	}
	cp := *m.b
	return &cp, nil
}
func (m *memBookings) ListForCustomer(context.Context, string, ports.Page) ([]*domain.Booking, error) { return nil, nil }
func (m *memBookings) ListForVendor(context.Context, string, ports.Page) ([]*domain.Booking, error)   { return nil, nil }
func (m *memBookings) ListAll(context.Context, ports.Page) ([]*domain.Booking, error)                  { return nil, nil }
func (m *memBookings) UpdateStatus(context.Context, string, domain.BookingStatus) error               { return nil }
func (m *memBookings) SetPayment(context.Context, string, string) error                                { return nil }
func (m *memBookings) MarkPaid(context.Context, string, string) error                                  { return nil }
func (m *memBookings) Stats(context.Context) (ports.BookingStats, error)                               { return ports.BookingStats{}, nil }

type memVendors struct {
	mu  sync.Mutex
	avg float64
	cnt int
}

func (m *memVendors) Upsert(context.Context, string, domain.VendorInput) (*domain.Vendor, error) { return nil, nil }
func (m *memVendors) FindByID(context.Context, string) (*domain.Vendor, error)                   { return nil, errs.ErrNotFound }
func (m *memVendors) FindByUserID(context.Context, string) (*domain.Vendor, error)               { return nil, errs.ErrNotFound }
func (m *memVendors) FindByIDs(context.Context, []string) ([]*domain.Vendor, error)              { return nil, nil }
func (m *memVendors) Search(context.Context, ports.VendorQuery) ([]*domain.Vendor, int, error)   { return nil, 0, nil }
func (m *memVendors) UpdateStatus(context.Context, string, domain.VendorStatus) error            { return nil }
func (m *memVendors) UpdateRating(_ context.Context, _ string, avg float64, count int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.avg, m.cnt = avg, count
	return nil
}
func (m *memVendors) CountByStatus(context.Context) (map[domain.VendorStatus]int, error) { return nil, nil }

func TestSubmit_HappyPath_UpdatesRating(t *testing.T) {
	bk := &memBookings{b: &domain.Booking{ID: "b1", CustomerID: "c", VendorID: "v1", Status: domain.BookingCompleted}}
	rev := &memReviews{}
	vs := &memVendors{}
	svc := review.New(review.Deps{Reviews: rev, Bookings: bk, Vendors: vs})
	r, err := svc.Submit(context.Background(), "c", "b1", 5, "great")
	if err != nil {
		t.Fatal(err)
	}
	if r.Rating != 5 {
		t.Fatalf("rating not stored: %+v", r)
	}
	if vs.avg != 5 || vs.cnt != 1 {
		t.Fatalf("vendor rating not updated: avg=%v cnt=%v", vs.avg, vs.cnt)
	}
}

func TestSubmit_OnlyOwner(t *testing.T) {
	bk := &memBookings{b: &domain.Booking{ID: "b1", CustomerID: "owner", VendorID: "v1", Status: domain.BookingCompleted}}
	svc := review.New(review.Deps{Reviews: &memReviews{}, Bookings: bk, Vendors: &memVendors{}})
	_, err := svc.Submit(context.Background(), "stranger", "b1", 5, "")
	if !errors.Is(err, errs.ErrForbidden) {
		t.Fatalf("want forbidden, got %v", err)
	}
}

func TestSubmit_RatingRange(t *testing.T) {
	svc := review.New(review.Deps{Reviews: &memReviews{}, Bookings: &memBookings{}, Vendors: &memVendors{}})
	if _, err := svc.Submit(context.Background(), "c", "b", 0, ""); !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("rating 0 invalid, got %v", err)
	}
	if _, err := svc.Submit(context.Background(), "c", "b", 6, ""); !errors.Is(err, errs.ErrInvalidInput) {
		t.Fatalf("rating 6 invalid, got %v", err)
	}
}

func TestSubmit_BookingMustBeCompletedOrPaid(t *testing.T) {
	bk := &memBookings{b: &domain.Booking{ID: "b1", CustomerID: "c", VendorID: "v1", Status: domain.BookingPending}}
	svc := review.New(review.Deps{Reviews: &memReviews{}, Bookings: bk, Vendors: &memVendors{}})
	_, err := svc.Submit(context.Background(), "c", "b1", 5, "")
	if !errors.Is(err, errs.ErrConflict) {
		t.Fatalf("want conflict for pending booking, got %v", err)
	}
}

func TestSubmit_DuplicateBookingRejected(t *testing.T) {
	bk := &memBookings{b: &domain.Booking{ID: "b1", CustomerID: "c", VendorID: "v1", Status: domain.BookingPaid}}
	rev := &memReviews{}
	svc := review.New(review.Deps{Reviews: rev, Bookings: bk, Vendors: &memVendors{}})
	if _, err := svc.Submit(context.Background(), "c", "b1", 5, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Submit(context.Background(), "c", "b1", 4, ""); !errors.Is(err, errs.ErrAlreadyExists) {
		t.Fatalf("expected duplicate rejected, got %v", err)
	}
}
