package tests

import (
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func seedApprovedVendor(t *testing.T, e *env, email, name, cat, city string, price int64) string {
	t.Helper()
	tok := e.signup(email, "password123", name, "vendor")
	_, body := e.do("POST", "/api/vendor", tok, map[string]any{
		"name": name, "category": cat, "city": city, "priceFrom": price, "description": "Pro " + cat,
	})
	v := decode[domain.Vendor](t, body)
	_, _ = e.do("PATCH", "/api/admin/vendors/"+v.ID, e.adminTok, map[string]any{"status": "approved"})
	return v.ID
}

func TestSearchByCategory(t *testing.T) {
	e := newEnv(t)
	seedApprovedVendor(t, e, "a@b.kz", "Rixos", "Venue", "Almaty", 1_500_000)
	seedApprovedVendor(t, e, "b@b.kz", "Aitu", "Photo", "Almaty", 450_000)

	_, body := e.do("GET", "/api/vendors?category=Photo", "", nil)
	r := decode[struct {
		Items []domain.Vendor `json:"items"`
		Total int             `json:"total"`
	}](t, body)
	if r.Total != 1 || len(r.Items) != 1 || r.Items[0].Name != "Aitu" {
		t.Errorf("category filter: %+v", r)
	}
}

func TestSearchByPriceRange(t *testing.T) {
	e := newEnv(t)
	seedApprovedVendor(t, e, "a@b.kz", "Cheap", "Photo", "Almaty", 100_000)
	seedApprovedVendor(t, e, "b@b.kz", "Mid", "Photo", "Almaty", 500_000)
	seedApprovedVendor(t, e, "c@b.kz", "Lux", "Photo", "Almaty", 2_000_000)

	_, body := e.do("GET", "/api/vendors?price_min=200000&price_max=1000000", "", nil)
	r := decode[struct {
		Items []domain.Vendor `json:"items"`
	}](t, body)
	if len(r.Items) != 1 || r.Items[0].Name != "Mid" {
		t.Errorf("price range: %+v", r.Items)
	}
}

func TestSearchByFullText(t *testing.T) {
	e := newEnv(t)
	seedApprovedVendor(t, e, "a@b.kz", "Rixos Almaty Ballroom", "Venue", "Almaty", 1_000_000)
	seedApprovedVendor(t, e, "b@b.kz", "Studio Aitu", "Photo", "Almaty", 500_000)

	_, body := e.do("GET", "/api/vendors?q=Rixos", "", nil)
	r := decode[struct {
		Items []domain.Vendor `json:"items"`
	}](t, body)
	if len(r.Items) != 1 || r.Items[0].Name != "Rixos Almaty Ballroom" {
		t.Errorf("fts search: %+v", r.Items)
	}
}

func TestSearchSortByPrice(t *testing.T) {
	e := newEnv(t)
	seedApprovedVendor(t, e, "a@b.kz", "B", "Photo", "Almaty", 500_000)
	seedApprovedVendor(t, e, "b@b.kz", "A", "Photo", "Almaty", 200_000)
	seedApprovedVendor(t, e, "c@b.kz", "C", "Photo", "Almaty", 800_000)

	_, body := e.do("GET", "/api/vendors?sort=price_asc", "", nil)
	r := decode[struct {
		Items []domain.Vendor `json:"items"`
	}](t, body)
	if len(r.Items) != 3 || r.Items[0].PriceFrom != 200_000 || r.Items[2].PriceFrom != 800_000 {
		t.Errorf("sort: %+v", r.Items)
	}
}

func TestSearchPagination(t *testing.T) {
	e := newEnv(t)
	for i := 0; i < 5; i++ {
		letter := string(rune('a' + i))
		seedApprovedVendor(t, e, letter+"@b.kz", "V"+letter, "Photo", "Almaty", 100_000)
	}
	_, body := e.do("GET", "/api/vendors?page=2&limit=2", "", nil)
	r := decode[struct {
		Items []domain.Vendor `json:"items"`
		Total int             `json:"total"`
		Page  int             `json:"page"`
		Limit int             `json:"limit"`
	}](t, body)
	if r.Total != 5 || len(r.Items) != 2 || r.Page != 2 || r.Limit != 2 {
		t.Errorf("pagination: %+v", r)
	}
}
