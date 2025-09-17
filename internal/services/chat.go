package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	chatpb "github.com/studyguides-com/study-guides-api/api/v1/chat"
	"github.com/studyguides-com/study-guides-api/internal/lib/ai"
	"github.com/studyguides-com/study-guides-api/internal/mcp"
	"github.com/studyguides-com/study-guides-api/internal/mcp/indexing"
	"github.com/studyguides-com/study-guides-api/internal/mcp/kpi"
	"github.com/studyguides-com/study-guides-api/internal/mcp/tag"
	indexingcore "github.com/studyguides-com/study-guides-api/internal/core/indexing"
	"github.com/studyguides-com/study-guides-api/internal/store"

	"github.com/sashabaranov/go-openai"
)

const (
	ToolChoiceTypeAuto = "auto"
	// Conversation management constants
	MaxConversationMessages = 10   // Maximum number of messages to keep in history
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

	case "Deploy":
		return "Deployment initiated successfully"

	case "Rollback":
		return "Rollback initiated successfully"

	case "ListDeployments":
		// Count lines to estimate number of deployments
		lines := strings.Split(answer, "\n")
		deploymentCount := 0
		for _, line := range lines {
			if strings.Contains(line, ". ") || strings.Contains(line, "|") || strings.Contains(line, "{") {
				deploymentCount++
			}
		}

		format := "list"
		if f, ok := params["format"]; ok {
			format = f
		}

		return fmt.Sprintf("Retrieved %d deployments in %s format", deploymentCount, format)

	case "GetDeploymentStatus":
		return "Retrieved deployment status information"

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
	mcpProcessor *mcp.MCPProcessor
	store        store.Store
}

func NewChatService(store store.Store, ai ai.AiClient) *ChatService {
	// Create MCP processor
	mcpProcessor := mcp.NewMCPProcessor(ai)
	
	// Register tag repository
	tagRepo := tag.NewTagRepositoryAdapter(store.TagStore())
	mcpProcessor.Register(tag.ResourceName, tagRepo, tag.GetResourceSchema())
	
	// Register KPI repository
	kpiRepo := kpi.NewKPIRepositoryAdapter(store.KPIStore())
	mcpProcessor.Register(kpi.ResourceName, kpiRepo, kpi.GetResourceSchema())
	
	// Register indexing repository
	indexingBusiness := indexingcore.NewBusinessService(store)
	indexingRepo := indexing.NewIndexingRepositoryAdapter(indexingBusiness)
	mcpProcessor.Register(indexing.ResourceName, indexingRepo, indexing.GetResourceSchema())
	
	return &ChatService{
		mcpProcessor: mcpProcessor,
		store:        store,
	}
}

// convertMCPResponseToLegacyFormat converts MCP response to match the old API format
func convertMCPResponseToLegacyFormat(mcpResp *mcp.Response) (string, string, map[string]string) {
	// Extract operation and parameters from MCP response
	operation := "MCP_Operation"
	parameters := make(map[string]string)
	
	// The MCP response contains the formatted data we need
	answer := mcpResp.Message
	if !mcpResp.Success {
		answer = fmt.Sprintf("Error: %s", mcpResp.Error)
	}
	
	// If we have data, include it in the answer
	if mcpResp.Data != nil {
		if _, err := json.Marshal(mcpResp.Data); err == nil {
			// For backwards compatibility, we format the data nicely
			answer = mcpResp.Message
		}
	}
	
	return answer, operation, parameters
}

func (s *ChatService) Chat(ctx context.Context, req *chatpb.ChatRequest) (*chatpb.ChatResponse, error) {
	// Debug: Print incoming request details
	fmt.Printf("DEBUG: === MCP CHAT REQUEST ===\n")
	fmt.Printf("DEBUG: User ID: %s\n", req.Context.GetUserId())
	fmt.Printf("DEBUG: Session ID: %s\n", req.Context.GetSessionId())
	fmt.Printf("DEBUG: Message: %s\n", req.Message)
	fmt.Printf("DEBUG: Context Metadata: %+v\n", req.Context.GetMetadata())

	resp, err := PublicBaseHandler(ctx, func(ctx context.Context) (interface{}, error) {
		// Get conversation history from context metadata
		conversationHistory := FromContextMetadata(req.Context.GetMetadata())

		// Add current user message to history
		currentUserMessage := openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Message,
		}
		conversationHistory.AddMessage(currentUserMessage)

		// Debug: Print conversation history
		fmt.Printf("DEBUG: === CONVERSATION HISTORY ===\n")
		fmt.Printf("DEBUG: History length: %d messages\n", len(conversationHistory.Messages))
		for i, msg := range conversationHistory.Messages {
			fmt.Printf("DEBUG: Message %d - Role: %s, Content: %s\n", i+1, msg.Role, msg.Content)
		}

		// Process request through MCP system
		fmt.Printf("DEBUG: === MCP PROCESSING ===\n")
		mcpResponse, err := s.mcpProcessor.ProcessRequest(ctx, req.Message)
		if err != nil {
			fmt.Printf("DEBUG: MCP processing failed: %v\n", err)
			return nil, fmt.Errorf("failed to process request: %w", err)
		}

		// Debug: Print MCP response
		fmt.Printf("DEBUG: === MCP RESPONSE ===\n")
		fmt.Printf("DEBUG: Success: %v\n", mcpResponse.Success)
		fmt.Printf("DEBUG: Message: %s\n", mcpResponse.Message)
		fmt.Printf("DEBUG: Error: %s\n", mcpResponse.Error)
		if mcpResponse.Count != nil {
			fmt.Printf("DEBUG: Count: %d\n", *mcpResponse.Count)
		}

		// Convert MCP response to legacy format for backwards compatibility
		answer, operation, parameters := convertMCPResponseToLegacyFormat(mcpResponse)

		// Add assistant response to conversation history
		conversationHistory.AddAssistantResponse(answer, operation, parameters)

		// Create updated context with conversation history
		updatedContext := &chatpb.Context{
			UserId:    req.Context.GetUserId(),
			SessionId: req.Context.GetSessionId(),
			Metadata:  conversationHistory.ToContextMetadata(),
		}

		// Debug: Print final response
		fmt.Printf("DEBUG: === FINAL RESPONSE ===\n")
		fmt.Printf("DEBUG: Answer: %s\n", answer)
		fmt.Printf("DEBUG: Operation: %s\n", operation)
		fmt.Printf("DEBUG: Parameters: %+v\n", parameters)

		return &chatpb.ChatResponse{
			Answer:     answer,
			Context:    updatedContext,
			Operation:  operation,
			Parameters: parameters,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return resp.(*chatpb.ChatResponse), nil
}
