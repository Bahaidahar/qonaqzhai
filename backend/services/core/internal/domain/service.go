package domain

import (
	"strings"
	"time"
)

// ServiceUnit is the pricing dimension shown next to the price.
type ServiceUnit string

const (
	UnitFixed  ServiceUnit = "fixed"
	UnitHour   ServiceUnit = "hour"
	UnitItem   ServiceUnit = "item"
	UnitPerson ServiceUnit = "person"
	UnitDay    ServiceUnit = "day"
)

// Valid reports whether the unit is one of the defined values.
func (u ServiceUnit) Valid() bool {
	switch u {
	case UnitFixed, UnitHour, UnitItem, UnitPerson, UnitDay:
		return true
	}
	return false
}

// Service is one offering on a vendor's menu.
type Service struct {
	ID          string      `json:"id"`
	VendorID    string      `json:"vendorId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       int64       `json:"price"`
	Unit        ServiceUnit `json:"unit"`
	IsActive    bool        `json:"isActive"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// ServiceInput captures user-supplied fields for create / update.
type ServiceInput struct {
	Name        string
	Description string
	Price       int64
	Unit        ServiceUnit
	IsActive    *bool // nil → default true on create, no-op on update
}

// Normalize trims whitespace + applies defaults.
func (in *ServiceInput) Normalize() {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	if in.Unit == "" {
		in.Unit = UnitFixed
	}
}
