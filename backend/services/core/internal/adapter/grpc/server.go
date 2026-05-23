// Package grpc exposes core service operations over gRPC. Other services call
// these instead of touching core's database.
package grpc

import (
	"context"

	corev1 "qonaqzhai-backend/gen/proto/core/v1"
	"qonaqzhai-backend/pkg/grpcutil"

	"qonaqzhai-backend/services/core/internal/domain"
	"qonaqzhai-backend/services/core/internal/usecase/admin"
	"qonaqzhai-backend/services/core/internal/usecase/booking"
	"qonaqzhai-backend/services/core/internal/usecase/vendor"
)

// Server implements corev1.CoreServiceServer.
type Server struct {
	corev1.UnimplementedCoreServiceServer
	vendors  *vendor.Service
	bookings *booking.Service
	admin    *admin.Service
}

// New constructs the gRPC server.
func New(vendors *vendor.Service, bookings *booking.Service, ad *admin.Service) *Server {
	return &Server{vendors: vendors, bookings: bookings, admin: ad}
}

// GetVendor returns the vendor by id (no public-only filter).
func (s *Server) GetVendor(ctx context.Context, req *corev1.GetVendorRequest) (*corev1.Vendor, error) {
	v, err := s.vendors.FindByID(ctx, req.GetVendorId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return vendorProto(v), nil
}

// GetVendorByUser returns the vendor owned by userID.
func (s *Server) GetVendorByUser(ctx context.Context, req *corev1.GetVendorByUserRequest) (*corev1.Vendor, error) {
	v, err := s.vendors.FindByUserID(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return vendorProto(v), nil
}

// ListVendorsByIDs returns vendors for the supplied ids.
func (s *Server) ListVendorsByIDs(ctx context.Context, req *corev1.ListVendorsByIDsRequest) (*corev1.ListVendorsByIDsResponse, error) {
	vs, err := s.vendors.ListByIDs(ctx, req.GetVendorIds())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	out := make([]*corev1.Vendor, len(vs))
	for i, v := range vs {
		out[i] = vendorProto(v)
	}
	return &corev1.ListVendorsByIDsResponse{Vendors: out}, nil
}

// GetBooking returns a booking by id.
func (s *Server) GetBooking(ctx context.Context, req *corev1.GetBookingRequest) (*corev1.Booking, error) {
	b, err := s.bookings.GetRaw(ctx, req.GetBookingId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return bookingProto(b), nil
}

// IsBookingAccepted answers realtime's thread-creation precheck.
func (s *Server) IsBookingAccepted(ctx context.Context, req *corev1.IsBookingAcceptedRequest) (*corev1.IsBookingAcceptedResponse, error) {
	b, vendorUserID, accepted, err := s.bookings.IsAccepted(ctx, req.GetBookingId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	out := &corev1.IsBookingAcceptedResponse{Accepted: accepted}
	if b != nil {
		out.CustomerId = b.CustomerID
		out.VendorId = b.VendorID
	}
	out.VendorUserId = vendorUserID
	return out, nil
}

// MarkBookingPaid is the payment-svc callback.
func (s *Server) MarkBookingPaid(ctx context.Context, req *corev1.MarkBookingPaidRequest) (*corev1.Booking, error) {
	b, err := s.bookings.MarkPaid(ctx, req.GetBookingId(), req.GetPaymentId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return bookingProto(b), nil
}

// AdminStats fan-outs core's own metrics. Cross-service totals (users, GMV
// breakdown) are aggregated by the caller (gateway / dashboard).
func (s *Server) AdminStats(ctx context.Context, _ *corev1.AdminStatsRequest) (*corev1.AdminStatsResponse, error) {
	st, err := s.admin.Compute(ctx)
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return &corev1.AdminStatsResponse{
		VendorsTotal:     int32(st.VendorsTotal),
		VendorsPending:   int32(st.VendorsPending),
		VendorsApproved:  int32(st.VendorsApproved),
		BookingsTotal:    int32(st.BookingsTotal),
		BookingsPending:  int32(st.BookingsPending),
		BookingsAccepted: int32(st.BookingsAccept),
		BookingsPaid:     int32(st.BookingsPaid),
		Gmv:              st.GMV,
	}, nil
}

func vendorProto(v *domain.Vendor) *corev1.Vendor {
	if v == nil {
		return nil
	}
	return &corev1.Vendor{
		Id:          v.ID,
		UserId:      v.UserID,
		Name:        v.Name,
		Category:    v.Category,
		City:        v.City,
		Description: v.Description,
		PriceFrom:   v.PriceFrom,
		Status:      string(v.Status),
		RatingAvg:   v.RatingAvg,
		RatingCount: int32(v.RatingCount),
		PhotoIds:    v.PhotoIDs,
		CreatedAt:   v.CreatedAt.Unix(),
		UpdatedAt:   v.UpdatedAt.Unix(),
	}
}

func bookingProto(b *domain.Booking) *corev1.Booking {
	if b == nil {
		return nil
	}
	return &corev1.Booking{
		Id:         b.ID,
		CustomerId: b.CustomerID,
		VendorId:   b.VendorID,
		ServiceId:  b.ServiceID,
		EventDate:  b.EventDate,
		GuestCount: int32(b.GuestCount),
		Note:       b.Note,
		Status:     string(b.Status),
		Amount:     b.Amount,
		PaymentId:  b.PaymentID,
		CreatedAt:  b.CreatedAt.Unix(),
	}
}
