package tests

import (
	"net/http"
	"testing"

	"qonaqzhai-backend/internal/domain"
)

func TestSignupValidation(t *testing.T) {
	e := newEnv(t)
	cases := []struct {
		name string
		body map[string]any
		want int
	}{
		{"invalid email", map[string]any{"email": "no-at", "password": "password123", "role": "customer"}, 400},
		{"short password", map[string]any{"email": "a@b.kz", "password": "short", "role": "customer"}, 400},
		{"missing email", map[string]any{"password": "password123", "role": "customer"}, 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, _ := e.do("POST", "/api/signup", "", tc.body)
			if res.StatusCode != tc.want {
				t.Errorf("want %d got %d", tc.want, res.StatusCode)
			}
		})
	}
}

func TestSignupDuplicate(t *testing.T) {
	e := newEnv(t)
	_ = e.signup("a@b.kz", "password123", "A", "customer")
	res, _ := e.do("POST", "/api/signup", "", map[string]any{
		"email": "a@b.kz", "password": "password123", "name": "B", "role": "customer",
	})
	if res.StatusCode != http.StatusConflict {
		t.Errorf("want 409 got %d", res.StatusCode)
	}
}

func TestSignupDefaultsToCustomer(t *testing.T) {
	e := newEnv(t)
	res, body := e.do("POST", "/api/signup", "", map[string]any{
		"email": "a@b.kz", "password": "password123", "role": "garbage",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("want 201 got %d: %s", res.StatusCode, body)
	}
	r := decode[authResp](t, body)
	if r.User.Role != domain.RoleCustomer {
		t.Errorf("want customer fallback got %s", r.User.Role)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	e := newEnv(t)
	_ = e.signup("a@b.kz", "password123", "A", "customer")
	res, _ := e.do("POST", "/api/login", "", map[string]any{
		"email": "a@b.kz", "password": "wrong",
	})
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401 got %d", res.StatusCode)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("POST", "/api/login", "", map[string]any{
		"email": "nope@x.kz", "password": "password123",
	})
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401 got %d", res.StatusCode)
	}
}

func TestSuspendedUserCannotLogin(t *testing.T) {
	e := newEnv(t)
	_ = e.signup("c@b.kz", "password123", "C", "customer")
	uid := e.userIDByEmail("c@b.kz")
	res, _ := e.do("PATCH", "/api/admin/users/"+uid, e.adminTok, map[string]any{"status": "suspended"})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("suspend failed: %d", res.StatusCode)
	}
	res, _ = e.do("POST", "/api/login", "", map[string]any{
		"email": "c@b.kz", "password": "password123",
	})
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("want 403 got %d", res.StatusCode)
	}
}

func TestMeRequiresAuth(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("GET", "/api/me", "", nil)
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401 got %d", res.StatusCode)
	}
}

func TestMeReturnsUser(t *testing.T) {
	e := newEnv(t)
	tok := e.signup("a@b.kz", "password123", "Aigerim", "customer")
	res, body := e.do("GET", "/api/me", tok, nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("want 200 got %d", res.StatusCode)
	}
	u := decode[domain.User](t, body)
	if u.Email != "a@b.kz" || u.Name != "Aigerim" || u.Role != domain.RoleCustomer {
		t.Errorf("unexpected user: %+v", u)
	}
}

func TestInvalidToken(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("GET", "/api/me", "garbage", nil)
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("want 401 got %d", res.StatusCode)
	}
}

func TestRefreshTokenRotates(t *testing.T) {
	e := newEnv(t)
	a := e.signupFull("a@b.kz", "password123", "A", "customer")
	if a.RefreshToken == "" {
		t.Fatal("signup did not return refresh token")
	}

	res, body := e.do("POST", "/api/auth/refresh", "", map[string]any{"refreshToken": a.RefreshToken})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("refresh: %d %s", res.StatusCode, body)
	}
	r := decode[authResp](t, body)
	if r.RefreshToken == "" || r.RefreshToken == a.RefreshToken {
		t.Errorf("refresh did not rotate: %s == %s", r.RefreshToken, a.RefreshToken)
	}
	// old refresh must be revoked
	res, _ = e.do("POST", "/api/auth/refresh", "", map[string]any{"refreshToken": a.RefreshToken})
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("revoked refresh accepted: %d", res.StatusCode)
	}
}

func TestLogoutRevokesRefresh(t *testing.T) {
	e := newEnv(t)
	a := e.signupFull("a@b.kz", "password123", "A", "customer")
	res, _ := e.do("POST", "/api/auth/logout", "", map[string]any{"refreshToken": a.RefreshToken})
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("logout: %d", res.StatusCode)
	}
	res, _ = e.do("POST", "/api/auth/refresh", "", map[string]any{"refreshToken": a.RefreshToken})
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("refresh after logout: %d", res.StatusCode)
	}
}

func TestForgotPasswordSilentForUnknown(t *testing.T) {
	e := newEnv(t)
	res, _ := e.do("POST", "/api/auth/forgot-password", "", map[string]any{"email": "unknown@x.kz"})
	if res.StatusCode != http.StatusOK {
		t.Errorf("forgot password leaks for unknown email: %d", res.StatusCode)
	}
}
