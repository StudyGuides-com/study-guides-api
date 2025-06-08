package search

import (
	"context"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
)

// AlgoliaStore represents an Algolia search client
type AlgoliaStore struct {
	client *search.Client
}

// NewAlgoliaSearchClient creates a new Algolia search client
func NewAlgoliaStore(appID, apiKey string) *AlgoliaStore {
	client := search.NewClient(appID, apiKey)
	
	return &AlgoliaStore{
		client: client,
	}
}

// GetIndex returns an Algolia index by name
func (c *AlgoliaStore) GetIndex(indexName string) *search.Index {
	return c.client.InitIndex(indexName)
}

// SearchTags searches for tags using Algolia
func (c *AlgoliaStore) SearchTags(ctx context.Context, contextType types.ContextType, query string) ([]*sharedpb.TagSearchResult, error) {
	index := c.GetIndex("tags")
	
	// Create search parameters with context filter
	opts := []interface{}{
		opt.Filters("context:" + string(contextType)),
	}

	// Perform the search
	res, err := index.Search(query, opts...)
	if err != nil {
		return nil, err
	}

	// Convert Algolia results to TagSearchResult
	results := make([]*sharedpb.TagSearchResult, 0, len(res.Hits))
	for _, hit := range res.Hits {
		// Create tag hierarchy from the tags array
		tagHierarchy := make([]*sharedpb.TagSearchPath, 0)
		if tags, ok := hit["tags"].([]interface{}); ok {
			for _, tag := range tags {
				if tagMap, ok := tag.(map[string]interface{}); ok {
					tagPath := &sharedpb.TagSearchPath{
						Id:   tagMap["ID"].(string),
						Name: tagMap["Name"].(string),
						Type: tagMap["Type"].(string),
					}
					tagHierarchy = append(tagHierarchy, tagPath)
				}
			}
		}

		// Create the tag search result
		tag := &sharedpb.TagSearchResult{
			Id:                 hit["id"].(string),
			Name:               hit["name"].(string),
			Description:        hit["description"].(string),
			Type:               hit["type"].(string),
			HasQuestions:       hit["hasQuestions"].(bool),
			ContentRating:      hit["contentRating"].(string),
			MetaTags:           convertToStringSlice(hit["metaTags"]),
			ContentDescriptors: convertToStringSlice(hit["contentDescriptors"]),
			TagHierarchy:       tagHierarchy,
		}
		results = append(results, tag)
	}

	return results, nil
}

// convertToStringSlice converts an interface{} to []string
func convertToStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	
	slice, ok := v.([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if str, ok := item.(string); ok {
			result = append(result, str)
		}
	}
	return result
}

