package search

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SearchStore interface {
	SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error)
	SearchUsers(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.UserSearchResult, error)
}

func NewAlgoliaSearchStore(ctx context.Context, appID, apiKey string) SearchStore {
	return NewAlgoliaStore(appID, apiKey)
}
