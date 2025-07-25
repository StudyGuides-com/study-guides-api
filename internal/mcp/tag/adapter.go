package tag

import (
	"context"
	"fmt"

	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"github.com/studyguides-com/study-guides-api/internal/store/tag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TagRepositoryAdapter adapts the existing SqlTagStore to implement TagRepository
type TagRepositoryAdapter struct {
	store tag.TagStore
}

// NewTagRepositoryAdapter creates a new adapter for the existing tag store
func NewTagRepositoryAdapter(store tag.TagStore) TagRepository {
	return &TagRepositoryAdapter{
		store: store,
	}
}

// Generic CRUD operations

// Find implements the generic find operation with rich filtering
func (a *TagRepositoryAdapter) Find(ctx context.Context, filter TagFilter) ([]Tag, error) {
	// Convert our rich filter to the existing store's parameter format
	params := a.filterToParams(filter)
	
	// Use existing store method
	pbTags, err := a.store.ListTagsWithFilters(ctx, params)
	if err != nil {
		return nil, err
	}
	
	// Convert protobuf tags to domain tags
	domainTags := make([]Tag, len(pbTags))
	for i, pbTag := range pbTags {
		domainTag, err := a.pbTagToDomainTag(pbTag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag %s: %w", pbTag.Id, err)
		}
		domainTags[i] = *domainTag
	}
	
	return domainTags, nil
}

// FindByID implements the generic findById operation
func (a *TagRepositoryAdapter) FindByID(ctx context.Context, id string) (*Tag, error) {
	pbTag, err := a.store.GetTagByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return a.pbTagToDomainTag(pbTag)
}

// Create implements the generic create operation
func (a *TagRepositoryAdapter) Create(ctx context.Context, entity Tag) (*Tag, error) {
	// For now, return error since we don't have create in existing store
	// TODO: Implement when needed
	return nil, fmt.Errorf("create operation not yet implemented")
}

// Update implements the generic update operation
func (a *TagRepositoryAdapter) Update(ctx context.Context, id string, update TagUpdate) (*Tag, error) {
	// For now, return error since we don't have update in existing store
	// TODO: Implement when needed
	return nil, fmt.Errorf("update operation not yet implemented")
}

// Delete implements the generic delete operation
func (a *TagRepositoryAdapter) Delete(ctx context.Context, id string) error {
	// For now, return error since we don't have delete in existing store
	// TODO: Implement when needed
	return fmt.Errorf("delete operation not yet implemented")
}

// Count implements the generic count operation
func (a *TagRepositoryAdapter) Count(ctx context.Context, filter TagFilter) (int, error) {
	params := a.filterToParams(filter)
	return a.store.CountTags(ctx, params)
}

// Tag-specific operations

// FindByParent returns all direct children of a parent tag
func (a *TagRepositoryAdapter) FindByParent(ctx context.Context, parentID string) ([]Tag, error) {
	pbTags, err := a.store.ListTagsByParent(ctx, parentID)
	if err != nil {
		return nil, err
	}
	
	domainTags := make([]Tag, len(pbTags))
	for i, pbTag := range pbTags {
		domainTag, err := a.pbTagToDomainTag(pbTag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag %s: %w", pbTag.Id, err)
		}
		domainTags[i] = *domainTag
	}
	
	return domainTags, nil
}

// FindRoots returns all root tags (no parent) with optional filtering
func (a *TagRepositoryAdapter) FindRoots(ctx context.Context, filter TagFilter) ([]Tag, error) {
	params := a.filterToParams(filter)
	pbTags, err := a.store.ListRootTags(ctx, params)
	if err != nil {
		return nil, err
	}
	
	domainTags := make([]Tag, len(pbTags))
	for i, pbTag := range pbTags {
		domainTag, err := a.pbTagToDomainTag(pbTag)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tag %s: %w", pbTag.Id, err)
		}
		domainTags[i] = *domainTag
	}
	
	return domainTags, nil
}

// UniqueTypes returns all unique tag types in the system
func (a *TagRepositoryAdapter) UniqueTypes(ctx context.Context) ([]sharedv1.TagType, error) {
	return a.store.UniqueTagTypes(ctx)
}

