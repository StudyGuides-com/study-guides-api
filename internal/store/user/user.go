package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &SqlUserStore{db: db}, nil
}