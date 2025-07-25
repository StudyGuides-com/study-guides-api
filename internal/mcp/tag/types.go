package tag

import (
	"encoding/json"
	"strings"
	"time"

	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// AccessType represents access levels (not in shared protos)
type AccessType string

const (
	AccessTypePublic    AccessType = "Public"
	AccessTypePrivate   AccessType = "Private"
	AccessTypeReadOnly  AccessType = "ReadOnly"
	AccessTypeReadWrite AccessType = "ReadWrite"
)

// Tag represents the domain entity for tags
type Tag struct {
	ID                 string                      `json:"id"`
	BatchID            *string                     `json:"batchId,omitempty"`
	Hash               string                      `json:"hash"`
	Name               string                      `json:"name"`
	Description        *string                     `json:"description,omitempty"`
	Type               sharedv1.TagType            `json:"type"`
	Context            *sharedv1.ContextType       `json:"context,omitempty"`
	ParentTagID        *string                     `json:"parentTagId,omitempty"`
	HasChildren        bool                        `json:"hasChildren"`
	HasQuestions       bool                        `json:"hasQuestions"`
	ContentRating      sharedv1.ContentRating      `json:"contentRating"`
	ContentDescriptors []string                    `json:"contentDescriptors"`
	MetaTags           []string                    `json:"metaTags"`
	Public             bool                        `json:"public"`
	OwnerID            *string                     `json:"ownerId,omitempty"`
	AccessCount        int                         `json:"accessCount"`
	Metadata           map[string]interface{}      `json:"metadata,omitempty"`
	CreatedAt          time.Time                   `json:"createdAt"`
	UpdatedAt          time.Time                   `json:"updatedAt"`
}

// TagFilter defines the available filters for querying tags
type TagFilter struct {
	// Basic filters
	Type        *sharedv1.TagType     `json:"type,omitempty"`
	Context     *sharedv1.ContextType `json:"context,omitempty"`
	Name        *string               `json:"name,omitempty"`          // Partial name search
	ParentID    *string               `json:"parentId,omitempty"`      // Find children of parent
	OwnerID     *string               `json:"ownerId,omitempty"`       // Filter by owner
	Public      *bool                 `json:"public,omitempty"`        // Public/private filter
	
	// Structure filters
	HasChildren  *bool `json:"hasChildren,omitempty"`   // Tags with/without children
	HasQuestions *bool `json:"hasQuestions,omitempty"`  // Tags with/without questions
	IsRoot       *bool `json:"isRoot,omitempty"`        // Root tags (no parent)
	
	// Content filters
	ContentRating      *sharedv1.ContentRating `json:"contentRating,omitempty"`
	ContentDescriptors []string                    `json:"contentDescriptors,omitempty"`
	MetaTags           []string                    `json:"metaTags,omitempty"`
	
	// Pagination
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
	
	// Sorting
	SortBy    *string `json:"sortBy,omitempty"`    // "name", "createdAt", "updatedAt", "accessCount"
	SortOrder *string `json:"sortOrder,omitempty"` // "asc", "desc"
}

// TagUpdate defines the fields that can be updated for a tag
type TagUpdate struct {
	Name               *string                     `json:"name,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	Type               *sharedv1.TagType           `json:"type,omitempty"`
	Context            *sharedv1.ContextType       `json:"context,omitempty"`
	ParentTagID        *string                     `json:"parentTagId,omitempty"`
	ContentRating      *sharedv1.ContentRating `json:"contentRating,omitempty"`
	ContentDescriptors []string                    `json:"contentDescriptors,omitempty"`
	MetaTags           []string                    `json:"metaTags,omitempty"`
	Public             *bool                       `json:"public,omitempty"`
	Metadata           map[string]interface{}      `json:"metadata,omitempty"`
}

// TagCreate defines the fields required to create a new tag
type TagCreate struct {
	Name               string                      `json:"name" validate:"required"`
	Description        *string                     `json:"description,omitempty"`
	Type               sharedv1.TagType            `json:"type" validate:"required"`
	Context            *sharedv1.ContextType       `json:"context,omitempty"`
	ParentTagID        *string                     `json:"parentTagId,omitempty"`
	ContentRating      *sharedv1.ContentRating `json:"contentRating,omitempty"`
	ContentDescriptors []string                    `json:"contentDescriptors,omitempty"`
	MetaTags           []string                    `json:"metaTags,omitempty"`
	Public             *bool                       `json:"public,omitempty"`
	OwnerID            *string                     `json:"ownerId,omitempty"`
	Metadata           map[string]interface{}      `json:"metadata,omitempty"`
}

// TagStats represents aggregated statistics about tags
type TagStats struct {
	TotalTags       int                                        `json:"totalTags"`
	PublicTags      int                                        `json:"publicTags"`
	PrivateTags     int                                        `json:"privateTags"`
	TagsByType      map[sharedv1.TagType]int                   `json:"tagsByType"`
	TagsByContext   map[sharedv1.ContextType]int               `json:"tagsByContext"`
	TagsByRating    map[sharedv1.ContentRating]int         `json:"tagsByRating"`
}

// TagHierarchy represents a tag with its children for tree operations
type TagHierarchy struct {
	Tag      Tag             `json:"tag"`
	Children []TagHierarchy  `json:"children,omitempty"`
	Depth    int             `json:"depth"`
}

// TagSearchOptions defines advanced search options
type TagSearchOptions struct {
	Query           string   `json:"query,omitempty"`           // Full-text search
	Fields          []string `json:"fields,omitempty"`          // Fields to search in
	BoostFields     []string `json:"boostFields,omitempty"`     // Fields to boost in relevance
	FacetFilters    []string `json:"facetFilters,omitempty"`    // Algolia-style facet filters
	NumericFilters  []string `json:"numericFilters,omitempty"`  // Numeric range filters
	HighlightFields []string `json:"highlightFields,omitempty"` // Fields to highlight in results
}

// UnmarshalJSON provides custom JSON unmarshaling for TagFilter to handle string-to-enum conversion
func (tf *TagFilter) UnmarshalJSON(data []byte) error {
	// Define a temporary struct with string fields
	type tagFilterAlias struct {
		// Basic filters
		Type        *string `json:"type,omitempty"`
		Context     *string `json:"context,omitempty"`
		Name        *string `json:"name,omitempty"`
		ParentID    *string `json:"parentId,omitempty"`
		OwnerID     *string `json:"ownerId,omitempty"`
		Public      *bool   `json:"public,omitempty"`
		
		// Structure filters
		HasChildren  *bool `json:"hasChildren,omitempty"`
		HasQuestions *bool `json:"hasQuestions,omitempty"`
		IsRoot       *bool `json:"isRoot,omitempty"`
		
		// Content filters
		ContentRating      *string  `json:"contentRating,omitempty"`
		ContentDescriptors []string `json:"contentDescriptors,omitempty"`
		MetaTags           []string `json:"metaTags,omitempty"`
		
		// Pagination and sorting
		Limit     *int    `json:"limit,omitempty"`
		Offset    *int    `json:"offset,omitempty"`
		SortBy    *string `json:"sortBy,omitempty"`
		SortOrder *string `json:"sortOrder,omitempty"`
	}
	
	var temp tagFilterAlias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	
	// Copy all direct fields
	tf.Name = temp.Name
	tf.ParentID = temp.ParentID
	tf.OwnerID = temp.OwnerID
	tf.Public = temp.Public
	tf.HasChildren = temp.HasChildren
	tf.HasQuestions = temp.HasQuestions
	tf.IsRoot = temp.IsRoot
	tf.ContentDescriptors = temp.ContentDescriptors
	tf.MetaTags = temp.MetaTags
	tf.Limit = temp.Limit
	tf.Offset = temp.Offset
	tf.SortBy = temp.SortBy
	tf.SortOrder = temp.SortOrder
	
	// Convert Type string to enum
	if temp.Type != nil {
		typeStr := strings.ToUpper(*temp.Type)
		if val, ok := sharedv1.TagType_value[typeStr]; ok {
			enumVal := sharedv1.TagType(val)
			tf.Type = &enumVal
		}
	}
	
	// Convert Context string to enum
	if temp.Context != nil {
		contextStr := strings.ToUpper(*temp.Context)
		if val, ok := sharedv1.ContextType_value[contextStr]; ok {
			enumVal := sharedv1.ContextType(val)
			tf.Context = &enumVal
		}
	}
	
	// Convert ContentRating string to enum
	if temp.ContentRating != nil {
		ratingStr := strings.ToUpper(*temp.ContentRating)
		if val, ok := sharedv1.ContentRating_value[ratingStr]; ok {
			enumVal := sharedv1.ContentRating(val)
			tf.ContentRating = &enumVal
		}
	}
	
	return nil
}