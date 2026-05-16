package booking_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/booking"
	"qonaqzhai-backend/internal/usecase/inmem"
)

func newSvc(t *testing.T) (*booking.Service, *inmem.VendorRepo, *inmem.BookingRepo) {
	t.Helper()
	vrepo := inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)
	brepo := inmem.NewBookingRepo(inmem.NewSeqIDGen("b-").New)
	return booking.New(booking.Deps{
		Bookings: brepo,
		Vendors:  vrepo,
		Clock:    &inmem.FixedClock{T: time.Now()},
	}), vrepo, brepo
}

func TestCreateRequiresExistingVendor(t *testing.T) {
	t.Parallel()
	svc, _, _ := newSvc(t)
	_, err := svc.Create(context.Background(), "cust", booking.CreateInput{VendorID: "missing", EventDate: "2026-06-12"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("err=%v want ErrNotFound", err)
	}
}

func TestCreateValidatesFields(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vend-user", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})

	cases := []struct {
		name      string
		in        booking.CreateInput
		wantErr   error
	}{
		{"empty vendor", booking.CreateInput{EventDate: "2026-06-12"}, domain.ErrInvalidInput},
		{"empty date", booking.CreateInput{VendorID: v.ID}, domain.ErrInvalidInput},
		{"negative guests", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-12", GuestCount: -1}, domain.ErrInvalidInput},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := svc.Create(ctx, "cust", c.in); !errors.Is(err, c.wantErr) {
				t.Errorf("err=%v want %v", err, c.wantErr)
			}
		})
	}
}

func TestCreateStartsAsPending(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, err := svc.Create(ctx, "cust", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-12", GuestCount: 150})
	if err != nil {
		t.Fatal(err)
	}
	if b.Status != domain.BookingPending {
		t.Errorf("status=%s", b.Status)
	}
	if b.CustomerID != "cust" {
		t.Errorf("customerId=%s", b.CustomerID)
	}
}

func TestVendorAcceptDeclineFlow(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, _ := svc.Create(ctx, "cust", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-12"})

	if _, err := svc.VendorTransition(ctx, "vu", b.ID, domain.BookingAccepted); err != nil {
		t.Fatalf("accept: %v", err)
	}
	// further vendor transitions: only completed allowed
	if _, err := svc.VendorTransition(ctx, "vu", b.ID, domain.BookingDeclined); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("vendor decline after accept allowed? err=%v", err)
	}
	if _, err := svc.VendorTransition(ctx, "vu", b.ID, domain.BookingCompleted); err != nil {
		t.Fatalf("complete: %v", err)
	}
}

func TestVendorCannotTouchOtherVendorBooking(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v1, _ := vrepo.Upsert(ctx, "vu1", domain.VendorInput{Name: "1", Category: "Y", City: "Z"})
	_, _ = vrepo.Upsert(ctx, "vu2", domain.VendorInput{Name: "2", Category: "Y", City: "Z"})
	b, _ := svc.Create(ctx, "cust", booking.CreateInput{VendorID: v1.ID, EventDate: "2026-06-12"})

	if _, err := svc.VendorTransition(ctx, "vu2", b.ID, domain.BookingAccepted); !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("err=%v want ErrForbidden", err)
	}
}

func TestCustomerCancel(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, _ := svc.Create(ctx, "cust", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-12"})

	if _, err := svc.CustomerTransition(ctx, "other", b.ID, domain.BookingCancelled); !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("non-owner cancel allowed: %v", err)
	}
	if _, err := svc.CustomerTransition(ctx, "cust", b.ID, domain.BookingAccepted); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("customer accept allowed: %v", err)
	}
	if _, err := svc.CustomerTransition(ctx, "cust", b.ID, domain.BookingCancelled); err != nil {
		t.Fatalf("cancel: %v", err)
	}
}

func TestListByRole(t *testing.T) {
	t.Parallel()
	svc, vrepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	_, _ = svc.Create(ctx, "cust1", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-12"})
	_, _ = svc.Create(ctx, "cust2", booking.CreateInput{VendorID: v.ID, EventDate: "2026-06-13"})

	list, _ := svc.ListForCustomer(ctx, "cust1")
	if len(list) != 1 {
		t.Errorf("customer list = %d, want 1", len(list))
	}
	list, _ = svc.ListForVendor(ctx, "vu")
	if len(list) != 2 {
		t.Errorf("vendor list = %d, want 2", len(list))
	}
	list, _ = svc.ListAll(ctx)
	if len(list) != 2 {
		t.Errorf("all = %d, want 2", len(list))
	}
}
