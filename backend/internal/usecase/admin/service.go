// Package admin implements admin moderation and analytics use cases.
package admin

import (
	"context"
	"fmt"
	"sort"
	"time"

	"qonaqzhai-backend/internal/domain"
	"qonaqzhai-backend/internal/usecase"
)

// Notifier emits admin-driven notifications. Optional.
type Notifier interface {
	Enqueue(ctx context.Context, n *domain.Notification) error
}

// Deps bundles admin service collaborators.
type Deps struct {
	Users    usecase.UserRepo
	Vendors  usecase.VendorRepo
	Bookings usecase.BookingRepo
	Reviews  usecase.ReviewRepo
	Audit    usecase.AuditRepo // optional
	Notifier Notifier          // optional
}

// Service exposes admin moderation + stats operations.
type Service struct{ d Deps }

// New constructs an admin Service.
func New(d Deps) *Service { return &Service{d: d} }

// audit best-effort logs an admin action (silently swallows errors).
func (s *Service) audit(ctx context.Context, actorID, actorEmail, action, targetType, targetID, meta string) {
	if s.d.Audit == nil || actorID == "" {
		return
	}
	_ = s.d.Audit.Create(ctx, &domain.AuditEntry{
		ActorID: actorID, ActorEmail: actorEmail,
		Action: action, TargetType: targetType, TargetID: targetID, Meta: meta,
	})
}

// AuditLog returns the latest admin actions.
func (s *Service) AuditLog(ctx context.Context, limit int) ([]*domain.AuditEntry, error) {
	if s.d.Audit == nil {
		return []*domain.AuditEntry{}, nil
	}
	return s.d.Audit.List(ctx, limit)
}

// ListUsers returns every account.
func (s *Service) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.d.Users.List(ctx)
}

// SetUserStatus suspends or restores a user account.
func (s *Service) SetUserStatus(ctx context.Context, actorID, actorEmail, id string, status domain.UserStatus) (*domain.User, error) {
	if status != domain.UserActive && status != domain.UserSuspended {
		return nil, fmt.Errorf("status: %w", domain.ErrInvalidInput)
	}
	if err := s.d.Users.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	s.audit(ctx, actorID, actorEmail, "user.set_status", "user", id, string(status))
	return s.d.Users.FindByID(ctx, id)
}

// SetVendorStatus moderates a vendor profile.
func (s *Service) SetVendorStatus(ctx context.Context, actorID, actorEmail, id string, status domain.VendorStatus) (*domain.Vendor, error) {
	if !status.Valid() {
		return nil, fmt.Errorf("status: %w", domain.ErrInvalidInput)
	}
	if err := s.d.Vendors.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	s.audit(ctx, actorID, actorEmail, "vendor.set_status", "vendor", id, string(status))
	v, err := s.d.Vendors.FindByID(ctx, id)
	if err == nil && s.d.Notifier != nil {
		var notifType domain.NotificationType
		var title, body string
		switch status {
		case domain.VendorApproved:
			notifType = domain.NotifVendorApproved
			title = "Vendor profile approved"
			body = "<p>Your vendor profile is live in the catalog.</p>"
		case domain.VendorRejected:
			notifType = domain.NotifVendorRejected
			title = "Vendor profile rejected"
			body = "<p>Your vendor profile was not approved. Please contact support.</p>"
		}
		if notifType != "" {
			_ = s.d.Notifier.Enqueue(ctx, &domain.Notification{
				UserID:  v.UserID,
				Type:    notifType,
				Channel: domain.ChannelBoth,
				Title:   title,
				Body:    body,
			})
		}
	}
	return v, err
}

// Stats aggregates platform-wide KPIs.
type Stats struct {
	Users            int `json:"users"`
	Customers        int `json:"customers"`
	Vendors          int `json:"vendors"`
	Admins           int `json:"admins"`
	VendorProfiles   int `json:"vendor_profiles"`
	VendorsPending   int `json:"vendors_pending"`
	VendorsApproved  int `json:"vendors_approved"`
	VendorsRejected  int `json:"vendors_rejected"`
	BookingsTotal    int `json:"bookings_total"`
	BookingsPending  int `json:"bookings_pending"`
	BookingsAccepted int `json:"bookings_accepted"`
	BookingsPaid     int `json:"bookings_paid"`
	BookingsRevenue  int64 `json:"bookings_revenue"`
}

