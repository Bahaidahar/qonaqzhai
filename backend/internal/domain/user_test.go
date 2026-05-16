package domain

import "testing"

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"  Foo@Bar.KZ ": "foo@bar.kz",
		"x@y.com":       "x@y.com",
		"":              "",
	}
	for in, want := range cases {
		if got := NormalizeEmail(in); got != want {
			t.Errorf("NormalizeEmail(%q)=%q want %q", in, got, want)
		}
	}
}

func TestValidEmail(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in   string
		want bool
	}{
		{"foo@bar.kz", true},
		{"a@b.c", true},
		{"", false},
		{"no-at-sign", false},
		{"@nolocal.com", false},
		{"nolocal@", false},
		{"missing-dot@host", false},
	}
	for _, c := range cases {
		if got := ValidEmail(c.in); got != c.want {
			t.Errorf("ValidEmail(%q)=%v want %v", c.in, got, c.want)
		}
	}
}

func TestValidPassword(t *testing.T) {
	t.Parallel()
	if ValidPassword("short") {
		t.Error("short password accepted")
	}
	if !ValidPassword("password123") {
		t.Error("8+ password rejected")
	}
}

func TestDefaultName(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name, email, want string
	}{
		{"Aigerim", "x@y.com", "Aigerim"},
		{"  ", "aigerim@example.kz", "aigerim"},
		{"", "noattag", "noattag"},
		{"  Spaced  ", "x@y.com", "Spaced"},
	}
	for _, c := range cases {
		if got := DefaultName(c.name, c.email); got != c.want {
			t.Errorf("DefaultName(%q,%q)=%q want %q", c.name, c.email, got, c.want)
		}
	}
}

func TestUserIsActive(t *testing.T) {
	t.Parallel()
	if !(&User{Status: UserActive}).IsActive() {
		t.Error("active user reported inactive")
	}
	if (&User{Status: UserSuspended}).IsActive() {
		t.Error("suspended user reported active")
	}
}
