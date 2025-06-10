package search

import (
	"context"
	"os"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)

// SearchOptions contains all options for performing a search
type SearchOptions struct {
	// UserID is the ID of the user performing the search, can be nil for anonymous users
	UserID *string
	// ContextType is the type of context to filter by, can be empty for no context filter
	ContextType types.ContextType
}

type SearchStore interface {
	SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error)
}

func NewAlgoliaSearchStore() SearchStore {
	return NewAlgoliaStore(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_ADMIN_API_KEY"))
}
