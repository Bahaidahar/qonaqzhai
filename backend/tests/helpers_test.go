package tests

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"qonaqzhai-backend/internal/app"
	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/infra/config"
	"qonaqzhai-backend/internal/infra/db/testpg"
)

type env struct {
	t        *testing.T
	app      *app.App
	ts       *httptest.Server
	adminTok string
}

func newEnv(t *testing.T) *env {
	t.Helper()
	dsn := testpg.Start(t)
	cfg := config.Config{
		DatabaseURL:       dsn,
		JWTSecret:         "test-secret-test-secret-test-secret",
		CORSOrigin:        "*",
		AccessTTL:         time.Hour,
		RefreshTTL:        24 * time.Hour,
		ResetTTL:          time.Hour,
		RateLimitDisabled: true,
		BcryptCost:        4, // fast hashing in tests; production uses 12
	}
	a, err := app.New(context.Background(), cfg)
	if err != nil {
		t.Fatalf("app: %v", err)
	}
	t.Cleanup(func() { _ = a.Close() })

	ts := httptest.NewServer(a.Handler)
	t.Cleanup(ts.Close)

	e := &env{t: t, app: a, ts: ts}
	e.adminTok = e.login("admin@qonaqzhai.kz", "admin12345")
	return e
}

func (e *env) do(method, path, token string, body any) (*http.Response, []byte) {
	e.t.Helper()
	var rdr io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			e.t.Fatalf("marshal: %v", err)
		}
		rdr = bytes.NewReader(buf)
	}
	req, err := http.NewRequest(method, e.ts.URL+path, rdr)
	if err != nil {
		e.t.Fatalf("req: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatalf("do: %v", err)
	}
	defer res.Body.Close()
	out, _ := io.ReadAll(res.Body)
	return res, out
}

func (e *env) postMultipart(path, token, field, filename, mime string, data []byte) (*http.Response, []byte) {
	e.t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{fmt.Sprintf(`form-data; name=%q; filename=%q`, field, filename)}
	h["Content-Type"] = []string{mime}
	part, err := w.CreatePart(h)
	if err != nil {
		e.t.Fatalf("createpart: %v", err)
	}
	if _, err := part.Write(data); err != nil {
		e.t.Fatalf("write: %v", err)
	}
	w.Close()

	req, _ := http.NewRequest("POST", e.ts.URL+path, &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e.t.Fatalf("do mp: %v", err)
	}
	defer res.Body.Close()
	out, _ := io.ReadAll(res.Body)
	return res, out
}

func decode[T any](t *testing.T, body []byte) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("decode: %v body=%s", err, string(body))
	}
	return v
}

type authResp struct {
	Token        string       `json:"token"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *domain.User `json:"user"`
}

func (e *env) signup(email, password, name, role string) string {
	e.t.Helper()
	res, body := e.do("POST", "/api/signup", "", map[string]any{
		"email": email, "password": password, "name": name, "role": role,
	})
	if res.StatusCode != http.StatusCreated {
		e.t.Fatalf("signup %s want 201 got %d: %s", email, res.StatusCode, body)
	}
	return decode[authResp](e.t, body).Token
}

func (e *env) signupFull(email, password, name, role string) authResp {
	e.t.Helper()
	res, body := e.do("POST", "/api/signup", "", map[string]any{
		"email": email, "password": password, "name": name, "role": role,
	})
	if res.StatusCode != http.StatusCreated {
		e.t.Fatalf("signup %s want 201 got %d: %s", email, res.StatusCode, body)
	}
	return decode[authResp](e.t, body)
}

func (e *env) login(email, password string) string {
	e.t.Helper()
	res, body := e.do("POST", "/api/login", "", map[string]any{
		"email": email, "password": password,
	})
	if res.StatusCode != http.StatusOK {
		e.t.Fatalf("login %s want 200 got %d: %s", email, res.StatusCode, body)
	}
	return decode[authResp](e.t, body).Token
}

func (e *env) loginFull(email, password string) authResp {
	e.t.Helper()
	res, body := e.do("POST", "/api/login", "", map[string]any{
		"email": email, "password": password,
	})
	if res.StatusCode != http.StatusOK {
		e.t.Fatalf("login %s want 200 got %d: %s", email, res.StatusCode, body)
	}
	return decode[authResp](e.t, body)
}

func pngFixture() []byte {
	const b64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="
	dec, _ := base64.StdEncoding.DecodeString(b64)
	return dec
}

// userIDByEmail fetches the user id via the admin /api/admin/users endpoint.
// Useful for tests that need to refer to a freshly-created user without exposing the DB.
func (e *env) userIDByEmail(email string) string {
	e.t.Helper()
	_, body := e.do("GET", "/api/admin/users", e.adminTok, nil)
	r := decode[struct {
		Items []domain.User `json:"items"`
	}](e.t, body)
	for _, u := range r.Items {
		if u.Email == email {
			return u.ID
		}
	}
	e.t.Fatalf("user %s not found via admin list", email)
	return ""
}
