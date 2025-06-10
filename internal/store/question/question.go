package question

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sharedpb "github.com/studyguides-com/study-guides-api/api/v1/shared"
)

type QuestionStore interface {
	GetQuestionsByTagID(ctx context.Context, tagID string) ([]*sharedpb.Question, error)
	Report(ctx context.Context, questionID string, userId string, reportType sharedpb.ReportType, reason string) error
}

func NewSqlQuestionStore(ctx context.Context, dbURL string) (*SqlQuestionStore, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to connect to postgres: "+err.Error())
	}
	return &SqlQuestionStore{db: db}, nil
}
