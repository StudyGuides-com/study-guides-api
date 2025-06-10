package middleware

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// SessionDetails contains information about the current user session
type SessionDetails struct {
	UserID    *string
	UserRoles *[]sharedpb.UserRole
	IsAuth    bool
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
	userRoles, _ := UserRolesFromContext(ctx)
	
	// Always return a valid SessionDetails, even if no JWT
	return &SessionDetails{
		UserID:    &userID,  // Will be empty string if not found
		UserRoles: &userRoles,  // Will be empty slice if not found
		IsAuth:    ok,  // true if we got a valid userID from context
	}
}

// HasRole checks if the user has the specified role
func (s *SessionDetails) HasRole(role sharedpb.UserRole) bool {
	if s.UserRoles == nil {
		return false
	}
	
	for _, userRole := range *s.UserRoles {
		if userRole == role {
			return true
		}
	}
	return false
}