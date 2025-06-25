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

const (
	ToolChoiceTypeAuto = "auto"
	// Conversation management constants
	MaxConversationMessages = 10  // Maximum number of messages to keep in history
	MaxMessageLength        = 1000 // Maximum characters per message to store
)

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

// createResponseSummary creates a summary of the response for conversation history
func createResponseSummary(answer string, operation string, params map[string]string) string {
	// If the response is short, keep it as is
	if len(answer) <= MaxMessageLength {
		return answer
	}
	
	// For long responses, create a summary based on the operation
	switch operation {
	case "ListTags":
		// Count lines to estimate number of items
		lines := strings.Split(answer, "\n")
		itemCount := 0
		for _, line := range lines {
			if strings.Contains(line, ". ") || strings.Contains(line, "|") || strings.Contains(line, "{") {
				itemCount++
			}
		}
		
		format := "list"
		if f, ok := params["format"]; ok {
			format = f
		}
		
		tagType := ""
		if t, ok := params["type"]; ok {
			tagType = t
		}
		
		if tagType != "" {
			return fmt.Sprintf("Retrieved %d %s tags in %s format", itemCount, tagType, format)
		} else {
			return fmt.Sprintf("Retrieved %d root tags in %s format", itemCount, format)
		}
		
	case "TagCount":
		return "Retrieved tag count information"
		
	case "UniqueTagTypes":
		// Count the number of tag types
		lines := strings.Split(answer, "\n")
		typeCount := 0
		for _, line := range lines {
			if strings.Contains(line, ". ") {
				typeCount++
			}
		}
		return fmt.Sprintf("Retrieved %d unique tag types", typeCount)
		
	default:
		// For unknown operations, truncate and add ellipsis
		if len(answer) > MaxMessageLength {
			return answer[:MaxMessageLength] + "..."
		}
		return answer
	}
}

// AddMessage adds a message to the conversation history with management
func (ch *ConversationHistory) AddMessage(message openai.ChatCompletionMessage) {
	// Truncate message content if it's too long
	if len(message.Content) > MaxMessageLength {
		message.Content = message.Content[:MaxMessageLength] + "..."
	}
	
	ch.Messages = append(ch.Messages, message)
	
	// Keep only the most recent messages
	if len(ch.Messages) > MaxConversationMessages {
		// Remove oldest messages, keeping the most recent ones
		ch.Messages = ch.Messages[len(ch.Messages)-MaxConversationMessages:]
	}
}

// AddAssistantResponse adds an assistant response with smart summarization
func (ch *ConversationHistory) AddAssistantResponse(answer string, operation string, params map[string]string) {
	// Create a summary for conversation history
	summary := createResponseSummary(answer, operation, params)
	
	assistantMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: summary,
	}
	
	ch.AddMessage(assistantMessage)
}

// TruncateHistory removes old messages to prevent token limit issues
func (ch *ConversationHistory) TruncateHistory() {
	if len(ch.Messages) > MaxConversationMessages {
		// Keep only the most recent messages
		ch.Messages = ch.Messages[len(ch.Messages)-MaxConversationMessages:]
	}
}

// GetMessagesForAI returns messages optimized for AI consumption
func (ch *ConversationHistory) GetMessagesForAI() []openai.ChatCompletionMessage {
	// Create a copy and truncate if needed
	messages := make([]openai.ChatCompletionMessage, len(ch.Messages))
	copy(messages, ch.Messages)
	
	// If we have too many messages, keep only the most recent ones
	if len(messages) > MaxConversationMessages {
		messages = messages[len(messages)-MaxConversationMessages:]
	}
	
	return messages
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
	
	// Add format selection guidance
	formatGuidance := `
	When using ListTags or TagCount, choose the appropriate format based on user intent:
	- Use 'list' (default) for human-readable responses, general browsing, or when user asks to "show", "list", or "see" tags
	- Use 'json' when user mentions "export", "data", "API", "programmatic", "download", or machine-readable needs
	- Use 'csv' when user mentions "spreadsheet", "Excel", "analysis", "compare", or data import needs
	- Use 'table' when user specifically asks for table format or markdown
	- If unsure, default to 'list' for human interaction
	`
	
	return fmt.Sprintf(`
	You are an intent router. Allowed operations: %s.
	%s
	%s
	%s
	If none apply, call Unknown.
	Always pick exactly one.
	Please respond using the provided tool to return your response in JSON format.
	`, operationsList, detailsList, tagTypeGuidance, formatGuidance)
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

		// Use optimized messages for AI to prevent token limit issues
		aiMessages := conversationHistory.GetMessagesForAI()
		raw, err := s.ai.ChatCompletionWithHistory(ctx, systemPrompt, aiMessages, tools, nil)
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
		conversationHistory.AddAssistantResponse(answer, plan.Operation, plan.Parameters)
		
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
			Answer:  answer,
			Context: updatedContext,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return resp.(*chatpb.ChatResponse), nil
}

func mustJson(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
