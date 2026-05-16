package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func TestReviewSubmitRequiresCompletedBooking(t *testing.T) {
	e := newEnv(t)
	v, _, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 50,
	})
	b := decode[domain.Booking](t, body)

	res, _ := e.do("POST", "/api/reviews", custTok, map[string]any{
		"bookingId": b.ID, "rating": 5, "text": "fine",
	})
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("review on pending booking allowed: %d", res.StatusCode)
	}
}

func TestReviewFullFlow(t *testing.T) {
	e := newEnv(t)
	v, vendTok, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 50,
	})
	b := decode[domain.Booking](t, body)

	// vendor accepts → completed (vendor allowed to complete from accepted/paid)
	_, _ = e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "accepted"})
	res, body := e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "completed"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("complete: %d %s", res.StatusCode, body)
	}

	// submit review
	res, body = e.do("POST", "/api/reviews", custTok, map[string]any{
		"bookingId": b.ID, "rating": 5, "text": "great",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("submit: %d %s", res.StatusCode, body)
	}
	rv := decode[domain.Review](t, body)
	if rv.Rating != 5 || rv.VendorID != v.ID {
		t.Errorf("review wrong: %+v", rv)
	}

	// listing
	_, body = e.do("GET", "/api/vendors/"+v.ID+"/reviews", "", nil)
	r := decode[struct {
		Items []domain.Review `json:"items"`
	}](t, body)
	if len(r.Items) != 1 {
		t.Errorf("list=%d", len(r.Items))
	}

	// duplicate prevented
	res, _ = e.do("POST", "/api/reviews", custTok, map[string]any{
		"bookingId": b.ID, "rating": 4, "text": "again",
	})
	if res.StatusCode != http.StatusConflict {
		t.Errorf("dup review: %d", res.StatusCode)
	}

	// vendor rating recomputed
	_, body = e.do("GET", "/api/vendors/"+v.ID, "", nil)
	got := decode[domain.Vendor](t, body)
	if got.RatingCount != 1 || got.RatingAvg != 5 {
		t.Errorf("rating not propagated: %+v", got)
	}

	// admin delete
	res, _ = e.do("DELETE", "/api/admin/reviews/"+rv.ID, e.adminTok, nil)
	if res.StatusCode != http.StatusNoContent {
		t.Errorf("admin delete: %d", res.StatusCode)
	}
	_, body = e.do("GET", "/api/vendors/"+v.ID, "", nil)
	got = decode[domain.Vendor](t, body)
	if got.RatingCount != 0 {
		t.Errorf("rating not reset after delete: %+v", got)
	}
}

func TestReviewInvalidRating(t *testing.T) {
	e := newEnv(t)
	v, vendTok, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12",
	})
	b := decode[domain.Booking](t, body)
	_, _ = e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "accepted"})
	_, _ = e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "completed"})
	res, _ := e.do("POST", "/api/reviews", custTok, map[string]any{
		"bookingId": b.ID, "rating": 0,
	})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}
