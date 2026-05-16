package domain

import "time"

// NotificationChannel selects the delivery path for an event.
type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
	ChannelBoth  NotificationChannel = "email+push"
)

// NotificationType is the semantic event kind.
type NotificationType string

const (
	NotifSignupWelcome   NotificationType = "signup.welcome"
	NotifPasswordReset   NotificationType = "auth.password_reset"
	NotifBookingCreated  NotificationType = "booking.created"
	NotifBookingAccepted NotificationType = "booking.accepted"
	NotifBookingDeclined NotificationType = "booking.declined"
	NotifBookingPaid     NotificationType = "booking.paid"
	NotifVendorApproved  NotificationType = "vendor.approved"
	NotifVendorRejected  NotificationType = "vendor.rejected"
)

// Notification is an in-app + delivery record.
type Notification struct {
	ID        string              `json:"id"`
	UserID    string              `json:"userId"`
	Type      NotificationType    `json:"type"`
	Channel   NotificationChannel `json:"channel"`
	Title     string              `json:"title"`
	Body      string              `json:"body"`
	Status    string              `json:"status"` // queued | sent | failed
	CreatedAt time.Time           `json:"createdAt"`
}
