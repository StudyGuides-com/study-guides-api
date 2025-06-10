package interaction

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	interactionpb "github.com/studyguides-com/study-guides-api/api/v1/interaction"
	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type InteractionStore interface {
	// Methods that modify the question return the updated question
	AnswerCorrectly(ctx context.Context, req *interactionpb.InteractRequest) (*sharedpb.Question, error)
	AnswerIncorrectly(ctx context.Context, req *interactionpb.InteractRequest) (*sharedpb.Question, error)
	AnswerEasy(ctx context.Context, req *interactionpb.InteractRequest) (*sharedpb.Question, error)
	AnswerHard(ctx context.Context, req *interactionpb.InteractRequest) (*sharedpb.Question, error)
	
	// View-only methods don't return a question
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

