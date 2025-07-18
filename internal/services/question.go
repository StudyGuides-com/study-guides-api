package services

import (
	"context"

	questionpb "github.com/studyguides-com/study-guides-api/api/v1/question"
	"github.com/studyguides-com/study-guides-api/internal/middleware"
	"github.com/studyguides-com/study-guides-api/internal/store"
)

type QuestionService struct {
	questionpb.UnimplementedQuestionServiceServer
	store store.Store
}

func NewQuestionService(store store.Store) *QuestionService {
	return &QuestionService{
		store: store,
	}
}

func (s *QuestionService) ForTag(ctx context.Context, req *questionpb.ForTagRequest) (*questionpb.QuestionsResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		questions, err := s.store.QuestionStore().GetQuestionsByTagID(ctx, req.TagId)
		if err != nil {
			return nil, err
		}
		return &questionpb.QuestionsResponse{
			Questions: questions,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*questionpb.QuestionsResponse), nil
}

func (s *QuestionService) Report(ctx context.Context, req *questionpb.ReportQuestionRequest) (*questionpb.ReportQuestionResponse, error) {
	resp, err := AuthBaseHandler(ctx, func(ctx context.Context, session *middleware.SessionDetails) (interface{}, error) {
		return &questionpb.ReportQuestionResponse{
			Success: true,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return resp.(*questionpb.ReportQuestionResponse), nil
}
