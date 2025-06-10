// internal/service/base.go

package services

import (
	"context"

	"github.com/studyguides-com/study-guides-api/internal/middleware"
)

// PublicBaseHandler is used when the request may or may not be authenticated.
func PublicBaseHandler(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	return fn(ctx)
}

// AuthBaseHandler is used when authentication is required.
func AuthBaseHandler(ctx context.Context, fn func(ctx context.Context, userID *string, userRoles *[]string) (interface{}, error)) (interface{}, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	userRoles, rolesOk := middleware.UserRolesFromContext(ctx)
	if ok && rolesOk {
		return fn(ctx, &userID, &userRoles)
	}
	return fn(ctx, nil, nil)
}
