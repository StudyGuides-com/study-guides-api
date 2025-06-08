package search

import (
	"context"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)




type SearchStore interface {
	SearchTags(ctx context.Context, contextType types.ContextType, query string) ([]*sharedpb.TagSearchResult, error)
}