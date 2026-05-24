// Package idgen returns opaque UUIDv4 identifiers.
package idgen

import "github.com/google/uuid"

// UUID emits canonical hyphenated UUIDv4 strings.
type UUID struct{}

// New returns a UUID generator.
func New() UUID { return UUID{} }

// New returns a fresh UUIDv4 string.
func (UUID) New() string { return uuid.NewString() }
