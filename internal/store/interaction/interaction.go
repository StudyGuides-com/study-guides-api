package interaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

)

type InteractionStore interface {
	Interact(ctx context.Context) error
}



func NewSqlInteractionStore(ctx context.Context, dbURL string) (*SqlInteractionStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlInteractionStore{db: db}, nil
}

