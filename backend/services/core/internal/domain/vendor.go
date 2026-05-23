// Package domain holds the core service's pure business entities. They have no
// dependencies on any project package outside this one.
package domain

import (
	"strings"
	"time"
)

// VendorStatus is the moderation state of a vendor profile.
type VendorStatus string

const (
	VendorPending  VendorStatus = "pending"
	VendorApproved VendorStatus = "approved"
	VendorRejected VendorStatus = "rejected"
)

// Valid reports whether the status is a known value.
func (s VendorStatus) Valid() bool {
	switch s {
	case VendorPending, VendorApproved, VendorRejected:
		return true
	}
	return false
}

// Vendor is the business profile owned by a vendor-role user living in
// auth-svc. UserID is therefore a foreign UUID with no DB constraint.
type Vendor struct {
	ID          string       `json:"id"`
	UserID      string       `json:"userId"`
	Name        string       `json:"name"`
	Category    string       `json:"category"`
	City        string       `json:"city"`
	Description string       `json:"description"`
	PriceFrom   int64        `json:"priceFrom"`
	Status      VendorStatus `json:"status"`
	RatingAvg   float64      `json:"ratingAvg"`
	RatingCount int          `json:"ratingCount"`
	PhotoIDs    []string     `json:"photoIds"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// IsPublic reports whether the vendor should be visible in the public catalog.
func (v *Vendor) IsPublic() bool { return v.Status == VendorApproved }

// VendorInput captures user-supplied fields for create / update operations.
type VendorInput struct {
	Name        string
	Category    string
	City        string
	Description string
	PriceFrom   int64
}

// Normalize trims whitespace on text fields.
func (in *VendorInput) Normalize() {
	in.Name = strings.TrimSpace(in.Name)
	in.Category = strings.TrimSpace(in.Category)
	in.City = strings.TrimSpace(in.City)
	in.Description = strings.TrimSpace(in.Description)
}
