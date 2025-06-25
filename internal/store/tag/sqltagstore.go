package tag

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SqlTagStore struct {
	db *pgxpool.Pool
}

type tagRow struct {
	ID                 string            `db:"id"`
	BatchID            *string           `db:"batchId"`
	Hash               string            `db:"hash"`
	Name               string            `db:"name"`
	Description        *string           `db:"description"`
	Type               string            `db:"type"`
	Context            string            `db:"context"`
	ParentTagID        *string           `db:"parentTagId"`
	ContentRating      string            `db:"contentRating"`
	ContentDescriptors []string          `db:"contentDescriptors"`
	MetaTags           []string          `db:"metaTags"`
	Public             bool              `db:"public"`
	AccessCount        int64             `db:"accessCount"`
	Metadata           json.RawMessage   `db:"metadata"`
	CreatedAt          time.Time         `db:"createdAt"`
	UpdatedAt          time.Time         `db:"updatedAt"`
	OwnerID            *string           `db:"ownerId"`
	HasQuestions       bool              `db:"hasQuestions"`
	HasChildren        bool              `db:"hasChildren"`
}

func (s *SqlTagStore) GetTagByID(ctx context.Context, id string) (*sharedpb.Tag, error) {
	var row tagRow

	err := pgxscan.Get(ctx, s.db, &row, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get tag by id: "+err.Error())
	}

	metadata := make(map[string]string)
	if row.Metadata != nil {
		// First try to unmarshal as map[string]interface{} to handle mixed types
		var rawMetadata map[string]interface{}
		if err := json.Unmarshal(row.Metadata, &rawMetadata); err == nil {
			// Convert all values to strings
			for k, v := range rawMetadata {
				switch val := v.(type) {
				case string:
					metadata[k] = val
				case float64:
					metadata[k] = fmt.Sprintf("%v", val)
				case bool:
					metadata[k] = fmt.Sprintf("%v", val)
				default:
					// For any other type, convert to JSON string
					if jsonBytes, err := json.Marshal(val); err == nil {
						metadata[k] = string(jsonBytes)
					} else {
						metadata[k] = fmt.Sprintf("%v", val)
					}
				}
			}
		} else {
			// Fallback: try direct unmarshaling to map[string]string
			if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
				return nil, status.Error(codes.Internal, "failed to unmarshal metadata: "+err.Error())
			}
		}
	}

	return &sharedpb.Tag{
		Id:                 row.ID,
		BatchId:            row.BatchID,
		Hash:               row.Hash,
		Name:               row.Name,
		Description:        row.Description,
		Type:               sharedpb.TagType(sharedpb.TagType_value[row.Type]),
		Context:            row.Context,
		ParentTagId:        row.ParentTagID,
		ContentRating:      sharedpb.ContentRating(sharedpb.ContentRating_value[row.ContentRating]),
		ContentDescriptors: row.ContentDescriptors,
		MetaTags:           row.MetaTags,
		Public:             row.Public,
		AccessCount:        int32(row.AccessCount),
		Metadata:           metadata,
		CreatedAt:          timestamppb.New(row.CreatedAt),
		UpdatedAt:          timestamppb.New(row.UpdatedAt),
		OwnerId:            row.OwnerID,
		HasQuestions:       row.HasQuestions,
		HasChildren:        row.HasChildren,
	}, nil
}

