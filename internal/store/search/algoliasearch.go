package search

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AlgoliaStore represents an Algolia search client
type AlgoliaStore struct {
	client *search.Client
}

// NewAlgoliaStore creates a new Algolia search client
func NewAlgoliaStore(appID, apiKey string) (*AlgoliaStore, error) {
	if appID == "" || apiKey == "" {
		return nil, status.Error(codes.InvalidArgument, "appID and apiKey are required")
	}

	client := search.NewClient(appID, apiKey)
	if client == nil {
		return nil, status.Error(codes.Internal, "failed to create Algolia client")
	}

	return &AlgoliaStore{
		client: client,
	}, nil
}

// GetIndex returns an Algolia index by name
func (s *AlgoliaStore) GetIndex(indexName string) *search.Index {
	return s.client.InitIndex(indexName)
}

// buildTagFilters builds the search filters based on the search options
func (c *AlgoliaStore) buildTagFilters(opts *SearchOptions) []interface{} {
	var filters []string

	// Add context filter if not "all"
	if opts.ContextType != "" && opts.ContextType != types.ContextTypeAll {
		filters = append(filters, fmt.Sprintf("context:%s", opts.ContextType))
	}

	// Add type filter
	if opts.Type != sharedpb.TagType_TAG_TYPE_UNSPECIFIED {
		filters = append(filters, fmt.Sprintf("type:%s", opts.Type))
	}

	// Handle public/private content based on context and user ID
	if opts.ContextType == types.ContextTypeUserGeneratedContent {
		if opts.UserID != nil && *opts.UserID != "" {
			filters = append(filters, fmt.Sprintf("public:false AND ownerId:%s", *opts.UserID))
		}
	} else if opts.UserID != nil && *opts.UserID != "" {
		filters = append(filters, fmt.Sprintf("(public:true OR ownerId:%s)", *opts.UserID))
	} else {
		filters = append(filters, "public:true")
	}

	// Join all filters with AND
	if len(filters) > 0 {
		return []interface{}{
			opt.Filters(strings.Join(filters, " AND ")),
		}
	}

	return nil
}

func (s *AlgoliaStore) SearchUsers(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.UserSearchResult, error) {
	// If user is authenticated and has admin role, return full results
	if opts.HasRole(sharedpb.UserRole_USER_ROLE_ADMIN) {
		log.Printf("Admin user search with query: %s", query)
		index := s.GetIndex("users")
		res, err := index.Search(query)
		if err != nil {
			log.Printf("Error searching for users: %v", err)
			return nil, err
		}
		return NewUserSearchResults(res.Hits), nil
	}

	// Non-admin users can only search for their own user record
	if opts.UserID != nil && *opts.UserID != "" {
		log.Printf("Non-admin user search with query: %s, userID: %s", query, *opts.UserID)
		index := s.GetIndex("users")
		res, err := index.Search(query, opt.Filters(fmt.Sprintf("id:%s", *opts.UserID)))
		if err != nil {
			log.Printf("Error searching for user: %v", err)
			return nil, err
		}
		return NewUserSearchResults(res.Hits), nil
	}

	return nil, status.Error(codes.PermissionDenied, "you must be an administrator to search users")
}

// SearchTags searches for tags using Algolia
func (s *AlgoliaStore) SearchTags(ctx context.Context, query string, opts *SearchOptions) ([]*sharedpb.TagSearchResult, error) {
	log.Printf("Searching for tags with query: %s, context: %v, userID: %v", query, opts.ContextType, opts.UserID)
	index := s.GetIndex("tags")

	searchOpts := s.buildTagFilters(opts)

	// Perform the search
	res, err := index.Search(query, searchOpts...)
	if err != nil {
		log.Printf("Error searching for tags: %v", err)
		return nil, err
	}

	return NewTagSearchResults(res.Hits), nil
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
		Id:                        id,
		Name:                      name,
		Type:                      sharedpb.TagType(sharedpb.TagType_value[tagType]),
		ContentRating:             sharedpb.ContentRating(sharedpb.ContentRating_value[contentRating]),
		MetaTags:                  convertToStringSlice(metaTags),
		ContentDescriptors:        convertToStringSlice(contentDescriptors),
		Tags:                      tags,
		Context:                   context,
		Public:                    public,
		HasQuestions:              hasQuestions,
		HasChildren:               hasChildren,
		BatchId:                   batchId,
		MissingMetaTags:           missingMetaTags,
		MissingContentRating:      missingContentRating,
		MissingContentDescriptors: missingContentDescriptors,
		ObjectId:                  objectID,
	}
}

// NewTagSearchResults converts Algolia search results to TagSearchResults
func NewTagSearchResults(hits []map[string]interface{}) []*sharedpb.TagSearchResult {
	results := make([]*sharedpb.TagSearchResult, 0, len(hits))
	log.Printf("Found %d tags", len(hits))
	for _, hit := range hits {
		tag := NewTagSearchResult(hit)
		results = append(results, tag)
	}
	return results
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
		Type: sharedpb.TagType(sharedpb.TagType_value[tagType]),
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

func NewUserSearchResult(hit map[string]interface{}) *sharedpb.UserSearchResult {
	id, _ := hit["id"].(string)
	name, _ := hit["name"].(string)
	email, _ := hit["email"].(string)
	gamerTag, _ := hit["gamerTag"].(string)
	createdAt, _ := hit["createdAt"].(string)
	stripeCustomerId, _ := hit["stripeCustomerId"].(string)
	objectID, _ := hit["objectID"].(string)

	return &sharedpb.UserSearchResult{
		Id:               id,
		Name:             name,
		Email:            email,
		GamerTag:         gamerTag,
		CreatedAt:        createdAt,
		StripeCustomerId: stripeCustomerId,
		ObjectId:         objectID,
	}
}

func NewUserSearchResults(hits []map[string]interface{}) []*sharedpb.UserSearchResult {
	results := make([]*sharedpb.UserSearchResult, 0, len(hits))
	log.Printf("Found %d users", len(hits))
	for _, hit := range hits {
		user := NewUserSearchResult(hit)
		results = append(results, user)
	}
	return results
}
