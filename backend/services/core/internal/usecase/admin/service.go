// Package admin aggregates platform-wide stats by fanning out to local repos.
package admin

import (
	"context"

	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/ports"
)

// Deps bundles admin collaborators.
type Deps struct {
	Vendors  ports.VendorRepo
	Bookings ports.BookingRepo
}

// Service computes core-side admin metrics.
type Service struct{ d Deps }

// New constructs an admin Service.
func New(d Deps) *Service { return &Service{d: d} }

// Stats returns the aggregate dashboard payload.
type Stats struct {
	VendorsTotal    int
	VendorsPending  int
	VendorsApproved int
	BookingsTotal   int
	BookingsPending int
	BookingsAccept  int
	BookingsPaid    int
	GMV             int64
}

// Compute returns the dashboard payload. Each repo issues exactly one query.
func (s *Service) Compute(ctx context.Context) (Stats, error) {
	vc, err := s.d.Vendors.CountByStatus(ctx)
	if err != nil {
		return Stats{}, err
	}
	bs, err := s.d.Bookings.Stats(ctx)
	if err != nil {
		return Stats{}, err
	}
	total := 0
	for _, n := range vc {
		total += n
	}
	return Stats{
		VendorsTotal:    total,
		VendorsPending:  vc[domain.VendorPending],
		VendorsApproved: vc[domain.VendorApproved],
		BookingsTotal:   bs.Total,
		BookingsPending: bs.Pending,
		BookingsAccept:  bs.Accepted,
		BookingsPaid:    bs.Paid,
		GMV:             bs.GMV,
	}, nil
}
