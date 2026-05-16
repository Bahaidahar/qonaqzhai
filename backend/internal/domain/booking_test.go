package domain

import "testing"

func TestBookingStatusValid(t *testing.T) {
	t.Parallel()
	valid := []BookingStatus{
		BookingPending, BookingAccepted, BookingDeclined,
		BookingCancelled, BookingCompleted, BookingPaid,
	}
	for _, s := range valid {
		if !s.Valid() {
			t.Errorf("%q expected valid", s)
		}
	}
	if BookingStatus("zombie").Valid() {
		t.Error("unknown status accepted")
	}
}

func TestVendorMayTransition(t *testing.T) {
	t.Parallel()
	cases := []struct {
		from, to BookingStatus
		ok       bool
	}{
		{BookingPending, BookingAccepted, true},
		{BookingPending, BookingDeclined, true},
		{BookingPending, BookingCompleted, false},
		{BookingPaid, BookingAccepted, true},
		{BookingAccepted, BookingCompleted, true},
		{BookingDeclined, BookingAccepted, false},
		{BookingCancelled, BookingAccepted, false},
		{BookingPending, BookingCancelled, false},
		{BookingAccepted, BookingDeclined, false},
	}
	for _, c := range cases {
		b := &Booking{Status: c.from}
		if got := b.VendorMayTransition(c.to); got != c.ok {
			t.Errorf("%s→%s: got %v want %v", c.from, c.to, got, c.ok)
		}
	}
}

func TestCustomerMayTransition(t *testing.T) {
	t.Parallel()
	cases := []struct {
		from, to BookingStatus
		ok       bool
	}{
		{BookingPending, BookingCancelled, true},
		{BookingAccepted, BookingCancelled, true},
		{BookingCompleted, BookingCancelled, false},
		{BookingPending, BookingAccepted, false},
		{BookingCancelled, BookingCancelled, false},
	}
	for _, c := range cases {
		b := &Booking{Status: c.from}
		if got := b.CustomerMayTransition(c.to); got != c.ok {
			t.Errorf("%s→%s: got %v want %v", c.from, c.to, got, c.ok)
		}
	}
}
