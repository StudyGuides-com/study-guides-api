package tag

import (
	"context"
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
	Metadata           map[string]string `db:"metadata"`
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
		FROM "Tag"
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, status.Error(codes.Internal, "get tag by id: "+err.Error())
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
		Metadata:           row.Metadata,
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
		FROM "Tag"
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
		FROM "Tag"
		WHERE type = $1
	`, tagType.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "list tags by type: "+err.Error())
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListRootTags(ctx context.Context) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, "batchId", hash, name, description, type, context, "parentTagId",
		       "contentRating", "contentDescriptors", "metaTags", public, "accessCount",
		       metadata, "createdAt", "updatedAt", "ownerId", "hasQuestions", "hasChildren"
		FROM "Tag"
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
			Metadata:           row.Metadata,
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
