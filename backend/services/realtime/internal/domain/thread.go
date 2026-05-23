// Package domain holds realtime service entities.
package domain

import "time"

// Thread is a DM channel attached to a booking. CustomerID + VendorID are
// foreign UUIDs from auth-svc (vendor's user id, not the vendor row id).
type Thread struct {
	ID         string    `json:"id"`
	BookingID  string    `json:"bookingId"`
	CustomerID string    `json:"customerId"`
	VendorID   string    `json:"vendorId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Message is one chat line in a thread.
type Message struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"threadId"`
	SenderID  string    `json:"senderId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
