// Package grpc exposes realtime operations over gRPC.
package grpc

import (
	"context"

	realtimev1 "qonaqzhai-backend/gen/proto/realtime/v1"
	"qonaqzhai-backend/pkg/grpcutil"

	"qonaqzhai-backend/services/realtime/internal/ports"
	"qonaqzhai-backend/services/realtime/internal/usecase/thread"
)

// Server implements realtimev1.RealtimeServiceServer.
type Server struct {
	realtimev1.UnimplementedRealtimeServiceServer
	threads   *thread.Service
	publisher ports.Publisher
}

// New constructs the gRPC server.
func New(threads *thread.Service, publisher ports.Publisher) *Server {
	return &Server{threads: threads, publisher: publisher}
}

// EnsureThread idempotently creates the chat thread for a booking.
func (s *Server) EnsureThread(ctx context.Context, req *realtimev1.EnsureThreadRequest) (*realtimev1.Thread, error) {
	t, err := s.threads.Ensure(ctx, req.GetBookingId(), req.GetCustomerId(), req.GetVendorId())
	if err != nil {
		return nil, grpcutil.ToStatus(err)
	}
	return &realtimev1.Thread{
		Id:         t.ID,
		BookingId:  t.BookingID,
		CustomerId: t.CustomerID,
		VendorId:   t.VendorID,
		CreatedAt:  t.CreatedAt.Unix(),
		UpdatedAt:  t.UpdatedAt.Unix(),
	}, nil
}

// PublishEvent fans an event out to the hub.
func (s *Server) PublishEvent(_ context.Context, req *realtimev1.PublishEventRequest) (*realtimev1.PublishEventResponse, error) {
	if s.publisher != nil {
		s.publisher.Publish(req.GetEvent(), req.GetPayloadJson(), req.GetUserIds()...)
	}
	return &realtimev1.PublishEventResponse{Delivered: int32(len(req.GetUserIds()))}, nil
}
