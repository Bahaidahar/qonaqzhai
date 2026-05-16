package domain

import "testing"

func TestRoleValid(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   Role
		want bool
	}{
		{RoleCustomer, true},
		{RoleVendor, true},
		{RoleAdmin, true},
		{Role(""), false},
		{Role("guest"), false},
	}
	for _, c := range cases {
		if got := c.in.Valid(); got != c.want {
			t.Errorf("Role(%q).Valid() = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseRole(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want Role
	}{
		{"customer", RoleCustomer},
		{"vendor", RoleVendor},
		{"admin", RoleCustomer}, // admin not self-signupable → fallback
		{"", RoleCustomer},
		{"hacker", RoleCustomer},
	}
	for _, c := range cases {
		if got := ParseRole(c.in); got != c.want {
			t.Errorf("ParseRole(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
