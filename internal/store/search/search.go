package search

import (
	"context"
	"os"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SearchStore interface {
	SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error)
	SearchUsers(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.UserSearchResult, error)
}

func NewAlgoliaSearchStore() SearchStore {
	return NewAlgoliaStore(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_ADMIN_API_KEY"))
}
