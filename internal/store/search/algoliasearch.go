package search

import (
	"context"
	"fmt"
	"log"

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
func (c *AlgoliaStore) SearchTagsForContext(ctx context.Context, contextType types.ContextType, query string) ([]*sharedpb.TagSearchResult, error) {
	
	log.Printf("Searching for tags with query: %s, contextType: %s", query, contextType)
	index := c.GetIndex("tags")
	
	// Create search parameters with context filter
	opts := []interface{}{
		opt.Filters("context:" + string(contextType)),
	}

	// Perform the search
	res, err := index.Search(query, opts...)
	if err != nil {
		log.Printf("Error searching for tags: %v", err)
		return nil, err
	}

	// Convert Algolia results to TagSearchResult
	results := make([]*sharedpb.TagSearchResult, 0, len(res.Hits))
	log.Printf("Found %d tags", len(res.Hits))
	for _, hit := range res.Hits {
		tag := NewTagSearchResult(hit)
		results = append(results, tag)
	}

	return results, nil
}

func (c *AlgoliaStore) SearchTags(ctx context.Context, tagID string, query string) ([]*sharedpb.TagSearchResult, error) {
	log.Printf("Searching for tags with query: %s, tagID: %s", query, tagID)
	index := c.GetIndex("tags")

	opts := []interface{}{
		opt.Filters("id:" + tagID),
	}
}

// NewTagSearchResult creates a TagSearchResult from a hit
func NewTagSearchResult(hit map[string]interface{}) *sharedpb.TagSearchResult {
	// Safely extract fields with nil checks
	id, _ := hit["id"].(string)
	name, _ := hit["name"].(string)
	metaTags, _ := hit["metaTags"].([]interface{})
	tagType, _ := hit["type"].(string)
	contentRating, _ := hit["contentRating"].(string)
	tags := NewTagSearchPaths(hit)
	public, _ := hit["public"].(bool)
	contentDescriptors, _ := hit["contentDescriptors"].([]interface{})
	context, _ := hit["context"].(string)
	hasQuestions, _ := hit["hasQuestions"].(bool)
	hasChildren, _ := hit["hasChildren"].(bool)
	batchId, _ := hit["batchId"].(string)
	missingMetaTags, _ := hit["missingMetaTags"].(bool)
	missingContentRating, _ := hit["missingContentRating"].(bool)
	missingContentDescriptors, _ := hit["missingContentDescriptors"].(bool)
	objectID, _ := hit["objectID"].(string)

	// Create the tag search result
	return &sharedpb.TagSearchResult{
		Id:                      id,
		Name:                    name,
		Type:                    tagType,
		ContentRating:           contentRating,
		MetaTags:                convertToStringSlice(metaTags),
		ContentDescriptors:      convertToStringSlice(contentDescriptors),
		Tags:                    tags,
		Context:                 context,
		Public:                  public,
		HasQuestions:            hasQuestions,
		HasChildren:             hasChildren,
		BatchId:                 batchId,
		MissingMetaTags:         missingMetaTags,
		MissingContentRating:    missingContentRating,
		MissingContentDescriptors: missingContentDescriptors,
		ObjectId:                objectID,
	}
}

// NewTagSearchPath creates a TagSearchPath from a map of tag data
func NewTagSearchPath(tagMap map[string]interface{}, index int) (*sharedpb.TagSearchPath, error) {
	// Safely extract fields with nil checks
	id, idOk := tagMap["id"].(string)
	name, nameOk := tagMap["name"].(string)
	tagType, typeOk := tagMap["type"].(string)

	// Log if any required fields are missing
	if !idOk || !nameOk || !typeOk {
		log.Printf("Warning: Missing required fields in tag object at index %d. ID: %v, Name: %v, Type: %v", 
			index, idOk, nameOk, typeOk)
		return nil, fmt.Errorf("missing required fields in tag object")
	}
	
	return &sharedpb.TagSearchPath{
		Id:   id,
		Name: name,
		Type: tagType,
	}, nil
}

// NewTagSearchPaths processes the tags array from a hit and returns a slice of TagSearchPath
func NewTagSearchPaths(hit map[string]interface{}) []*sharedpb.TagSearchPath {
	tagHierarchy := make([]*sharedpb.TagSearchPath, 0)
	if tags, ok := hit["tags"].([]interface{}); ok {
		for i, tag := range tags {
			if tagMap, ok := tag.(map[string]interface{}); ok {
				// Check if the map is empty
				if len(tagMap) == 0 {
					log.Printf("Warning: Empty tag object found at index %d in tag hierarchy", i)
					continue
				}

				tagPath, err := NewTagSearchPath(tagMap, i)
				if err != nil {
					continue
				}
				tagHierarchy = append(tagHierarchy, tagPath)
			} else {
				log.Printf("Warning: Invalid tag object type at index %d: %T", i, tag)
			}
		}
	} else {
		log.Printf("Warning: 'tags' field is not an array or is missing in search result")
	}
	return tagHierarchy
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

