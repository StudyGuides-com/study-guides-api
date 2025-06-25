package tag

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type TagStore interface {
	GetTagByID(ctx context.Context, id string) (*sharedpb.Tag, error)
	ListTagsByParent(ctx context.Context, parentID string) ([]*sharedpb.Tag, error)
	ListTagsByType(ctx context.Context, tagType sharedpb.TagType) ([]*sharedpb.Tag, error)
	ListTagsByContext(ctx context.Context, context string) ([]*sharedpb.Tag, error)
	ListTagsWithFilters(ctx context.Context, params map[string]string) ([]*sharedpb.Tag, error)
	ListRootTags(ctx context.Context) ([]*sharedpb.Tag, error)
	UniqueTagTypes(ctx context.Context) ([]sharedpb.TagType, error)
	UniqueContextTypes(ctx context.Context) ([]string, error)
	Report(ctx context.Context, tagID string, userId string, reportType sharedpb.ReportType, reason string) error
	Favorite(ctx context.Context, tagID string, userId string) error
	Unfavorite(ctx context.Context, tagID string, userId string) error
	CountTags(ctx context.Context, params map[string]string) (int, error)
}

func NewSqlTagStore(ctx context.Context, dbURL string) (*SqlTagStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlTagStore{db: db}, nil
}
