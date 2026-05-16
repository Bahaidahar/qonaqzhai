package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
	"qonaqzhai-backend/internal/usecase/auth"
	"qonaqzhai-backend/internal/usecase/inmem"
)

func newSvc(t *testing.T) (*auth.Service, *inmem.UserRepo, *inmem.RefreshTokenRepo, *inmem.PasswordResetRepo, *inmem.FixedClock) {
	t.Helper()
	users := inmem.NewUserRepo(inmem.NewSeqIDGen("u-").New)
	rts := inmem.NewRefreshTokenRepo(inmem.NewSeqIDGen("rt-").New)
	prs := inmem.NewPasswordResetRepo(inmem.NewSeqIDGen("pr-").New)
	clock := &inmem.FixedClock{T: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	svc := auth.New(auth.Deps{
		Users:          users,
		Refresh:        rts,
		PasswordResets: prs,
		Hasher:         inmem.PlainHasher{},
		Tokens:         inmem.NewFakeIssuer(),
		Clock:          clock,
		IDs:            inmem.NewSeqIDGen("id-"),
		AccessTTL:      time.Hour,
		RefreshTTL:     24 * time.Hour,
		ResetTTL:       time.Hour,
	})
	return svc, users, rts, prs, clock
}

func TestSignupValidatesInput(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()

	cases := []struct {
		name             string
		email, pw, dname string
		role             string
		wantErr          error
	}{
		{"invalid email", "not-an-email", "password123", "Bob", "customer", domain.ErrInvalidInput},
		{"short password", "x@y.com", "short", "Bob", "customer", domain.ErrInvalidInput},
		{"ok customer", "x@y.com", "password123", "Bob", "customer", nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := svc.Signup(ctx, auth.SignupInput{
				Email: c.email, Password: c.pw, Name: c.dname, Role: c.role,
			})
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("err=%v want %v", err, c.wantErr)
			}
		})
	}
}