// TimePoint is one bucket on a time-series chart.
type TimePoint struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// CategoryCount aggregates per-category vendor totals.
type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

// FunnelStage is one step in the vendor approval funnel.
type FunnelStage struct {
	Stage string `json:"stage"`
	Count int    `json:"count"`
}

// BookingsTimeseries aggregates bookings per day between [from, to).
// Empty days are omitted.
func (s *Service) BookingsTimeseries(ctx context.Context, from, to time.Time) ([]TimePoint, error) {
	all, err := s.d.Bookings.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	buckets := map[string]int64{}
	for _, b := range all {
		if !from.IsZero() && b.CreatedAt.Before(from) {
			continue
		}
		if !to.IsZero() && b.CreatedAt.After(to) {
			continue
		}
		day := b.CreatedAt.UTC().Format("2006-01-02")
		buckets[day]++
	}
	out := make([]TimePoint, 0, len(buckets))
	for day, n := range buckets {
		out = append(out, TimePoint{Date: day, Value: n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Date < out[j].Date })
	return out, nil
}

// TopCategories returns the most popular vendor categories, capped by limit.
func (s *Service) TopCategories(ctx context.Context, limit int) ([]CategoryCount, error) {
	vendors, _, err := s.d.Vendors.Search(ctx, usecase.VendorQuery{Limit: 10_000})
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	for _, v := range vendors {
		counts[v.Category]++
	}
	out := make([]CategoryCount, 0, len(counts))
	for c, n := range counts {
		out = append(out, CategoryCount{Category: c, Count: n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}

// ApprovalFunnel returns counts at each stage of the vendor moderation flow.
func (s *Service) ApprovalFunnel(ctx context.Context) ([]FunnelStage, error) {
	vendors, _, err := s.d.Vendors.Search(ctx, usecase.VendorQuery{Limit: 10_000})
	if err != nil {
		return nil, err
	}
	var pending, approved, rejected int
	for _, v := range vendors {
		switch v.Status {
		case domain.VendorPending:
			pending++
		case domain.VendorApproved:
			approved++
		case domain.VendorRejected:
			rejected++
		}
	}
	return []FunnelStage{
		{Stage: "submitted", Count: pending + approved + rejected},
		{Stage: "pending", Count: pending},
		{Stage: "approved", Count: approved},
		{Stage: "rejected", Count: rejected},
	}, nil
}

// Stats computes a snapshot for the admin dashboard.
func (s *Service) Stats(ctx context.Context) (*Stats, error) {
	users, err := s.d.Users.List(ctx)
	if err != nil {
		return nil, err
	}
	vendors, _, err := s.d.Vendors.Search(ctx, usecase.VendorQuery{Limit: 10_000})
	if err != nil {
		return nil, err
	}
	bookings, err := s.d.Bookings.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	st := &Stats{
		Users:          len(users),
		VendorProfiles: len(vendors),
		BookingsTotal:  len(bookings),
	}
	for _, u := range users {
		switch u.Role {
		case domain.RoleCustomer:
			st.Customers++
		case domain.RoleVendor:
			st.Vendors++
		case domain.RoleAdmin:
			st.Admins++
		}
	}
	for _, v := range vendors {
		switch v.Status {
		case domain.VendorPending:
			st.VendorsPending++
		case domain.VendorApproved:
			st.VendorsApproved++
		case domain.VendorRejected:
			st.VendorsRejected++
		}
	}
	for _, b := range bookings {
		switch b.Status {
		case domain.BookingPending:
			st.BookingsPending++
		case domain.BookingAccepted:
			st.BookingsAccepted++
		case domain.BookingPaid:
			st.BookingsPaid++
			st.BookingsRevenue += b.Amount
		}
	}
	return st, nil
}
