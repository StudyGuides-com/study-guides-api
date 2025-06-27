package roland

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type RolandStore interface {
	HealthCheck(ctx context.Context) error
	Bundles(ctx context.Context) ([]*sharedpb.Bundle, error)
	BundlesByParserType(ctx context.Context, parserType sharedpb.ParserType) ([]*sharedpb.Bundle, error)
	SaveBundle(ctx context.Context, bundle *sharedpb.Bundle, force bool) (bool, error)
	UpdateGob(ctx context.Context, id string, gobPayload []byte, force bool) (bool, error)
	DeleteAllBundles(ctx context.Context) error
	DeleteBundleByID(ctx context.Context, id string) error
	DeleteBundlesByShortID(ctx context.Context, prefix string) (int, error)
	MarkBundleExported(ctx context.Context, id string, exportType sharedpb.ExportType) (bool, error)
}

func NewSqlRolandStore(ctx context.Context, dbURL string) (*SqlRolandStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlRolandStore{db: db}, nil
}