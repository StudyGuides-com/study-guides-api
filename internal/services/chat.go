package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"github.com/studyguides-com/study-guides-api/internal/errors"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/lib/router"
	"github.com/studyguides-com/study-guides-api/internal/lib/tools"

	"github.com/sashabaranov/go-openai"
)

const ToolChoiceTypeAuto = "auto"

// ConversationHistory manages conversation history using Context metadata
type ConversationHistory struct {
	Messages []openai.ChatCompletionMessage
}

// NewConversationHistory creates a new conversation history
func NewConversationHistory() *ConversationHistory {
	return &ConversationHistory{
		Messages: []openai.ChatCompletionMessage{},
	}
}

// FromContextMetadata loads conversation history from context metadata
func FromContextMetadata(metadata map[string]string) *ConversationHistory {
	history := NewConversationHistory()
	
	if historyJSON, exists := metadata["conversation_history"]; exists && historyJSON != "" {
		var messages []openai.ChatCompletionMessage
		if err := json.Unmarshal([]byte(historyJSON), &messages); err == nil {
			history.Messages = messages
		}
	}
	
	return history
}

// ToContextMetadata saves conversation history to context metadata
func (ch *ConversationHistory) ToContextMetadata() map[string]string {
	metadata := make(map[string]string)
	
	if len(ch.Messages) > 0 {
		if historyJSON, err := json.Marshal(ch.Messages); err == nil {
			metadata["conversation_history"] = string(historyJSON)
		}
	}
	
	return metadata
}

// AddMessage adds a message to the conversation history
func (ch *ConversationHistory) AddMessage(message openai.ChatCompletionMessage) {
	ch.Messages = append(ch.Messages, message)
}

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

// buildSystemPrompt dynamically creates a system prompt based on available tools
func buildSystemPrompt() string {
	toolDefinitions := tools.GetClassificationDefinitions()
	
	var operations []string
	var operationDetails []string
	
	for _, toolDef := range toolDefinitions {
		operationName := toolDef.Name
		operations = append(operations, operationName)
		
		// Build parameter list for this operation
		var params []string
		for _, param := range toolDef.Parameters {
			params = append(params, param.Name)
		}
		
		if len(params) > 0 {
			operationDetails = append(operationDetails, fmt.Sprintf("For %s, allowed parameters: %s.", operationName, strings.Join(params, ", ")))
		} else {
			operationDetails = append(operationDetails, fmt.Sprintf("For %s, no parameters required.", operationName))
		}
	}
	
	operationsList := strings.Join(operations, ", ")
	detailsList := strings.Join(operationDetails, "\n")
	
	// Add tag type guidance
	tagTypeGuidance := `
	When using ListTags with a type parameter, use the exact tag types that exist in the system.
	Common tag types include: Category, UserContent, UserTopic, Branch, Instruction_Type, 
	Instruction_Group, Instruction, Chapter, Section, etc.
	
	Use the exact type name as it appears in the system. Do NOT use synonyms or variations.
	`
	
	return fmt.Sprintf(`
	You are an intent router. Allowed operations: %s.
	%s
	%s
	If none apply, call Unknown.
	Always pick exactly one.
	Please respond using the provided tool to return your response in JSON format.
	`, operationsList, detailsList, tagTypeGuidance)
}

func (s *ChatService) Chat(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	// Debug: Print incoming request details
	fmt.Printf("DEBUG: === CHAT REQUEST ===\n")
	fmt.Printf("DEBUG: User ID: %s\n", req.Context.GetUserId())
	fmt.Printf("DEBUG: Session ID: %s\n", req.Context.GetSessionId())
	fmt.Printf("DEBUG: Message: %s\n", req.Message)
	fmt.Printf("DEBUG: Context Metadata: %+v\n", req.Context.GetMetadata())
	
	resp, err := PublicBaseHandler(ctx, func(ctx context.Context) (interface{}, error) {

		systemPrompt := buildSystemPrompt()

		// Add current user message to history
		currentUserMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Message,
		}
		// Get conversation history from context metadata
		conversationHistory := FromContextMetadata(req.Context.GetMetadata())
		
		// Debug: Print conversation history
		fmt.Printf("DEBUG: === CONVERSATION HISTORY ===\n")
		fmt.Printf("DEBUG: History length: %d messages\n", len(conversationHistory.Messages))
		for i, msg := range conversationHistory.Messages {
			fmt.Printf("DEBUG: Message %d - Role: %s, Content: %s\n", i+1, msg.Role, msg.Content)
		}
		
		conversationHistory.AddMessage(currentUserMessage)
		
		// Debug: Print updated history
		fmt.Printf("DEBUG: === UPDATED HISTORY ===\n")
		fmt.Printf("DEBUG: Updated history length: %d messages\n", len(conversationHistory.Messages))
		for i, msg := range conversationHistory.Messages {
			fmt.Printf("DEBUG: Message %d - Role: %s, Content: %s\n", i+1, msg.Role, msg.Content)
		}

		tools := tools.GetClassificationTools()

		raw, err := s.ai.ChatCompletionWithHistory(ctx, systemPrompt, conversationHistory.Messages, tools, nil)
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
		
		// Debug: Print AI response and plan
		fmt.Printf("DEBUG: === AI RESPONSE ===\n")
		fmt.Printf("DEBUG: Operation: %s\n", plan.Operation)
		fmt.Printf("DEBUG: Parameters: %+v\n", plan.Parameters)

		answer, err := s.router.Route(ctx, plan.Operation, plan.Parameters)
		if err != nil {
			return nil, err
		}
		
		// Debug: Print router answer
		fmt.Printf("DEBUG: === ROUTER ANSWER ===\n")
		fmt.Printf("DEBUG: Answer: %s\n", answer)

		// Add assistant response to conversation history
		assistantMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: answer,
		}
		conversationHistory.AddMessage(assistantMessage)
		
		// Debug: Print final history
		fmt.Printf("DEBUG: === FINAL HISTORY ===\n")
		fmt.Printf("DEBUG: Final history length: %d messages\n", len(conversationHistory.Messages))
		for i, msg := range conversationHistory.Messages {
			fmt.Printf("DEBUG: Message %d - Role: %s, Content: %s\n", i+1, msg.Role, msg.Content)
		}

		// Create updated context with conversation history
		updatedContext := &chatpb.Context{
			UserId:    req.Context.GetUserId(),
			SessionId: req.Context.GetSessionId(),
			Metadata:  conversationHistory.ToContextMetadata(),
		}
		
		// Debug: Print updated context
		fmt.Printf("DEBUG: === UPDATED CONTEXT ===\n")
		fmt.Printf("DEBUG: Updated Context Metadata: %+v\n", updatedContext.GetMetadata())

		return &chatpb.ChatResponse{
			Context:    updatedContext,
			Operation:  plan.Operation,
			Parameters: plan.Parameters,
			Answer:     answer,
			PlanJson:   mustJson(plan),
		}, nil
	})
	
	// Check if resp is nil before type assertion
	if resp == nil {
		return nil, err
	}
	
	return resp.(*chatpb.ChatResponse), err
}

func mustJson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