func TestSignupCreatesUserAndReturnsTokens(t *testing.T) {
	t.Parallel()
	svc, users, _, _, _ := newSvc(t)
	ctx := context.Background()

	out, err := svc.Signup(ctx, auth.SignupInput{
		Email: "aigerim@example.kz", Password: "password123", Name: "Aigerim", Role: "customer",
	})
	if err != nil {
		t.Fatalf("signup err: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatal("expected both tokens issued")
	}
	if out.User.Email != "aigerim@example.kz" {
		t.Errorf("email=%q", out.User.Email)
	}
	if out.User.Role != domain.RoleCustomer {
		t.Errorf("role=%q", out.User.Role)
	}

	stored, err := users.FindByEmail(ctx, "aigerim@example.kz")
	if err != nil {
		t.Fatalf("user not stored: %v", err)
	}
	if stored.PasswordHash == "password123" {
		t.Fatal("password stored in plain text!")
	}
}

func TestSignupDefaultsToCustomerWhenRoleInvalid(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	out, err := svc.Signup(context.Background(), auth.SignupInput{
		Email: "x@y.com", Password: "password123", Role: "admin", // not allowed at signup
	})
	if err != nil {
		t.Fatalf("signup err: %v", err)
	}
	if out.User.Role != domain.RoleCustomer {
		t.Errorf("role=%q want customer", out.User.Role)
	}
}

func TestSignupRejectsDuplicateEmail(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()
	_, err := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	if !errors.Is(err, domain.ErrAlreadyExists) {
		t.Errorf("err=%v want ErrAlreadyExists", err)
	}
}

func TestLoginVerifiesCredentials(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()

	_, err := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := svc.Login(ctx, "x@y.com", "WRONG"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Errorf("expected bad credentials, got %v", err)
	}
	if _, err := svc.Login(ctx, "nobody@y.com", "password123"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Errorf("expected bad credentials for unknown email, got %v", err)
	}
	out, err := svc.Login(ctx, "x@y.com", "password123")
	if err != nil {
		t.Fatalf("login err: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Error("expected both tokens")
	}
}

func TestLoginRejectsSuspendedUsers(t *testing.T) {
	t.Parallel()
	svc, users, _, _, _ := newSvc(t)
	ctx := context.Background()
	out, _ := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	_ = users.UpdateStatus(ctx, out.User.ID, domain.UserSuspended)
	if _, err := svc.Login(ctx, "x@y.com", "password123"); !errors.Is(err, domain.ErrSuspended) {
		t.Errorf("err=%v want ErrSuspended", err)
	}
}

func TestRefreshIssuesNewTokensAndRevokesOld(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()
	out, err := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	r1, err := svc.Refresh(ctx, out.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if r1.RefreshToken == out.RefreshToken {
		t.Error("refresh token should rotate")
	}
	if r1.AccessToken == "" {
		t.Error("expected access token")
	}
	// old refresh must no longer work
	if _, err := svc.Refresh(ctx, out.RefreshToken); !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("old refresh err=%v want ErrUnauthorized", err)
	}
}

func TestLogoutRevokesRefreshToken(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()
	out, _ := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})
	if err := svc.Logout(ctx, out.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Refresh(ctx, out.RefreshToken); !errors.Is(err, domain.ErrUnauthorized) {
		t.Error("revoked refresh accepted")
	}
}

type captureMailer struct {
	calls []struct{ To, Subject, Body string }
}

func (c *captureMailer) Send(_ context.Context, to, subj, body string) error {
	c.calls = append(c.calls, struct{ To, Subject, Body string }{to, subj, body})
	return nil
}

func TestForgotPasswordIssuesTokenAndEmail(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()
	_, _ = svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})

	mailer := &captureMailer{}
	svc.SetMailer(mailer)

	tok, err := svc.ForgotPassword(ctx, "x@y.com")
	if err != nil {
		t.Fatalf("forgot: %v", err)
	}
	if tok == "" {
		t.Error("expected raw token returned for direct delivery")
	}
	if len(mailer.calls) != 1 {
		t.Fatalf("expected 1 email, got %d", len(mailer.calls))
	}
	if mailer.calls[0].To != "x@y.com" {
		t.Errorf("recipient=%q", mailer.calls[0].To)
	}
}

func TestForgotPasswordSilentForUnknownEmail(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	mailer := &captureMailer{}
	svc.SetMailer(mailer)
	// must NOT leak existence — no error
	if _, err := svc.ForgotPassword(context.Background(), "nobody@y.com"); err != nil {
		t.Errorf("forgot for unknown email returned err: %v (should silently succeed)", err)
	}
	if len(mailer.calls) != 0 {
		t.Error("should not email unknown address")
	}
}

func TestResetPasswordUpdatesHashAndRevokesRefresh(t *testing.T) {
	t.Parallel()
	svc, _, _, _, _ := newSvc(t)
	ctx := context.Background()
	out, _ := svc.Signup(ctx, auth.SignupInput{Email: "x@y.com", Password: "password123"})

	tok, err := svc.ForgotPassword(ctx, "x@y.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.ResetPassword(ctx, tok, "newpassword456"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := svc.Login(ctx, "x@y.com", "password123"); !errors.Is(err, domain.ErrBadCredentials) {
		t.Error("old password still works")
	}
	if _, err := svc.Login(ctx, "x@y.com", "newpassword456"); err != nil {
		t.Errorf("new password rejected: %v", err)
	}
	// reset token must be single-use
	if err := svc.ResetPassword(ctx, tok, "another789x"); !errors.Is(err, domain.ErrUnauthorized) {
		t.Errorf("token reused: err=%v", err)
	}
	// old refresh revoked
	if _, err := svc.Refresh(ctx, out.RefreshToken); !errors.Is(err, domain.ErrUnauthorized) {
		t.Error("refresh after password reset still valid")
	}
}

// keep usecase pkg referenced (avoids unused import on partial builds)
var _ = usecase.Claims{}
