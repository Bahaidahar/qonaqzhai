// Package usecase defines application-level service ports (interfaces).
// Concrete usecases live in subpackages and depend only on these ports + domain.
package usecase

import (
	"context"
	"time"

	"qonaqzhai-backend/internal/domain"
)

// UserRepo is the persistence boundary for users.
type UserRepo interface {
	Create(ctx context.Context, u *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	UpdateStatus(ctx context.Context, id string, status domain.UserStatus) error
	UpdatePasswordHash(ctx context.Context, id, hash string) error
}

// VendorRepo persists vendor profiles.
type VendorRepo interface {
	Upsert(ctx context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error)
	FindByID(ctx context.Context, id string) (*domain.Vendor, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Vendor, error)
	Search(ctx context.Context, q VendorQuery) ([]*domain.Vendor, int, error)
	UpdateStatus(ctx context.Context, id string, status domain.VendorStatus) error
	UpdateRating(ctx context.Context, id string, avg float64, count int) error
}

// VendorQuery captures all catalog filters and pagination options.
type VendorQuery struct {
	Q          string
	Category   string
	City       string
	MinPrice   int64
	MaxPrice   int64
	MinRating  float64
	Status     domain.VendorStatus // empty → no filter
	Sort       string              // price_asc | price_desc | rating_desc | newest
	Page       int                 // 1-based
	Limit      int
}

// ServiceRepo persists per-vendor services (menu items).
type ServiceRepo interface {
	Create(ctx context.Context, vendorID string, in domain.ServiceInput) (*domain.Service, error)
	Update(ctx context.Context, id string, in domain.ServiceInput) (*domain.Service, error)
	FindByID(ctx context.Context, id string) (*domain.Service, error)
	ListByVendor(ctx context.Context, vendorID string, activeOnly bool) ([]*domain.Service, error)
	Delete(ctx context.Context, id string) error
	MinActivePrice(ctx context.Context, vendorID string) (int64, error)
}

// PhotoRepo persists vendor photos.
type PhotoRepo interface {
	Create(ctx context.Context, vendorID, mime string, data []byte) (*domain.Photo, error)
	Find(ctx context.Context, id string) (*domain.Photo, error)
	Delete(ctx context.Context, id string) error
	ListIDs(ctx context.Context, vendorID string) ([]string, error)
}

// BookingRepo persists bookings.
type BookingRepo interface {
	Create(ctx context.Context, b *domain.Booking) (*domain.Booking, error)
	Find(ctx context.Context, id string) (*domain.Booking, error)
	ListForCustomer(ctx context.Context, customerID string) ([]*domain.Booking, error)
	ListForVendor(ctx context.Context, vendorID string) ([]*domain.Booking, error)
	ListAll(ctx context.Context) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error
	SetPayment(ctx context.Context, id, paymentID string) error
	SetService(ctx context.Context, id, serviceID string) error
}

// ReviewRepo persists vendor reviews.
type ReviewRepo interface {
	Create(ctx context.Context, r *domain.Review) (*domain.Review, error)
	FindByID(ctx context.Context, id string) (*domain.Review, error)
	FindByBooking(ctx context.Context, bookingID string) (*domain.Review, error)
	ListByVendor(ctx context.Context, vendorID string) ([]*domain.Review, error)
	Delete(ctx context.Context, id string) error
	AggregateForVendor(ctx context.Context, vendorID string) (avg float64, count int, err error)
}

// RefreshTokenRepo persists refresh tokens (hashed).
// FindActiveByHash takes a `now` argument so that "active" is defined by the
// service clock, not the repo wall-clock. SQL implementations should embed
// `now` in a `WHERE expires_at > ? AND revoked_at IS NULL` clause.
type RefreshTokenRepo interface {
	Create(ctx context.Context, t *domain.RefreshToken) error
	FindActiveByHash(ctx context.Context, hash string, now time.Time) (*domain.RefreshToken, error)
	Revoke(ctx context.Context, id string, at time.Time) error
	RevokeAllForUser(ctx context.Context, userID string, at time.Time) error
}

// PasswordResetRepo persists password reset tokens (hashed).
type PasswordResetRepo interface {
	Create(ctx context.Context, t *domain.PasswordResetToken) error
	FindByHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string, at time.Time) error
}

// AuditRepo persists admin audit log entries.
type AuditRepo interface {
	Create(ctx context.Context, e *domain.AuditEntry) error
	List(ctx context.Context, limit int) ([]*domain.AuditEntry, error)
}

// NotificationRepo persists in-app notifications and delivery records.
type NotificationRepo interface {
	Create(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
	ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Notification, error)
	MarkSent(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string) error
}

// PasswordHasher abstracts bcrypt for testability.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Verify(hash, plain string) error
}

// TokenIssuer abstracts JWT access token issuance + parsing.
type TokenIssuer interface {
	Issue(u *domain.User, ttl time.Duration) (string, error)
	Parse(raw string) (Claims, error)
}

// Claims is the parsed JWT payload.
type Claims struct {
	UserID string
	Email  string
	Role   domain.Role
}

// Clock returns the current time. Inject for deterministic tests.
type Clock interface {
	Now() time.Time
}

// IDGen generates new opaque identifiers (UUIDv4 in production).
type IDGen interface {
	New() string
}

// AIClient is the chat / planner port.
type AIClient interface {
	Generate(ctx context.Context, userMessage string, vendors []VendorRef) (*ChatReply, error)
}

// VendorRef is the slim catalog reference passed to the AI planner.
type VendorRef struct {
	ID        string
	Name      string
	Category  string
	City      string
	PriceFrom int64
}

// ChatReply is the AI-generated structured response.
type ChatReply struct {
	Reply  string
	Blocks []ChatBlock
}

// ChatBlock is a typed UI block returned alongside the reply text.
type ChatBlock struct {
	Type string
	Data map[string]any
}

// Mailer sends transactional email.
type Mailer interface {
	Send(ctx context.Context, to, subject, htmlBody string) error
}

// Pusher sends push notifications via FCM.
type Pusher interface {
	Send(ctx context.Context, tokens []string, title, body string, data map[string]string) error
}

// PaymentGateway creates payment intents and verifies callbacks.
type PaymentGateway interface {
	CreatePayment(ctx context.Context, in PaymentIntent) (PaymentRedirect, error)
	VerifyCallback(form map[string]string) (CallbackResult, error)
}

// PaymentIntent captures all fields needed to initiate a payment.
type PaymentIntent struct {
	OrderID     string
	Amount      int64
	Currency    string
	Description string
	CustomerEmail string
	SuccessURL  string
	FailureURL  string
}

// PaymentRedirect is the gateway-provided customer redirect.
type PaymentRedirect struct {
	URL           string
	TransactionID string
}

// CallbackResult is the verified outcome of a payment callback.
type CallbackResult struct {
	OrderID       string
	TransactionID string
	Success       bool
	Amount        int64
}
