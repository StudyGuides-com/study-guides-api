package middleware

import "context"

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func UserRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(userRoleKey).([]string)
	return roles, ok
}