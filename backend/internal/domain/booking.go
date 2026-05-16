package domain

import "time"

// BookingStatus is the lifecycle state of a booking.
type BookingStatus string

const (
	BookingPending   BookingStatus = "pending"
	BookingAccepted  BookingStatus = "accepted"
	BookingDeclined  BookingStatus = "declined"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
	BookingPaid      BookingStatus = "paid"
)

// Valid reports whether the status is a known value.
func (s BookingStatus) Valid() bool {
	switch s {
	case BookingPending, BookingAccepted, BookingDeclined,
		BookingCancelled, BookingCompleted, BookingPaid:
		return true
	}
	return false
}

// Booking is a customer's reservation of a vendor's service.
type Booking struct {
	ID         string        `json:"id"`
	CustomerID string        `json:"customerId"`
	VendorID   string        `json:"vendorId"`
	EventDate  string        `json:"eventDate"`
	GuestCount int           `json:"guestCount"`
	Note       string        `json:"note"`
	Status     BookingStatus `json:"status"`
	Amount     int64         `json:"amount"`
	PaymentID  string        `json:"paymentId,omitempty"`
	CreatedAt  time.Time     `json:"createdAt"`
}

// VendorMayTransition reports whether a vendor may move the booking
// from its current status to the requested next status.
func (b *Booking) VendorMayTransition(next BookingStatus) bool {
	switch next {
	case BookingAccepted, BookingDeclined:
		return b.Status == BookingPending || b.Status == BookingPaid
	case BookingCompleted:
		return b.Status == BookingAccepted || b.Status == BookingPaid
	}
	return false
}

// CustomerMayTransition reports whether the booking's customer may move
// the booking to the requested status.
func (b *Booking) CustomerMayTransition(next BookingStatus) bool {
	if next != BookingCancelled {
		return false
	}
	return b.Status == BookingPending || b.Status == BookingAccepted
}
