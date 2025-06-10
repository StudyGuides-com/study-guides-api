package tag

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type SqlTagStore struct {
	db *pgxpool.Pool
}

type tagRow struct {
	ID            string `db:"id"`
	Name          string `db:"name"`
	Type          string `db:"type"`
	Context       string `db:"context"`
	ContentRating string `db:"content_rating"`
	Public        bool   `db:"public"`
	AccessCount   int64  `db:"access_count"`
}

func (s *SqlTagStore) GetTagByID(ctx context.Context, id string) (*sharedpb.Tag, error) {
	var row tagRow

	err := pgxscan.Get(ctx, s.db, &row, `
		SELECT id, name, type, context, content_rating, public, access_count
		FROM tags
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get tag by id: %w", err)
	}

	return &sharedpb.Tag{
		Id:            row.ID,
		Name:          row.Name,
		Type:          sharedpb.TagType(sharedpb.TagType_value[row.Type]),
		Context:       row.Context,
		ContentRating: sharedpb.ContentRating(sharedpb.ContentRating_value[row.ContentRating]),
		Public:        row.Public,
		AccessCount:   int32(row.AccessCount),
	}, nil
}

func (s *SqlTagStore) ListTagsByParent(ctx context.Context, parentID string) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, name, type, context, content_rating, public, access_count
		FROM tags
		WHERE parent_id = $1
	`, parentID)
	if err != nil {
		return nil, fmt.Errorf("list tags by parent: %w", err)
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListTagsByType(ctx context.Context, tagType sharedpb.TagType) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, name, type, context, content_rating, public, access_count
		FROM tags
		WHERE type = $1
	`, tagType.String())
	if err != nil {
		return nil, fmt.Errorf("list tags by type: %w", err)
	}

	return mapRowsToTags(rows), nil
}

func (s *SqlTagStore) ListRootTags(ctx context.Context) ([]*sharedpb.Tag, error) {
	var rows []tagRow

	err := pgxscan.Select(ctx, s.db, &rows, `
		SELECT id, name, type, context, content_rating, public, access_count
		FROM tags
		WHERE parent_id IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("list root tags: %w", err)
	}

	return mapRowsToTags(rows), nil
}

func mapRowsToTags(rows []tagRow) []*sharedpb.Tag {
	var tags []*sharedpb.Tag
	for _, row := range rows {
		tags = append(tags, &sharedpb.Tag{
			Id:            row.ID,
			Name:          row.Name,
			Type:          sharedpb.TagType(sharedpb.TagType_value[row.Type]),
			Context:       row.Context,
			ContentRating: sharedpb.ContentRating(sharedpb.ContentRating_value[row.ContentRating]),
			Public:        row.Public,
			AccessCount:   int32(row.AccessCount),
		})
	}
	return tags
}
