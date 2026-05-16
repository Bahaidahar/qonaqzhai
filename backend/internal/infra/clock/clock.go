// Package clock provides the production Clock implementation.
package clock

import "time"

// System implements usecase.Clock against the OS wall clock.
type System struct{}

// New returns a System clock.
func New() System { return System{} }

// Now returns time.Now() in UTC.
func (System) Now() time.Time { return time.Now().UTC() }
