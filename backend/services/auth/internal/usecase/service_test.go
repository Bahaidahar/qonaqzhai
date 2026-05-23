package usecase_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/errs"

	"qonaqzhai-backend/services/auth/internal/domain"
	"qonaqzhai-backend/services/auth/internal/ports"
	"qonaqzhai-backend/services/auth/internal/usecase"
)

// --- in-memory fakes ---------------------------------------------------------

type memUsers struct {
	mu       sync.Mutex
	byID     map[string]*domain.User
	byEmail  map[string]*domain.User
	pwHashes map[string]string
}

func newMemUsers() *memUsers {
	return &memUsers{
		byID:     map[string]*domain.User{},
		byEmail:  map[string]*domain.User{},
		pwHashes: map[string]string{},
	}
}

func (m *memUsers) Create(_ context.Context, u *domain.User) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.byEmail[u.Email]; exists {
		return nil, errs.ErrAlreadyExists
	}
	if u.ID == "" {
		u.ID = "u-" + strconv.Itoa(len(m.byID)+1)
	}
	cp := *u
	m.byID[cp.ID] = &cp
	m.byEmail[cp.Email] = &cp
	m.pwHashes[cp.ID] = cp.PasswordHash
	return &cp, nil
}

func (m *memUsers) FindByID(_ context.Context, id string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byID[id]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *u
	cp.PasswordHash = m.pwHashes[id]
	return &cp, nil
}

func (m *memUsers) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byEmail[email]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *u
	cp.PasswordHash = m.pwHashes[u.ID]
	return &cp, nil
}

func (m *memUsers) FindByIDs(_ context.Context, ids []string) ([]*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.User{}
	for _, id := range ids {
		if u, ok := m.byID[id]; ok {
			cp := *u
			out = append(out, &cp)
		}
	}
	return out, nil
}

func (m *memUsers) List(_ context.Context, opts ports.ListUsersOpts) ([]*domain.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := []*domain.User{}
	for _, u := range m.byID {
		if opts.Role != "" && string(u.Role) != opts.Role {
			continue
		}
		cp := *u
		out = append(out, &cp)
	}
	return out, nil
}

func (m *memUsers) UpdateStatus(_ context.Context, id string, status domain.UserStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	u, ok := m.byID[id]
	if !ok {
		return errs.ErrNotFound
	}
	u.Status = status
	return nil
}

func (m *memUsers) UpdatePasswordHash(_ context.Context, id, hash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.byID[id]; !ok {
		return errs.ErrNotFound
	}
	m.pwHashes[id] = hash
	return nil
}

type memRefresh struct {
	mu     sync.Mutex
	rows   map[string]*domain.RefreshToken
	byHash map[string]*domain.RefreshToken
}

func newMemRefresh() *memRefresh {
	return &memRefresh{rows: map[string]*domain.RefreshToken{}, byHash: map[string]*domain.RefreshToken{}}
}

func (r *memRefresh) Create(_ context.Context, t *domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.ID == "" {
		t.ID = "rt-" + strconv.Itoa(len(r.rows)+1)
	}
	cp := *t
	r.rows[cp.ID] = &cp
	r.byHash[cp.TokenHash] = &cp
	return nil
}

func (r *memRefresh) FindActiveByHash(_ context.Context, hash string, now time.Time) (*domain.RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.byHash[hash]
	if !ok || !t.Active(now) {
		return nil, errs.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *memRefresh) Revoke(_ context.Context, id string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	t.RevokedAt = &at
	return nil
}

func (r *memRefresh) RevokeAllForUser(_ context.Context, userID string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.rows {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &at
		}
	}
	return nil
}

type memResets struct {
	mu     sync.Mutex
	rows   map[string]*domain.PasswordResetToken
	byHash map[string]*domain.PasswordResetToken
}

func newMemResets() *memResets {
	return &memResets{rows: map[string]*domain.PasswordResetToken{}, byHash: map[string]*domain.PasswordResetToken{}}
}

func (r *memResets) Create(_ context.Context, t *domain.PasswordResetToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.ID == "" {
		t.ID = "pr-" + strconv.Itoa(len(r.rows)+1)
	}
	cp := *t
	r.rows[cp.ID] = &cp
	r.byHash[cp.TokenHash] = &cp
	return nil
}

func (r *memResets) FindByHash(_ context.Context, hash string) (*domain.PasswordResetToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.byHash[hash]
	if !ok {
		return nil, errs.ErrNotFound
	}
	cp := *t
	return &cp, nil
}

func (r *memResets) MarkUsed(_ context.Context, id string, at time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.rows[id]
	if !ok {
		return errs.ErrNotFound
	}
	t.UsedAt = &at
	return nil
}

type plainHasher struct{}

func (plainHasher) Hash(p string) (string, error) { return "h:" + p, nil }
func (plainHasher) Verify(hash, p string) error {
	if hash == "h:"+p {
		return nil
	}
	return errors.New("mismatch")
}

type fixedClock struct{ t time.Time }

func (c *fixedClock) Now() time.Time { return c.t }

