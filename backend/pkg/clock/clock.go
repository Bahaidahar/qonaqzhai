// Package clock provides a System clock implementation that every service can
// embed via DI. Tests substitute a fixed clock.
package clock

import "time"

// System reports the OS wall clock in UTC.
type System struct{}

// New returns a System clock.
func New() System { return System{} }

// Now returns time.Now() in UTC.
func (System) Now() time.Time { return time.Now().UTC() }