// UniqueContexts returns all unique context types in the system
func (a *TagRepositoryAdapter) UniqueContexts(ctx context.Context) ([]sharedv1.ContextType, error) {
	contextStrings, err := a.store.UniqueContextTypes(ctx)
	if err != nil {
		return nil, err
	}
	
	// Convert strings to enum values
	contexts := make([]sharedv1.ContextType, 0, len(contextStrings))
	for _, contextStr := range contextStrings {
		if contextType := sharedv1.ContextType(sharedv1.ContextType_value[contextStr]); contextType != 0 {
			contexts = append(contexts, contextType)
		}
	}
	
	return contexts, nil
}

// Stub implementations for methods not yet needed
func (a *TagRepositoryAdapter) FindHierarchy(ctx context.Context, tagID string, maxDepth int) (*TagHierarchy, error) {
	return nil, fmt.Errorf("FindHierarchy not yet implemented")
}

func (a *TagRepositoryAdapter) FindAncestors(ctx context.Context, tagID string) ([]Tag, error) {
	return nil, fmt.Errorf("FindAncestors not yet implemented")
}

func (a *TagRepositoryAdapter) GetStats(ctx context.Context, filter TagFilter) (*TagStats, error) {
	return nil, fmt.Errorf("GetStats not yet implemented")
}

func (a *TagRepositoryAdapter) Search(ctx context.Context, options TagSearchOptions) ([]Tag, error) {
	return nil, fmt.Errorf("Search not yet implemented")
}

func (a *TagRepositoryAdapter) ValidateHierarchy(ctx context.Context, childID, parentID string) error {
	return nil // No validation for now
}

func (a *TagRepositoryAdapter) UpdateHierarchyFlags(ctx context.Context, tagID string) error {
	return nil // No-op for now
}

func (a *TagRepositoryAdapter) BulkUpdateParent(ctx context.Context, tagIDs []string, newParentID *string) error {
	return fmt.Errorf("BulkUpdateParent not yet implemented")
}

// Helper methods

// filterToParams converts TagFilter to the map[string]string format expected by existing store
func (a *TagRepositoryAdapter) filterToParams(filter TagFilter) map[string]string {
	params := make(map[string]string)
	
	if filter.Type != nil {
		params["type"] = filter.Type.String()
	}
	
	if filter.Context != nil {
		params["contextType"] = filter.Context.String()
	}
	
	if filter.Name != nil {
		params["name"] = *filter.Name
	}
	
	if filter.Public != nil {
		if *filter.Public {
			params["public"] = "true"
		} else {
			params["public"] = "false"
		}
	}
	
	if filter.Limit != nil {
		params["limit"] = fmt.Sprintf("%d", *filter.Limit)
	}
	
	if filter.Offset != nil {
		params["offset"] = fmt.Sprintf("%d", *filter.Offset)
	}
	
	if filter.SortBy != nil {
		params["sortBy"] = *filter.SortBy
	}
	
	if filter.SortOrder != nil {
		params["sortOrder"] = *filter.SortOrder
	}
	
	// Handle IsRoot filter by setting/unsetting parentId
	if filter.IsRoot != nil && *filter.IsRoot {
		params["parentId"] = "" // Empty string indicates root tags
	}
	
	return params
}

