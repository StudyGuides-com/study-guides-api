package indexing

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

// AlgoliaTagRecord represents the exact structure for Algolia indexing
type AlgoliaTagRecord struct {
	ObjectID                  string    `json:"objectID"`
	ID                        string    `json:"id"`
	Name                      string    `json:"name"`
	Description               *string   `json:"description,omitempty"`
	HasQuestions              bool      `json:"hasQuestions"`
	Type                      string    `json:"type"`
	Context                   string    `json:"context"`
	Public                    bool      `json:"public"`
	BatchID                   *string   `json:"batchId,omitempty"`
	MetaTags                  []string  `json:"metaTags"`
	ContentRating             string    `json:"contentRating"`
	ContentDescriptors        []string  `json:"contentDescriptors"`
	MissingContentDescriptors bool      `json:"missingContentDescriptors"`
	MissingContentRating      bool      `json:"missingContentRating"`
	MissingMetaTags           bool      `json:"missingMetaTags"`
	OwnerID                   *string   `json:"ownerId,omitempty"`
	AccessList                []string  `json:"accessList"`
	Tags                      []TagInfo `json:"tags"`
}

// TagInfo represents tag ancestry information
type TagInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	ParentTagID  *string `json:"parentTagId,omitempty"`
	HasQuestions bool    `json:"hasQuestions"`
	HasChildren  bool    `json:"hasChildren"`
}

// processTagOperation handles indexing operations for tags
func (s *SqlIndexingStore) processTagOperation(ctx context.Context, op IndexOperation, force bool) error {
	if op.Action == "delete" {
		// Delete from Algolia
		_, err := s.tagIndex.DeleteObject(op.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to delete from Algolia: %w", err)
		}
		return nil
	}
	
	// Get tag data
	tag, err := s.tagStore.GetTagByID(ctx, op.ObjectID)
	if err != nil {
		return fmt.Errorf("failed to get tag: %w", err)
	}
	
	// Get access list
	accessList, err := s.getTagAccessList(ctx, op.ObjectID)
	if err != nil {
		return fmt.Errorf("failed to get access list: %w", err)
	}
	
	// Get ancestry
	ancestry, err := s.getTagAncestry(ctx, op.ObjectID)
	if err != nil {
		return fmt.Errorf("failed to get ancestry: %w", err)
	}
	
	// Transform to Algolia record
	record := transformTagToAlgoliaRecord(tag, accessList, ancestry)
	
	// Compute hash of the record
	recordJSON, _ := json.Marshal(record)
	hash := sha256.Sum256(recordJSON)
	
	if !force {
		// Check if content has changed
		state, _ := s.GetIndexState(ctx, "Tag", op.ObjectID)
		if state != nil && state.LastIndexedHash != nil {
			if bytes.Equal(state.LastIndexedHash, hash[:]) {
				// Content hasn't changed, skip indexing
				return nil
			}
		}
	}
	
	// Push to Algolia
	_, err = s.tagIndex.SaveObject(record)
	if err != nil {
		return fmt.Errorf("failed to save to Algolia: %w", err)
	}
	
	// Update state
	if err := s.UpdateIndexState(ctx, "Tag", op.ObjectID, hash[:]); err != nil {
		return fmt.Errorf("failed to update index state: %w", err)
	}
	
	return nil
}

// getTagAccessList retrieves the list of users with access to a tag
func (s *SqlIndexingStore) getTagAccessList(ctx context.Context, tagID string) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT "userId" FROM "TagAccess" WHERE "tagId" = $1
	`, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tag access: %w", err)
	}
	defer rows.Close()
	
	var accessList []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("failed to scan user ID: %w", err)
		}
		accessList = append(accessList, userID)
	}
	
	return accessList, nil
}

// tagAncestryRow represents a row from the ancestry query
type tagAncestryRow struct {
	ID           string  `db:"id"`
	Name         string  `db:"name"`
	Type         string  `db:"type"`
	ParentTagID  *string `db:"parentTagId"`
	HasQuestions bool    `db:"hasQuestions"`
	HasChildren  bool    `db:"hasChildren"`
}

// getTagAncestry retrieves the ancestry chain for a tag
func (s *SqlIndexingStore) getTagAncestry(ctx context.Context, tagID string) ([]TagInfo, error) {
	// Use recursive CTE to climb the ancestry tree
	query := `
		WITH RECURSIVE ancestry AS (
			SELECT id, name, type, "parentTagId", "hasQuestions", "hasChildren"
			FROM "Tag" 
			WHERE id = $1
			
			UNION ALL
			
			SELECT t.id, t.name, t.type, t."parentTagId", t."hasQuestions", t."hasChildren"
			FROM "Tag" t
			INNER JOIN ancestry a ON t.id = a."parentTagId"
		)
		SELECT id, name, type, "parentTagId", "hasQuestions", "hasChildren"
		FROM ancestry
		ORDER BY id
	`
	
	var rows []tagAncestryRow
	err := pgxscan.Select(ctx, s.pool, &rows, query, tagID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tag ancestry: %w", err)
	}
	
	// Convert to TagInfo slice
	var ancestry []TagInfo
	for _, row := range rows {
		ancestry = append(ancestry, TagInfo{
			ID:           row.ID,
			Name:         row.Name,
			Type:         row.Type,
			ParentTagID:  row.ParentTagID,
			HasQuestions: row.HasQuestions,
			HasChildren:  row.HasChildren,
		})
	}
	
	return ancestry, nil
}

// transformTagToAlgoliaRecord transforms a tag into the Algolia record format
func transformTagToAlgoliaRecord(tag *sharedpb.Tag, accessList []string, ancestry []TagInfo) AlgoliaTagRecord {
	// Handle nil/empty slices
	if tag.MetaTags == nil {
		tag.MetaTags = []string{}
	}
	if tag.ContentDescriptors == nil {
		tag.ContentDescriptors = []sharedpb.ContentDescriptorType{}
	}
	if accessList == nil {
		accessList = []string{}
	}
	if ancestry == nil {
		ancestry = []TagInfo{}
	}
	
	// Convert content descriptors to strings
	contentDescriptors := make([]string, len(tag.ContentDescriptors))
	for i, cd := range tag.ContentDescriptors {
		contentDescriptors[i] = cd.String()
	}
	
	return AlgoliaTagRecord{
		ObjectID:                  tag.Id,
		ID:                        tag.Id,
		Name:                      tag.Name,
		Description:               tag.Description,
		HasQuestions:              tag.HasQuestions,
		Type:                      tag.Type.String(),
		Context:                   tag.Context.String(),
		Public:                    tag.Public,
		BatchID:                   tag.BatchId,
		MetaTags:                  tag.MetaTags,
		ContentRating:             tag.ContentRating.String(),
		ContentDescriptors:        contentDescriptors,
		MissingContentDescriptors: len(contentDescriptors) == 0,
		MissingContentRating:      tag.ContentRating == sharedpb.ContentRating_RatingPending,
		MissingMetaTags:           len(tag.MetaTags) == 0,
		OwnerID:                   tag.OwnerId,
		AccessList:                accessList,
		Tags:                      ancestry,
	}
}