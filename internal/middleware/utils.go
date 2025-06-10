package middleware

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func UserRolesFromContext(ctx context.Context) ([]sharedpb.UserRole, bool) {
	roles, ok := ctx.Value(userRoleKey).([]sharedpb.UserRole)
	return roles, ok
}