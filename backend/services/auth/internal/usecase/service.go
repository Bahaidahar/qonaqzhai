// Package usecase implements the auth application service: signup, login,
// refresh, logout, and password recovery.
package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	pkgauth "qonaqzhai-backend/pkg/auth"
	"qonaqzhai-backend/pkg/errs"

	"qonaqzhai-backend/services/auth/internal/domain"
	"qonaqzhai-backend/services/auth/internal/ports"
)

// Deps bundles all collaborators required by the Service.
type Deps struct {
	Users          ports.UserRepo
	Refresh        ports.RefreshTokenRepo
	PasswordResets ports.PasswordResetRepo
	Hasher         ports.PasswordHasher
	Signer         *pkgauth.JWTSigner
	Clock          ports.Clock
	IDs            ports.IDGen
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	ResetTTL       time.Duration
}

// Service is the auth application service.
type Service struct {
	d      Deps
	mu     sync.RWMutex
	mailer ports.Mailer
}

// New constructs a Service. Zero TTLs fall back to sane defaults.
func New(d Deps) *Service {
	if d.AccessTTL == 0 {
		d.AccessTTL = 15 * time.Minute
	}
	if d.RefreshTTL == 0 {
		d.RefreshTTL = 30 * 24 * time.Hour
	}
	if d.ResetTTL == 0 {
		d.ResetTTL = time.Hour
	}
	return &Service{d: d}
}

// SetMailer wires a mailer for password-reset delivery. Without one,
// ForgotPassword still issues a token but does not email it.
func (s *Service) SetMailer(m ports.Mailer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mailer = m
}

func (s *Service) currentMailer() ports.Mailer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mailer
}

// SignupInput captures fields needed to create an account.
type SignupInput struct {
	Email    string
	Password string
	Name     string
	Role     string
}

// AuthOutput is returned from signup / login / refresh.
type AuthOutput struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
}

// Signup creates a new account and issues both tokens. Admin accounts are
// seeded out-of-band, never through public signup.
func (s *Service) Signup(ctx context.Context, in SignupInput) (*AuthOutput, error) {
	in.Email = domain.NormalizeEmail(in.Email)
	if !domain.ValidEmail(in.Email) {
		return nil, fmt.Errorf("email: %w", errs.ErrInvalidInput)
	}
	if !domain.ValidPassword(in.Password) {
		return nil, fmt.Errorf("password too short: %w", errs.ErrInvalidInput)
	}

	hash, err := s.d.Hasher.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	u := &domain.User{
		Email:        in.Email,
		Name:         domain.DefaultName(in.Name, in.Email),
		PasswordHash: hash,
		Role:         domain.ParseRole(in.Role),
		Status:       domain.UserActive,
		CreatedAt:    s.d.Clock.Now(),
	}
	created, err := s.d.Users.Create(ctx, u)
	if err != nil {
		return nil, err
	}
	return s.issueTokens(ctx, created)
}

// Login verifies credentials and returns fresh tokens.
func (s *Service) Login(ctx context.Context, email, password string) (*AuthOutput, error) {
	email = domain.NormalizeEmail(email)
	u, err := s.d.Users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, errs.ErrBadCredentials
		}
		return nil, err
	}
	if !u.IsActive() {
		return nil, errs.ErrSuspended
	}
	if err := s.d.Hasher.Verify(u.PasswordHash, password); err != nil {
		return nil, errs.ErrBadCredentials
	}
	return s.issueTokens(ctx, u)
}

// Refresh rotates the refresh token: the supplied one is revoked and a new pair
// is issued.
func (s *Service) Refresh(ctx context.Context, refresh string) (*AuthOutput, error) {
	if refresh == "" {
		return nil, errs.ErrUnauthorized
	}
	hash := hashToken(refresh)
	now := s.d.Clock.Now()
	t, err := s.d.Refresh.FindActiveByHash(ctx, hash, now)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}
	if err := s.d.Refresh.Revoke(ctx, t.ID, now); err != nil {
		return nil, err
	}
	u, err := s.d.Users.FindByID(ctx, t.UserID)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}
	if !u.IsActive() {
		return nil, errs.ErrSuspended
	}
	return s.issueTokens(ctx, u)
}

// Logout revokes the supplied refresh token. Missing tokens are silently
// ignored.
func (s *Service) Logout(ctx context.Context, refresh string) error {
	if refresh == "" {
		return nil
	}
	hash := hashToken(refresh)
	now := s.d.Clock.Now()
	t, err := s.d.Refresh.FindActiveByHash(ctx, hash, now)
	if err != nil {
		return nil
	}
	return s.d.Refresh.Revoke(ctx, t.ID, now)
}