type seqIDs struct {
	mu sync.Mutex
	n  int
}

func (s *seqIDs) New() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	return "id-" + strconv.Itoa(s.n)
}

// --- helpers -----------------------------------------------------------------

func newSvc(t *testing.T) (*usecase.Service, *memUsers, *fixedClock) {
	t.Helper()
	clk := &fixedClock{t: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	users := newMemUsers()
	svc := usecase.New(usecase.Deps{
		Users:          users,
		Refresh:        newMemRefresh(),
		PasswordResets: newMemResets(),
		Hasher:         plainHasher{},
		Signer:         pkgauth.NewJWTSigner([]byte("test-secret"), "test"),
		Clock:          clk,
		IDs:            &seqIDs{},
		AccessTTL:      15 * time.Minute,
		RefreshTTL:     time.Hour,
		ResetTTL:       30 * time.Minute,
	})
	return svc, users, clk
}

// --- tests -------------------------------------------------------------------

func TestSignupLoginRefresh(t *testing.T) {
	svc, _, _ := newSvc(t)
	ctx := context.Background()

	out, err := svc.Signup(ctx, usecase.SignupInput{
		Email: "a@b.kz", Password: "password123", Name: "Aigerim",
	})
	if err != nil {
		t.Fatalf("signup: %v", err)
	}
	if out.AccessToken == "" || out.RefreshToken == "" {
		t.Fatalf("expected both tokens")
	}
	if out.User.Role != domain.RoleCustomer {
		t.Fatalf("expected default role customer, got %q", out.User.Role)
	}

	login, err := svc.Login(ctx, "A@B.KZ", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if login.User.ID != out.User.ID {
		t.Fatalf("user id mismatch after login")
	}

	if _, err := svc.Login(ctx, "a@b.kz", "wrong"); !errors.Is(err, errs.ErrBadCredentials) {
		t.Fatalf("expected bad credentials, got %v", err)
	}

	refreshed, err := svc.Refresh(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if refreshed.RefreshToken == login.RefreshToken {
		t.Fatalf("refresh token should rotate")
	}
	if _, err := svc.Refresh(ctx, login.RefreshToken); !errors.Is(err, errs.ErrUnauthorized) {
		t.Fatalf("old refresh token must be revoked, got %v", err)
	}
}

func TestSuspendedUserCannotLogin(t *testing.T) {
	svc, users, _ := newSvc(t)
	ctx := context.Background()
	out, err := svc.Signup(ctx, usecase.SignupInput{Email: "x@y.kz", Password: "password123"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.SetUserStatus(ctx, out.User.ID, domain.UserSuspended); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Login(ctx, "x@y.kz", "password123"); !errors.Is(err, errs.ErrSuspended) {
		t.Fatalf("expected suspended, got %v", err)
	}
	// Sanity check: the underlying record indeed flipped.
	got, _ := users.FindByID(ctx, out.User.ID)
	if got.Status != domain.UserSuspended {
		t.Fatalf("status not persisted")
	}
}

func TestPasswordResetFlow(t *testing.T) {
	svc, _, _ := newSvc(t)
	ctx := context.Background()
	if _, err := svc.Signup(ctx, usecase.SignupInput{Email: "p@q.kz", Password: "password123"}); err != nil {
		t.Fatal(err)
	}
	token, err := svc.ForgotPassword(ctx, "p@q.kz")
	if err != nil || token == "" {
		t.Fatalf("forgot password: token=%q err=%v", token, err)
	}
	if err := svc.ResetPassword(ctx, token, "newpassword12"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if _, err := svc.Login(ctx, "p@q.kz", "newpassword12"); err != nil {
		t.Fatalf("login with new password: %v", err)
	}
	if _, err := svc.Login(ctx, "p@q.kz", "password123"); !errors.Is(err, errs.ErrBadCredentials) {
		t.Fatalf("old password should not work, got %v", err)
	}
	// Reset token is single-use.
	if err := svc.ResetPassword(ctx, token, "another1234"); !errors.Is(err, errs.ErrUnauthorized) {
		t.Fatalf("token must be single-use, got %v", err)
	}
}

func TestForgotPasswordUnknownEmailNoLeak(t *testing.T) {
	svc, _, _ := newSvc(t)
	token, err := svc.ForgotPassword(context.Background(), "nobody@nowhere.kz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if token != "" {
		t.Fatalf("token must be empty for unknown email, got %q", token)
	}
}

func TestEnsureAdminIdempotent(t *testing.T) {
	svc, _, _ := newSvc(t)
	ctx := context.Background()
	u1, err := svc.EnsureAdmin(ctx, "admin@qonaqzhai.kz", "password123", "Admin")
	if err != nil {
		t.Fatal(err)
	}
	if u1.Role != domain.RoleAdmin {
		t.Fatalf("expected admin, got %q", u1.Role)
	}
	u2, err := svc.EnsureAdmin(ctx, "admin@qonaqzhai.kz", "different456", "Admin")
	if err != nil {
		t.Fatal(err)
	}
	if u2.ID != u1.ID {
		t.Fatalf("second call should return existing admin, got new id")
	}
}
