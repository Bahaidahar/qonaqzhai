package domain

import (
	"errors"
	"testing"
)

func TestVendorStatusValid(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   VendorStatus
		want bool
	}{
		{VendorPending, true},
		{VendorApproved, true},
		{VendorRejected, true},
		{VendorStatus(""), false},
		{VendorStatus("ban"), false},
	}
	for _, c := range cases {
		if got := c.in.Valid(); got != c.want {
			t.Errorf("VendorStatus(%q).Valid()=%v want %v", c.in, got, c.want)
		}
	}
}

func TestVendorIsPublic(t *testing.T) {
	t.Parallel()
	if !(&Vendor{Status: VendorApproved}).IsPublic() {
		t.Error("approved should be public")
	}
	if (&Vendor{Status: VendorPending}).IsPublic() {
		t.Error("pending must not be public")
	}
	if (&Vendor{Status: VendorRejected}).IsPublic() {
		t.Error("rejected must not be public")
	}
}

func TestVendorInputNormalize(t *testing.T) {
	t.Parallel()
	in := &VendorInput{Name: "  Rixos ", Category: " Venue ", City: " Almaty ", Description: " desc "}
	in.Normalize()
	if in.Name != "Rixos" || in.Category != "Venue" || in.City != "Almaty" || in.Description != "desc" {
		t.Errorf("normalize failed: %+v", in)
	}
}

func TestVendorInputValidate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   VendorInput
		err  error
	}{
		{"ok", VendorInput{Name: "X", Category: "Y", City: "Z", PriceFrom: 100}, nil},
		{"empty name", VendorInput{Category: "Y", City: "Z"}, ErrInvalidInput},
		{"empty cat", VendorInput{Name: "X", City: "Z"}, ErrInvalidInput},
		{"empty city", VendorInput{Name: "X", Category: "Y"}, ErrInvalidInput},
		{"negative price", VendorInput{Name: "X", Category: "Y", City: "Z", PriceFrom: -1}, ErrInvalidInput},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.in.Validate(); !errors.Is(err, c.err) {
				t.Errorf("err=%v want %v", err, c.err)
			}
		})
	}
}
