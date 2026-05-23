// Package domain holds the auth service's pure business entities and rules.
// Domain types depend on no other project packages so they can be tested in
// isolation and exposed at adapter boundaries without leaks.
package domain

// Role enumerates the kinds of principals the platform recognises.
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

// ParseRole defaults to RoleCustomer for empty / unknown / admin input.
// Public signup endpoints must never let callers self-elevate to admin.
func ParseRole(s string) Role {
	r := Role(s)
	if r != RoleCustomer && r != RoleVendor {
		return RoleCustomer
	}
	return r
}
