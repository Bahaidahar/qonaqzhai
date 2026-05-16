package review_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase/inmem"
	"qonaqzhai-backend/internal/usecase/review"
)

func newSvc(t *testing.T) (*review.Service, *inmem.VendorRepo, *inmem.BookingRepo, *inmem.ReviewRepo) {
	t.Helper()
	vrepo := inmem.NewVendorRepo(inmem.NewSeqIDGen("v-").New)
	brepo := inmem.NewBookingRepo(inmem.NewSeqIDGen("b-").New)
	rrepo := inmem.NewReviewRepo(inmem.NewSeqIDGen("r-").New)
	return review.New(review.Deps{
		Reviews:  rrepo,
		Bookings: brepo,
		Vendors:  vrepo,
		Clock:    &inmem.FixedClock{T: time.Now()},
	}), vrepo, brepo, rrepo
}

func seedCompletedBooking(t *testing.T, vrepo *inmem.VendorRepo, brepo *inmem.BookingRepo, customer, vendorUser string) (vendorID, bookingID string) {
	t.Helper()
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, vendorUser, domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, _ := brepo.Create(ctx, &domain.Booking{
		CustomerID: customer,
		VendorID:   v.ID,
		EventDate:  "2026-06-12",
		Status:     domain.BookingCompleted,
	})
	return v.ID, b.ID
}

func TestSubmitRequiresCompletedBooking(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, _ := brepo.Create(ctx, &domain.Booking{CustomerID: "cust", VendorID: v.ID, Status: domain.BookingPending})

	if _, err := svc.Submit(ctx, "cust", review.SubmitInput{BookingID: b.ID, Rating: 5, Text: "ok"}); !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("submit on non-completed booking allowed: %v", err)
	}
}

func TestSubmitChecksCustomerOwnership(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	_, bID := seedCompletedBooking(t, vrepo, brepo, "cust", "vu")
	if _, err := svc.Submit(context.Background(), "other", review.SubmitInput{BookingID: bID, Rating: 5, Text: "ok"}); !errors.Is(err, domain.ErrForbidden) {
		t.Errorf("non-owner review allowed: %v", err)
	}
}

func TestSubmitValidatesRating(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	_, bID := seedCompletedBooking(t, vrepo, brepo, "cust", "vu")
	for _, bad := range []int{0, -1, 6, 99} {
		if _, err := svc.Submit(context.Background(), "cust", review.SubmitInput{BookingID: bID, Rating: bad, Text: "x"}); !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("rating %d accepted", bad)
		}
	}
}

func TestSubmitOneReviewPerBooking(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	_, bID := seedCompletedBooking(t, vrepo, brepo, "cust", "vu")
	ctx := context.Background()
	if _, err := svc.Submit(ctx, "cust", review.SubmitInput{BookingID: bID, Rating: 5, Text: "first"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Submit(ctx, "cust", review.SubmitInput{BookingID: bID, Rating: 4, Text: "second"}); !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("second review allowed: %v", err)
	}
}

func TestSubmitUpdatesVendorRatingAggregate(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b1, _ := brepo.Create(ctx, &domain.Booking{CustomerID: "c1", VendorID: v.ID, Status: domain.BookingCompleted})
	b2, _ := brepo.Create(ctx, &domain.Booking{CustomerID: "c2", VendorID: v.ID, Status: domain.BookingCompleted})

	if _, err := svc.Submit(ctx, "c1", review.SubmitInput{BookingID: b1.ID, Rating: 5, Text: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Submit(ctx, "c2", review.SubmitInput{BookingID: b2.ID, Rating: 3, Text: "b"}); err != nil {
		t.Fatal(err)
	}
	fresh, _ := vrepo.FindByID(ctx, v.ID)
	if fresh.RatingCount != 2 {
		t.Errorf("count=%d", fresh.RatingCount)
	}
	if fresh.RatingAvg != 4 {
		t.Errorf("avg=%v want 4", fresh.RatingAvg)
	}
}

func TestDeleteRecomputesAggregate(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	ctx := context.Background()
	v, _ := vrepo.Upsert(ctx, "vu", domain.VendorInput{Name: "X", Category: "Y", City: "Z"})
	b, _ := brepo.Create(ctx, &domain.Booking{CustomerID: "c", VendorID: v.ID, Status: domain.BookingCompleted})
	rv, _ := svc.Submit(ctx, "c", review.SubmitInput{BookingID: b.ID, Rating: 5, Text: "a"})
	if err := svc.AdminDelete(ctx, rv.ID); err != nil {
		t.Fatal(err)
	}
	fresh, _ := vrepo.FindByID(ctx, v.ID)
	if fresh.RatingCount != 0 || fresh.RatingAvg != 0 {
		t.Errorf("aggregate after delete: count=%d avg=%v", fresh.RatingCount, fresh.RatingAvg)
	}
}

func TestListByVendor(t *testing.T) {
	t.Parallel()
	svc, vrepo, brepo, _ := newSvc(t)
	ctx := context.Background()
	_, bID := seedCompletedBooking(t, vrepo, brepo, "c", "vu")
	v, _ := vrepo.FindByUserID(ctx, "vu")
	_, _ = svc.Submit(ctx, "c", review.SubmitInput{BookingID: bID, Rating: 5, Text: "ok"})
	list, err := svc.ListByVendor(ctx, v.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Errorf("len=%d", len(list))
	}
}
