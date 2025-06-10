package interaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	interactionpb "github.com/studyguides-com/study-guides-api/api/v1/interaction"
)

type InteractionStore interface {
	AnswerCorrectly(ctx context.Context, req *interactionpb.InteractRequest) error
	AnswerIncorrectly(ctx context.Context, req *interactionpb.InteractRequest) error
	AnswerEasy(ctx context.Context, req *interactionpb.InteractRequest) error
	AnswerHard(ctx context.Context, req *interactionpb.InteractRequest) error
	Reveal(ctx context.Context, req *interactionpb.InteractRequest) error
	ViewLearnMore(ctx context.Context, req *interactionpb.InteractRequest) error
	ViewPassage(ctx context.Context, req *interactionpb.InteractRequest) error
}



func NewSqlInteractionStore(ctx context.Context, dbURL string) (*SqlInteractionStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlInteractionStore{db: db}, nil
}

