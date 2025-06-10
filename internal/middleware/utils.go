package middleware

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// SessionDetails contains information about the current user session
type SessionDetails struct {
	UserID    *string
	UserRoles *[]sharedpb.UserRole
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}

func UserRolesFromContext(ctx context.Context) ([]sharedpb.UserRole, bool) {
	roles, ok := ctx.Value(userRoleKey).([]sharedpb.UserRole)
	return roles, ok
}

// GetSessionDetails extracts session information from the context
func GetSessionDetails(ctx context.Context) *SessionDetails {
	userID, ok := UserIDFromContext(ctx)
	userRoles, rolesOk := UserRolesFromContext(ctx)
	
	if !ok || !rolesOk {
		return &SessionDetails{
			UserID:    nil,
			UserRoles: nil,
		}
	}
	
	return &SessionDetails{
		UserID:    &userID,
		UserRoles: &userRoles,
	}
}