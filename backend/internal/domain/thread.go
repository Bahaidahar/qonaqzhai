package domain

import "time"

// BookingThread is a private DM between the customer and the vendor of a single booking.
// Created automatically when the vendor accepts the booking.
type BookingThread struct {
	ID         string    `json:"id"`
	BookingID  string    `json:"bookingId"`
	CustomerID string    `json:"customerId"`
	VendorID   string    `json:"vendorId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// ThreadMessage is one message in a booking thread.
type ThreadMessage struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"threadId"`
	SenderID  string    `json:"senderId"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
