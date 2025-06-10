package search

import (
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
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