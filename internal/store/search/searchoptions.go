package search

import (
	"context"

	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/types"
)

// SearchOptions contains all options for performing a search
type SearchOptions struct {
	// UserID is the ID of the user performing the search, can be nil for anonymous users
	UserID *string
	// ContextType is the type of context to filter by, can be empty for no context filter
	ContextType types.ContextType
	// UserRoles is the roles of the user performing the search, can be nil for anonymous users
	UserRoles *[]sharedpb.UserRole
	// Type is the type of tag to filter by
	Type sharedpb.TagType
}

// NewSearchOptions creates a new SearchOptions with optional fields
func NewSearchOptions(opts ...func(*SearchOptions)) *SearchOptions {
	options := &SearchOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithUserID sets the UserID field
func WithUserID(userID *string) func(*SearchOptions) {
	return func(o *SearchOptions) {
		o.UserID = userID
	}
}

// WithContextType sets the ContextType field
func WithContextType(contextType types.ContextType) func(*SearchOptions) {
	return func(o *SearchOptions) {
		o.ContextType = contextType
	}
}

// WithUserRoles sets the UserRoles field
func WithUserRoles(userRoles *[]sharedpb.UserRole) func(*SearchOptions) {
	return func(o *SearchOptions) {
		o.UserRoles = userRoles
	}
}

// WithType sets the Type field
func WithType(tagType sharedpb.TagType) func(*SearchOptions) {
	return func(o *SearchOptions) {
		o.Type = tagType
	}
}

// FromProtoContextType converts a proto ContextType to our internal types.ContextType
func FromProtoContextType(ct sharedpb.ContextType) types.ContextType {
	switch ct {
	case sharedpb.ContextType_CONTEXT_TYPE_COLLEGES:
		return types.ContextTypeColleges
	case sharedpb.ContextType_CONTEXT_TYPE_CERTIFICATIONS:
		return types.ContextTypeCertifications
	case sharedpb.ContextType_CONTEXT_TYPE_ENTRANCE_EXAMS:
		return types.ContextTypeEntranceExams
	case sharedpb.ContextType_CONTEXT_TYPE_AP_EXAMS:
		return types.ContextTypeAPExams
	case sharedpb.ContextType_CONTEXT_TYPE_DOD:
		return types.ContextTypeDoD
	default:
		return types.ContextTypeAll
	}
}

// NewSearchOptionsFromRequest creates a new SearchOptions from the context and request
func NewSearchOptionsFromRequest(ctx context.Context, req *searchpb.SearchTagsRequest) *SearchOptions {
	session := middleware.GetSessionDetails(ctx)
	return NewSearchOptions(
		WithUserID(session.UserID),
		WithUserRoles(session.UserRoles),
		WithContextType(FromProtoContextType(req.Context)),
	)
}

// NewSearchOptionsFrom creates a new SearchOptions from the context for user search
func NewSearchOptionsFrom(ctx context.Context) *SearchOptions {
	session := middleware.GetSessionDetails(ctx)
	return NewSearchOptions(
		WithUserID(session.UserID),
		WithUserRoles(session.UserRoles),
	)
}

// HasRole checks if the user has the specified role
func (s *SearchOptions) HasRole(role sharedpb.UserRole) bool {
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