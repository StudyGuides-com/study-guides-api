package search

import (
	"context"

	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SearchStore interface {
	SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error)
	SearchUsers(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.UserSearchResult, error)
	ListIndexes(ctx context.Context) *searchpb.ListIndexesResponse
}

func NewAlgoliaSearchStore(ctx context.Context, appID, apiKey string) (SearchStore, error) {
	return NewAlgoliaStore(appID, apiKey)
}
