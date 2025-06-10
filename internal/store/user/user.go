package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type UserStore interface {
	UserByID(ctx context.Context, userID string) (*sharedpb.User, error)
	UserByEmail(ctx context.Context, email string) (*sharedpb.User, error)
	Profile(ctx context.Context, userID string) (*sharedpb.User, error)
}

func NewSqlUserStore(ctx context.Context, dbURL string) (*SqlUserStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlUserStore{db: db}, nil
}
