package services

import (
	"context"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"github.com/studyguides-com/study-guides-api/internal/lib/router"
)

type ChatService struct {
	chatpb.UnimplementedChatServiceServer
	router router.Router
}

func NewChatService(router router.Router) *ChatService {
	return &ChatService{
		router: router,
	}
}

func (s *ChatService) Chat(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context) (interface{}, error) {
		// 👇 In real flow: you'd call LLM here to classify
    // For now, hardcode a fake plan:
    planOperation := "GetTagCount"
    planParameters := map[string]string{
      "type": "Course",
      "contextType": "College",
    }

	answer, err := s.router.Route(ctx, planOperation, planParameters)
    if err != nil {
      return nil, err
    }
		if err != nil {
			return nil, err
		}
		return &chatpb.ChatResponse{
			Context: req.Context,
			Operation: planOperation,
			Parameters: planParameters,
			Answer: answer,
			PlanJson: "{}",
		}, nil
	})
	return resp.(*chatpb.ChatResponse), err
}
