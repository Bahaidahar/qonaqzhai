// Package ports defines the interfaces the core service usecases depend on.
// Adapter packages provide the concrete implementations.
package ports

import (
	"context"
	"time"

	"qonaqzhai-backend/services/core/internal/domain"
)

// VendorQuery is the search/filter set for the public catalog.
type VendorQuery struct {
	Q         string
	Category  string
	City      string
	Status    domain.VendorStatus
	MinPrice  int64
	MaxPrice  int64
	MinRating float64
	Sort      string // "" | price_asc | price_desc | rating_desc
	Page      int
	Limit     int
}

// VendorRepo persists vendor profiles.
type VendorRepo interface {
	Upsert(ctx context.Context, userID string, in domain.VendorInput) (*domain.Vendor, error)
	FindByID(ctx context.Context, id string) (*domain.Vendor, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Vendor, error)
	FindByIDs(ctx context.Context, ids []string) ([]*domain.Vendor, error)
	Search(ctx context.Context, q VendorQuery) ([]*domain.Vendor, int, error)
	UpdateStatus(ctx context.Context, id string, status domain.VendorStatus) error
	UpdateRating(ctx context.Context, id string, avg float64, count int) error
	CountByStatus(ctx context.Context) (map[domain.VendorStatus]int, error)
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
	Stats(ctx context.Context) (BookingStats, error)
}

// BookingStats is a single aggregated row returned by the admin dashboard.
type BookingStats struct {
	Total    int
	Pending  int
	Accepted int
	Paid     int
	GMV      int64 // sum(amount) of paid bookings, in minor units
}

// ServiceRepo persists vendor service menus.
type ServiceRepo interface {
	Create(ctx context.Context, s *domain.Service) (*domain.Service, error)
	Update(ctx context.Context, id string, in domain.ServiceInput) (*domain.Service, error)
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (*domain.Service, error)
	ListByVendor(ctx context.Context, vendorID string) ([]*domain.Service, error)
}

// PhotoRepo persists vendor photos.
type PhotoRepo interface {
	Insert(ctx context.Context, p *domain.Photo) (*domain.Photo, error)
	Find(ctx context.Context, id string) (*domain.Photo, error)
	Delete(ctx context.Context, id string) error
	ListByVendor(ctx context.Context, vendorID string) ([]*domain.Photo, error)
}

// ReviewRepo persists reviews.
type ReviewRepo interface {
	Create(ctx context.Context, r *domain.Review) (*domain.Review, error)
	ListForVendor(ctx context.Context, vendorID string) ([]*domain.Review, error)
	FindByBooking(ctx context.Context, bookingID string) (*domain.Review, error)
	AggregateForVendor(ctx context.Context, vendorID string) (float64, int, error)
}

// NotificationRepo persists in-app notifications.
type NotificationRepo interface {
	Enqueue(ctx context.Context, n *domain.Notification) (*domain.Notification, error)
	ListForUser(ctx context.Context, userID string, limit int) ([]*domain.Notification, error)
	MarkSent(ctx context.Context, id string) error
}

// FCMTokenRepo persists device push tokens.
type FCMTokenRepo interface {
	Upsert(ctx context.Context, t *domain.FCMToken) error
	Delete(ctx context.Context, token string) error
	ListByUsers(ctx context.Context, userIDs []string) ([]*domain.FCMToken, error)
}

// AuthClient lets the core service ask auth-svc for user data without taking a
// hard dependency on the JWT secret or the users table.
type AuthClient interface {
	GetUser(ctx context.Context, userID string) (*ExternalUser, error)
	GetUsersBatch(ctx context.Context, userIDs []string) ([]*ExternalUser, error)
}

// ExternalUser is the subset of auth-svc User we use here.
type ExternalUser struct {
	ID    string
	Email string
	Name  string
	Role  string
}

// PaymentClient charges a saved card on payment-svc.
type PaymentClient interface {
	Charge(ctx context.Context, in ChargeRequest) (*PaymentResult, error)
}

// ChargeRequest is the payload sent to payment-svc.
type ChargeRequest struct {
	BookingID string
	UserID    string
	CardID    string
	Amount    int64
	Currency  string
}

// PaymentResult is the payment-svc response.
type PaymentResult struct {
	ID     string
	Status string
}

// RealtimeClient lets core trigger thread creation + event fan-out on
// realtime-svc.
type RealtimeClient interface {
	EnsureThread(ctx context.Context, bookingID, customerID, vendorUserID string) error
	Publish(ctx context.Context, event string, payloadJSON []byte, userIDs ...string) error
}

// Pusher delivers push notifications. nil-safe.
type Pusher interface {
	Push(ctx context.Context, tokens []string, title, body string) error
}

// Clock abstracts the current time so tests can pin it.
type Clock interface{ Now() time.Time }

// IDGen generates new opaque entity IDs.
type IDGen interface{ New() string }
