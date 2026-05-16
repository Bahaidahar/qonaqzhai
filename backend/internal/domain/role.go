// Package domain holds pure business entities and rules.
// Domain types depend on no other project packages.
package domain

// Role is the user role enum.
type Role string

const (
	RoleCustomer Role = "customer"
	RoleVendor   Role = "vendor"
	RoleAdmin    Role = "admin"
)

// Valid reports whether r is one of the defined roles.
func (r Role) Valid() bool {
	switch r {
	case RoleCustomer, RoleVendor, RoleAdmin:
		return true
	}
	return false
}

// ParseRole returns the role with default fallback to RoleCustomer
// when the input is empty or unknown. Used at HTTP/DTO boundary.
func ParseRole(s string) Role {
	r := Role(s)
	if r != RoleCustomer && r != RoleVendor {
		return RoleCustomer
	}
	return r
}
