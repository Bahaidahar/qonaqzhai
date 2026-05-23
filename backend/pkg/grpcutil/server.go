package grpcutil

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingUnaryInterceptor logs every unary RPC with method + duration + code.
func LoggingUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)
		log.Info("grpc",
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Duration("dur", time.Since(start)),
		)
		return resp, err
	}
}

// RecoverUnaryInterceptor turns panics into Internal errors.
func RecoverUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("grpc panic",
					slog.String("method", info.FullMethod),
					slog.Any("err", rec),
					slog.String("stack", string(debug.Stack())),
				)
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}