// pbTagToDomainTag converts a protobuf Tag to a domain Tag
func (a *TagRepositoryAdapter) pbTagToDomainTag(pbTag *sharedv1.Tag) (*Tag, error) {
	if pbTag == nil {
		return nil, fmt.Errorf("nil protobuf tag")
	}
	
	// Convert ContentDescriptors from enum slice to string slice
	contentDescriptors := make([]string, len(pbTag.ContentDescriptors))
	for i, cd := range pbTag.ContentDescriptors {
		contentDescriptors[i] = cd.String()
	}

	domainTag := &Tag{
		ID:                 pbTag.Id,
		Hash:               pbTag.Hash,
		Name:               pbTag.Name,
		Type:               pbTag.Type,
		HasChildren:        pbTag.HasChildren,
		HasQuestions:       pbTag.HasQuestions,
		ContentRating:      pbTag.ContentRating,
		ContentDescriptors: contentDescriptors,
		MetaTags:           pbTag.MetaTags,
		Public:             pbTag.Public,
		AccessCount:        int(pbTag.AccessCount),
	}
	
	// Handle optional fields (they're already pointers in protobuf)
	if pbTag.BatchId != nil && *pbTag.BatchId != "" {
		domainTag.BatchID = pbTag.BatchId
	}
	
	if pbTag.Description != nil && *pbTag.Description != "" {
		domainTag.Description = pbTag.Description
	}
	
	// Context is not optional in protobuf, but we treat it as optional in domain
	domainTag.Context = &pbTag.Context
	
	if pbTag.ParentTagId != nil && *pbTag.ParentTagId != "" {
		domainTag.ParentTagID = pbTag.ParentTagId
	}
	
	if pbTag.OwnerId != nil && *pbTag.OwnerId != "" {
		domainTag.OwnerID = pbTag.OwnerId
	}
	
	// Handle timestamps
	if pbTag.CreatedAt != nil {
		domainTag.CreatedAt = pbTag.CreatedAt.AsTime()
	}
	
	if pbTag.UpdatedAt != nil {
		domainTag.UpdatedAt = pbTag.UpdatedAt.AsTime()
	}
	
	// Handle metadata - convert from protobuf Metadata to map[string]interface{}
	if pbTag.Metadata != nil && pbTag.Metadata.Metadata != nil {
		metadata := make(map[string]interface{})
		for k, v := range pbTag.Metadata.Metadata {
			metadata[k] = v
		}
		domainTag.Metadata = metadata
	}
	
	return domainTag, nil
}

// domainTagToPbTag converts a domain Tag back to protobuf (for create/update operations)
func (a *TagRepositoryAdapter) domainTagToPbTag(domainTag *Tag) (*sharedv1.Tag, error) {
	if domainTag == nil {
		return nil, fmt.Errorf("nil domain tag")
	}
	
	// Convert ContentDescriptors from string slice to enum slice
	contentDescriptors := make([]sharedv1.ContentDescriptorType, len(domainTag.ContentDescriptors))
	for i, cd := range domainTag.ContentDescriptors {
		if val, ok := sharedv1.ContentDescriptorType_value[cd]; ok {
			contentDescriptors[i] = sharedv1.ContentDescriptorType(val)
		}
	}

	pbTag := &sharedv1.Tag{
		Id:                 domainTag.ID,
		Hash:               domainTag.Hash,
		Name:               domainTag.Name,
		Type:               domainTag.Type,
		HasChildren:        domainTag.HasChildren,
		HasQuestions:       domainTag.HasQuestions,
		ContentRating:      domainTag.ContentRating,
		ContentDescriptors: contentDescriptors,
		MetaTags:           domainTag.MetaTags,
		Public:             domainTag.Public,
		AccessCount:        int32(domainTag.AccessCount),
	}
	
	// Handle optional fields (protobuf expects pointers)
	if domainTag.BatchID != nil {
		pbTag.BatchId = domainTag.BatchID
	}
	
	if domainTag.Description != nil {
		pbTag.Description = domainTag.Description
	}
	
	if domainTag.Context != nil {
		pbTag.Context = *domainTag.Context
	}
	
	if domainTag.ParentTagID != nil {
		pbTag.ParentTagId = domainTag.ParentTagID
	}
	
	if domainTag.OwnerID != nil {
		pbTag.OwnerId = domainTag.OwnerID
	}
	
	// Handle timestamps
	if !domainTag.CreatedAt.IsZero() {
		pbTag.CreatedAt = timestamppb.New(domainTag.CreatedAt)
	}
	
	if !domainTag.UpdatedAt.IsZero() {
		pbTag.UpdatedAt = timestamppb.New(domainTag.UpdatedAt)
	}
	
	// Handle metadata - convert from map[string]interface{} to protobuf Metadata
	if domainTag.Metadata != nil {
		pbMetadata := &sharedv1.Metadata{
			Metadata: make(map[string]string),
		}
		for k, v := range domainTag.Metadata {
			// Convert interface{} to string
			pbMetadata.Metadata[k] = fmt.Sprintf("%v", v)
		}
		pbTag.Metadata = pbMetadata
	}
	
	return pbTag, nil
}