// ForgotPassword issues a single-use reset token bound to the user's email.
// For unknown emails it returns ("", nil) — never leaking account existence.
func (s *Service) ForgotPassword(ctx context.Context, email string) (string, error) {
	email = domain.NormalizeEmail(email)
	u, err := s.d.Users.FindByEmail(ctx, email)
	if err != nil {
		return "", nil
	}
	raw, err := randomToken(32)
	if err != nil {
		return "", err
	}
	t := &domain.PasswordResetToken{
		UserID:    u.ID,
		TokenHash: hashToken(raw),
		ExpiresAt: s.d.Clock.Now().Add(s.d.ResetTTL),
	}
	if err := s.d.PasswordResets.Create(ctx, t); err != nil {
		return "", err
	}
	if m := s.currentMailer(); m != nil {
		body := fmt.Sprintf(
			"<p>Reset your password using this token (valid %v):</p><pre>%s</pre>",
			s.d.ResetTTL, raw,
		)
		_ = m.Send(ctx, u.Email, "Reset your Qonaqzhai password", body)
	}
	return raw, nil
}

// ResetPassword consumes a reset token and updates the user's password hash.
// All refresh tokens for the user are revoked as a defensive measure.
func (s *Service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	if !domain.ValidPassword(newPassword) {
		return fmt.Errorf("password too short: %w", errs.ErrInvalidInput)
	}
	hash := hashToken(rawToken)
	t, err := s.d.PasswordResets.FindByHash(ctx, hash)
	if err != nil {
		return errs.ErrUnauthorized
	}
	now := s.d.Clock.Now()
	if !t.Valid(now) {
		return errs.ErrUnauthorized
	}
	pwHash, err := s.d.Hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	if err := s.d.Users.UpdatePasswordHash(ctx, t.UserID, pwHash); err != nil {
		return err
	}
	if err := s.d.PasswordResets.MarkUsed(ctx, t.ID, now); err != nil {
		return err
	}
	return s.d.Refresh.RevokeAllForUser(ctx, t.UserID, now)
}

// VerifyAccessToken decodes a bearer token. Used by the gRPC verifier endpoint
// other services call.
func (s *Service) VerifyAccessToken(_ context.Context, raw string) (pkgauth.Claims, time.Time, error) {
	return s.d.Signer.Parse(raw)
}

// FindUser returns a user by id. Used by the gRPC GetUser endpoint.
func (s *Service) FindUser(ctx context.Context, id string) (*domain.User, error) {
	return s.d.Users.FindByID(ctx, id)
}

// FindUsers returns multiple users by id, preserving order roughly. Missing IDs
// are silently omitted.
func (s *Service) FindUsers(ctx context.Context, ids []string) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	return s.d.Users.FindByIDs(ctx, ids)
}

// ListUsersForAdmin returns paginated users.
func (s *Service) ListUsersForAdmin(ctx context.Context, opts ports.ListUsersOpts) ([]*domain.User, error) {
	if opts.Limit <= 0 || opts.Limit > 200 {
		opts.Limit = 50
	}
	return s.d.Users.List(ctx, opts)
}

// EnsureAdmin idempotently creates an admin account. Returns the existing
// user when the email is already taken (without touching the password). Called
// once at service startup from main.go when ADMIN_EMAIL / ADMIN_PASSWORD are
// set.
func (s *Service) EnsureAdmin(ctx context.Context, email, password, name string) (*domain.User, error) {
	email = domain.NormalizeEmail(email)
	if existing, err := s.d.Users.FindByEmail(ctx, email); err == nil {
		return existing, nil
	}
	if !domain.ValidEmail(email) {
		return nil, fmt.Errorf("admin email: %w", errs.ErrInvalidInput)
	}
	if !domain.ValidPassword(password) {
		return nil, fmt.Errorf("admin password too short: %w", errs.ErrInvalidInput)
	}
	hash, err := s.d.Hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	u := &domain.User{
		Email:        email,
		Name:         domain.DefaultName(name, email),
		PasswordHash: hash,
		Role:         domain.RoleAdmin,
		Status:       domain.UserActive,
		CreatedAt:    s.d.Clock.Now(),
	}
	return s.d.Users.Create(ctx, u)
}

// SetUserStatus is an admin operation.
func (s *Service) SetUserStatus(ctx context.Context, id string, status domain.UserStatus) (*domain.User, error) {
	if status != domain.UserActive && status != domain.UserSuspended {
		return nil, errs.ErrInvalidInput
	}
	if err := s.d.Users.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.d.Users.FindByID(ctx, id)
}

func (s *Service) issueTokens(ctx context.Context, u *domain.User) (*AuthOutput, error) {
	access, err := s.d.Signer.Issue(pkgauth.Claims{
		UserID: u.ID, Email: u.Email, Role: string(u.Role), Status: string(u.Status),
	}, s.d.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("issue access: %w", err)
	}
	raw, err := randomToken(32)
	if err != nil {
		return nil, err
	}
	rt := &domain.RefreshToken{
		UserID:    u.ID,
		TokenHash: hashToken(raw),
		ExpiresAt: s.d.Clock.Now().Add(s.d.RefreshTTL),
	}
	if err := s.d.Refresh.Create(ctx, rt); err != nil {
		return nil, err
	}
	return &AuthOutput{User: u, AccessToken: access, RefreshToken: raw}, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
