// internal/service/base.go

package services

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PublicBaseHandler is used when the request may or may not be authenticated.
func PublicBaseHandler(ctx context.Context, fn func(ctx context.Context, userID *string) (interface{}, error)) (interface{}, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if ok {
		return fn(ctx, &userID)
	}
	return fn(ctx, nil)
}

// AuthBaseHandler is used when authentication is required.
func AuthBaseHandler(ctx context.Context, fn func(ctx context.Context, userID string) (interface{}, error)) (interface{}, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}
	return fn(ctx, userID)
}
