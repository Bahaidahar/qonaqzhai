package domain

import "time"

// NotificationChannel selects the delivery path for an event.
type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelPush  NotificationChannel = "push"
	ChannelBoth  NotificationChannel = "email+push"
)

// Notification is an in-app + delivery record.
type Notification struct {
	ID        string              `json:"id"`
	UserID    string              `json:"userId"`
	Type      string              `json:"type"`
	Channel   NotificationChannel `json:"channel"`
	Title     string              `json:"title"`
	Body      string              `json:"body"`
	Status    string              `json:"status"` // queued | sent | failed
	CreatedAt time.Time           `json:"createdAt"`
}

// FCMToken records a device token for push delivery.
type FCMToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	Platform  string    `json:"platform"`
	CreatedAt time.Time `json:"createdAt"`
}
