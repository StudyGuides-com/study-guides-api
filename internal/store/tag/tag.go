package tag

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

	type tagStore struct {
}

type TagStore interface {
	GetTagByID(ctx context.Context, id string) (*sharedpb.Tag, error)
	ListTagsByParent(ctx context.Context, parentID string) ([]*sharedpb.Tag, error)
	ListTagsByType(ctx context.Context, tagType sharedpb.TagType) ([]*sharedpb.Tag, error)
	ListRootTags(ctx context.Context) ([]*sharedpb.Tag, error)
}



func NewSqlTagStore(ctx context.Context, dbURL string) (*SqlTagStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &SqlTagStore{db: db}, nil
}