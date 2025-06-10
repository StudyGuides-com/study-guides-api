package search

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	searchpb "github.com/studyguides-com/study-guides-api/api/v1/search"
)

type SearchStore interface {
	SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error)
	SearchUsers(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.UserSearchResult, error)
	ListIndexes(ctx context.Context) *searchpb.ListIndexesResponse
}

func NewAlgoliaSearchStore(ctx context.Context, appID, apiKey string) (SearchStore, error) {
	return NewAlgoliaStore(appID, apiKey)
}
