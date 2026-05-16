package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func TestVendorUpsertAndStatusBadge(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "Rixos", "vendor")

	res, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Rixos Almaty", "category": "Venue", "city": "Almaty",
		"description": "Premier", "priceFrom": 1500000,
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("upsert: %d %s", res.StatusCode, body)
	}
	v := decode[domain.Vendor](t, body)
	if v.Status != domain.VendorPending {
		t.Errorf("new vendor should be pending, got %s", v.Status)
	}

	res, _ = e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Rixos Updated", "category": "Venue", "city": "Almaty",
	})
	if res.StatusCode != http.StatusOK {
		t.Errorf("second upsert: %d", res.StatusCode)
	}
	_, body = e.do("GET", "/api/vendor", vend, nil)
	v2 := decode[domain.Vendor](t, body)
	if v2.ID != v.ID || v2.Name != "Rixos Updated" {
		t.Errorf("upsert idempotency broken: %+v", v2)
	}
}

func TestVendorRequiredFields(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	res, _ := e.do("POST", "/api/vendor", vend, map[string]any{"name": "Only name"})
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}

func TestPendingVendorHiddenFromPublic(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Pending Co", "category": "Venue", "city": "Almaty",
	})
	_, body := e.do("GET", "/api/vendors", "", nil)
	r := decode[struct {
		Items []*domain.Vendor `json:"items"`
	}](t, body)
	if len(r.Items) != 0 {
		t.Errorf("pending should be hidden, got %d", len(r.Items))
	}
}

func TestAdminApprovesVendor(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})
	v := decode[domain.Vendor](t, body)

	res, body := e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{
		"status": "approved",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("approve: %d %s", res.StatusCode, body)
	}
	updated := decode[domain.Vendor](t, body)
	if updated.Status != domain.VendorApproved {
		t.Errorf("not approved: %s", updated.Status)
	}

	_, body = e.do("GET", "/api/vendors", "", nil)
	r := decode[struct {
		Items []*domain.Vendor `json:"items"`
	}](t, body)
	if len(r.Items) != 1 {
		t.Errorf("want 1 approved, got %d", len(r.Items))
	}
}

func TestPhotoUploadAndServe(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})

	png := pngFixture()
	res, body := e.postMultipart("/api/vendor/photos", vend, "photo", "test.png", "image/png", png)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("upload: %d %s", res.StatusCode, body)
	}
	p := decode[domain.Photo](t, body)
	if p.ID == "" || p.Size != int64(len(png)) {
		t.Errorf("bad photo: %+v", p)
	}

	res, raw := e.do("GET", "/api/photos/"+p.ID, "", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("serve: %d", res.StatusCode)
	}
	if got := res.Header.Get("Content-Type"); got != "image/png" {
		t.Errorf("content-type: %s", got)
	}
	if len(raw) != len(png) {
		t.Errorf("bytes len %d want %d", len(raw), len(png))
	}
}

func TestPhotoUploadRejectsNonImage(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})
	res, _ := e.postMultipart("/api/vendor/photos", vend, "photo", "x.txt", "text/plain", []byte("hello"))
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("want 400 got %d", res.StatusCode)
	}
}

func TestPhotoUploadRequiresVendorProfile(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	res, _ := e.postMultipart("/api/vendor/photos", vend, "photo", "x.png", "image/png", pngFixture())
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("want 404 got %d", res.StatusCode)
	}
}

func TestVendorDetailShowsPhotoIds(t *testing.T) {
	e := newEnv(t)
	vend := e.signup("v@b.kz", "password123", "V", "vendor")
	_, body := e.do("POST", "/api/vendor", vend, map[string]any{
		"name": "Co", "category": "Venue", "city": "Almaty",
	})
	v := decode[domain.Vendor](t, body)
	_, _ = e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{"status": "approved"})

	_, pbody := e.postMultipart("/api/vendor/photos", vend, "photo", "x.png", "image/png", pngFixture())
	p := decode[domain.Photo](t, pbody)

	_, dbody := e.do("GET", "/api/vendors/"+v.ID, "", nil)
	got := decode[domain.Vendor](t, dbody)
	if len(got.PhotoIDs) != 1 || got.PhotoIDs[0] != p.ID {
		t.Errorf("photoIds mismatch: %v want [%s]", got.PhotoIDs, p.ID)
	}
}
