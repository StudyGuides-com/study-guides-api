package search

import (
	"context"
	"os"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)




type SearchStore interface {
	SearchTags(ctx context.Context, contextType types.ContextType, query string) ([]*sharedpb.TagSearchResult, error)
}


func NewAlgoliaSearchStore() SearchStore {
	return NewAlgoliaStore(os.Getenv("ALGOLIA_APP_ID"), os.Getenv("ALGOLIA_ADMIN_API_KEY"))
}
