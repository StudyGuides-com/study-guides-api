package services

import (
	"context"
	"encoding/json"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"github.com/studyguides-com/study-guides-api/internal/errors"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/lib/router"
	"github.com/studyguides-com/study-guides-api/internal/lib/tools"

	"github.com/sashabaranov/go-openai"
)

const ToolChoiceTypeAuto = "auto"

type ChatService struct {
	chatpb.UnimplementedChatServiceServer
	router router.Router
	ai     ai.AiClient
}

func NewChatService(router router.Router, ai ai.AiClient) *ChatService {
	return &ChatService{
		router: router,
		ai:     ai,
	}
}

func (s *ChatService) Chat(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context) (interface{}, error) {

		systemPrompt := `
	You are an intent router. Allowed operations: GetTagCount.
	For GetTagCount, allowed parameters: type, contextType.
	If none apply, call Unknown.
	Always pick exactly one.
	Please respond using the provided tool to return your response in JSON format.
	`

		tools := tools.GetClassificationTools()

		raw, err := s.ai.ChatCompletionWithTools(ctx, systemPrompt, req.Message, tools, &openai.ToolChoice{
			Type: ToolChoiceTypeAuto,
		})
		if err != nil {
			return nil, err
		}

		var plan struct {
			Operation  string
			Parameters map[string]string
		}

		// parse from tool call JSON:
		var chatResp openai.ChatCompletionResponse
		if err := json.Unmarshal([]byte(raw), &chatResp); err != nil {
			return nil, err
		}
		if len(chatResp.Choices) == 0 || len(chatResp.Choices[0].Message.ToolCalls) == 0 {
			return nil, errors.ErrToolNotFound
		}
		toolCall := chatResp.Choices[0].Message.ToolCalls[0]
		plan.Operation = toolCall.Function.Name
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &plan.Parameters); err != nil {
			return nil, err
		}

		answer, err := s.router.Route(ctx, plan.Operation, plan.Parameters)
		if err != nil {
			return nil, err
		}

		return &chatpb.ChatResponse{
			Context:    req.Context,
			Operation:  plan.Operation,
			Parameters: plan.Parameters,
			Answer:     answer,
			PlanJson:   mustJson(plan),
		}, nil
	})
	return resp.(*chatpb.ChatResponse), err
}

func mustJson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
