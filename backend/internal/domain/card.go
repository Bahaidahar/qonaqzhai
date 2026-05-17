package domain

import (
	"strings"
	"time"
)

// PaymentCard is a mock saved card. We never store the PAN — only last4 + brand.
type PaymentCard struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Brand     string    `json:"brand"`
	Last4     string    `json:"last4"`
	ExpMonth  int       `json:"expMonth"`
	ExpYear   int       `json:"expYear"`
	Holder    string    `json:"holder"`
	IsDefault bool      `json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
}

// CardInput captures fields supplied to add a card. Number is hashed away —
// only last4 + brand are persisted (MVP mock, never used against a real PSP).
type CardInput struct {
	Number   string
	ExpMonth int
	ExpYear  int
	Holder   string
}

// Normalize strips spaces from the PAN + trims fields.
func (in *CardInput) Normalize() {
	in.Number = strings.ReplaceAll(strings.TrimSpace(in.Number), " ", "")
	in.Holder = strings.TrimSpace(in.Holder)
}

// Validate enforces lightweight Luhn-style checks (length + expiry).
func (in *CardInput) Validate(now time.Time) error {
	if l := len(in.Number); l < 13 || l > 19 {
		return ErrInvalidInput
	}
	for _, r := range in.Number {
		if r < '0' || r > '9' {
			return ErrInvalidInput
		}
	}
	if in.ExpMonth < 1 || in.ExpMonth > 12 {
		return ErrInvalidInput
	}
	yy := in.ExpYear
	if yy < 100 {
		yy += 2000
	}
	if yy < now.Year() || (yy == now.Year() && in.ExpMonth < int(now.Month())) {
		return ErrInvalidInput
	}
	return nil
}

// Last4 extracts the last four digits of the PAN.
func (in *CardInput) Last4() string {
	if len(in.Number) < 4 {
		return in.Number
	}
	return in.Number[len(in.Number)-4:]
}

// DetectBrand sniffs the card brand from the BIN.
func DetectBrand(pan string) string {
	if len(pan) == 0 {
		return "unknown"
	}
	switch pan[0] {
	case '4':
		return "visa"
	case '5':
		return "mastercard"
	case '3':
		return "amex"
	case '6':
		return "discover"
	}
	return "unknown"
}
