// internal/middleware/ratelimit.go
package middleware

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type limiterStore struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	r        rate.Limit
	burst    int
}

func newLimiterStore(r rate.Limit, burst int) *limiterStore {
	return &limiterStore{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		burst:    burst,
	}
}

func (s *limiterStore) getLimiter(key string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	limiter, ok := s.limiters[key]
	if !ok {
		limiter = rate.NewLimiter(s.r, s.burst)
		s.limiters[key] = limiter
	}
	return limiter
}

func RateLimitUnaryInterceptor(r rate.Limit, burst int) grpc.UnaryServerInterceptor {
	store := newLimiterStore(r, burst)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Try to extract userID from JWT first
		var key string
		if userID, ok := UserIDFromContext(ctx); ok {
			key = userID
		} else {
			// fallback to IP address
			if p, ok := peer.FromContext(ctx); ok {
				key = p.Addr.String()
			} else {
				key = "unknown"
			}
		}

		limiter := store.getLimiter(key)
		if !limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}
