package tests

import (
	"net/http"
	"testing"
)

func TestCustomerCannotAccessVendorEndpoints(t *testing.T) {
	e := newEnv(t)
	cust := e.signup("c@b.kz", "password123", "C", "customer")
	endpoints := []struct{ method, path string }{
		{"POST", "/api/vendor"},
		{"GET", "/api/vendor"},
		{"POST", "/api/vendor/photos"},
	}
	for _, ep := range endpoints {
		res, _ := e.do(ep.method, ep.path, cust, map[string]any{})
		if res.StatusCode != http.StatusForbidden {
			t.Errorf("%s %s want 403 got %d", ep.method, ep.path, res.StatusCode)
		}
	}
}

func TestVendorCannotCreateBooking(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	res, _ := e.do("POST", "/api/bookings", vend, map[string]any{
		"vendorId": "x", "eventDate": "2026-12-01",
	})
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("want 403 got %d", res.StatusCode)
	}
}

func TestCustomerForbiddenOnAdmin(t *testing.T) {
	e := newEnv(t)
	cust := e.signup("c@b.kz", "password123", "C", "customer")
	res, _ := e.do("GET", "/api/admin/users", cust, nil)
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("want 403 got %d", res.StatusCode)
	}
}

func TestVendorForbiddenOnAdmin(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	res, _ := e.do("GET", "/api/admin/users", vend, nil)
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("want 403 got %d", res.StatusCode)
	}
}

func TestAdminAccessAllowed(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("GET", "/api/admin/users", e.adminTok, nil)
	if res.StatusCode != http.StatusOK {
		t.Errorf("want 200 got %d", res.StatusCode)
	}
}
