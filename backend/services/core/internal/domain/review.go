package domain

import "time"

// Review is a customer's evaluation of a vendor after a completed booking.
type Review struct {
	ID         string    `json:"id"`
	BookingID  string    `json:"bookingId"`
	CustomerID string    `json:"customerId"`
	VendorID   string    `json:"vendorId"`
	Rating     int       `json:"rating"`
	Text       string    `json:"text"`
	CreatedAt  time.Time `json:"createdAt"`
}

// MinRating and MaxRating bound the 5-star scale.
const (
	MinRating = 1
	MaxRating = 5
)

// ValidRating reports whether the rating falls within the allowed range.
func ValidRating(r int) bool { return r >= MinRating && r <= MaxRating }
