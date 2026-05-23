// Package clock provides the system clock implementation.
package clock

import "time"

// System implements ports.Clock against the OS wall clock.
type System struct{}

// New returns a System clock.
func New() System { return System{} }

// Now returns time.Now() in UTC.
func (System) Now() time.Time { return time.Now().UTC() }
