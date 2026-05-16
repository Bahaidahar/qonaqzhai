package admin_test

import (
	"context"
	"errors"
	"testing"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/admin"
	"qonaqzhai-backend/internal/usecase/inmem"
)

func newSvc(t *testing.T) (*admin.Service, *inmem.UserRepo, *inmem.VendorRepo, *inmem.BookingRepo) {
	t.Helper()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	vendors := inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)
	bookings := inmem.NewBookingRepo(inmem.NewSeqIDGen("b-").New)
	reviews := inmem.NewReviewRepo(inmem.NewSeqIDGen("r-").New)
	return admin.New(admin.Deps{
		Users:    users,
		Vendors:  vendors,
		Bookings: bookings,
		Reviews:  reviews,
	}), users, vendors, bookings
}

func TestSetUserStatusValidates(t *testing.T) {
	t.Parallel()
	svc, _, _, _ := newSvc(t)
	if _, err := svc.SetUserStatus(context.Background(), "admin", "admin@x", "x", "zombie"); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err=%v want ErrInvalidInput", err)
	}
}

func TestSetVendorStatusValidates(t *testing.T) {
	t.Parallel()
	svc, _, _, _ := newSvc(t)
	if _, err := svc.SetVendorStatus(context.Background(), "admin", "admin@x", "x", "zombie"); !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("err=%v want ErrInvalidInput", err)
	}
}

func TestStatsAggregates(t *testing.T) {
	t.Parallel()
	svc, users, vendors, bookings := newSvc(t)
	ctx := context.Background()

	cu, _ := users.Create(ctx, &domain.User{Email: "a@b", Role: domain.RoleCustomer})
	_, _ = users.Create(ctx, &domain.User{Email: "v@b", Role: domain.RoleVendor})
	_, _ = users.Create(ctx, &domain.User{Email: "ad@b", Role: domain.RoleAdmin})

	v, _ := vendors.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	_ = vendors.UpdateStatus(ctx, v.ID, domain.VendorApproved)

	_, _ = bookings.Create(ctx, &domain.Booking{CustomerID: cu.ID, VendorID: v.ID, Status: domain.BookingPending})
	_, _ = bookings.Create(ctx, &domain.Booking{CustomerID: cu.ID, VendorID: v.ID, Status: domain.BookingPaid, Amount: 100})

	st, err := svc.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if st.Users != 3 {
		t.Errorf("users=%d", st.Users)
	}
	if st.Customers != 1 || st.Vendors != 1 || st.Admins != 1 {
		t.Errorf("role split: %+v", st)
	}
	if st.VendorsApproved != 1 {
		t.Errorf("approved=%d", st.VendorsApproved)
	}
	if st.BookingsTotal != 2 || st.BookingsPaid != 1 || st.BookingsRevenue != 100 {
		t.Errorf("bookings: %+v", st)
	}
}
