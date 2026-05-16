// Package idgen provides the production IDGen implementation using UUIDv4.
package idgen

import "github.com/google/uuid"

// UUID implements usecase.IDGen via google/uuid v4 strings.
type UUID struct{}

// New returns a UUID generator.
func New() UUID { return UUID{} }

// New returns a fresh UUIDv4 in canonical hyphenated string form.
func (UUID) New() string { return uuid.NewString() }
