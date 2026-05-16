package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func TestAdminStats(t *testing.T) {
	e := newEnv(t)
	_ = e.signup("c@b.kz", "password123", "C", "customer")
	_ = e.signup("v@b.kz", "password123", "V", "vendor")

	res, body := e.do("GET", "/api/admin/stats", e.adminTok, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("stats: %d %s", res.StatusCode, body)
	}
	r := decode[map[string]any](t, body)
	if int(r["users"].(float64)) != 3 {
		t.Errorf("users want 3 got %v", r["users"])
	}
	if int(r["customers"].(float64)) != 1 || int(r["vendors"].(float64)) != 1 || int(r["admins"].(float64)) != 1 {
		t.Errorf("role counts: %+v", r)
	}
}

func TestAdminApproveAndStatsReflect(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})
	v := decode[domain.Vendor](t, body)

	_, body = e.do("GET", "/api/admin/stats", e.adminTok, nil)
	r := decode[map[string]any](t, body)
	if int(r["vendors_pending"].(float64)) != 1 || int(r["vendors_approved"].(float64)) != 0 {
		t.Errorf("before approve: %+v", r)
	}

	_, _ = e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{"status": "approved"})
	_, body = e.do("GET", "/api/admin/stats", e.adminTok, nil)
	r = decode[map[string]any](t, body)
	if int(r["vendors_pending"].(float64)) != 0 || int(r["vendors_approved"].(float64)) != 1 {
		t.Errorf("after approve: %+v", r)
	}
}

func TestAdminInvalidVendorStatus(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})
	v := decode[domain.Vendor](t, body)
	res, _ := e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{
		"status": "garbage",
	})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}

func TestAdminListUsers(t *testing.T) {
	e := newEnv(t)
	_ = e.signup("c@b.kz", "password123", "C", "customer")
	res, body := e.do("GET", "/api/admin/users", e.adminTok, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("users: %d", res.StatusCode)
	}
	r := decode[struct {
		Items []domain.User `json:"items"`
	}](t, body)
	if len(r.Items) < 2 {
		t.Errorf("want at least 2 users, got %d", len(r.Items))
	}
}
