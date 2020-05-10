package ratelimiter

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	NameKey = "name"
)

type RateLimiter interface {
	Wait(name string) bool
}

func UnaryServerInterceptor(rl RateLimiter) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (
		interface{},
		error,
	) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return "", status.Error(codes.Unauthenticated, "name requried")
		}

		names := md.Get(NameKey)
		if len(names) == 0 {
			return "", status.Error(codes.Unauthenticated, "name requried")
		}

		if !rl.Wait(names[0]) {
			return "", status.Error(codes.ResourceExhausted, "exceeded rate limit")
		}

		return handler(ctx, req)
	}
}
