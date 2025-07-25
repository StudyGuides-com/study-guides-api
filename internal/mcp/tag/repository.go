package tag

import (
	"context"
	"time"

	"github.com/studyguides-com/study-guides-api/internal/repository"
	sharedv1 "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// TagRepository extends the generic repository with tag-specific operations
type TagRepository interface {
	// Generic CRUD operations
	repository.Repository[Tag, TagFilter, TagUpdate]

	// Tag-specific operations
	
	// FindByParent returns all direct children of a parent tag
	FindByParent(ctx context.Context, parentID string) ([]Tag, error)
	
	// FindRoots returns all root tags (no parent) with optional filtering
	FindRoots(ctx context.Context, filter TagFilter) ([]Tag, error)
	
	// FindHierarchy returns a tag with all its descendants in a tree structure
	FindHierarchy(ctx context.Context, tagID string, maxDepth int) (*TagHierarchy, error)
	
	// FindAncestors returns all ancestor tags up to the root
	FindAncestors(ctx context.Context, tagID string) ([]Tag, error)
	
	// UniqueTypes returns all unique tag types in the system
	UniqueTypes(ctx context.Context) ([]sharedv1.TagType, error)
	
	// UniqueContexts returns all unique context types in the system
	UniqueContexts(ctx context.Context) ([]sharedv1.ContextType, error)
	
	// GetStats returns aggregated statistics about tags
	GetStats(ctx context.Context, filter TagFilter) (*TagStats, error)
	
	// Search performs full-text search with advanced options
	Search(ctx context.Context, options TagSearchOptions) ([]Tag, error)
	
	// ValidateHierarchy checks if a parent-child relationship would create a cycle
	ValidateHierarchy(ctx context.Context, childID, parentID string) error
	
	// UpdateHierarchyFlags updates hasChildren flags after hierarchy changes
	UpdateHierarchyFlags(ctx context.Context, tagID string) error
	
	// BulkUpdateParent moves multiple tags under a new parent
	BulkUpdateParent(ctx context.Context, tagIDs []string, newParentID *string) error
}

// TagAccessRepository handles tag access control operations
type TagAccessRepository interface {
	// GrantAccess gives a user access to a tag
	GrantAccess(ctx context.Context, tagID, userID string, accessType AccessType) error
	
	// RevokeAccess removes a user's access to a tag
	RevokeAccess(ctx context.Context, tagID, userID string) error
	
	// GetUserAccess returns a user's access level for a tag
	GetUserAccess(ctx context.Context, tagID, userID string) (*AccessType, error)
	
	// ListUserTags returns all tags a user has access to
	ListUserTags(ctx context.Context, userID string, accessType *AccessType) ([]Tag, error)
	
	// ListTagUsers returns all users with access to a tag
	ListTagUsers(ctx context.Context, tagID string) ([]TagAccess, error)
}

// TagAccess represents a user's access to a tag
type TagAccess struct {
	TagID      string                `json:"tagId"`
	UserID     string                `json:"userId"`
	AccessType AccessType   `json:"accessType"`
	CreatedAt  time.Time             `json:"createdAt"`
}

// TagInteractionRepository handles user interactions with tags
type TagInteractionRepository interface {
	// AddToFavorites adds a tag to user's favorites
	AddToFavorites(ctx context.Context, userID, tagID string) error
	
	// RemoveFromFavorites removes a tag from user's favorites
	RemoveFromFavorites(ctx context.Context, userID, tagID string) error
	
	// GetFavorites returns a user's favorite tags
	GetFavorites(ctx context.Context, userID string) ([]Tag, error)
	
	// AddToRecent adds a tag to user's recent history
	AddToRecent(ctx context.Context, userID, tagID string) error
	
	// GetRecent returns a user's recently viewed tags
	GetRecent(ctx context.Context, userID string, limit int) ([]Tag, error)
	
	// RateTag allows a user to rate a tag
	RateTag(ctx context.Context, userID, tagID string, rating int) error
	
	// GetRating returns a user's rating for a tag
	GetRating(ctx context.Context, userID, tagID string) (*int, error)
	
	// GetAverageRating returns the average rating for a tag
	GetAverageRating(ctx context.Context, tagID string) (*float64, error)
	
	// ReportTag allows a user to report a tag
	ReportTag(ctx context.Context, userID, tagID string, reportType sharedv1.ReportType, reason string) error
}

// ResourceConstants for registration
const (
	ResourceName = "tag"
)

// GetResourceSchema returns the schema definition for the tag resource
func GetResourceSchema() repository.ResourceSchema {
	return repository.ResourceSchema{
		EntityType: Tag{},
		FilterType: TagFilter{},
		UpdateType: TagUpdate{},
	}
}