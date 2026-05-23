package domain

import "time"

// PaymentStatus is the lifecycle state of a charge attempt.
type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentCaptured PaymentStatus = "captured"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"
)

// Payment is one charge attempt against a booking.
type Payment struct {
	ID          string        `json:"id"`
	BookingID   string        `json:"bookingId"`
	UserID      string        `json:"userId"`
	CardID      string        `json:"cardId"`
	Amount      int64         `json:"amount"`
	Currency    string        `json:"currency"`
	Status      PaymentStatus `json:"status"`
	ProviderRef string        `json:"providerRef"`
	CreatedAt   time.Time     `json:"createdAt"`
}