func (s *SqlTagStore) ListTagsByParent(ctx context.Context, parentID string) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE "parentTagId" = $1
	`, parentID)
	if err != nil {
		return nil, status.Error(codes.Internal, "list tags by parent: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListTagsByType(ctx context.Context, tagType sharedpb.TagType) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE type = $1
	`, tagType.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "list tags by type: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListTagsByContext(ctx context.Context, context string) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE context = $1
	`, context)
	if err != nil {
		return nil, status.Error(codes.Internal, "list tags by context: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListTagsWithFilters(ctx context.Context, params map[string]string) ([]*sharedpb.Tag, error) {
	query := `SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag" WHERE TRUE`
	args := []interface{}{}

	if tagType, ok := params["type"]; ok && tagType != "" {
		query += fmt.Sprintf(` AND type = $%d`, len(args)+1)
		args = append(args, tagType)
	}

	if contextType, ok := params["contextType"]; ok && contextType != "" {
		query += fmt.Sprintf(` AND context = $%d`, len(args)+1)
		args = append(args, contextType)
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		// Parse the boolean string
		var public bool
		if publicStr == "true" {
			public = true
		} else if publicStr == "false" {
			public = false
		} else {
			return nil, status.Error(codes.InvalidArgument, "public parameter must be 'true' or 'false'")
		}
		query += fmt.Sprintf(` AND public = $%d`, len(args)+1)
		args = append(args, public)
	}

	// Handle rootOnly parameter for filtering root tags
	if rootOnly, ok := params["rootOnly"]; ok && rootOnly == "true" {
		query += ` AND "parentTagId" IS NULL`
	}

	var rows []tagRow
	err := pgxscan.Select(ctx, s.db, &rows, query, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, "list tags with filters: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListRootTags(ctx context.Context) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM public."Tag"
		WHERE "parentTagId" IS NULL
	`)
	if err != nil {
		return nil, status.Error(codes.Internal, "list root tags: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func mapRowsToTags(rows []tagRow) []*sharedpb.Tag {
	var tags []*sharedpb.Tag
	for _, row := range rows {
		metadata := make(map[string]string)
		if row.Metadata != nil {
			// First try to unmarshal as map[string]interface{} to handle mixed types
			var rawMetadata map[string]interface{}
			if err := json.Unmarshal(row.Metadata, &rawMetadata); err == nil {
				// Convert all values to strings
				for k, v := range rawMetadata {
					switch val := v.(type) {
					case string:
						metadata[k] = val
					case float64:
						metadata[k] = fmt.Sprintf("%v", val)
					case bool:
						metadata[k] = fmt.Sprintf("%v", val)
					default:
						// For any other type, convert to JSON string
						if jsonBytes, err := json.Marshal(val); err == nil {
							metadata[k] = string(jsonBytes)
						} else {
							metadata[k] = fmt.Sprintf("%v", val)
						}
					}
				}
			} else {
				// Fallback: try direct unmarshaling to map[string]string
				if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
					// Log error but continue with empty metadata
					fmt.Printf("failed to unmarshal metadata for tag %s: %v\n", row.ID, err)
				}
			}
		}

		tags = append(tags, &sharedpb.Tag{
			Id:                 row.ID,
			BatchId:            row.BatchID,
			Hash:               row.Hash,
			Name:               row.Name,
			Description:        row.Description,
			Type:               sharedpb.TagType(sharedpb.TagType_value[row.Type]),
			Context:            row.Context,
			ParentTagId:        row.ParentTagID,
			ContentRating:      sharedpb.ContentRating(sharedpb.ContentRating_value[row.ContentRating]),
			ContentDescriptors: row.ContentDescriptors,
			MetaTags:           row.MetaTags,
			Public:             row.Public,
			AccessCount:        int32(row.AccessCount),
			Metadata:           metadata,
			CreatedAt:          timestamppb.New(row.CreatedAt),
			UpdatedAt:          timestamppb.New(row.UpdatedAt),
			OwnerId:            row.OwnerID,
			HasQuestions:       row.HasQuestions,
			HasChildren:        row.HasChildren,
		})
	}
	return tags
}

func (s *SqlTagStore) Report(ctx context.Context, tagID string, userId string, reportType sharedpb.ReportType, reason string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO "UserTagReport" ("userId", "tagId", report)
		VALUES ($1, $2, $3)
		ON CONFLICT ("userId", "tagId") 
		DO UPDATE SET report = $3, "createdAt" = CURRENT_TIMESTAMP
	`, userId, tagID, reportType)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to report tag: %v", err)
	}
	return nil
}

func (s *SqlTagStore) Favorite(ctx context.Context, tagID string, userId string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO "UserTagFavorite" ("userId", "tagId", "createdAt")
		VALUES ($1, $2, CURRENT_TIMESTAMP)
	`, userId, tagID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to favorite tag: %v", err)
	}
	return nil
}

func (s *SqlTagStore) Unfavorite(ctx context.Context, tagID string, userId string) error {
	_, err := s.db.Exec(ctx, `
		DELETE FROM "UserTagFavorite" WHERE "userId" = $1 AND "tagId" = $2
	`, userId, tagID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to unfavorite tag: %v", err)
	}
	return nil
}

func (s *SqlTagStore) CountTags(ctx context.Context, params map[string]string) (int, error) {
	query := `SELECT COUNT(*) FROM public."Tag" WHERE TRUE`
	args := []interface{}{}

	if tagType, ok := params["type"]; ok && tagType != "" {
		query += fmt.Sprintf(` AND type = $%d`, len(args)+1)
		args = append(args, tagType)
	}

	if contextType, ok := params["contextType"]; ok && contextType != "" {
		query += fmt.Sprintf(` AND context = $%d`, len(args)+1)
		args = append(args, contextType)
	}

	if publicStr, ok := params["public"]; ok && publicStr != "" {
		// Parse the boolean string
		var public bool
		if publicStr == "true" {
			public = true
		} else if publicStr == "false" {
			public = false
		} else {
			return 0, status.Error(codes.InvalidArgument, "public parameter must be 'true' or 'false'")
		}
		query += fmt.Sprintf(` AND public = $%d`, len(args)+1)
		args = append(args, public)
	}

	var count int
	if err := s.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *SqlTagStore) UniqueTagTypes(ctx context.Context) ([]sharedpb.TagType, error) {
	var typeStrings []string

	err := pgxscan.Select(ctx, s.db, &typeStrings, `
		SELECT DISTINCT type 
		FROM public."Tag" 
		WHERE type IS NOT NULL 
		ORDER BY type
	`)
	if err != nil {
		return nil, status.Error(codes.Internal, "get unique tag types: "+err.Error())
	}

	var tagTypes []sharedpb.TagType
	for _, typeStr := range typeStrings {
		if tagType, exists := sharedpb.TagType_value[typeStr]; exists {
			tagTypes = append(tagTypes, sharedpb.TagType(tagType))
		}
	}

	return tagTypes, nil
}
