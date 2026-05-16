package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func setupApprovedVendor(t *testing.T, e *env) (*domain.Vendor, string, string) {
	t.Helper()
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty", "priceFrom": 500000,
	})
	v := decode[domain.Vendor](t, body)
	_, _ = e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{"status": "approved"})

	cust := e.signup("c@b.kz", "password123", "C", "customer")
	return &v, vend, cust
}

func TestBookingFullLifecycle(t *testing.T) {
	e := newEnv(t)
	v, vendTok, custTok := setupApprovedVendor(t, e)

	res, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 150, "note": "Toi",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create booking: %d %s", res.StatusCode, body)
	}
	b := decode[domain.Booking](t, body)
	if b.Status != domain.BookingPending {
		t.Errorf("new should be pending, got %s", b.Status)
	}

	_, body = e.do("GET", "/api/bookings", vendTok, nil)
	r := decode[struct {
		Items []domain.Booking `json:"items"`
	}](t, body)
	if len(r.Items) != 1 || r.Items[0].ID != b.ID {
		t.Fatalf("vendor incoming: %+v", r.Items)
	}

	res, body = e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "accepted"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("accept: %d %s", res.StatusCode, body)
	}
	upd := decode[domain.Booking](t, body)
	if upd.Status != domain.BookingAccepted {
		t.Errorf("not accepted: %s", upd.Status)
	}
}

func TestBookingVendorCannotCancel(t *testing.T) {
	e := newEnv(t)
	v, vendTok, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 100,
	})
	b := decode[domain.Booking](t, body)
	res, _ := e.do("PATCH", "/api/bookings/"+b.ID, vendTok, map[string]any{"status": "cancelled"})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("vendor cancel should fail, got %d", res.StatusCode)
	}
}

func TestBookingCustomerCannotAccept(t *testing.T) {
	e := newEnv(t)
	v, _, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 100,
	})
	b := decode[domain.Booking](t, body)
	res, _ := e.do("PATCH", "/api/bookings/"+b.ID, custTok, map[string]any{"status": "accepted"})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("customer accept should fail, got %d", res.StatusCode)
	}
}

func TestBookingCustomerCancels(t *testing.T) {
	e := newEnv(t)
	v, _, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 100,
	})
	b := decode[domain.Booking](t, body)
	res, body := e.do("PATCH", "/api/bookings/"+b.ID, custTok, map[string]any{"status": "cancelled"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("cancel: %d %s", res.StatusCode, body)
	}
	upd := decode[domain.Booking](t, body)
	if upd.Status != domain.BookingCancelled {
		t.Errorf("not cancelled: %s", upd.Status)
	}
}

func TestBookingOtherCustomerCannotCancel(t *testing.T) {
	e := newEnv(t)
	v, _, custTok := setupApprovedVendor(t, e)
	_, body := e.do("POST", "/api/bookings", custTok, map[string]any{
		"vendorId": v.ID, "eventDate": "2026-08-12", "guestCount": 100,
	})
	b := decode[domain.Booking](t, body)
	other := e.signup("c2@b.kz", "password123", "C2", "customer")
	res, _ := e.do("PATCH", "/api/bookings/"+b.ID, other, map[string]any{"status": "cancelled"})
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("want 403 got %d", res.StatusCode)
	}
}

func TestBookingMissingFields(t *testing.T) {
	e := newEnv(t)
	_, _, custTok := setupApprovedVendor(t, e)
	res, _ := e.do("POST", "/api/bookings", custTok, map[string]any{
		"eventDate": "2026-08-12",
	})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}
